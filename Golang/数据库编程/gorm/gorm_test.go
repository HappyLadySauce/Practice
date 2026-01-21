package gorm_test

import (
	"testing"
	mygorm "happyladysauce/gorm"
)

func TestGorm(t *testing.T) {
	// mygorm.GormQuickStart()
	db := mygorm.InitGormDB()
	mygorm.CreateData(db)
}