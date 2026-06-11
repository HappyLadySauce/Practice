package main

import (
	"context"
	"fmt"
	"sync"

	openai "github.com/sashabaranov/go-openai"

	"github.com/milvus-io/milvus/client/v2/entity"
	// "github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

const (
	collectionName	= "go_knowledge"
	embeddingDim	= 768 // nomic-embed-text-v1.5 输出维度
)

// 知识库文档
var documents = []string{
	"Go语言的并发模型基于CSP理论，goroutine是轻量级协程，初始栈空间仅2KB，通过channel进行类型安全的通信。",
	"Go语言的GC采用三色标记清除算法，从1.5版本开始STW时间控制在毫秒级，可通过GOGC环境变量调整触发频率。",
	"Go语言的接口是隐式实现的，只要类型实现了接口的所有方法就自动满足该接口，空接口interface{}可用any代替。",
	"Go Module是官方依赖管理方案，go.mod记录模块路径和依赖版本，go.sum保存哈希校验值，go mod tidy清理依赖。",
	"Go语言的错误处理采用显式返回error的方式，errors.Is和errors.As用于判断错误类型，支持%w格式化动词包装错误。",
}

const (
	Model	= "text-embedding-nomic-embed-text-v1.5"
	APIURL	= "http://100.100.100.254:11434/v1"
)

func main() {
	wg := sync.WaitGroup{}
	ctx := context.Background()

	fmt.Println("开始运行")
	milvusClient, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address: "100.100.100.254:19530",
	})
	if err != nil {
		panic(fmt.Errorf("连接milvus失败: ", err))
	}
	defer milvusClient.Close(ctx)
	fmt.Println("milvus连接成功")

	aiConfig := openai.DefaultConfig("")
	aiConfig.BaseURL = APIURL
	aiClient := openai.NewClientWithConfig(aiConfig)

	// 如果 Collection 存在先删除
	has, _ := milvusClient.HasCollection(ctx, milvusclient.NewHasCollectionOption(collectionName))
	if has {
		milvusClient.DropCollection(ctx, milvusclient.NewDropCollectionOption(collectionName))
	}

	// 定义 Collection 的 Schema
	schema := entity.NewSchema().WithName(collectionName).WithDescription("Go语言知识库")
	schema.WithField(entity.NewField().WithName("id").WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true).WithIsAutoID(true)).
	WithField(entity.NewField().WithName("content").WithDataType(entity.FieldTypeVarChar).WithMaxLength(2000)).
	WithField(entity.NewField().WithName("embedding").WithDataType(entity.FieldTypeFloatVector).WithDim(embeddingDim))

	// 创建 Collection
	err = milvusClient.CreateCollection(ctx, milvusclient.NewCreateCollectionOption(collectionName, schema))
	if err != nil {
		fmt.Printf("创建 Collection 失败: %v\n", err)
		return
	}

	fmt.Println("✅ Collection 创建成功")


	fmt.Println("\n📚 正在并行处理文档向量化...")
	contents := make([]string, len(documents))
	embeddings := make([][]float32, len(documents))

	for i, doc := range documents {
		wg.Add(1)
		go func() {
			contents[i] = doc
			emb, err := getEmbedding(ctx, aiClient, doc)
			if err != nil {
				fmt.Printf("文档 %d 向量化失败: %v\n", i, err)
				wg.Done()
				return
			}
			embeddings[i] = emb
			fmt.Printf("文档 %v 已向量化(%d维)\n", i, len(emb))
			wg.Done()
		}()
	}

	wg.Wait()
}


func getEmbedding(ctx context.Context, client *openai.Client, text string) ([]float32, error) {
	resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.EmbeddingModel(Model),
		Input: []string{text},
	})
	if err != nil {
		return nil, err
	}
	return resp.Data[0].Embedding, nil
}



