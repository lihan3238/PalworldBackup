package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"github.com/robfig/cron"
)

const backupDir = "D:\\backup\\"
const sourceDir = "D:\\steamcmd\\steamapps\\common\\PalServer\\Pal\\Saved\\"
const maxBackups = 10

func main() {
	c := cron.New()

	// 每十分钟执行一次备份任务
	c.AddFunc("*/60 * * * *", func() {
		err := backup()
		if err != nil {
			fmt.Println("备份失败:", err)
		}
	})

	c.Start()

	// 等待程序运行
	select {}
}

func backup() error {
	// 创建备份目录
	err := os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		return err
	}

	// 生成时间戳
	timestamp := time.Now().Format("20060102150405")

	// 备份目标路径
	backupPath := filepath.Join(backupDir, "backup_"+timestamp) + "\\"

	// 执行备份命令
	// 执行备份命令
	cmd := exec.Command("cmd", "/c", "xcopy", sourceDir, backupPath, "/E")

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
