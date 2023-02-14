package main

import (
	"Ecommerce/internal/cards"
	"Ecommerce/internal/models"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type jsonResponse struct {
	OK      bool   `json:"ok""`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Println(err)
	}
	amount, err := strconv.Atoi(payload.Amount)
	fmt.Println(amount)
	if err != nil {
		app.errorLog.Println(err)
	}
	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: payload.Currency,
	}
	okay := true

	pi, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		okay = false
	}
	if okay {
		out, err := json.MarshalIndent(pi, "", "  ")
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		j := jsonResponse{
			OK:      false,
			Message: msg,
			Content: "",
		}
		out, err := json.MarshalIndent(j, "", "  ")
		if err != nil {
			app.errorLog.Println(err)

		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}

//Customers
func (app *application) createCustomer(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
	}
	var customer models.Customer
	err = json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := app.DB.InsertCustomer(customer)
	if err != nil {
		j := jsonResponse{
			OK:      false,
			Message: fmt.Sprintf("error", err),
		}
		out, _ := json.MarshalIndent(j, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}

	j := jsonResponse{
		OK: true,
		ID: id,
	}
	out, err := json.MarshalIndent(j, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

//Orders
func (app *application) createOrder(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
	}
	var order models.Order
	err = json.NewDecoder(r.Body).Decode(&order)

	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := app.DB.InsertOrder(order)
	if err != nil {
		j := jsonResponse{
			OK:      false,
			Message: fmt.Sprintf("error", err),
		}
		out, _ := json.MarshalIndent(j, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		j := jsonResponse{
			OK: true,
			ID: id,
		}
		out, err := json.MarshalIndent(j, "", "  ")
		if err != nil {
			app.errorLog.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}

}

//transactions
func (app *application) createTransaction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
	}
	var transaction models.Transactions
	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := app.DB.InsertTransaction(transaction)
	if err != nil {
		j := jsonResponse{
			OK:      false,
			Message: fmt.Sprintf("error", err),
		}
		out, _ := json.MarshalIndent(j, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}

	j := jsonResponse{
		OK: true,
		ID: id,
	}
	out, err := json.MarshalIndent(j, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

//items
func (app *application) getItemByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	itemID, _ := strconv.Atoi(id)
	item, err := app.DB.GetItem(itemID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	out, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) getAllProducts(w http.ResponseWriter, r *http.Request) {
	items, err := app.DB.GetAllItems()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	out, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
