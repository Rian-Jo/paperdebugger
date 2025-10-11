package services

import (
    "context"
    "crypto/sha1"
    "encoding/hex"
    "fmt"
)

type mockProvider struct{}

func NewMockBillingProvider() BillingProvider { return &mockProvider{} }

func (m *mockProvider) CreateCheckout(ctx context.Context, userID string, plan BillingPlan, successURL, cancelURL string) (string, string, error) {
    // Deterministic pseudo token
    h := sha1.Sum([]byte(userID + string(plan)))
    token := hex.EncodeToString(h[:])
    return fmt.Sprintf("https://mock.billing/checkout/%s", token), fmt.Sprintf("mock_cus_%s", token[:12]), nil
}

func (m *mockProvider) CreatePortal(ctx context.Context, userID string) (string, error) {
    h := sha1.Sum([]byte(userID))
    token := hex.EncodeToString(h[:])
    return fmt.Sprintf("https://mock.billing/portal/%s", token), nil
}

func (m *mockProvider) ParseWebhook(ctx context.Context, body []byte, headers map[string]string) (*WebhookEvent, error) {
    // For mock, accept everything and synthesize an event id from body
    h := sha1.Sum(body)
    id := hex.EncodeToString(h[:])
    // Simple convention: the first header X-Mock-Event-Type decides type; fallback to subscription_updated
    t := headers["X-Mock-Event-Type"]
    if t == "" { t = "subscription_updated" }
    uid := headers["X-Mock-User-Id"]
    return &WebhookEvent{ID: id, Type: t, UserID: uid}, nil
}
