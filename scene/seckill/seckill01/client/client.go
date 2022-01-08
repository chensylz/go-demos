package main

import (
	"fmt"
	"github.com/1278651995/go-demos/scene/seckill/seckill01/server/models"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-resty/resty/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	client := resty.New()
	client2 := resty.New()
	client.SetBaseURL("http://127.0.0.1:8001").
		SetTimeout(3 * time.Second)
	client2.SetBaseURL("http://127.0.0.1:8002").
		SetTimeout(3 * time.Second)
	go func() {
		for {
			ShowDBData()
			time.Sleep(100 * time.Millisecond)
		}
	}()
	go func() {
		for {
			CreateOrder(client)
			time.Sleep(time.Millisecond * 1)
		}
	}()
	go func() {
		for {
			CreateOrder(client2)
			time.Sleep(time.Millisecond * 1)
		}
	}()
	select {
	case <-time.After(30 * time.Second):
		return
	}
}

func CreateOrder(client *resty.Client) {
	body := map[string]interface{}{
		"username":       gofakeit.Name(),
		"merchandise_id": 1,
	}
	resp, err := client.R().SetBody(body).Post("/order")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp.String())
}

func ShowDBData() {
	db := NewMysql()
	var merchandise models.Merchandise
	db.Where("deleted_at is null").First(&merchandise)
	var count int64
	db.Model(&models.Order{}).Where("merchandise_id = ?", merchandise.ID).Count(&count)
	fmt.Printf("当前%s商品已售出%d， 库存还剩%d \n", merchandise.Name, count, merchandise.Stock)
}

func NewMysql() *gorm.DB {
	dsn := "root:123456@tcp(127.0.0.1:3306)/chensy?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
