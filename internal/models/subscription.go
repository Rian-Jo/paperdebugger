package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Subscription struct {
    BaseModel `bson:",inline"`

    UserID                  bson.ObjectID `bson:"user_id"`
    Provider               string        `bson:"provider"`
    ProviderCustomerID     string        `bson:"provider_customer_id"`
    ProviderSubscriptionID string        `bson:"provider_subscription_id"`

    Plan   string `bson:"plan"`   // free, pro
    Status string `bson:"status"` // none, trialing, active, past_due, canceled

    CurrentPeriodEnd   bson.DateTime `bson:"current_period_end"`
    CancelAtPeriodEnd  bool          `bson:"cancel_at_period_end"`
}

func (Subscription) CollectionName() string { return "subscriptions" }

type BillingEvent struct {
    ID        bson.ObjectID `bson:"_id"`
    EventID   string        `bson:"event_id,unique"`
    Type      string        `bson:"type"`
    Processed bool          `bson:"processed"`
}

func (BillingEvent) CollectionName() string { return "billing_events" }

