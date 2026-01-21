package redis_test

import (
    "context"
    "strconv"
    "testing"
    "time"

    "github.com/redis/go-redis/v9"
    myRedis "happyladysauce/redis"
)

func TestInitRedis(t *testing.T) {
    myRedis.InitRedis()
    if myRedis.Client == nil {
        t.Fatalf("Client should not be nil after InitRedis")
    }
    if err := myRedis.Client.Ping(context.Background()).Err(); err != nil {
        t.Fatalf("Ping failed: %v", err)
    }
}

func TestStringValue(t *testing.T) {
    myRedis.InitRedis()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    myRedis.StringValue(ctx, myRedis.Client)

    if _, err := myRedis.Client.Get(ctx, "name").Result(); err != redis.Nil {
        t.Fatalf("expected 'name' to be deleted after StringValue; got err=%v", err)
    }
    if _, err := myRedis.Client.Get(ctx, "age").Result(); err != redis.Nil {
        t.Fatalf("expected 'age' to be deleted after StringValue; got err=%v", err)
    }
}

func TestSetGetStruct(t *testing.T) {
    myRedis.InitRedis()
    ctx := context.Background()

    stu := &myRedis.Student{Id: 101, Name: "Alice", Sex: "F", Age: 20}
    if err := myRedis.SetStruct(ctx, myRedis.Client, stu); err != nil {
        t.Fatalf("SetStruct failed: %v", err)
    }

    got := myRedis.GetStruct(ctx, myRedis.Client, stu.Id)
    if got == nil {
        t.Fatalf("GetStruct returned nil for id=%d", stu.Id)
    }
    if got.Name != stu.Name || got.Sex != stu.Sex || got.Age != stu.Age || got.Id != stu.Id {
        t.Fatalf("GetStruct mismatch: got=%+v want=%+v", got, stu)
    }

    // cleanup
    myRedis.Client.Del(ctx, "STU_"+strconv.Itoa(stu.Id))
}

func TestDeleteKey(t *testing.T) {
    myRedis.InitRedis()
    ctx := context.Background()
    // Ensure the function runs without panic; underlying DEL on a non-existent key should be safe
    myRedis.DeleteKey(ctx, myRedis.Client)
}

func TestXxx(t *testing.T) {
    myRedis.InitRedis()
    ctx := context.Background()
    myRedis.GetStudentFromRedis(ctx, myRedis.Client)
}