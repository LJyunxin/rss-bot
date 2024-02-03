package services

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/LJ-WorkSpace/feishu-RSS-bot/models"
	"github.com/LJ-WorkSpace/feishu-RSS-bot/models/dao"
	"github.com/mmcdole/gofeed"
)

func PushSubscription() error {
	origins, err := models.GetDataSlice()
	if err != nil {
		return err
	}

	fp := gofeed.NewParser()
	for index, data := range origins {
		err = pushFeed(&data, fp)
		if err != nil {
			log.Println("push feed failed:", data.Name, err)
		}
		origins[index].UpdatedAt = time.Now()
	}

	err = models.UpdatedDataFile(origins)
	if err != nil {
		return err
	}

	return nil
}

func pushFeed(data *dao.Subscription, fp *gofeed.Parser) error {
	feed, err := fp.ParseURL(data.Url)
	if err != nil {
		return err
	}

	for _, item := range feed.Items {
		if item.PublishedParsed.Unix() <= data.UpdatedAt.Unix() {
			break
		}

		err = pushToLark(data.Name, item)
		if err != nil {
			log.Println("push to lark failed:", data.Name, err)
		}
	}

	return nil
}

func pushToLark(name string, item *gofeed.Item) error {
	temp := fmt.Sprintf(dao.Template, name, item.Title, item.Link)

	resp, err := http.Post(dao.Webhook, "application/json", strings.NewReader(temp))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintln("push failed:http response code", resp.StatusCode))
	}

	return nil
}
