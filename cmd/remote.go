package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path"
)

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "远程服务器地址管理",
	Long:  `远程服务器地址管理`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列举所有",
	Long:  `列举所有`,
	Run: func(cmd *cobra.Command, args []string) {
		remotes := viper.Get("remote")
		remoteDefault := viper.GetString("remoteDefault")
		var s = remotes.(map[string]interface{})
		for k, v := range s {
			if k == remoteDefault {
				color.Red(fmt.Sprintf("* %s %s (current)", k, v))
			} else {
				fmt.Println(fmt.Sprintf("  %s %s", k, v))
			}
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "新增远程服务器地址",
	Long:  `新增远程服务器地址`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("请按照格式添加：gg-cli remote add #key #url")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		remotes := viper.Get("remote")
		var s = remotes.(map[string]interface{})
		s[args[0]] = args[1]

		home, _ := homedir.Dir()

		viper.Set("remote", s)
		if err := viper.WriteConfigAs(path.Join(home, ".gg-cli.yaml")); err == nil {
			fmt.Printf("新增远程服务地址成功 #key = %s, #url = %s \n", args[0], args[1])
		} else {
			fmt.Printf("新增远程服务地址失败 %v\n", err)
		}
	},
}

var delCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除远程服务器地址",
	Long:  `删除远程服务器地址`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("请按照格式删除：gg-cli remote del #key")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		remotes := viper.Get("remote")
		var s = remotes.(map[string]interface{})
		before := len(s)
		delete(s, args[0])
		end := len(s)
		if before == end {
			return
		}
		home, _ := homedir.Dir()
		viper.Set("remote", s)
		if err := viper.WriteConfigAs(path.Join(home, ".gg-cli.yaml")); err == nil {
			fmt.Printf("删除远程服务地址成功 #key = %s\n", args[0])
		} else {
			fmt.Printf("删除远程服务地址失败 %v\n", err)
		}
	},
}

var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "变更当前远程服务器地址",
	Long:  `变更当前远程服务器地址`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("请按照格式变更：gg-cli remote change #key")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		remotes := viper.Get("remote")
		var s = remotes.(map[string]interface{})
		_, ok := s[args[0]]
		if !ok {
			fmt.Printf("没有 %s 的地址别名\n", args[0])
			return
		}

		home, _ := homedir.Dir()
		viper.Set("remoteDefault", args[0])
		if err := viper.WriteConfigAs(path.Join(home, ".gg-cli.yaml")); err == nil {
			fmt.Printf("变更远程服务地址成功 #key = %s\n", args[0])
		} else {
			fmt.Printf("变更远程服务地址失败 %v\n", err)
		}
	},
}

func GetRemoteUrl() string {
	remotes := viper.Get("remote")
	remoteDefault := viper.GetString("remoteDefault")
	var s = remotes.(map[string]interface{})
	return s[remoteDefault].(string)
}

func init() {
	remoteCmd.AddCommand(lsCmd)
	remoteCmd.AddCommand(addCmd)
	remoteCmd.AddCommand(delCmd)
	remoteCmd.AddCommand(changeCmd)
}
