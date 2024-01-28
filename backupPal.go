package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"github.com/robfig/cron"
)

const configFile = "config.json"
const maxBackups = 100

// Config 结构体表示配置信息
type Config struct {
	BackupDir string `json:"backupDir"`
	SourceDir string `json:"sourceDir"`
}

func main() {
	// 读取配置文件，如果不存在则创建默认配置
	config, err := readOrCreateConfig()
	if err != nil {
		fmt.Println("读取或创建配置文件失败:", err)
		return
	}

	c := cron.New()

	// 每分钟执行一次备份任务
	c.AddFunc("*/60 * * * *", func() {
		err := backup(config.BackupDir, config.SourceDir)
		if err != nil {
			fmt.Println("备份失败:", err)
		}
	})

	c.Start()

	// 等待程序运行
	select {}
}

func backup(backupDir, sourceDir string) error {
	// 创建备份目录
	err := os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		return err
	}

	// 生成时间戳
	timestamp := time.Now().Format("20060102150405")

	// 备份目标路径
	backupPath := filepath.Join(backupDir, "backup_"+timestamp) + "\\"
	sourcePath := filepath.Join(sourceDir)
	// 执行备份命令
	cmd := exec.Command("xcopy", sourcePath, backupPath, "/E", "/I", "/H", "/Y")
	fmt.Println("执行命令:", cmd.String())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("备份失败: %w", err)
	}

	// 删除最老的备份，保留最多maxBackups个备份
	backups, err := getBackupList(backupDir)
	if err != nil {
		return err
	}

	if len(backups) > maxBackups {
		sort.Sort(ByTime(backups))
		for i := 0; i < len(backups)-maxBackups; i++ {
			err := os.RemoveAll(filepath.Join(backupDir, backups[i]))
			if err != nil {
				return err
			}
			fmt.Println("删除备份:", backups[i])
		}
	}

	fmt.Println("备份成功:", backupPath)
	return nil
}

func readOrCreateConfig() (Config, error) {
	config, err := readConfig()
	if err != nil && os.IsNotExist(err) {
		// 如果配置文件不存在，创建默认配置并写入文件
		defaultConfig := Config{BackupDir: "C:/backup/", SourceDir: "C:/steamcmd/steamapps/common/PalServer/Pal/Saved/"}
		err := writeConfig(defaultConfig)
		if err != nil {
			return config, err
		}
		return defaultConfig, nil
	}
	return config, err
}

func readConfig() (Config, error) {
	var config Config

	file, err := os.Open(configFile)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func writeConfig(config Config) error {
	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	return err
}

// 获取备份目录下的备份列表
func getBackupList(dir string) ([]string, error) {
	var backups []string

	files, err := os.ReadDir(dir)
	if err != nil {
		return backups, err
	}

	for _, file := range files {
		if file.IsDir() {
			backups = append(backups, file.Name())
		}
	}

	return backups, nil
}

// ByTime 实现备份列表按时间排序
type ByTime []string

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i] < a[j] }
