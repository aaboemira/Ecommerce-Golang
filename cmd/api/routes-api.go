package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-type", "X-Csrf-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	mux.Post("/api/payment-intent", app.GetPaymentIntent)

	mux.Post("/api/customers", app.createCustomer)

	mux.Post("/api/orders", app.createOrder)

	mux.Post("/api/transactions", app.createTransaction)

	mux.Get("/api/getItem/{id}", app.getItemByID)
	mux.Get("/api/products", app.getAllProducts)

	return mux
}
