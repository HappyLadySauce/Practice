package sqlbuilder

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"

	gsb "github.com/huandu/go-sqlbuilder"
)

// SQL 插入
func SqlInsert() {
	insertBuilder := gsb.NewInsertBuilder()
	insertBuilder = insertBuilder.InsertInto("user").Cols("name", "password").Values("testbuilder", "MD5(testbuilder)")
	for i :=0; i < 3; i++ {
		randName := rand.Text()
		insertBuilder = insertBuilder.Values(randName[:4], "builderpassword")
	}
	sql, args := insertBuilder.Build()	// 构建 sql 语句
	log.Println(sql)	// INSERT INTO user (name, password) VALUES (?, ?), (?, ?), (?, ?), (?, ?)
	log.Println(args...)
}

// SQL 删除
func SqlDelete() {
	deleteBuilder := gsb.NewDeleteBuilder()
	deleteBuilder = deleteBuilder.DeleteFrom("user").Where(
		deleteBuilder.Equal("name", "test1"),
	)
	sql, args := deleteBuilder.Build()
	log.Println(sql)
	log.Println(args...)
}

// 数据读取
func SqlRead() {
	selectBuilder := gsb.NewSelectBuilder()
	selectBuilder.SetFlavor(gsb.MySQL)	// 不同的数据库sql语法会有差异,用过Flavor指定使用哪种数据库的语法
	selectBuilder = selectBuilder.Select("name").From("user")
	sql, args := selectBuilder.Build()
	log.Println(sql)
	log.Println(args...)
}


// 数据更新
func SqlUpdate() {
	updateBuilder := gsb.NewUpdateBuilder()
	updateBuilder = updateBuilder.Update("user").Set("name", "testupdate").Where(
		updateBuilder.Equal("name", "testbuilder"),
	)
	sql, args := updateBuilder.Build()
	log.Println(sql)
	log.Println(args...)
}

func InitDataBase() (*sql.DB, error) {
	// 连接数据库
	mysql, err := sql.Open("mysql", "test:test@tcp(100.100.100.2:3306)/test?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")
	if err != nil {
		return nil, err
	}
	return mysql, nil
}

// 大量数据插入
func MassInsertStmt(db *sql.DB) {
	insertBuilder := gsb.NewInsertBuilder()
	insertBuilder = insertBuilder.InsertInto("user")
	insertBuilder = insertBuilder.Cols("name", "password")
	randTest := rand.Text()
	insertBuilder = insertBuilder.Values(
		randTest[:8],
		"TestPassword",
	)

	sql, args := insertBuilder.Build()
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Printf("stmt create err: %v\n", err)
	}
	stmt.Exec(args...)
	stmt.Close()
}

// 定义多表查询结构体
type User struct {
	// 公共主键
	Id 		int
	// 表1中查找name
	Name 	string
	// 表2中查找sex
	Sex		string
}

// 多表联合查
func Query(db *sql.DB) map[int]User {
	rows, err := db.Query("select id, name from user")
	if err != nil {
		log.Printf("err: %v\n", err)
	}
	defer rows.Close()

	rect := make(map[int]User, 10)
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		rect[id] = User{
			Id: id,
			Name: name,
		}
	}
	return rect
}

type UserPage struct {
	Id		int
	Name	string
}

// 分页查询
func QueryByPage(db *sql.DB, pageSize, page int) (total int, data []*UserPage) {
	rows, err := db.Query("select conut(*) from user")
	if err != nil {
		log.Printf("query total from user err: %v\n", err)
		return 0, nil
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&total)

	offset := pageSize * (page - 1)
	rows2, err := db.Query(fmt.Sprint("select id, name from user limit %d, %d", offset, pageSize))
	if err != nil {
		log.Printf("query data from user err: %v\n", err)
		return total, nil
	}
	for rows2.Next() {
		var id int
		var name string
		rows2.Scan(&id, &name)
		data = append(data, &UserPage{Id: id, Name: name})
	}
	return
}












