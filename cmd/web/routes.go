package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(SessionLoad)
	mux.Get("/", app.appHome)
	mux.Get("/virtual-terminal", app.VirtualTerminal)
	mux.Get("/order/{id}", app.ChargeProduct)
	mux.Get("/shop", app.ProductsPageHandler)
	mux.Post("/payment-succeeded", app.PaymentSuccess)
	//mux.Post("/checkout", app.ShoppingCartHandler)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
