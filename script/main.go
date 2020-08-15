// create by d1y<chenhonzhou@gmail.com>

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/d1y/b23/ffmpeg"
	"github.com/d1y/plist2json"
)

// 当前路径
var curr = ""

// 输出结果路径
var resultPath = ""

func main() {
	fmt.Println("auto run")
	auto()
}

func init() {
	x, err := os.Getwd()
	if err != nil {
		log.Fatalln("get current path is error")
	}
	curr = x
	resultPath = path.Join(x, "../results")
	ensureDir(resultPath)
}

// readDataDir 读取数据文件夹
func readDataDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path)
}

func auto() {
	var dataPath = path.Join(curr, "../data")
	l, e := readDataDir(dataPath)
	if e != nil {
		log.Fatalln("read dir is error")
	}
	for _, item := range l {
		var filename = item.Name()
		var ext = filepath.Ext(filename)
		var rawfullpath = path.Join(dataPath, filename)
		if !item.IsDir() {
			switch ext {
			case ".caf":
				var Nfile = strings.Replace(filename, ".caf", "", len(filename)-3)
				var R = path.Join(resultPath, Nfile)
				flag := ffmpeg.ConvertFormat2mp3(rawfullpath, R)
				var text = "成功"
				if !flag {
					text = "失败"
				}
				fmt.Println(fmt.Sprintf("转换结果: %v, 输出结果: %v", text, R))
			case ".plist":
				var Nfile = strings.Replace(filename, "plist", "json", len(filename)-4)
				var R = path.Join(resultPath, Nfile)
				e := plist2json.Easy(rawfullpath, R)
				if e != nil {
					log.Fatalln(e)
				}
				fmt.Println(fmt.Sprintf("转换成功: %v", R))
			default:
				var R = path.Join(resultPath, filename)
				fmt.Println(fmt.Sprintf("原样复制: %v", R))
				copy(rawfullpath, R)
			}
		}
	}
}

// ensureDir 自动创建文件
func ensureDir(fileName string) {
	dirName := fileName
	// dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, 0755)
		if merr != nil {
			panic(merr)
		}
	}
}

// copy 复制文件
func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// Exists 判断所给路径文件/文件夹是否存在
// https://www.php.cn/be/go/446030.html
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
// https://www.php.cn/be/go/446030.html
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
// https://www.php.cn/be/go/446030.html
func IsFile(path string) bool {
	return !IsDir(path)
}
