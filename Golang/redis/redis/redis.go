package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// 定义全局 redis 客户端
var (
	Client *redis.Client
)

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:		"100.100.100.5:6379",
		DB:			0,
		Username:	"",
		Password:	"@Ssddffqxc547",
	})

	if err := client.Ping(context.Background()).Err(); err !=nil {
		slog.Error("connect to redis failed. error: ", "error", err)
	} else {
		slog.Info("connect to redis ok.")
	}

	Client = client
}

func StringValue(ctx context.Context, client *redis.Client) {
	key := "name"
	value := "HappyLadySauce"

	if err := client.Set(ctx, key, value, 0).Err(); err != nil {	// 插入数据, 并设置超时时间, 当设置为0时表示永不超时
		fmt.Printf("redis insert data err: %v\n", err)
	}
	defer client.Del(ctx, key)	// 函数结束时删除本次插入的数据, 不影响下次允许演示

	client.Expire(ctx, key, 3 * time.Second)	// 显示设置超时时间
	time.Sleep(2 * time.Second)

	v, err := client.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("get redis key: %v data err: %v\n",  key, err)
	}
	fmt.Printf("redis data: %v:%v\n", key, v)

	if err = client.Set(ctx, "age", 18, 0).Err(); err != nil {	// 默认插入redis中的数据会转为string类型
		fmt.Printf("redis insert data err: %v\n", err)
	}
	defer client.Del(ctx, "age")	// 函数结束时删除本次插入的数据, 不影响下次允许演示

	v1, err := client.Get(ctx, "age").Int()	// 将redis中的string类型转为int类型接收
	if err != nil {
		fmt.Printf("get redis key: %v data err: %v\n",  key, err)
	}
	fmt.Printf("redis data: %v:%v\n", key, v1)
}

func DeleteKey(ctx context.Context, client *redis.Client) {
	n, err := client.Del(ctx, "not_exissts").Result()
	if err == nil {
		fmt.Printf("删除%d个key\n", n)
	}
}

// 在 redis 中存储结构体, 最好是先把结构体转为字符串
type Student struct{
	Id		int
	Name	string
	Sex		string
	Age 	int
}

func SetStruct(ctx context.Context, client *redis.Client, data *Student) error {
	if data == nil {
		return nil
	}

	key := "STU_" + strconv.Itoa(data.Id)	// 添加学生业务前缀, 避免业务混淆
	v, err := json.Marshal(data)	// 将结构体转为 JSON 字符串
	if err != nil {
		return err
	}
	err = client.Set(ctx, key, string(v), 0).Err()
	return err
}

func GetStruct(ctx context.Context, client *redis.Client, id int) *Student {
	key := "STU_" + strconv.Itoa(id)
	v, err := client.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("get redis key: %v data err: %v\n",  key, err)
		return nil
	}

	var data Student
	if err = json.Unmarshal([]byte(v), &data); err != nil {
		fmt.Printf("Student JSON data Unmarshal err: %v\n", err)
		return nil
	}
	return &data
}

// redis 的复杂类型操作

// List 数组操作 使用 List 可以天然的保证顺序
func ValueList(ctx context.Context, client *redis.Client) {
	key := "id"
	value := []interface{}{1, 2, 4, 1, "你好"}

	if err := client.RPush(ctx, key, value...).Err(); err != nil {	// RPush(Right)表示尾插, LPush(Left)表示头插
		fmt.Printf("redis insert data err: %v\n", err)
	}
	defer client.Del(ctx, key)

	v1, err := client.LRange(ctx, key, 0, -1).Result()	// 截取, 双闭区间. LRange表示List
	if err != nil {
		fmt.Printf("get redis key: %v data err: %v\n",  key, err)
	}
	fmt.Printf("redis data: %v:%v\n", key, v1)	// {1, 2, 4, 1, "你好"} 全都是字符串
}

// Set 数据操作, Set 与 List 很相似, 只是 Set 不能有重复的数据, 并且不保证数据顺序
// 可以使用 redis 中的 Set 进行排重
// 求交集/差集
func ValueSet(ctx context.Context, client *redis.Client) {
	key1 := "id1"
	value := []interface{}{1, 2, 4, 1, "你好"}

	if err := client.SAdd(ctx, key1, value...).Err(); err != nil {
		fmt.Printf("redis insert data err: %v\n", err)
	}
	defer client.Del(ctx, key1)

	// 判断 Set 中是否包含指定元素
	var index any
	index = 1
	if client.SIsMember(ctx, key1, value).Val() {
		fmt.Printf("Set 中包含元素 %v\n", index)
	} else {
		fmt.Printf("Set 中不包含元素 %v\n", index)
	}

	index = 3
	if client.SIsMember(ctx, key1, value).Val() {
		fmt.Printf("Set 中包含元素 %v\n", index)
	} else {
		fmt.Printf("Set 中不包含元素 %v\n", index)
	}

	// 遍历 Set
	for _, ele := range client.SMembers(ctx, key1).Val() {	// SMembers 表示获取 Set 中的所有元素
		fmt.Printf("Set 元素: %v\n", ele)
	}

	key2 := "id2"
	value2 := []interface{}{1, 2, 3, 1, "你好"}
	if err := client.SAdd(ctx, key2, value2...).Err(); err != nil {
		fmt.Printf("redis insert data err: %v\n", err)
	}
	defer client.Del(ctx, key2)

	// Set 之间求差集
	fmt.Println("key1 - key2 的差集")
	for _, ele := range client.SDiff(ctx, key1, key2).Val() {
		fmt.Println(ele)
	}
	fmt.Println("key2 - key1 的差集")
	for _, ele := range client.SDiff(ctx, key2, key1).Val() {
		fmt.Println(ele)
	}

	// Set 之间的交集
	fmt.Println("key1 与 key2 的交集")
	for _, ele := range client.SInter(ctx, key1, key2).Val() {
		fmt.Println(ele)
	}
}

// Redis 中的 ZSet 是有序的 Set, 可以排重且保证顺序
func ValueZSet(ctx context.Context, client *redis.Client) {
	key := "id"
	value := []redis.Z{{Member: "张珊", Score: 1}, {Member: "哈哈", Score: 3}, {Member: "你好", Score: 2}}	// 按照 Score 用来排序, 比如把时间戳给 Score

	if err := client.ZAdd(ctx, key, value...).Err(); err != nil {
		fmt.Printf("redis insert data err: %v\n", err)
	}
	defer client.Del(ctx, key)

	// 遍历 ZSet, 按 Score 有序输出 Member
	for _, ele := range client.ZRange(ctx, key, 0, -1).Val() {
		fmt.Println(ele)
	}
}

// 哈希表 Value
func ValueHashtable(ctx context.Context, client *redis.Client) {
	student1 := map[string]interface{}{"Name": "张珊", "Age": 18, "Sex": "男"}
	if err := client.HMSet(ctx, "学生1", student1).Err(); err != nil {
		fmt.Printf("redis insert data err: %v\n", err)
	}
	defer client.Del(ctx, "学生1")

	student2 := map[string]interface{}{"Name": "张三", "Age": 18, "Sex": "男"}
	if err := client.HMSet(ctx, "学生2", student2).Err(); err != nil {
		fmt.Printf("redis insert data err: %v\n", err)
	}
	defer client.Del(ctx, "学生2")

	age, err := client.HGet(ctx, "学生1", "Age").Int()	// 通过指定学生字段, 再指定结构体中的字段获取数据, 并且存储在 redis 中的数据均为 string 需要进行转换
	if err != nil {
		fmt.Printf("get redis key: %v data err: %v\n",  "学生1", err)
		return
	}
	fmt.Printf("学生1 的年龄为: %v\n", age)

	for key, value := range client.HGetAll(ctx, "学生2").Val() {	// HGetAll 表示获取哈希表中的所有字段和值
		fmt.Printf("学生2 的 %v 为: %v\n", key, value)
	}
}

// 在 redis 中存储结构体, 最好是先把结构体转为字符串
// type Student struct{
// 	Id		int
// 	Name	string
// 	Sex		string
// 	Age 	int
// }

// 批量遍历 redis
func GetStudentFromRedis(ctx context.Context, client *redis.Client) {
	const P = 100
	// 构造学生数据
	for i := 0; i < P; i++ {
		key := "STU_" + strconv.Itoa(i)	
		value := &Student{
			Id:		i,
			Name: 	"HappyLadySauce",
			Sex: 	"男",
			Age:	18,
		}

		data, err := json.Marshal(value)
		if err != nil {
			fmt.Printf("marshal student err: %v\n", err)
		}

		if err := client.Set(ctx, key, string(data), 5 * time.Second).Err(); err != nil {
			fmt.Printf("set redis key: %v data err: %v\n", key, err)
		}
	} 

	// 函数结束时删除模拟数据
	defer func(ctx context.Context, client *redis.Client) {
		for i := 0; i < P; i++ {
		key := "STU_" + strconv.Itoa(i)	
		client.Del(ctx, key)
		}
	}(ctx, client)

	// 遍历 redis
	const MID = "STU_"
	var cursor uint64 = 0

	for {
		keys, c, err := client.Scan(ctx, cursor, MID+"*", 3).Result()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("cursor %d keys count %d\n", c, len(keys))
		if c == 0 {
			break
		}
		cursor = c	// 为游标重新赋值
	}
}







