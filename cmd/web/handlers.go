package main

import (
	"Ecommerce/internal/models"
	"net/http"
)

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["StripeKey"] = app.config.stripe.key
	if err := app.renderTemplate(w, r, "terminal", &templateData{StringMap: data}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}
func (app *application) PaymentSuccess(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
	}
	data := make(map[string]interface{})
	data["name"] = r.Form.Get("cardholder-name")
	data["email"] = r.Form.Get("cardholder-email")
	data["paymentIntent"] = r.Form.Get("payment_intent")
	data["paymentMethod"] = r.Form.Get("payment_method")
	data["paymentAmount"] = r.Form.Get("payment_amount")
	data["paymentCurrency"] = r.Form.Get("payment_currency")

	err = app.renderTemplate(w, r, "paymentSuccess", &templateData{Data: data})
	if err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) OrderPage(w http.ResponseWriter, r *http.Request) {
	fidgets := &models.Item{ID: 1, Name: "Fidget Spinner", Description: "Fidget spinner for kids", InventoryLevel: 100, Price: 1000}
	data := make(map[string]interface{})
	data["item"] = fidgets

	if err := app.renderTemplate(w, r, "order", &templateData{Data: data}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}
