package adapters

import (
	"gorm.io/gorm"
	"stori/internal/core/domain"
	"stori/internal/core/ports"
	"stori/orm"
)

func ProvideTransactionRepository(db *gorm.DB) ports.TransactionRepository {
	return &mySQLTransactionRepository{
		db: db,
	}
}

type mySQLTransactionRepository struct {
	db *gorm.DB
}

func (u mySQLTransactionRepository) Save(transaction domain.Transaction) (domain.Transaction, error) {
	ormAccount := orm.Transaction{
		ID:    transaction.Id,
		Value: transaction.Value,
		Date:  transaction.Date,
	}
	if err := u.db.Create(&ormAccount).Error; err != nil {
		return domain.Transaction{}, err
	}

	return transaction, nil
}

func (u mySQLTransactionRepository) GetTransactionById(id int) (domain.Transaction, error) {
	var orm orm.Transaction
	if err := u.db.First(&orm, "id = ?", id).Error; err != nil {
		return domain.Transaction{}, err
	}
	return domain.Transaction{
		Id:    orm.ID,
		Value: orm.Value,
		Date:  orm.Date,
	}, nil
}

func (u mySQLTransactionRepository) GetAll() *[]domain.Transaction {
	var orm []orm.Transaction
	if err := u.db.Raw("select * from `transactions`").Scan(&orm).Error; err != nil {
		return nil
	}
	result := make([]domain.Transaction, len(orm))
	for i := 0; i < len(orm); i++ {
		result[i] = domain.Transaction{
			Id:    orm[i].ID,
			Value: orm[i].Value,
			Date:  orm[i].Date,
		}
	}
	return &result
}
