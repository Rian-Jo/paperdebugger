package services

import (
    "context"
    "strings"
    "testing"
)

func TestMockProviderCheckoutAndPortal(t *testing.T) {
    p := NewMockBillingProvider()
    url, cus, err := p.CreateCheckout(context.Background(), "507f1f77bcf86cd799439011", PlanPro, "https://ok", "https://cancel")
    if err != nil { t.Fatalf("CreateCheckout error: %v", err) }
    if !strings.Contains(url, "/checkout/") { t.Fatalf("unexpected checkout url: %s", url) }
    if !strings.HasPrefix(cus, "mock_cus_") { t.Fatalf("unexpected customer id: %s", cus) }

    portal, err := p.CreatePortal(context.Background(), "507f1f77bcf86cd799439011")
    if err != nil { t.Fatalf("CreatePortal error: %v", err) }
    if !strings.Contains(portal, "/portal/") { t.Fatalf("unexpected portal url: %s", portal) }
}

func TestMockProviderParseWebhook(t *testing.T) {
    p := NewMockBillingProvider()
    ev, err := p.ParseWebhook(context.Background(), []byte("abc"), map[string]string{"X-Mock-Event-Type": "subscription_created", "X-Mock-User-Id": "507f1f77bcf86cd799439011"})
    if err != nil { t.Fatalf("ParseWebhook error: %v", err) }
    if ev.ID == "" || ev.Type != "subscription_created" || ev.UserID == "" { t.Fatalf("unexpected event: %+v", ev) }
}

