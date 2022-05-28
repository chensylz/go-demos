package main

import (
	"context"
	"log"

	"github.com/1278651995/go-demos/scene/seckill/seckill03/ent"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Server *gin.Engine
	Ent    *ent.Client
}

func NewApp() *App {
	return &App{
		Ent:    NewMysql(),
		Server: newEngine(),
	}
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

func newEngine() *gin.Engine {
	return gin.Default()
}
