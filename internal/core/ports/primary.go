package ports

import (
	"stori/internal/core/domain"
)

type TransactionUseCase interface {
	CreateTransaction(id, date, value string) (domain.Transaction, error)
	GetTransactionById(id int) (domain.Transaction, error)
	SendEmail(to, from, pass string) (*domain.Response, error)
}
