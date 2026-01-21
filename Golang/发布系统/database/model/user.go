package model

type User struct {
	Id			int `gorm:"primaryKey"`
	Name 		string
	PassWord	string	`gorm:"column:password"`	// column:password 表示数据库中的字段名
}