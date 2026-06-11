package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	einoIndexer "github.com/cloudwego/eino-ext/components/indexer/milvus2"
	einoRetriever "github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2/search_mode"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

const (
	collectionName = "go_knowledge"
	embeddingDim   = 768 // text-embedding-nomic-embed-text-v1.5 输出维度
	vectorField    = "vector"
)

// 知识库文档
var documents = []*schema.Document{
	{ID: "doc_1", Content: "Go语言的并发模型基于CSP理论，goroutine是轻量级协程，初始栈空间仅2KB，通过channel进行类型安全的通信。"},
	{ID: "doc_2", Content: "Go语言的GC采用三色标记清除算法，从1.5版本开始STW时间控制在毫秒级，可通过GOGC环境变量调整触发频率。"},
	{ID: "doc_3", Content: "Go语言的接口是隐式实现的，只要类型实现了接口的所有方法就自动满足该接口，空接口interface{}可用any代替。"},
	{ID: "doc_4", Content: "Go Module是官方依赖管理方案，go.mod记录模块路径和依赖版本，go.sum保存哈希校验值，go mod tidy清理依赖。"},
	{ID: "doc_5", Content: "Go语言的错误处理采用显式返回error的方式，errors.Is和errors.As用于判断错误类型，支持%w格式化动词包装错误。"},
}

const (
	// Embedding 专用模型；对话生成可另用 harrier-oss-v1-0.6b
	Model	= "text-embedding-nomic-embed-text-v1.5"
	APIURL	= "http://100.100.100.254:11434/v1"
)


func main() {
	ctx := context.Background()

	dim := embeddingDim
	embedder, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
		APIKey: "",
		Model: Model,
		BaseURL: APIURL,
		Dimensions: &dim,
	})
	if err != nil {
		fmt.Printf("创建 Embedder 失败: %w\n", err)
		return
	}
	// fmt.Println("✅ Embedder 初始化成功")

	vectors, err := embedder.EmbedStrings(ctx, []string{"Hello, Go Agent!"})
	if err != nil {
		fmt.Printf("Embedding 测试失败: %w", err)
		return
	}
	fmt.Printf("测试文本向量: %v\n", vectors)

	milvusClientConfig := &milvusclient.ClientConfig{
		Address: "100.100.100.254:19530",
	}

	// Drop legacy collection if present (e.g. created by openai/main.go with incompatible schema).
	// 若存在旧版 Collection（如 openai/main.go 创建的 schema 不兼容），先删除再重建。
	milvusClient, err := milvusclient.New(ctx, milvusClientConfig)
	if err != nil {
		fmt.Printf("连接 Milvus 失败: %v\n", err)
		return
	}
	defer milvusClient.Close(ctx)

	has, err := milvusClient.HasCollection(ctx, milvusclient.NewHasCollectionOption(collectionName))
	if err != nil {
		fmt.Printf("检查 Collection 失败: %v\n", err)
		return
	}
	if has {
		if err = milvusClient.DropCollection(ctx, milvusclient.NewDropCollectionOption(collectionName)); err != nil {
			fmt.Printf("删除旧 Collection 失败: %v\n", err)
			return
		}
		fmt.Printf("已删除旧 Collection: %s\n", collectionName)
	}

	indexer, err := einoIndexer.NewIndexer(ctx, &einoIndexer.IndexerConfig{
		ClientConfig: milvusClientConfig,
		Collection:   collectionName,
		Vector: &einoIndexer.VectorConfig{
			VectorField: vectorField,
			Dimension:   int64(embeddingDim),
			MetricType:  einoIndexer.COSINE,
			IndexBuilder: einoIndexer.NewHNSWIndexBuilder().
				WithM(16).
				WithEfConstruction(200),
		},
		Embedding: embedder,
	})
	if err != nil {
		fmt.Printf("创建 Indexer 失败: %v\n", err)
		return
	}
	fmt.Println("Indexer 初始化完成")

	fmt.Println("正在存入文档")
	ids, err := indexer.Store(ctx, documents)
	if err != nil {
		fmt.Printf("存储文档失败: %v\n", err)
		return
	}
	fmt.Printf("成功存储 %d 篇文档, ID: %v\n", len(ids), ids)


	retriever, err := einoRetriever.NewRetriever(ctx, &einoRetriever.RetrieverConfig{
		ClientConfig: milvusClientConfig,
		Collection:   collectionName,
		VectorField:  vectorField,
		OutputFields: []string{"id", "content"},
		TopK:         3,
		SearchMode:   search_mode.NewApproximate(einoRetriever.COSINE),
		Embedding:    embedder,
	})
	if err != nil {
		fmt.Printf("创建 Retriever 失败: %v\n", err)
		return
	}
	fmt.Println("Retriever 初始化成功")

	queries := []string{
		"Go语言怎么做并发编程？",
		"Go的垃圾回收机制是怎样的？",
		"怎么管理Go项目的依赖？",
	}

	for _, query := range queries {
		fmt.Printf("\n 查询: %s\n", query)
		docs, err := retriever.Retrieve(ctx, query)
		if err != nil {
			fmt.Printf("检索失败: %v\n", err)
			continue
		}

		for i, doc := range docs {
			fmt.Printf("[%d]相似度: %.4f\n内容: %s\n", i+1, doc.Score(), doc.Content)
		}
	}


}





