Cobra 是一个 [Go 语言](https://zhida.zhihu.com/search?content_id=227653816&content_type=Article&match_order=1&q=Go+%E8%AF%AD%E8%A8%80&zd_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ6aGlkYV9zZXJ2ZXIiLCJleHAiOjE3NDY1ODk2OTYsInEiOiJHbyDor63oqIAiLCJ6aGlkYV9zb3VyY2UiOiJlbnRpdHkiLCJjb250ZW50X2lkIjoyMjc2NTM4MTYsImNvbnRlbnRfdHlwZSI6IkFydGljbGUiLCJtYXRjaF9vcmRlciI6MSwiemRfdG9rZW4iOm51bGx9.Id47rm2HKIg4x2UxIBftkAWcvJslBvDJIdgZCMlCJ_U&zhida_source=entity)开发的命令行（CLI）框架，它提供了简洁、灵活且强大的方式来创建命令行程序。它包含一个用于创建命令行程序的库（Cobra 库），以及一个用于快速生成基于 Cobra 库的命令行程序工具（Cobra 命令）。Cobra 是由 Go 团队成员 [spf13](https://link.zhihu.com/?target=https%3A//spf13.com/) 为 [Hugo](https://link.zhihu.com/?target=https%3A//gohugo.io/) 项目创建的，并已被许多流行的 Go 项目所采用，如 [Kubernetes](https://zhida.zhihu.com/search?content_id=227653816&content_type=Article&match_order=1&q=Kubernetes&zd_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ6aGlkYV9zZXJ2ZXIiLCJleHAiOjE3NDY1ODk2OTYsInEiOiJLdWJlcm5ldGVzIiwiemhpZGFfc291cmNlIjoiZW50aXR5IiwiY29udGVudF9pZCI6MjI3NjUzODE2LCJjb250ZW50X3R5cGUiOiJBcnRpY2xlIiwibWF0Y2hfb3JkZXIiOjEsInpkX3Rva2VuIjpudWxsfQ.vaCk0FjNmW_TpQBvwuA6lyaQKgiFiz0mgvV--6KO1Iw&zhida_source=entity)、[Helm](https://zhida.zhihu.com/search?content_id=227653816&content_type=Article&match_order=1&q=Helm&zd_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ6aGlkYV9zZXJ2ZXIiLCJleHAiOjE3NDY1ODk2OTYsInEiOiJIZWxtIiwiemhpZGFfc291cmNlIjoiZW50aXR5IiwiY29udGVudF9pZCI6MjI3NjUzODE2LCJjb250ZW50X3R5cGUiOiJBcnRpY2xlIiwibWF0Y2hfb3JkZXIiOjEsInpkX3Rva2VuIjpudWxsfQ.ImezRpwM-P3tlVzQNZTLIxBDT_KX1GHA7Guyt3qJYwg&zhida_source=entity)、[Docker](https://zhida.zhihu.com/search?content_id=227653816&content_type=Article&match_order=1&q=Docker&zd_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ6aGlkYV9zZXJ2ZXIiLCJleHAiOjE3NDY1ODk2OTYsInEiOiJEb2NrZXIiLCJ6aGlkYV9zb3VyY2UiOiJlbnRpdHkiLCJjb250ZW50X2lkIjoyMjc2NTM4MTYsImNvbnRlbnRfdHlwZSI6IkFydGljbGUiLCJtYXRjaF9vcmRlciI6MSwiemRfdG9rZW4iOm51bGx9.oekQs3vtldXmLd_n3rKzme9Agh3hrfoNudrsSI3615Q&zhida_source=entity) (distribution)、[Etcd](https://zhida.zhihu.com/search?content_id=227653816&content_type=Article&match_order=1&q=Etcd&zd_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ6aGlkYV9zZXJ2ZXIiLCJleHAiOjE3NDY1ODk2OTYsInEiOiJFdGNkIiwiemhpZGFfc291cmNlIjoiZW50aXR5IiwiY29udGVudF9pZCI6MjI3NjUzODE2LCJjb250ZW50X3R5cGUiOiJBcnRpY2xlIiwibWF0Y2hfb3JkZXIiOjEsInpkX3Rva2VuIjpudWxsfQ.eCYc2Yu4VALH30UTMAUNKxBdvdi7cqzUsLK6mv28EeY&zhida_source=entity) 等。
[万字长文——Go 语言现代命令行框架 Cobra 详解 - 知乎](https://zhuanlan.zhihu.com/p/627848739)
[Go 命令行参数解析工具 pflag 使用 | 江湖十年 | 学而不思则罔，思而不学则殆。](https://jianghushinian.cn/2023/03/27/use-of-go-command-line-parameter-parsing-tool-pflag/)
bilbil 视频教程：[三步搞定 Cobra 命令行框架](https://www.bilibili.com/video/BV1uh4y13765/?vd_source=aec29d870700f0b3263a1f63df363d23)

1. cobra 命令行框架的基本项目结构逻辑
2. cobra 命令行的本地标志与持久标志
3. cobra 命令行参数与参数验证
## 1.Cobra 命令行框架的基本项目结构逻辑

Cobra 建立在**命令**、**参数**和**标志**这三个结构之上。要使用 Cobra 编写一个命令行程序，需要明确这三个概念。

- 命令（COMMAND）：命令表示要执行的操作。  
- 参数（ARG）：是命令的参数，一般用来表示操作的对象。  
- 标志（FLAG）：是命令的修饰，可以调整操作的行为。

> 一个好的命令行程序在使用时读起来像句子，用户会自然的理解并知道如何使用该程序。
> 要编写一个好的命令行程序，需要遵循的模式是：
> `APPNAME VERB NOUN --ADJECTIVE` 或 `APPNAME COMMAND ARG --FLAG`。
> 在这里 `VERB` 代表动词，`NOUN` 代表名词，`ADJECTIVE` 代表形容词。

```go
// 定义根命令
var rootCmd = &cobra.Command{
    // 命令行程序名称
    Use: "test",
    // 简短介绍
    Short: "A test CLI application",
    // 长介绍
    Long:  `This is a test CLI application built with Cobra.`,
    // RunE 是带 error 返回值的 Run
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("test CLI run begin.")
        fmt.Println("test CLI run end.")
    },
}

// 执行命令
func Execute() {
    rootCmd.Execute()
}
```

## 2.Cobra 命令行的本地标志与持久标志

使用 `init` 函数用来加载标志，并且加载配置文件。使用 `viper.BindPFlags(rootCmd.PersistentFlags())` 可一键根据标志名称进行变量赋值。

```go
var cfgFile string
var userLicense string

func init() {
	// 初始化配置
	cobra.OnInitialize(initConfig)
	// 持久标志
	rootCmd.PersistentFlags().Bool("viper", true, "")
	rootCmd.PersistentFlags().StringP("auther", "a", "HappyLadySauce", "")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "License", "l", "", "")

	// 本地标志
	rootCmd.Flags().StringP("source", "s", "", "")

	// 精确绑定viper某一标志
	// viper.BindPFlag("auther", rootCmd.PersistentFlags().Lookup("auther"))

	// 绑定所有持久标志
	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initConfig() {
	// 检查配置文件是否存在
	if cfgFile == "" {
		return
	}
	// 加载配置文件
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())
}
```

## 3.Cobra 命令行参数与参数验证

```go
var validString = []string{"123", "abc", "456"}

// 定义arg命令
var argCmd = &cobra.Command{
	Use:   "arg",
	Short: "A command to demonstrate argument handling",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("arg command run begin.")
		fmt.Printf("Arguments: %v\n", args)
		fmt.Println("arg command run end.")
	},
	// 参数验证函数
	ValidArgs: validString,
	Args:      combinedArgValidate,
}

// 初始化函数
func init() {
	// 将argCmd添加为rootCmd的子命令
	rootCmd.AddCommand(argCmd)
}

// 组合参数验证函数
func combinedArgValidate(cmd *cobra.Command, args []string) error {
	// 1. 首先检查参数数量
	if err := argValidate(cmd, args); err != nil {
		return err
	}
	// 2. 然后检查参数是否有效
	return cobra.OnlyValidArgs(cmd, args)
}

// 参数验证函数
func argValidate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("requires at least one arg")
	}
	return nil
}
```