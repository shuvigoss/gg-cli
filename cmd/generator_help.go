package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"gg-cli/common"
	"github.com/mitchellh/go-homedir"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

type GgResult struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

/**
 *	http://localhost:8080/query?name=gg-example
 */
func CompareAndSync(name, version string) bool {
	//查询远程服务器相关包信息
	result := queryResult(name, version)
	//确定需要下载的版本
	realVersion := getTarget(name, version, result)
	if realVersion == "" {
		return false
	}
	//查看本地是否已经缓存,并更新缓存
	return doCached(name, realVersion)
}

func getZipFile(children []os.FileInfo, name string) string {
	for _, f := range children {
		match, _ := regexp.Match("^"+name+".*\\.(zip|gz|tar)$", []byte(f.Name()))
		if match {
			return f.Name()
		}
	}
	return ""
}

func doCached(name, version string) bool {
	dir, _ := homedir.Dir()
	theDir := path.Join(dir, ".gg", name, version)
	var remoteSha256 = getRemoteSha256(name, version)
	if Exists(theDir) {
		readDir, _ := ioutil.ReadDir(theDir)
		zipFile := getZipFile(readDir, name)
		sha256File := path.Join(theDir, zipFile+".sha256")
		localSha256, err := ioutil.ReadFile(sha256File)
		if err != nil {
			fmt.Printf("读取 %s 文件异常 %v\n", sha256File, err)
			return false
		}

		if remoteSha256 == string(localSha256) {
			fmt.Println("文件未发生变化")
			return true
		}
		//不一致的话删除本地，重新下载
		if err := os.RemoveAll(theDir); err != nil {
			fmt.Printf("删除 %s 异常 %v\n", theDir, err)
			return false
		}
	}
	//下载
	return doDownload(name, version, remoteSha256, theDir)
}

/**
 * http://localhost:8080/download?name=gg-example&version=1.0.0
 */
func doDownload(name, version, remoteSha256, dir string) bool {
	remoteUrl, _ := url.Parse(GetRemoteUrl())
	remoteUrl.Path = path.Join(remoteUrl.Path, "download")
	query := remoteUrl.Query()
	query.Set("name", name)
	query.Set("version", version)
	remoteUrl.RawQuery = query.Encode()
	if err := downloadFile(dir, remoteUrl.String(), remoteSha256); err != nil {
		fmt.Printf("下载文件到本地异常 %v\n", err)
		return false
	}
	return true
}

func downloadFile(filepath, url, remoteSha256 string) error {
	fmt.Println("请求下载地址 : " + url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Header.Get("Content-Type") == "application/octet-stream" {
		//是下载,获取下载文件名
		content := resp.Header.Get("Content-Disposition")
		fileName := getFileName(content)
		_ = os.MkdirAll(filepath, os.ModePerm)
		realFullName := path.Join(filepath, fileName)
		ioutil.WriteFile(realFullName+".sha256", []byte(remoteSha256), os.ModePerm)
		out, err := os.Create(realFullName)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
	} else {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println("下载异常 " + string(bodyBytes))
		return errors.New("错误的下载结果")
	}

	return nil
}

//attachment; filename=gg-example.zip
func getFileName(content string) string {
	split1 := strings.Split(content, ";")
	split2 := strings.Split(split1[1], "=")
	return split2[1]
}

func getLocalSha256(filePath string) (string, error) {
	var hashValue string
	file, err := os.Open(filePath)
	if err != nil {
		return hashValue, err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return hashValue, err
	}
	hashInBytes := hash.Sum(nil)
	hashValue = hex.EncodeToString(hashInBytes)
	return hashValue, nil
}

func getRemoteSha256(name, version string) string {
	var sha256Str string
	remoteUrl, _ := url.Parse(GetRemoteUrl())
	remoteUrl.Path = path.Join(remoteUrl.Path, "check")
	query := remoteUrl.Query()
	query.Set("name", name)
	query.Set("version", version)
	remoteUrl.RawQuery = query.Encode()

	resp, err := http.Get(remoteUrl.String())
	fmt.Println("请求sha256接口 : " + remoteUrl.String())
	if err != nil {
		fmt.Printf("请求 %s 异常 %v\n", remoteUrl.String(), err)
		return sha256Str
	}
	defer resp.Body.Close()

	result, err := common.ToResult(resp)
	if err != nil {
		fmt.Printf("解析结果异常 %s %v\n", remoteUrl.String(), err)
		return sha256Str
	}

	if !common.IsSuccess(result) {
		fmt.Printf("请求失败 %s %v 请求结果：%s\n", remoteUrl.String(), result, result.Message)
		return sha256Str
	}

	if err := common.ToData(result, &sha256Str); err != nil {
		fmt.Printf("转换结果异常 %v\n", err)
	}
	return sha256Str
}

func getTarget(name string, version string, result []GgResult) string {
	var it interface{}
	for _, r := range result {
		if r.Name == name {
			it = r
			break
		}
	}
	if it == nil {
		fmt.Println("没有找到相关包 " + name)
		return ""
	}

	exactIt := it.(GgResult)
	sort.Strings(exactIt.Versions)

	find := false
	//取当前最大版本
	if version == "" {
		version = exactIt.Versions[len(exactIt.Versions)-1]
		find = true
	} else {
		for _, i := range exactIt.Versions {
			if i == version {
				find = true
				break
			}
		}
	}
	if !find {
		fmt.Printf("没有找到包 %s 所对应的版本 %s\n", name, version)
		return ""
	}
	return version
}

func queryResult(name string, version string) []GgResult {
	var findResults []GgResult
	remoteUrl, _ := url.Parse(GetRemoteUrl())
	remoteUrl.Path = path.Join(remoteUrl.Path, "query")
	query := remoteUrl.Query()
	query.Set("name", name)
	query.Set("version", version)
	remoteUrl.RawQuery = query.Encode()
	// Get the data
	urlStr := remoteUrl.String()
	fmt.Println("请求查询地址 : " + urlStr)
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Printf("请求 %s 异常 %v\n", urlStr, err)
		return findResults
	}
	defer resp.Body.Close()

	result, err := common.ToResult(resp)
	if err != nil {
		fmt.Printf("解析结果异常 %s %v\n", urlStr, err)
		return findResults
	}

	if !common.IsSuccess(result) {
		fmt.Printf("请求失败 %s %v 请求结果：%s\n", urlStr, result, result.Message)
		return findResults
	}

	if err := common.ToData(result, &findResults); err != nil {
		fmt.Printf("转换结果异常 %v\n", err)
		return findResults
	}
	return findResults
}
