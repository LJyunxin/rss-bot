package main

import (
	"log"
	"os"

	"github.com/LJ-WorkSpace/feishu-RSS-bot/handlers"
	"github.com/LJ-WorkSpace/feishu-RSS-bot/models/dao"
	"github.com/LJ-WorkSpace/feishu-RSS-bot/services"
	"github.com/robfig/cron/v3"
)

func init() {
	dao.FileMux.Lock()
	defer dao.FileMux.Unlock()

	if _, err := os.Stat(dao.FilePath); os.IsNotExist(err) {
		dataFile, err := os.Create(dao.FilePath)
		if err != nil {
			log.Panic("create data file failed ")
		}
		defer dataFile.Close()
		log.Println("create data file successfully")
	}

}

func main() {
	c := cron.New()
	_, err := c.AddFunc("@hourly", handlers.StartPushSubscription)
	if err != nil {
		log.Fatal("cronjob create failed")
	}
	c.Start()

	engine := handlers.GinStart()
	err = engine.Run(":8080")
	if err != nil {
		log.Fatal("gin run failed:", err)
	}

	err = services.PushSubscription()
	if err != nil {
		log.Println("push subscription failed:", err)
	}

	log.Println("system initialized")
}
