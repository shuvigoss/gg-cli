package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/mholt/archiver/v3"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

var offline bool
var dir string

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

	local := generateLocal(name, version)
	if local == "" {
		fmt.Println("创建目录异常")
		_ = os.RemoveAll(local)
		return
	}

	doGenerator(local)

}

func doGenerator(local string) {
	r := make(map[string]interface{})
	ggJson := path.Join(local, "gg.json")
	file, err := ioutil.ReadFile(ggJson)
	if err != nil {
		color.Red("读取gg.json %s 异常 %v", ggJson, err)
		return
	}
	_ = json.Unmarshal(file, &r)
	actions, ok := r["actions"]
	if !ok {
		fmt.Println("gg.json 没有actions, 无需处理")
		return
	}
	p := funk.Map(actions, func(m interface{}) map[string]interface{} {
		return m.(map[string]interface{})
	}).([]map[string]interface{})
	prompt := NewPrompt(p)
	prompt.Parse()
	if err = prompt.Run(); err != nil {
		color.Red("组织参数异常 %v", err)
		return
	}

	doReplace(prompt, local)
}

func doReplace(prompt Prompt, local string) {
	err := filepath.Walk(local, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("读取文件异常 %v", err)
			return err
		}
		parse, err := template.New("gg").Parse(string(f))
		if err != nil {
			fmt.Printf("解析模板异常 %v", err)
			return err
		}

		create, err := os.Create(path)
		if err != nil {
			log.Printf("创建文件异常 %v", err)
			return err
		}
		defer create.Close()

		if err = parse.Execute(create, prompt.Answer); err != nil {
			fmt.Printf("替换文件异常 %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("doReplace 异常 %v", err)
	} else {
		color.Blue("创建成功, 目录%s", local)
	}
}

func generateLocal(name, version string) string {
	//检查目录是否存在
	absolutePath := getAbsolutePath()
	if Exists(absolutePath) {
		color.Red("目标路径 %s 已经存在，请重新选择路径! ", absolutePath)
		return ""
	}

	homeDir, _ := homedir.Dir()
	d := path.Join(homeDir, ".gg", name)
	readDir, _ := ioutil.ReadDir(d)
	vs := funk.Map(readDir, func(fi os.FileInfo) string {
		return fi.Name()
	}).([]string)

	realVersion := GetRealVersion(name, version, vs)

	if realVersion == "" {
		color.Red("目标版本未找到 %s", name)
		return ""
	}

	_ = os.MkdirAll(absolutePath, os.ModePerm)
	localPath := path.Join(d, realVersion)
	infos, _ := ioutil.ReadDir(localPath)
	zipFile := GetZipFile(infos, name)
	if err := archiver.Unarchive(path.Join(localPath, zipFile), absolutePath); err != nil {
		color.Red("解压文件异常 %v", err)
		return ""
	}

	return absolutePath

}

func getAbsolutePath() string {
	var absolutePath string
	if strings.HasPrefix(dir, "/") {
		absolutePath = dir
	} else {
		wd, _ := os.Getwd()
		absolutePath = path.Join(wd, dir)
	}
	return absolutePath
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func init() {
	gCmd.PersistentFlags().BoolVarP(&offline, "offline", "o", false, "是否离线使用,默认为false")
	gCmd.Flags().StringVarP(&dir, "dir", "d", "", "生成文件夹名称(必填)")
	_ = gCmd.MarkFlagRequired("dir")
}
