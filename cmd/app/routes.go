package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/wallet", app.Transaction)
	mux.HandleFunc("/api/v1/wallets/", app.showBalance)

	return mux
}
