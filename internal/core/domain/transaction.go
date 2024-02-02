package domain

type Transaction struct {
	Id    string `json:"id"`
	Value string `json:"value"`
	Date  string `json:"date"`
}

func NewTransaction(id, date, value string) Transaction {
	return Transaction{
		Id:    id,
		Value: value,
		Date:  date,
	}
}
