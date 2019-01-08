package main

import (
	"flag"
	"strconv"
	"strings"

	"fmt"
	"log"
	"os"
	"time"

	"github.com/jimlawless/cfg"
	"io/ioutil"
	"sort"
)

var logger *log.Logger

var logfile *os.File

var cfg_file string

var cfg_map map[string]string

func main() {
	//load config
	flag.StringVar(&cfg_file, "conf", "./default.conf", "Input the config file use --conf=[file path]")
	flag.Parse()

	loadCfg(cfg_file)

	initLogger(cfg_map["log_file"])

	var pid = os.Getpid()

	logger.Println("config file:", cfg_file)
	logger.Println("LogFileDaemon process start,PID:", pid)

	maxFileSize := cfg_map["max_filesize"]
	maxSize, _ := strconv.ParseInt(maxFileSize, 0, 0)

	// input dir
	inputDir := cfg_map["input_dir"]
	dirAyy := strings.Split(inputDir, ",")

	// The time interval,unit second
	intervalStr := cfg_map["interval"]
	interval, _ := strconv.Atoi(intervalStr)
	if interval == 0 {
		interval = 10
	}
	for {
		for _, dir := range dirAyy {
			err := runDaemon(dir, maxSize*1024)
			//fmt.Println( files )
			if err != nil {
				logger.Println(dir+" runDeamon err:", err)
			}
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}

}

// 需要的文件信息字段
type FileInfoExt struct {
	FileSize int64  `json:"filesize"` //	文件大小
	FilePath string `json:"filepath"` //	文件路径
	FileDate int64  `json:"filedate"` //	 文件日期
}

func runDaemon(inputDir string, maxSize int64) (err error) {

	//遍历打印所有的文件名
	var files []FileInfoExt
	files, _ = GetAllFile(inputDir, files)

	//for i, v := range files {
	//	fmt.Println(fmt.Sprintf("index:%s , FilePath:%s , FileDate:%s , FileSize:%s", i, v.FilePath, v.FileDate, v.FileSize))
	//}

	var dirSize int64
	for _, v := range files {
		dirSize += v.FileSize
	}

	// 目录的文件总大小，大于设置值
	if dirSize > maxSize {
		// 按文件创建日期，从小到大排序
		sort.Slice(files, func(i, j int) bool {
			return files[i].FileDate < files[j].FileDate
		})
		subValue := dirSize - maxSize // 差值
		var deletedFileSize int64     // 已删除文件大小

		for _, v := range files {
			if deletedFileSize > subValue {
				break // 已删除文件的大小，大于差值
			}
			f, _ := os.Stat(v.FilePath)
			if !f.IsDir() {
				err := os.Remove(v.FilePath)
				if err != nil {
					//如果删除失败则输出 file remove Error!
					fmt.Println(v.FilePath + " file remove Error!")
					//输出错误详细信息
					fmt.Printf("%s", err)
				} else {
					//删除成功!
					deletedFileSize += v.FileSize
				}
			}
		}

		fmt.Println(fmt.Sprintf("file remove ok,total deleted size %s Kb!", deletedFileSize/1024))
	}

	return
}

func GetAllFile(pathname string, s []FileInfoExt) ([]FileInfoExt, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = GetAllFile(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			fileInfoExt := FileInfoExt{
				FilePath: fullName,
				FileDate: fi.ModTime().Unix(),
				FileSize: fi.Size(),
			}
			s = append(s, fileInfoExt)
		}
	}
	return s, nil
}

func initLogger(log_file string) {
	//logger init
	logfile, _ = os.OpenFile(log_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logger = log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
}

func loadCfg(filepath string) {
	cfg_map = make(map[string]string)
	cfg_err := cfg.Load(filepath, cfg_map)
	if cfg_err != nil {
		fmt.Println("load config file err")
	}
}

func inArray(array []string, findme string) bool {
	for _, v := range array {
		if strings.Contains(findme, v) {
			return true
		}
	}
	return false
}
