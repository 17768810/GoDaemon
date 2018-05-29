package main

import (
	"flag"
	"path/filepath"
	"strconv"
	"strings"

	"fmt"
	"log"
	"os"
	"time"

	"github.com/jimlawless/cfg"
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
	maxSize, _ := strconv.Atoi(maxFileSize)
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
			err := runDaemon(dir, maxSize)
			//fmt.Println( files )
			if err != nil {
				logger.Println(dir+" runDeamon err:", err)
			}
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}

}

func runDaemon(inputDir string, maxSize int) (err error) {
	files, err := getList(inputDir)
	//fmt.Println( files )
	if err == nil {
		delFile(files, maxSize)
	} else {
		logger.Println(err)
	}
	return
}

func getList(dir string) (files []string, err error) {
	var domains_filter []string
	var domains_filter_cfg string
	var gz_files []string
	match := fmt.Sprintf("%s/*.*", dir)
	gz_files, err = filepath.Glob(match)
	if err != nil {
		logger.Println("list gz files err:", gz_files)
		return
	}
	domains_filter_cfg = cfg_map["domains_filter"]
	domains_filter = strings.Split(domains_filter_cfg, ",")

	for _, fname := range gz_files {
		var in = inArray(domains_filter, filepath.Base(fname))
		//fmt.Println( "fname:",fname,"in?",in)
		if in {
			files = append(files, fname)
			//fmt.Println( "find one", files )
			continue
		}
	}
	return
}

// maxFileSize ,unit Kb
func delFile(files []string, maxFileSize int) {
	// today := time.Now().Format("2006-01-02")
	// goadf, err := os.OpenFile(cfg_map["output_dir"]+today+".log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	// if err != nil {
	// 	logger.Println("openfile err:", err)
	// 	return
	// }
	// defer goadf.Close()

	// var br *bufio.Reader

	for _, name := range files {
		fileSize, err := os.Stat(name)
		if err != nil {
			logger.Println("delFile err:", err)
		} else {
			//maxSize
			size := fileSize.Size()
			maxSize := int64(maxFileSize * 1024)
			if size > maxSize {
				err = os.Remove(name)
				if err != nil {
					logger.Println("delFile err:", err)
				} else {
					fmt.Println("fileName:"+name+" , fileSize:", size)
					logger.Println("fileName:"+name+" , fileSize:", size)
				}
			}
		}
	}
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
