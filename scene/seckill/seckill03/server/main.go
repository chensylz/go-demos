package main

import (
	"context"
	"github.com/1278651995/go-demos/scene/seckill/seckill03/ent/merchandise"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateOrderRequest struct {
	UserID        int  `json:"user_id"`
	MerchandiseID uint `json:"merchandise_id"`
}

func (s *App) CreateOrder(c *gin.Context) {
	ctx := context.Background()
	var request CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}
	effectRow, err := s.Ent.Merchandise.Update().Where(
		merchandise.ID(int(request.MerchandiseID)),
		merchandise.StockGT(0),
	).AddStock(-52).Save(ctx)
	if err != nil || effectRow == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "已售罄"})
		return
	}
	err = s.Ent.Order.Create().SetUserID(request.UserID).
		SetMerchandiseID(int(request.MerchandiseID)).Exec(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "抢购失败, 请重试!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "抢购成功"})
}

func (s *App) InitDBData() {
	s.Ent.Merchandise.Update().SetStock(100).ExecX(context.Background())
	s.Ent.Order.Delete().ExecX(context.Background())
}

func main() {
	app := NewApp()
	app.InitDBData()

	app.Server.POST("/order", app.CreateOrder)

	go func() {
		_ = app.Server.Run("127.0.0.1:8001")
	}()
	_ = app.Server.Run("127.0.0.1:8002")
}
