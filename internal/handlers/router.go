package handlers

import "github.com/gorilla/mux"

func ProvideRouter(
	transactionHandler *TransactionHttpHandler,
) *mux.Router {
	r := mux.NewRouter()
	router := r.PathPrefix("/v1").Subrouter()
	router.HandleFunc("/transaction", transactionHandler.createTransaction).Methods("POST")
	router.HandleFunc("/transaction/{id}", transactionHandler.getTransactionById).Methods("GET")
	router.HandleFunc("/transaction/send/email", transactionHandler.sendEmail).Methods("POST")
	return router
}
