package main

import (
	"log"
	"os"
	"time"

	"github.com/LJ-WorkSpace/feishu-RSS-bot/handlers"
	"github.com/LJ-WorkSpace/feishu-RSS-bot/models/dao"
	"github.com/robfig/cron/v3"
)

func init() {
	dao.FileMux.Lock()
	defer dao.FileMux.Unlock()

	if _, err := os.Stat(dao.DataFilePath); os.IsNotExist(err) {
		dataFile, err := os.Create(dao.DataFilePath)
		if err != nil {
			log.Panic("create data file failed ")
		}
		defer dataFile.Close()
		log.Println("create data file successfully")
	}

}

func loggerInit() {
	logFile, err := os.OpenFile(dao.LogFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Panicln("logger init fail:", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("=======================", time.Now(), "===========================")
}

func main() {
	loggerInit()
	c := cron.New()
	_, err := c.AddFunc("@hourly", handlers.PushSubscription)
	if err != nil {
		log.Fatal("cronjob create failed")
	}
	c.Start()

	engine := handlers.GinStart()
	err = engine.Run(":8080")
	if err != nil {
		log.Fatal("gin run failed:", err)
	}

	log.Println("system initialized")
}
