package usecases

import (
	"errors"
	"fmt"
	"math"
	"stori/internal/core/domain"
	"stori/internal/core/ports"
	"strconv"
	"strings"
)

func ProvideTransactionUseCase(
	transactionRepository ports.TransactionRepository,
) ports.TransactionUseCase {
	return &transactionUseCase{
		transactionRepository: transactionRepository,
	}
}

type transactionUseCase struct {
	transactionRepository ports.TransactionRepository
}

func (u transactionUseCase) SendEmail(to, from, pass string) (*domain.Response, error) {
	transactionInfo := u.transactionRepository.GetAll()
	calculations, err := getCalculations(transactionInfo)
	if err != nil {
		errors.New(fmt.Sprintf("%v %v", "hubo un error en los calculos", err))
	}
	err = sendEmail(calculations, to, from, pass)
	if err != nil {
		return nil, err
	}
	return calculations, nil
}

func (u transactionUseCase) CreateTransaction(id, date, value string) (domain.Transaction, error) {
	account := domain.NewTransaction(id, date, value)
	return u.transactionRepository.Save(account)
}

func (u transactionUseCase) GetTransactionById(id int) (domain.Transaction, error) {
	return u.transactionRepository.GetTransactionById(id)
}

func getCalculations(info *[]domain.Transaction) (*domain.Response, error) {
	if info == nil {
		return nil, errors.New("does not exists data")
	}
	averageDebit := 0.0
	averageCredit := 0.0
	totalCredit := 0.0
	totalDebit := 0.0
	quatityDebit := 0.0
	quantityCredit := 0.0
	indexMonth := 0
	numberOfTransactions := make(map[string]int, 12)
	for _, value := range *info {
		if value.Value[0:1] == "+" {
			// credit
			total, _ := strconv.ParseFloat(value.Value[1:len(value.Value)], 8)
			totalCredit = total + totalCredit
			quantityCredit++
		} else if value.Value[0:1] == "-" {
			// debit
			total, _ := strconv.ParseFloat(value.Value[1:len(value.Value)], 8)
			totalDebit = total + totalDebit
			quatityDebit++
		}
		if position := strings.Index(value.Date, "/"); position != -1 {
			indexMonth = position
		}
		month := value.Date[0:indexMonth]
		numberOfTransactions[month]++
	}
	averageDebit = totalDebit / quatityDebit
	averageCredit = totalCredit / quantityCredit
	response := domain.NewResponse(roundFloat(totalCredit-totalDebit, 2), roundFloat(averageDebit*-1, 2), roundFloat(averageCredit, 2), numberOfTransactions)
	return &response, nil
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
