package billing

import (
	"context"

	"paperdebugger/internal/libs/cfg"
	"paperdebugger/internal/libs/contextutil"
	"paperdebugger/internal/libs/logger"
	"paperdebugger/internal/libs/shared"
	"paperdebugger/internal/services"

	billingv1 "paperdebugger/pkg/gen/api/billing/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type BillingServer struct {
	billingv1.UnimplementedBillingServiceServer

	billingService *services.BillingService
	logger         *logger.Logger
	cfg            *cfg.Cfg
}

func NewBillingServer(billingService *services.BillingService, logger *logger.Logger, cfg *cfg.Cfg) billingv1.BillingServiceServer {
	return &BillingServer{billingService: billingService, logger: logger, cfg: cfg}
}

func toProtoStatus(dto *services.BillingStatusDTO) *billingv1.BillingStatus {
	if dto == nil {
		return &billingv1.BillingStatus{Plan: billingv1.BillingPlan_BILLING_PLAN_FREE, Status: billingv1.SubscriptionStatus_SUBSCRIPTION_STATUS_NONE}
	}
	var plan billingv1.BillingPlan
	switch dto.Plan {
	case services.PlanPro:
		plan = billingv1.BillingPlan_BILLING_PLAN_PRO
	default:
		plan = billingv1.BillingPlan_BILLING_PLAN_FREE
	}
	var status billingv1.SubscriptionStatus
	switch dto.Status {
	case services.StatusTrialing:
		status = billingv1.SubscriptionStatus_SUBSCRIPTION_STATUS_TRIALING
	case services.StatusActive:
		status = billingv1.SubscriptionStatus_SUBSCRIPTION_STATUS_ACTIVE
	case services.StatusPastDue:
		status = billingv1.SubscriptionStatus_SUBSCRIPTION_STATUS_PAST_DUE
	case services.StatusCanceled:
		status = billingv1.SubscriptionStatus_SUBSCRIPTION_STATUS_CANCELED
	default:
		status = billingv1.SubscriptionStatus_SUBSCRIPTION_STATUS_NONE
	}
	var ts *timestamppb.Timestamp
	if !dto.CurrentPeriodEnd.IsZero() {
		ts = timestamppb.New(dto.CurrentPeriodEnd)
	}
	return &billingv1.BillingStatus{Plan: plan, Status: status, CurrentPeriodEnd: ts, CancelAtPeriodEnd: dto.CancelAtPeriodEnd}
}

func (s *BillingServer) GetBillingStatus(ctx context.Context, req *billingv1.GetBillingStatusRequest) (*billingv1.GetBillingStatusResponse, error) {
	actor, err := contextutil.GetActor(ctx)
	if err != nil {
		return nil, err
	}
	status, err := s.billingService.GetStatus(ctx, actor.ID)
	if err != nil {
		return nil, err
	}
	return &billingv1.GetBillingStatusResponse{Status: toProtoStatus(status)}, nil
}

func (s *BillingServer) CreateCheckoutSession(ctx context.Context, req *billingv1.CreateCheckoutSessionRequest) (*billingv1.CreateCheckoutSessionResponse, error) {
	actor, err := contextutil.GetActor(ctx)
	if err != nil {
		return nil, err
	}
	var plan services.BillingPlan
	switch req.GetPlan() {
	case billingv1.BillingPlan_BILLING_PLAN_PRO:
		plan = services.PlanPro
	default:
		plan = services.PlanFree
	}
	url, err := s.billingService.StartCheckout(ctx, actor.ID, plan, req.GetSuccessUrl(), req.GetCancelUrl())
	if err != nil {
		return nil, err
	}
	return &billingv1.CreateCheckoutSessionResponse{RedirectUrl: url}, nil
}

func (s *BillingServer) CreatePortalSession(ctx context.Context, req *billingv1.CreatePortalSessionRequest) (*billingv1.CreatePortalSessionResponse, error) {
	actor, err := contextutil.GetActor(ctx)
	if err != nil {
		return nil, err
	}
	url, err := s.billingService.CreatePortal(ctx, actor.ID)
	if err != nil {
		return nil, err
	}
	return &billingv1.CreatePortalSessionResponse{PortalUrl: url}, nil
}

func (s *BillingServer) HandleWebhook(ctx context.Context, req *billingv1.HandleWebhookRequest) (*billingv1.HandleWebhookResponse, error) {
	// This method is unauthenticated; parse raw body + headers
	headers := map[string]string{}
	for k, v := range req.GetHeaders() {
		headers[k] = v
	}
	ev, err := s.billingService.Provider.ParseWebhook(ctx, req.GetBody(), headers)
	if err != nil {
		return nil, shared.ErrBadRequest(err)
	}
	if err := s.billingService.ApplyWebhookEvent(ctx, ev); err != nil {
		return nil, err
	}
	return &billingv1.HandleWebhookResponse{Ok: true}, nil
}
