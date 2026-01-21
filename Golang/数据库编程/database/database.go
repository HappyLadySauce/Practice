package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDataBase() (*sql.DB, error) {
	// 连接数据库
	mysql, err := sql.Open("mysql", "test:test@tcp(100.100.100.5:3306)/test?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")
	if err != nil {
		return nil, err
	}
	return mysql, nil
}

// Insert 数据插入
func Insert(db *sql.DB) {
	// 进行数据的插入
	// MD5 是一个加密算法，将字符串加密为一个 32 位的十六进制数
	// 这里使用 MD5 加密密码，是为了安全起见，不存储明文密码
	res, err := db.Exec("insert into user (name, password) values ('test1', 'MD5(test1)'), ('test2', 'MD5(test2)')")
	if err != nil {
		log.Printf("%v insert data err: %v\n", db.Driver(), err)
	}

	// 获取第一次插入数据的 ID
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("get insert last id err: %v\n", err)
	}
	log.Printf("after insert last id %d\n", lastId)

	// 获取总插入数据的行数
	rowsId, err := res.RowsAffected()
	if err != nil {
		log.Printf("get insert rows id err: %v\n", err)
	}
	log.Printf("after insert last id %d\n", rowsId)
}

// Replace 数据覆盖
func Replace(db *sql.DB) {
	// 进行数据的插入
	res, err := db.Exec("replace into user (name, password) values ('test1', 'MD5(test1)'), ('test2', 'MD5(test2)')")
	if err != nil {
		log.Printf("%v replace data err: %v\n", db.Driver(), err)
	}

	// 获取第一次插入数据的 ID
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("get insert last id err: %v\n", err)
	}
	log.Printf("after insert last id %d\n", lastId)

	// 获取总插入数据的行数
	rowsId, err := res.RowsAffected()
	if err != nil {
		log.Printf("get replace rows id err: %v\n", err)
	}
	log.Printf("after replace last id %d\n", rowsId)
}

// Delete 删除数据
func Delete(db *sql.DB) {
	res, err := db.Exec("delete from user where id > 10")
	if err != nil {
		log.Printf("delete user table id err: %v\n", err)
	}

	// 获取第一次插入数据的 ID
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("get insert last id err: %v\n", err)
	}
	log.Printf("after insert last id %d\n", lastId)

	// 获取总插入数据的行数
	rowsId, err := res.RowsAffected()
	if err != nil {
		log.Printf("get delete rows id err: %v\n", err)
	}
	log.Printf("after delete last id %d\n", rowsId)
}

// 更新数据
func Update(db *sql.DB) {
	res, err := db.Exec("update user set name = 'test1' where id = 1, password = 'MD5(test1)'")
	if err != nil {
		log.Printf("update user table id err: %v\n", err)
	}

	// 获取第一次插入数据的 ID
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("get insert last id err: %v\n", err)
	}
	log.Printf("after insert last id %d\n", lastId)

	// 获取总插入数据的行数
	rowsId, err := res.RowsAffected()
	if err != nil {
		log.Printf("get update rows id err: %v\n", err)
	}
	log.Printf("after update last id %d\n", rowsId)
}

// 查询数据
func Query(db *sql.DB) {
	rows, err := db.Query("select id, name from user")
	if err != nil {
		log.Printf("")
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var password string
		err := rows.Scan(&name, &password)
		if err != nil {
			log.Printf("scanf user table err: %v\n", err)
		}
		log.Printf("query data: name:%v password:%v from user table.", name, password)
	}
}

// 数据库事务
func Transaction(db *sql.DB) {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Printf("begin transaction err: %v\n", err)
	}
	_, err = tx.Exec("insert into user (name, password) values ('test3', 'MD5(test3)')")
	if err != nil {
		log.Printf("insert user table err: %v\n", err)
	}

	_, err = tx.Exec("insert into user (name, password) values ('test4', 'MD5(test4)')")
	if err != nil {
		log.Printf("insert user table err: %v\n", err)
	}

	tx.Commit()
}