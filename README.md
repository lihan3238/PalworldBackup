# PalworldBackup

## 简介

PalServer幻兽帕鲁的服务器优化很差，使用社区自发制作的优化程序可能会有风险，因此需要定期备份。

这是一个针对windows11上steamcmd运行的PalServer的存档定期备份小程序。

每分钟备份一次，只保留最新的100次备份

## 使用方法

下载`backupPal.exe`后，双击执行，在同一目录下自动创建配置文件`config.json`，修改为自己的存档路径与备份路径后重新运行即可。

修改格式:
```json
{
    "backupDir": "D:/backup/",
    "sourceDir": "D:/steamcmd/steamapps/common/PalServer/Pal/Saved/"
}
```

