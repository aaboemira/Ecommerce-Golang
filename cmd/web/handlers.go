package main

import (
	"Ecommerce/internal/cards"
	"Ecommerce/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

func (app *application) appHome(w http.ResponseWriter, r *http.Request) {

	if err := app.renderTemplate(w, r, "home", &templateData{}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

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
	card := cards.Card{
		Secret: app.config.stripe.secret,
		Key:    app.config.stripe.key,
	}
	pi, err := card.RetrivePaymentIntent(r.Form.Get("payment_intent"))
	if err != nil {
		app.errorLog.Fatal(err)
	}
	pm, err := card.RetrivePaymentMethod(r.Form.Get("payment_method"))

	customer := models.Customer{
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Email:     r.Form.Get("email"),
	}
	customerID, err := app.insertCustomer(customer)
	if err != nil {
		app.errorLog.Println("Couldn't insert", err)
	}

	amount, err := strconv.Atoi(r.Form.Get("payment_amount"))
	if err != nil {
		app.errorLog.Println(err)
	}
	currency := r.Form.Get("payment_currency")
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))
	lastFour := pm.Card.Last4
	expYear := pm.Card.ExpYear
	expMonth := pm.Card.ExpMonth
	bankCode := pi.LatestCharge.ID
	transaction := models.Transactions{
		Amount:         amount,
		Currency:       currency,
		LastFour:       lastFour,
		BankReturnCode: bankCode,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	trxID, err := app.insertTransaction(transaction)
	if err != nil {
		app.errorLog.Println(err)
	}
	itemID, _ := strconv.Atoi(r.Form.Get("item_id"))
	order := models.Order{
		TransactionID: trxID,
		ItemID:        itemID,
		StatusID:      1,
		Quantity:      quantity,
		Amount:        amount,
		CustomerID:    customerID,
		UpdatedAt:     time.Now(),
		CreatedAt:     time.Now(),
	}
	orderID, err := app.insertOrder(order)
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
	data["lastFour"] = lastFour
	data["expYear"] = expYear
	data["expMonth"] = expMonth
	data["bankCode"] = bankCode
	data["customerID"] = customerID
	data["transactionID"] = trxID
	data["transactionID"] = orderID

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

func (app *application) insertCustomer(cst models.Customer) (int, error) {

	var response models.ApiResponse
	api := fmt.Sprintf(app.config.api + "/api/customers")

	postBody, _ := json.Marshal(map[string]string{
		"first_name": cst.FirstName,
		"last_name":  cst.LastName,
		"email":      cst.Email,
	})
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(api, "application/json", reqBody)
	if err != nil {
		app.errorLog.Println(err)
		return 0, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		app.errorLog.Println(" an error occurred, please try again", err)
		return 0, err
	}

	return response.ID, nil
}

func (app *application) insertTransaction(trx models.Transactions) (int, error) {

	var response models.ApiResponse
	api := fmt.Sprintf(app.config.api + "/api/transactions")
	fmt.Println(trx.Amount)
	postBody, _ := json.Marshal(map[string]any{
		"amount":                trx.Amount,
		"currency":              trx.Currency,
		"last_four":             trx.LastFour,
		"bank_return_code":      trx.BankReturnCode,
		"transaction_status_id": 1,
		"created_at":            trx.CreatedAt,
		"updated_at":            trx.UpdatedAt,
	})
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(api, "application/json", reqBody)
	if err != nil {
		app.errorLog.Println(err)
		return 0, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	if !response.OK {
		return 0, errors.New(response.MSG)
	}

	return response.ID, nil
}

func (app *application) insertOrder(ord models.Order) (int, error) {

	var response models.ApiResponse
	api := fmt.Sprintf(app.config.api + "/api/orders")
	fmt.Println(ord.Amount)
	postBody, _ := json.Marshal(map[string]any{
		"item_id":        ord.ItemID,
		"transaction_id": ord.TransactionID,
		"status_id":      ord.StatusID,
		"quantity":       ord.Quantity,
		"amount":         ord.Amount,
		"customer_id":    ord.CustomerID,
		"created_at":     ord.CreatedAt,
		"updated_at":     ord.UpdatedAt,
	})
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(api, "application/json", reqBody)
	if err != nil {
		app.errorLog.Println(err)
		return 0, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}
	if !response.OK {
		return 0, errors.New(response.MSG)
	}
	fmt.Println(response)
	return response.ID, nil
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
