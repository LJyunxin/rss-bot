package dao

import (
	"os"
	"sync"
	"time"
)

const FilePath = "rss.json"
const Template = `{
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
