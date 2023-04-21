package payment

import (
	"encoding/json"
	"fmt"

	stripe "github.com/stripe/stripe-go/v74"

	session "github.com/stripe/stripe-go/v74/checkout/session"
	webhook "github.com/stripe/stripe-go/v74/webhook"
)

type Repository struct {
	stripeKey string
	siteName  string
	siteHost  string
}

func NewRepository(stripeKey, siteName, siteHost string) *Repository {
	return &Repository{
		stripeKey: stripeKey,
		siteName:  siteName,
		siteHost:  siteHost,
	}
}

func PlanTypeAndDurationToDescription(planType string, planDuration int64) string {
	return fmt.Sprintf("Plan x %d months", planDuration)
}

func (r *Repository) CreateSession(priceID string) (*stripe.CheckoutSession, error) {
	// Set the API key
	stripe.Key = r.stripeKey

	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(fmt.Sprintf("https://%s/x/j/p/1", r.siteHost)),
		CancelURL:  stripe.String(fmt.Sprintf("https://%s/x/j/p/0", r.siteHost)),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
	}

	session, err := session.New(params)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func HandleCheckoutSessionComplete(body []byte, endpointSecret, stripeSig string) (*stripe.CheckoutSession, error) {
	event, err := webhook.ConstructEvent(body, stripeSig, endpointSecret)
	if err != nil {
		return nil, fmt.Errorf("error verifying webhook signature: %v", err)
	}
	// Handle the checkout.session.completed event
	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			return nil, fmt.Errorf("error parsing webhook JSON: %v", err)
		}
		return &session, nil
	}
	return nil, nil
}
