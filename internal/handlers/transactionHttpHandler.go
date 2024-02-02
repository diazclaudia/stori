package handlers

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"stori/internal/core/ports"
	"strconv"
)

type body struct {
	CSV string `json:"csv"`
}

type createReq struct {
	Id    string `json:"id"`
	Value string `json:"value"`
	Date  string `json:"date"`
}

type Email struct {
	mail string
	pass string
}

func ProvideTransactionHttpHandler(
	uc ports.TransactionUseCase,
) *TransactionHttpHandler {
	return &TransactionHttpHandler{
		uc: uc,
	}
}

type TransactionHttpHandler struct {
	uc ports.TransactionUseCase
}

func (t *TransactionHttpHandler) sendEmail(w http.ResponseWriter, r *http.Request) {
	type emailReq struct {
		To   string `json:"to"`
		From string `json:"from"`
		Pass string `json:"pass"`
	}
	var req emailReq
	if err := decodeJSONBody(w, r, &req); err != nil {
		var mr *malformedRequestError
		log.Println(err)
		if errors.As(err, &mr) {
			http.Error(w, mr.Message, mr.Status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	fmt.Println("body ", req)
	transaction, err := t.uc.SendEmail(req.To, req.From, req.Pass)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Could not send email")
	}
	writeJson(w, transaction)
}

func (t *TransactionHttpHandler) getTransactionById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	transaction, err := t.uc.GetTransactionById(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Could not get transaction by id")
	}
	writeJson(w, transaction)
}

func (t *TransactionHttpHandler) createTransaction(w http.ResponseWriter, r *http.Request) {
	var req body
	if err := decodeJSONBody(w, r, &req); err != nil {
		var mr *malformedRequestError
		if errors.As(err, &mr) {
			http.Error(w, mr.Message, mr.Status)
		} else {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	f, err := os.Open(req.CSV)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	shoppingList := createTransactionList(data)
	jsonData, err := json.MarshalIndent(shoppingList, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
	for _, data := range shoppingList {
		user, err := t.uc.CreateTransaction(data.Id, data.Date, data.Value)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Could not create transaction")
		}
		writeJson(w, user)
	}
}

func createTransactionList(data [][]string) []createReq {
	// convert csv lines to array of structs
	var list []createReq
	for i, line := range data {
		if i > 0 { // omit header line
			var rec createReq
			for j, field := range line {
				if j == 0 {
					rec.Id = field
				} else if j == 1 {
					rec.Date = field
				} else if j == 2 {
					rec.Value = field
				}
			}
			list = append(list, rec)
		}
	}
	return list
}
