package cfg

import (
	"os"

	"github.com/joho/godotenv"
)

type Cfg struct {
    OpenAIAPIKey  string
    JwtSigningKey string

    MongoURI string

    BillingProvider    string
    StripeAPIKey       string
    StripeWebhookSecret string
}

var cfg *Cfg

func GetCfg() *Cfg {
    _ = godotenv.Load()
    cfg = &Cfg{
        OpenAIAPIKey:  os.Getenv("OPENAI_API_KEY"),
        JwtSigningKey: os.Getenv("JWT_SIGNING_KEY"),
        MongoURI:      mongoURI(),
        BillingProvider: os.Getenv("PD_BILLING_PROVIDER"),
        StripeAPIKey: os.Getenv("STRIPE_API_KEY"),
        StripeWebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
    }

    return cfg
}

func mongoURI() string {
	val := os.Getenv("PD_MONGO_URI")
	if val != "" {
		return val
	}

	return "mongodb://localhost:27017"
}
