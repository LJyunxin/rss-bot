package models

import (
	"encoding/json"
	"log"
	"os"

	"github.com/LJ-WorkSpace/feishu-RSS-bot/models/dao"
)

func GetDataSlice() ([]dao.Subscription, error) {
	var subscription []dao.Subscription
	dao.FileMux.RLock()
	defer dao.FileMux.RUnlock()

	data, err := os.ReadFile(dao.FilePath)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &subscription); err != nil {
		log.Println("json unmarshal failed")
	}

	return subscription, nil
}

func UpdatedDataFile(data []dao.Subscription) error {
	dao.FileMux.Lock()
	defer dao.FileMux.Unlock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err = os.WriteFile(dao.FilePath, jsonData, 0644); err != nil {
		return err
	}

	return nil
}
