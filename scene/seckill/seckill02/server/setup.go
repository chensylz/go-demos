package main

import (
	"github.com/1278651995/go-demos/scene/seckill/seckill01/server/models"
	"github.com/alicebob/miniredis/v2"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Server *gin.Engine
	Redis  *miniredis.Miniredis
}

func NewApp() *App {
	return &App{
		DB:     NewMysql(),
		Server: newEngine(),
		Redis:  newRedis(),
	}
}

func NewMysql() *gorm.DB {
	dsn := "root:123456@tcp(127.0.0.1:3306)/chensy?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	_ = db.AutoMigrate(&models.Order{}, &models.Merchandise{})
	return db
}

func newEngine() *gin.Engine {
	return gin.Default()
}

func newRedis() *miniredis.Miniredis {
	return miniredis.NewMiniRedis()
}
