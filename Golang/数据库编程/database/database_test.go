package database_test

import (
	"happyladysauce/database"
	"testing"
	"time"
)

func TestDatabase(t *testing.T) {
	// 测试数据库连接
	db, err := database.InitDataBase()
	if err != nil {
		t.Errorf("数据库连接失败: %v", err)
	}
	defer db.Close()

	// database.Insert(db)
	time.Sleep(100 * time.Millisecond)
	database.Replace(db)
	time.Sleep(100 * time.Millisecond)
	database.Query(db)
	time.Sleep(100 * time.Millisecond)
	// database.Update(db)
	time.Sleep(100 * time.Millisecond)
	database.Delete(db)
}