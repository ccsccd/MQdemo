package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main(){
	go httpPostForever()
	go httpPostForever()
	go httpPostForever()
	go httpPostForever()
	go httpPostForever()
	time.Sleep(5*time.Minute)
}
func httpPostForever(){
	for x := range time.Tick(1000*time.Millisecond) {
		httpPost(x)
	}
}
func httpPost(t time.Time) {
	resp, err := http.Post("http://localhost/shop/order",
		"application/x-www-form-urlencoded",
		strings.NewReader("buyer_id=1&good_id=10020&seller_id=66666&quantity=1"))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(body))
}