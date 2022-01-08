package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"github.com/1278651995/go-demos/scene/seckill/seckill01/server/models"
)

const MerchandiseKey = "merchandise:%d"

type CreateOrderRequest struct {
	UserName      string `json:"username"`
	MerchandiseID uint   `json:"merchandise_id"`
}

func (s *App) SyncStock(merchandiseID int) {
	var merchandise models.Merchandise
	s.DB.First(&merchandise, 1)
	_ = s.Redis.Set(fmt.Sprintf(MerchandiseKey, merchandiseID), strconv.Itoa(merchandise.Stock))
}

func (s *App) CreateOrder(c *gin.Context) {
	var request CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}
	stockStr, _ := s.Redis.Get(fmt.Sprintf(MerchandiseKey, request.MerchandiseID))
	stock, _ := strconv.Atoi(stockStr)
	if stock <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "已售罄"})
		return
	}
	stock, _ = s.Redis.Incr(fmt.Sprintf(MerchandiseKey, request.MerchandiseID), -1)
	fmt.Printf("当前还剩库存%d \n", stock)
	if stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "已售罄"})
		return
	}
	order := models.Order{
		UserName:      request.UserName,
		MerchandiseID: request.MerchandiseID,
	}
	resp := s.DB.Create(&order)
	if resp.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "抢购失败, 请重试!"})
		return
	}
	s.DB.Model(&models.Merchandise{}).Where("id = ?", 1).UpdateColumn("stock", stock)
	c.JSON(http.StatusOK, gin.H{"message": "抢购成功"})
}

func (s *App) InitDBData() {
	s.DB.Where("1 = 1").Delete(&models.Merchandise{})
	s.DB.Where("1 = 1").Delete(&models.Order{})
	merchandise := models.Merchandise{
		Name:  "原味吮指炸鸡块",
		Stock: 100,
	}
	s.DB.Create(&merchandise)
}

func main() {
	app := NewApp()
	app.SyncStock(1)
	//app.InitDBData()

	app.Server.POST("/order", app.CreateOrder)

	go func() {
		_ = app.Server.Run("127.0.0.1:8001")
	}()
	_ = app.Server.Run("127.0.0.1:8002")
}
