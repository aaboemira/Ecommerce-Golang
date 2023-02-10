package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/virtual-terminal", app.VirtualTerminal)
	mux.Get("/order", app.OrderPage)

	mux.Post("/payment-succeeded", app.PaymentSuccess)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
