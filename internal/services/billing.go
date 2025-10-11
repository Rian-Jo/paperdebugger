package services

import (
	"context"
	"strings"
	"time"

	"paperdebugger/internal/libs/cfg"
	"paperdebugger/internal/libs/db"
	"paperdebugger/internal/libs/logger"
	"paperdebugger/internal/libs/shared"
	"paperdebugger/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BillingService struct {
	BaseService
	Provider         BillingProvider
	subsCollection   *mongo.Collection
	eventsCollection *mongo.Collection
}

func NewBillingService(db *db.DB, cfg *cfg.Cfg, logger *logger.Logger) *BillingService {
	base := NewBaseService(db, cfg, logger)
	// For now we only support mock provider; selection by cfg can be added later
	provider := NewMockBillingProvider()
	return &BillingService{
		BaseService:      base,
		Provider:         provider,
		subsCollection:   base.db.Collection((models.Subscription{}).CollectionName()),
		eventsCollection: base.db.Collection((models.BillingEvent{}).CollectionName()),
	}
}

type BillingStatusDTO struct {
	Plan              BillingPlan
	Status            SubscriptionStatus
	CurrentPeriodEnd  time.Time
	CancelAtPeriodEnd bool
}

func parsePlan(p string) BillingPlan {
	switch strings.ToLower(p) {
	case string(PlanPro):
		return PlanPro
	default:
		return PlanFree
	}
}

func parseStatus(s string) SubscriptionStatus {
	switch strings.ToLower(s) {
	case string(StatusTrialing):
		return StatusTrialing
	case string(StatusActive):
		return StatusActive
	case string(StatusPastDue):
		return StatusPastDue
	case string(StatusCanceled):
		return StatusCanceled
	default:
		return StatusNone
	}
}

func (s *BillingService) GetStatus(ctx context.Context, userID bson.ObjectID) (*BillingStatusDTO, error) {
	var sub models.Subscription
	err := s.subsCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&sub)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &BillingStatusDTO{Plan: PlanFree, Status: StatusNone, CurrentPeriodEnd: time.Time{}, CancelAtPeriodEnd: false}, nil
		}
		return nil, err
	}
	return &BillingStatusDTO{
		Plan:              parsePlan(sub.Plan),
		Status:            parseStatus(sub.Status),
		CurrentPeriodEnd:  sub.CurrentPeriodEnd.Time(),
		CancelAtPeriodEnd: sub.CancelAtPeriodEnd,
	}, nil
}

func (s *BillingService) StartCheckout(ctx context.Context, userID bson.ObjectID, plan BillingPlan, successURL, cancelURL string) (string, error) {
	url, providerCustomerID, err := s.Provider.CreateCheckout(ctx, userID.Hex(), plan, successURL, cancelURL)
	if err != nil {
		return "", err
	}
	// Upsert a pending subscription record for tracking
	now := time.Now()
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{
		"provider":             "mock",
		"provider_customer_id": providerCustomerID,
		"plan":                 string(plan),
		"status":               string(StatusTrialing),
		"updated_at":           bson.NewDateTimeFromTime(now),
	}, "$setOnInsert": bson.M{
		"_id":        bson.NewObjectID(),
		"user_id":    userID,
		"created_at": bson.NewDateTimeFromTime(now),
	}}
	upsert := options.UpdateOne().SetUpsert(true)
	_, _ = s.subsCollection.UpdateOne(ctx, filter, update, upsert)
	return url, nil
}

func (s *BillingService) CreatePortal(ctx context.Context, userID bson.ObjectID) (string, error) {
	return s.Provider.CreatePortal(ctx, userID.Hex())
}

// ApplyWebhookEvent updates subscription state based on provider events.
func (s *BillingService) ApplyWebhookEvent(ctx context.Context, ev *WebhookEvent) error {
	if ev == nil {
		return shared.ErrBadRequest("webhook event is nil")
	}
	// idempotency: insert event doc; if exists, skip
	_, err := s.eventsCollection.InsertOne(ctx, models.BillingEvent{ID: bson.NewObjectID(), EventID: ev.ID, Type: ev.Type, Processed: true})
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}
	if ev.UserID == "" {
		return nil
	}
	// naive state change: set active unless canceled
	status := StatusActive
	if ev.Type == "subscription_canceled" {
		status = StatusCanceled
	}
	now := time.Now()
	uid, err := bson.ObjectIDFromHex(ev.UserID)
	if err != nil {
		return err
	}
	_, err = s.subsCollection.UpdateOne(ctx, bson.M{"user_id": uid}, bson.M{"$set": bson.M{
		"status":             string(status),
		"current_period_end": bson.NewDateTimeFromTime(now.Add(30 * 24 * time.Hour)),
		"updated_at":         bson.NewDateTimeFromTime(now),
	}})
	return err
}
