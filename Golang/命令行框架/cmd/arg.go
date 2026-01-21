package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
