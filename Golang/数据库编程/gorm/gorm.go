package gorm

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/hints"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// type User struct {
// 	Id			int
// 	Username	string
// 	Password	string
// }

// func GormQuickStart() {
// 	// 连接数据
// 	host := "100.100.100.5"
// 	port := 3306
// 	dbname := "test"
// 	user := "test"
// 	pass := "test"

// 	// data source name DSN 数据连接字符串
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
// 	db, err := gorm.Open(mysql.Open(dsn), nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 写入数据
// 	instance1 := User{
// 		Id: 1,
// 		Username: "HappyLadySauce",
// 		Password: "MD5(HappyladySauce)",
// 	}
// 	db.Create(&instance1)

// 	// 读取数据
// 	var instance2 User
// 	db.Find(&instance2)	// 读取全表
// 	fmt.Println(instance2)
// }

func InitGormDB() *gorm.DB {
	host := "100.100.100.5"
	port := 3306
	dbname := "test"
	user := "test"
	pass := "test"


	// 确保日志目录存在
	logDir := filepath.Join(".", "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("创建日志目录失败: %v", err))
	}
	logFile, err := os.OpenFile(filepath.Join(logDir, "gorm.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("打开日志文件失败: %v", err))
	}
	newLogger := logger.New(
		log.New(logFile, "\r\n", log.LstdFlags),	// io writer, 可以输出到文件, 也可以输出到终端
		logger.Config {
			SlowThreshold: 500 * time.Millisecond,	// 耗时超过此值定为慢查询
			LogLevel: logger.Info,	// 定义日志级别
			ParameterizedQueries: true,	// true 表示SQL日志里不包含参数
			Colorful: false,	// 禁用颜色
		},
	)

	// data source name DSN 数据连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,	// 使用Stmt优化性能
		SkipDefaultTransaction: true,	// 跳过默认的事务
		NamingStrategy: schema.NamingStrategy{	// 命名策略
			// TablePrefix: "t_",	// 增加表前缀
			// SingularTable: true,	// 表名映射时不加复数, 仅是驼峰-->蛇形
			// NoLowerCase: true,		// 不将字段名转换为小写
		},
		Logger: newLogger,
		DryRun: false,	// true 代表生成 SQL 但不执行
		DisableAutomaticPing: false,	// 在完成初始化之后 GORM 会自动 ping 数据库以检测数据库可用性
		DisableNestedTransaction: true,	// 禁用嵌套事务, 一般不会用到, 并提升性能
	})
	if err != nil {
		panic(err)
	}

	// 获取 gorm 连接池
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("获取数据库连接失败", "error", err)
	}
	sqlDB.SetMaxIdleConns(10)	// 设置连接池空闲连接的数量上限
	sqlDB.SetMaxOpenConns(100)	// 设置连接池的最大连接
	sqlDB.SetConnMaxIdleTime(time.Hour)	// 设置一个连接的最大生命周期
	return db
}

// 返回表名
func (*User) TableName() string {
	return "user"
}

// gorm.Model 的定义
// 默认情况下, GORM 使用 ID 作为主键, 使用结构体名的蛇形复数作为表名, 字段名的蛇形作为列名
// 不建议使用 Migrator(数据迁移)功能, 表的维护又DBA负责, 而不是开发人员
type User struct {	// 默认的表名: users, 使用 TableName 方法指定表名
	// ID		int	// 名为ID的字段默认为主键
	Id			int	`gorm:"primaryKey;column:id"`	// 显示指定主键, 显示指定表里的列名
	UserId		int	`gorm:"column:uid"`		// 显示指定列名
	Degree		string
	Keywords	[]string	`gorm:"json"`
	CreatedAt	time.Time	`gorm:"column:create_time"`	// gorm 会自动处理 CreatedAt 和 UpdatedAt
	UpdatedAt	time.Time	`gorm:"column:update_time"`
	Gender		string
	City		string
	Province	string	`gorm:"-"`	// 表示表里没有这一列, 但是这个结构体需要它
}

func CreateData(db *gorm.DB) {
	user := &User{UserId: rand.IntN(100), Degree: "本科", Gender: "男", City: "上海", Keywords: []string{"你好！", "世界"}}
	result := db.Create(&user)	// 必须传入指针, 因为需要给 user 结构体赋值
	if result.Error != nil {
		slog.Error("插入记录失败", "error", result.Error)
	}
	fmt.Printf("record id is %d\n", user.Id)
	fmt.Printf("影响行数%d\n", result.RowsAffected)

	// 会话模式, 开辟一个数据库连接会话, 并指定其特殊配置
	// tx := db.Session(&gorm.Session{SkipHooks: true})	// 不执行钩子
	// tx := db.Session(&gorm.Session{DryRUn: true})	// 生成SQL,但不执行
	// 一次性插入多条数据
	user1 := user
	user1.Id = 0	// 将 Id 主键设置为0,使用数据库的自增ID
	user1.UserId = rand.IntN(100)
	user2 := user
	user2.Id = 0
	user2.UserId = rand.IntN(100)
	users := []*User{user1, user2}
	result = db.Create(&users)
	fmt.Printf("影响行数%d\n", result.RowsAffected)

	// 当数据量太大时使用分配插入
	batchSize := 1
	user3 := user
	user3.Id = 0
	user3.UserId = rand.IntN(100)
	user4 := user
	user4.Id = 0
	user4.UserId = rand.IntN(100)
	// 一个批次一条SQL,且所有批次被放到一个事物中来执行
	db.CreateInBatches([]*User{user3, user4}, batchSize)
}

// 创建 map 进行数据插入
func CreateByMap(db *gorm.DB) {
	db.Model(User{}).Create(map[string]any {
		"uid": rand.IntN(100), "degree": "本科", "gender": "男", "city": "上海",
	})

	// 一次性插入多条数据使用 []map[string]any
	db.Model(User{}).Create([]map[string]any {
		{"uid": rand.IntN(100), "degree": "本科", "gender": "男", "city": "上海"},
		{"uid": rand.IntN(100), "degree": "本科", "gender": "男", "city": "上海"},
	})
}

func Delete(db *gorm.DB) {
	// 通过 Where 进行查询筛选, 通过 "degree==?" 防止 sql 注入攻击
	tx := db.Where("degree=?", "硕士").Delete(User{})
	fmt.Printf("删除%d行\n", tx.RowsAffected)

	var user User = User{Id: 10}
	db.Delete(user)	// 暗含Where条件id=10

	db.Delete(User{}, 1)	// 暗含id = 1
	db.Delete(User{}, []int{1,2,3})	// 暗含id = 1,2,3
}

// 更新
func Save(db *gorm.DB) {
	user := User{UserId: rand.IntN(100), Degree: "本科", Gender: "男", City: "上海", Keywords: []string{"你好！", "世界"}}
	db.Save(&user) // 主键为0时, Save 相当于 Create

	user1 := user
	user1.Degree = "硕士"
	db.Save(&user1)
}

func Update(db *gorm.DB) {
	// 根据 map 更新
	tx := db.Model(&User{}).Where("city=?", "北京").Updates(map[string]any{"degree": "硕士", "gender": "男"})
	fmt.Printf("更新了%d行\n", tx.RowsAffected)

	// 根据结构体更新, 只会更新结构体中的非0值
	db.Model(&User{}).Where("city=?", "北京").Updates(User{Degree: "硕士", Gender: "男"})
	fmt.Printf("更新了%d行\n", tx.RowsAffected)
}

// 查询数据
func Read(db *gorm.DB) {
	user := User{City: "HongKong"}
	tx := db.
	Select("uid, city, gender, keywords").	// 参数也可以这样传
	Where("uid>100 and degree='本科'").		// 没有Select时默认为select*
	Where("city in ?", []string{"北京", "上海"}).	// 多个Where之间是and关系
	Where("degree like ?", "%科").
	Or("gender=?", "男").	// Where 与 Or 是或关系  用?占位, 避免发生SQL注入攻击, ? 是后面的值的赋值
	Order("id desc, uid").	// 排序首先按照 id 进行排序, desc 表示倒序, 按照 uid 升序排列
	Order("city").	// 再按照 city 进行排序
	Offset(3).		// 跳过前面几条记录
	Limit(1).		// 读取几条结果
	Find(&user)	// Take、First、Last 查不到结果时会返回gorm.ErrRecordNotFound, 但Find不会,Find查无结果时就不去修改结构体
	// First、Last 分别表示取第一条和最后一条结果
	// Take 表示随机取一条结果
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("读DB失败", "error", tx.Error)
		} else {
			slog.Info("查无结果")
		}
	}

	user1 := User{Id: 47}
	tx = db.Find(&user1)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("读DB失败", "error", tx.Error)
		} else {
			slog.Info("查无结果")
		}
	} else {
		if tx.RowsAffected > 0 {
			fmt.Printf("read结果: %v\n", user1)
		} else {
			slog.Info("查无结果")
		}
	}

	var users []User
	// 智能判断是否使用索引, 若有索引, 则使用索引, 否则不使用索引
	db.Where("uid > 0").
		Clauses(hints.UseIndex("id", "id_uid")). // 智能判断是否使用索引, 若有索引, 则使用索引, 否则不使用索引
		Find(&users)
	db.Where("uid > 0").
		Clauses(hints.ForceIndex("id_uid")). // 强制使用索引 id_uid
		Find(&users)
}


func ReadWithStatistics(db *gorm.DB) {
	type Result struct {
		City	string
		Mid		float64
	}

	var results []Result
	// 统计每个城市的用户id的平均值, 并筛选出平均值大于0的城市
	db.Model(User{}).Select("city, avg(id) as mid").Group("city").Having("mid > 0").Find(&results)
	fmt.Println("group by having 查询结果:")
	for _, result := range results {
		fmt.Printf("%v\n", result)
	}

	// 统计每个城市的用户数量
	// Distinct 表示去重
	db.Table("user").Distinct("city").Find(results)	// db.Model(User{}) 同等于 db.Table("user")
	fmt.Println("distinct 查询结果:")
	for _, result := range results {
		fmt.Printf("%v\n", result)
	}

	var count int64
	db.Table("user").Distinct("city").Count(&count)	// 统计去重后的结果数量
	fmt.Printf("distinct 结果数量: %d\n", count)
}

// 事务
func Transaction(db gorm.DB) error {
	tx := db.Begin()	// 开始事务

	// 事务中进行数据库操作
	defer func ()  {
		if err := recover(); err != nil {
			tx.Rollback()	// 手动回滚
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	user := User{UserId: rand.IntN(100), Degree: "本科", Gender: "男", City: "上海", Keywords: []string{"你好！", "世界"} }
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()	// 手动回滚
		fmt.Println("第一次create回滚")
		return err
	}
	fmt.Printf("uid=%d\n", user.UserId)

	user.UserId = 0
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()	// 手动回滚
		fmt.Println("第一次create回滚")
		return err
	}
	return tx.Commit().Error
}

// 执行原生的SQL语句
func RawSQL(db *gorm.DB) {
	var user []User
	db.Raw("select * from user where id = ?", 1).Scan(&user)
	fmt.Printf("raw sql 查询结果: %v\n", user)

	rows, err := db.Raw("select * from user where id = ?", 1).Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.UserId, &user.Degree, &user.Gender, &user.City, &user.Keywords)
		fmt.Printf("%v\n", user)
	}
}

