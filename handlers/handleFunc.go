package handlers

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ljyunxin/rss-bot/models"
	"github.com/ljyunxin/rss-bot/models/dao"
	"github.com/ljyunxin/rss-bot/services"
	"github.com/mmcdole/gofeed"
)

func AddSubscription(ctx *gin.Context) {
	var newSubscription dao.Subscription
	err := ctx.ShouldBindBodyWith(&newSubscription, binding.JSON)
	if err != nil {
		ctx.JSON(400, gin.H{
			"msg": "binding body failed",
			"err": err,
		})
		return
	}

	data, err := models.GetDataSlice()
	if err != nil {
		ctx.JSON(500, gin.H{
			"msg": "get file data failed",
			"err": err,
		})
		return
	}

	newSubscription.UpdatedAt = time.Now()
	data = append(data, newSubscription)

	if err = models.UpdatedDataFile(data); err != nil {
		ctx.JSON(500, gin.H{
			"msg": "failed to update subscription file",
			"err": err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"msg": "add subscription success",
	})
}

func GetSubscription(ctx *gin.Context) {
	data, err := models.GetDataSlice()
	if err != nil {
		ctx.JSON(500, gin.H{
			"msg": "get data slice failed",
			"err": err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"data": data,
	})
}

func DeleteSubscription(ctx *gin.Context) {
	var item dao.Subscription
	err := ctx.ShouldBindBodyWith(&item, binding.JSON)
	if err != nil {
		ctx.JSON(400, gin.H{
			"msg": "binding json failed",
			"err": err,
		})
		return
	}

	data, err := models.GetDataSlice()
	if err != nil {
		ctx.JSON(500, gin.H{
			"msg": "get data slice failed",
			"err": err,
		})
		return
	}

	for index, i := range data {
		if i.Name == item.Name {
			data = append(data[:index], data[index+1:]...)
			err = models.UpdatedDataFile(data)
			if err != nil {
				ctx.JSON(500, gin.H{
					"msg": "update data file failed",
					"err": err,
				})
				return
			}

			ctx.JSON(200, gin.H{
				"msg": "delete successfully",
			})
			return
		}
	}

	ctx.JSON(400, gin.H{
		"msg": "data not found",
	})
}

func UpdateWebhook(ctx *gin.Context) {
	var webhookString dao.WebhookString
	err := ctx.ShouldBindJSON(&webhookString)
	if err != nil {
		ctx.JSON(400, gin.H{
			"msg": "bind json failed",
			"err": err,
		})
		return
	}

	dao.Webhook = webhookString.Webhook

	ctx.JSON(200, gin.H{
		"msg": "update webhook successfully",
	})
}

func PushSubscription() {
	origins, err := models.GetDataSlice()
	if err != nil {
		log.Println("get data slice failed", err)
		return
	}

	fp := gofeed.NewParser()
	for index, data := range origins {
		err = services.PushFeed(&data, fp)
		if err != nil {
			log.Println("push feed failed:", data.Name, err)
			continue
		}

		origins[index].UpdatedAt = time.Now()
	}

	err = models.UpdatedDataFile(origins)
	if err != nil {
		log.Println("updated data file failed:", err)
		return
	}
}
