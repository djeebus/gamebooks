package gorm

type UserStorage struct {
	Key   string `gorm:"primaryKey"`
	Value []byte
}
