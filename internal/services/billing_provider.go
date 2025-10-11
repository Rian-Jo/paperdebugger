package services

import (
    "context"
)

type BillingPlan string

const (
    PlanFree BillingPlan = "free"
    PlanPro  BillingPlan = "pro"
)

type SubscriptionStatus string

const (
    StatusNone     SubscriptionStatus = "none"
    StatusTrialing SubscriptionStatus = "trialing"
    StatusActive   SubscriptionStatus = "active"
    StatusPastDue  SubscriptionStatus = "past_due"
    StatusCanceled SubscriptionStatus = "canceled"
)

type WebhookEvent struct {
    ID     string
    Type   string
    UserID string // optional, provider-dependent
    // Provider-specific payload is omitted for now
}

// BillingProvider abstracts a payment provider like Stripe.
type BillingProvider interface {
    // CreateCheckout starts checkout and returns a redirect URL and optionally a provider customer id
    CreateCheckout(ctx context.Context, userID string, plan BillingPlan, successURL, cancelURL string) (redirectURL string, providerCustomerID string, err error)
    // CreatePortal returns the portal url for managing subscription
    CreatePortal(ctx context.Context, userID string) (portalURL string, err error)
    // ParseWebhook parses raw body and headers, then returns a normalized event
    ParseWebhook(ctx context.Context, body []byte, headers map[string]string) (*WebhookEvent, error)
}
