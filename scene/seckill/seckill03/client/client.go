package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql"

	"github.com/1278651995/go-demos/scene/seckill/seckill03/ent"
)

func main() {
	client := resty.New()
	client2 := resty.New()
	client.SetBaseURL("http://127.0.0.1:8001").
		SetTimeout(3 * time.Second)
	client2.SetBaseURL("http://127.0.0.1:8002").
		SetTimeout(3 * time.Second)
	mysqlClient := NewMysql()
	go func() {
		for {
			ShowDBData(mysqlClient)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	go func() {
		for {
			CreateOrder(client, 1)
			time.Sleep(time.Millisecond * 1)
		}
	}()
	go func() {
		for {
			CreateOrder(client2, 2)
			time.Sleep(time.Millisecond * 1)
		}
	}()
	select {
	case <-time.After(30 * time.Second):
		return
	}
}

func CreateOrder(client *resty.Client, id int) {
	body := map[string]interface{}{
		"user_id":        id,
		"merchandise_id": 1,
	}
	resp, err := client.R().SetBody(body).Post("/order")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp.String())
}

func ShowDBData(client *ent.Client) {
	m := client.Merchandise.Query().FirstX(context.Background())
	fmt.Printf("当前商品库存还剩%d \n", m.Stock)
}

func NewMysql() *ent.Client {
	dsn := "root:123456@tcp(127.0.0.1:3306)/chensy?charset=utf8mb4&parseTime=True&loc=Local"
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}
