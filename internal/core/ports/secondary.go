package ports

import (
	"stori/internal/core/domain"
)

type TransactionRepository interface {
	Save(transaction domain.Transaction) (domain.Transaction, error)
	GetTransactionById(id int) (domain.Transaction, error)
	GetAll() *[]domain.Transaction
}
