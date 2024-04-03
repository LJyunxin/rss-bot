package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ljyunxin/rss-bot/models/dao"
	"github.com/mmcdole/gofeed"
)

func PushFeed(data *dao.Subscription, fp *gofeed.Parser) error {
	feed, err := fp.ParseURL(data.Url)
	if err != nil && time.Since(data.UpdatedAt) > 24*time.Hour {
		pushErr := PushFailToLark(data.Name, err)
		return pushErr
	}
	if err != nil {
		return err
	}

	for _, item := range feed.Items {
		if item.PublishedParsed.Unix() <= data.UpdatedAt.Unix() {
			break
		}

		err = PushToLark(data.Name, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func PushFailToLark(name string, msg error) error {
	msgStr := strings.Replace(msg.Error(), `"`, `\"`, -1)
	temp := fmt.Sprintf(dao.MessageTemplate, name, "fetch feed failed", msgStr)

	resp, err := http.Post(dao.Webhook, "application/json", strings.NewReader(temp))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintln("push to lark failed:http response code", resp.StatusCode))
	}

	return nil
}

func PushToLark(name string, item *gofeed.Item) error {
	temp := fmt.Sprintf(dao.FeedTemplate, name, item.Title, item.Link)

	resp, err := http.Post(dao.Webhook, "application/json", strings.NewReader(temp))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintln("push to lark failed:http response code", resp.StatusCode))
	}

	return nil
}
