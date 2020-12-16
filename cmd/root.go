package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
)

var cfgFile string
var defaultConfig = `
remote:
  local: http://localhost:8080/
  remote: http://example.com/
remoteDefault: local
`

var rootCmd = &cobra.Command{
	Use:   "gg-cli",
	Short: "脚手架生成器(go generator)",
	Long:  `脚手架生成器(go generator)`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(remoteCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	createDefaultConfigFile(path.Join(home, ".gg-cli.yaml"))

	viper.AddConfigPath(home)
	viper.SetConfigName(".gg-cli")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("没有找到配置文件(.gg-cli.yaml) 请检查环境:" + cfgFile)
	}

}

func createDefaultConfigFile(f string) {
	if !Exists(f) {
		if err := ioutil.WriteFile(f, []byte(defaultConfig), os.ModePerm); err != nil {
			fmt.Printf("初始化 %s 文件异常", f)
		}
	}
}

func Exists(p string) bool {
	_, err := os.Stat(p)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
