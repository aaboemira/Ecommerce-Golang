package cards

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/paymentmethod"
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
func (c *Card) RetrivePaymentIntent(id string) (*stripe.PaymentIntent, error) {
	stripe.Key = c.Secret
	pi, err := paymentintent.Get(id, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("iam here")
	return pi, nil
}

func (c *Card) RetrivePaymentMethod(id string) (*stripe.PaymentMethod, error) {
	stripe.Key = c.Secret
	pm, err := paymentmethod.Get(id, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("iam here")
	return pm, nil
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
