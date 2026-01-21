package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 定义根命令
var rootCmd = &cobra.Command{
	// 命令行程序名称
	Use: "test",
	// 简短介绍
	Short: "A test CLI application",
	// 长介绍
	Long:  `This is a test CLI application built with Cobra.`,
	
	// RunE 是带 error 返回值的 Run
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test CLI run begin.")

		fmt.Println("test CLI run end.")
	},

}

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
	// 修改标志 设置无选项时的默认值
	// rootCmd.Flag("source").NoOptDefVal = "default"

	// 精确绑定viper某一标志
	// viper.BindPFlag("auther", rootCmd.PersistentFlags().Lookup("auther"))

	// 绑定所有持久标志
	viper.BindPFlags(rootCmd.PersistentFlags())
}

// 执行命令
func Execute() {
	rootCmd.Execute()
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