package api

import (
	"app/MQdemo/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status": 404,
		"error":  "页面不存在",
	})
}

func CheckError(err error, msg string) {
	if err != nil {
		log.Println(err)
		defer log.Println(msg)
		panic(err)
	}
}

func send(buyerId int, goodId int, sellerId int, quantity int) (rs bool) {
	rs = false

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	CheckError(err, "Can't connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	CheckError(err, "Can't create a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"order",
		true,
		false,
		false,
		false,
		nil)
	CheckError(err, "Can't declare a queue")

	order := model.Order{buyerId, goodId, sellerId, quantity}
	body, err := json.Marshal(order)
	CheckError(err, "Error encoding JSON")

	err = channel.Publish("",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	CheckError(err, "Can't publish a message")

	log.Printf("send: %s", string(body))

	rs = true
	return
}

func Order(c *gin.Context) {
	state := make(map[string]interface{})
	buyerId := c.Request.FormValue("buyer_id")
	goodId := c.Request.FormValue("good_id")
	sellerId := c.Request.FormValue("seller_id")
	quantity := c.Request.FormValue("quantity")
	//纯数字
	reg := "^\\d+$"
	result1, _ := regexp.MatchString(reg, buyerId)
	result2, _ := regexp.MatchString(reg, goodId)
	result3, _ := regexp.MatchString(reg, sellerId)
	result4, _ := regexp.MatchString(reg, quantity)
	if buyerId == "" || goodId == "" || sellerId == "" || quantity == "" {
		state["state"] = 2
		state["text"] = "参数不能有空！"
	} else if !result1 || !result2 || !result3 || !result4 {
		state["state"] = 3
		state["text"] = "参数格式有误！"
	} else {
		b, _ := strconv.Atoi(buyerId)
		g, _ := strconv.Atoi(goodId)
		s, _ := strconv.Atoi(sellerId)
		q, _ := strconv.Atoi(quantity)

		if send(b, g, s, q) {
			state["state"] = 1
			state["text"] = "入队成功！"
		} else {
			state["state"] = 0
			state["text"] = "入队失败！"
		}
	}
	sjson,_ :=json.Marshal(state)
	c.String(http.StatusOK, "%v", string(sjson))
}
