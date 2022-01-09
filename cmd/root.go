package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var (
	vers bool
)

var RootCmd = &cobra.Command{
	//节点名称
	Use: "restful-api",
	//节点短描述
	Short: "restful-api curd demo",
	//节点长描述
	Long: `restful-api`,
	//需要在该节点运行的内容
	RunE: func(cmd *cobra.Command, args []string) error {
		//如果vers为true则打印版本信息
		if vers {
			fmt.Println("version: 1.0.0")
			return nil
		}
		return fmt.Errorf("no flafs find")
	},
}

//cmd启动器
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

//命令行参数，赋值的定义和初始化
func init() {
	//全局标志PersistentFlags 所有的子节点都可以使用  Flags仅本节点可用
	//第一个形参是cmd 传参传给哪个变量接受，第二个代表 该 参数的名称，第三个代表 参数指令，第四个 代表参数的默认值，第五个代表参数的帮助描述
	RootCmd.Flags().BoolVarP(&vers, "version", "v", false, "the rest-api-demo version")
}