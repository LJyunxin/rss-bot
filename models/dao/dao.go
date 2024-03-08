package dao

import (
	"os"
	"sync"
	"time"
)

const DataFilePath = "rss.json"
const LogFilePath = "log.txt"
const FeedTemplate = `{
	"msg_type": "post",
	"content": {
		"post": {
			"zh_cn": {
				"title": "%s",
				"content": [
					[
						{
							"tag": "a",
							"text": "%s",
							"href": "%s"
						}
					]
				]
			}
		}
	}
}`

const MessageTemplate = `{
	"msg_type": "post",
	"content": {
		"post": {
			"zh_cn": {
				"title": "%s",
				"content": [
					[
						{
							"tag": "text",
							"text": "%s:\n"
						},
						{
							"tag": "text",
							"text": "%s"
						}
					]
				]
			}
		}
	}
}`

var Webhook = os.Getenv("webhook")
var FileMux = sync.RWMutex{}

type WebhookString struct {
	Webhook string `json:"webhook"`
}

type Subscription struct {
	Name      string `json:"name"`
	Url       string `json:"url"`
	UpdatedAt time.Time
}
