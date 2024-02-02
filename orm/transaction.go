package orm

type Transaction struct {
	ID    string `gorm:"type:bigint;primary_key"`
	Value string `gorm:"type:longtext"`
	Date  string `gorm:"type:longtext"`
}
