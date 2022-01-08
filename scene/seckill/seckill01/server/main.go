package main

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/1278651995/go-demos/scene/seckill/seckill01/server/models"
)

type CreateOrderRequest struct {
	UserName      string `json:"username"`
	MerchandiseID uint   `json:"merchandise_id"`
}

func (s *App) CreateOrder(c *gin.Context) {
	var request CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}
	var merchandise models.Merchandise
	s.DB.First(&merchandise, request.MerchandiseID)

	order := models.Order{
		UserName:      request.UserName,
		MerchandiseID: merchandise.ID,
	}
	merchandise.Stock--
	if merchandise.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "已售罄"})
		return
	}
	resp := s.DB.Create(&order)
	if resp.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "抢购失败, 请重试!"})
		return
	}
	s.DB.Save(&merchandise)
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
	//app.InitDBData()

	app.Server.POST("/order", app.CreateOrder)

	go func() {
		_ = app.Server.Run("127.0.0.1:8001")
	}()
	_ = app.Server.Run("127.0.0.1:8002")
}
