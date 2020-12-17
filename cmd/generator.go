package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var offline bool

var gCmd = &cobra.Command{
	Use:   "generator",
	Short: "生成模板",
	Long:  `生成模板`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("请正确使用 gg-cli generator #name[@#version] 命令")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		generate(cmd, args[0])
	},
}

func generate(cmd *cobra.Command, nameVersion string) {
	split := strings.Split(nameVersion, "@")
	name := split[0]
	version := ""
	if len(split) > 1 {
		version = split[1]
	}

	if !offline && !CompareAndSync(name, version) {
		fmt.Println("更新异常")
		return
	}
	fmt.Println("success")

}

func init() {
	gCmd.PersistentFlags().BoolVarP(&offline, "offline", "o", false, "是否离线使用,默认为false")
}
