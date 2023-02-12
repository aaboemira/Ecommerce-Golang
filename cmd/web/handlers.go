package main

import (
	"Ecommerce/internal/models"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
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
func (app *application) ProductsPageHandler(w http.ResponseWriter, r *http.Request) {
	api := fmt.Sprintf(app.config.api + "/api/products")
	resp, err := http.Get(api)
	if err != nil {
		app.errorLog.Println(err)
	}
	defer resp.Body.Close()
	var items []models.Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		app.errorLog.Println(" an error occurred, please try again", err)
	}
	fmt.Println(items)
	data := make(map[string]interface{})
	data["products"] = items
	if err := app.renderTemplate(w, r, "order", &templateData{Data: data}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

//func (app *application) ShoppingCartHandler(w http.ResponseWriter, r *http.Request) {
//	err := r.ParseForm()
//	if err != nil {
//		app.errorLog.Println(err)
//	}
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		panic(err)
//	}
//	log.Println(string(body))
//	data := make(map[string]interface{})
//
//	if err := app.renderTemplate(w, r, "order", &templateData{Data: data}, "stripe-js"); err != nil {
//		app.errorLog.Println(err)
//	}
//}
func (app *application) ChargeProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	api := fmt.Sprintf(app.config.api + "/api/getItem/" + id)
	resp, err := http.Get(api)
	if err != nil {
		app.errorLog.Println(err)
	}
	defer resp.Body.Close()
	var item models.Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		app.errorLog.Println(" an error occurred, please try again", err)
	}
	data := make(map[string]interface{})
	data["item"] = item
	if err := app.renderTemplate(w, r, "chargeProduct", &templateData{Data: data}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}
