package model

type User struct {
	Name		string	`form:"name" binding:"required, gte=2"`
	PassWord	string	`form:"password" binding:"required"`
}

type UpdatePassword struct {
	OldPassword	string	`form:"old_password" binding:"required"`
	NewPassword	string	`form:"new_password" binding:"required"`
}