package domain

type Response struct {
	TotalBalance       float64        `json:"total_balance"`
	NumberTransactions map[string]int `json:"number_of_transactions"`
	AverageDebit       float64        `json:"average_debit"`
	AverageCredit      float64        `json:"average_credit"`
}

func NewResponse(total, debit, credit float64, numberTransactions map[string]int) Response {
	return Response{
		TotalBalance:       total,
		NumberTransactions: numberTransactions,
		AverageDebit:       debit,
		AverageCredit:      credit,
	}
}
