package main

import (
	"context"
	"fmt"
	"sync"

	openai "github.com/sashabaranov/go-openai"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

const (
	collectionName	= "go_knowledge"
	embeddingDim	= 768 // text-embedding-nomic-embed-text-v1.5 输出维度
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
	// Embedding 专用模型；对话生成可另用 harrier-oss-v1-0.6b
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
		go func(idx int, text string) {
			defer wg.Done()
			contents[idx] = text
			emb, embErr := getEmbedding(ctx, aiClient, text)
			if embErr != nil {
				fmt.Printf("document %d embedding failed: %v\n", idx, embErr)
				return
			}
			embeddings[idx] = emb
			fmt.Printf("document %d embedded (%d dims)\n", idx, len(emb))
		}(i, doc)
	}
	wg.Wait()

	// Batch insert into Milvus (id is auto-generated).
	// 批量插入 Milvus（id 由数据库自动生成）
	validContents, validEmbeddings, ok := prepareInsertRows(contents, embeddings)
	if !ok {
		fmt.Println("batch insert aborted: some documents failed to embed")
		return
	}

	fmt.Printf("\ninserting %d documents into collection %q...\n", len(validContents), collectionName)
	insertResult, err := milvusClient.Insert(ctx, milvusclient.NewColumnBasedInsertOption(collectionName).
		WithVarcharColumn("content", validContents).
		WithFloatVectorColumn("embedding", embeddingDim, validEmbeddings),
	)
	if err != nil {
		fmt.Printf("batch insert failed: %v\n", err)
		return
	}
	fmt.Printf("batch insert succeeded, inserted %d rows\n", insertResult.InsertCount)

	flushTask, flushErr := milvusClient.Flush(ctx, milvusclient.NewFlushOption(collectionName))
	if flushErr != nil {
		fmt.Printf("flush collection failed: %v\n", flushErr)
		return
	}
	if err = flushTask.Await(ctx); err != nil {
		fmt.Printf("flush task failed: %v\n", err)
		return
	}

	hnswIndex := index.NewHNSWIndex(entity.COSINE, 16, 200)
	createIndexTask, err := milvusClient.CreateIndex(ctx, milvusclient.NewCreateIndexOption(collectionName, "embedding", hnswIndex))
	if err != nil {
		fmt.Printf("create index failed: %v\n", err)
		return
	}
	createIndexTask.Await(ctx)
	fmt.Printf("create index task completed\n")

	// 加载 Collection 到内存
	loadTask, err := milvusClient.LoadCollection(ctx, milvusclient.NewLoadCollectionOption(collectionName))
	if err != nil {
		fmt.Printf("load collection failed: %v\n", err)
		return
	}
	loadTask.Await(ctx)


	query := "Go语言怎么做并发编程？"
	fmt.Printf("query: %s\n", query)

	queryEmbedding, err := getEmbedding(ctx, aiClient, query)
	if err != nil {
		fmt.Printf("get embedding failed: %v\n", err)
		return
	}

	fmt.Println("query embedding:")
	searchResult, err := milvusClient.Search(ctx, milvusclient.NewSearchOption(collectionName, 3, []entity.Vector{entity.FloatVector(queryEmbedding)}).
		WithOutputFields("content"),
	)
	if err != nil {
		fmt.Printf("search failed: %v\n", err)
		return
	}

	fmt.Println("search result:")
	for _, result := range searchResult {
		contentCol := result.GetColumn("content")
		if contentCol == nil {
			fmt.Println("content field not found in search output")
			continue
		}
		for i := 0; i < result.ResultCount; i++ {
			content, colErr := contentCol.GetAsString(i)
			if colErr != nil {
				fmt.Printf("read content at index %d failed: %v\n", i, colErr)
				continue
			}
			score := result.Scores[i]
			fmt.Printf("content: %s, score: %f\n", content, score)
		}
	}
}

// prepareInsertRows filters rows with successful embeddings for batch insert.
// prepareInsertRows 过滤向量化成功的行，供批量插入使用
func prepareInsertRows(contents []string, embeddings [][]float32) ([]string, [][]float32, bool) {
	validContents := make([]string, 0, len(contents))
	validEmbeddings := make([][]float32, 0, len(contents))

	for i := range contents {
		if embeddings[i] == nil {
			return nil, nil, false
		}
		if len(embeddings[i]) != embeddingDim {
			fmt.Printf("document %d embedding dim mismatch: got %d, want %d\n", i, len(embeddings[i]), embeddingDim)
			return nil, nil, false
		}
		validContents = append(validContents, contents[i])
		validEmbeddings = append(validEmbeddings, embeddings[i])
	}
	return validContents, validEmbeddings, true
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



