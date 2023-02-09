package cards

import (
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}
type transaction struct {
	Currency            string
	Amount              int
	TransactionStatusID int
	LastFour            string
	BankReturnCode      string
}

func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(amount, currency)
}
func (c *Card) CreatePaymentIntent(amount int, currency string) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		msg := ""

		if stripErr, ok := err.(*stripe.Error); ok {

			msg = cardErrorMessage(stripErr.Code)
		}
		return nil, msg, err

	}
	return pi, "", nil
}
func cardErrorMessage(code stripe.ErrorCode) string {
	var msg = ""

	switch code {
	case stripe.ErrorCodeCardDeclined:
		msg = "Card was declined"
	case stripe.ErrorCodeExpiredCard:
		msg = "Card is invalid"
	default:
		msg = "Card was declined"
	}
	return msg
}
