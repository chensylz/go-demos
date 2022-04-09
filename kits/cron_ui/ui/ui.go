package ui

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

const DateTimeLayout = "2006-01-02 15:04:05"

type JobInfo struct {
	Name string
	Spec string
}

type CronUI struct {
	cron          *cron.Cron
	entriesInfo   map[cron.EntryID]JobInfo
	manualHistory map[cron.EntryID][]string
	engine        *gin.Engine
	host          string
}

func NewCronUI(c *cron.Cron, info map[cron.EntryID]JobInfo, host string) *CronUI {
	return &CronUI{
		cron:          c,
		host:          host,
		engine:        gin.Default(),
		entriesInfo:   info,
		manualHistory: make(map[cron.EntryID][]string),
	}
}

func (c *CronUI) Start() {
	c.setUp()
	err := c.engine.Run(c.host)
	log.Fatalln(err)
}

func (c *CronUI) Stop() {
}

func (c *CronUI) setUp() {
	c.engine.Use(CORS())
	c.engine.GET("/", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html;charset=UTF-8", []byte(template))
	})
	c.engine.GET("/cronjob", c.getCronJob)
	c.engine.POST("/cronjob/:entry_id", c.startCronJob)
}

func (c *CronUI) getCronJob(ctx *gin.Context) {
	entries := c.cron.Entries()
	result := make([]map[string]interface{}, len(entries))
	for i, entry := range entries {
		result[i] = make(map[string]interface{})
		result[i]["id"] = entry.ID
		result[i]["name"] = entry.ID
		if info, ok := c.entriesInfo[entry.ID]; ok {
			result[i]["name"] = info.Name
			result[i]["spec"] = info.Spec
		}
		if histories, ok := c.manualHistory[entry.ID]; ok {
			result[i]["manual_history"] = histories
		}
		result[i]["pre_time"] = entry.Prev.Format(DateTimeLayout)
		result[i]["next_time"] = entry.Next.Format(DateTimeLayout)

	}
	ctx.JSON(http.StatusOK, result)
}

func (c *CronUI) startCronJob(ctx *gin.Context) {
	entryIDStr := ctx.Param("entry_id")
	entryID, _ := strconv.Atoi(entryIDStr)
	entry := c.cron.Entry(cron.EntryID(entryID))
	if entry.Next.Before(time.Now().Add(10 * time.Second)) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "10秒后启动的任务无法手动启动",
		})
		return
	}
	go func(e *cron.Entry) {
		e.Job.Run()
	}(&entry)
	histories, ok := c.manualHistory[entry.ID]
	if !ok {
		histories = make([]string, 0, 1)
	}
	histories = append(histories, time.Now().Format(DateTimeLayout))
	c.manualHistory[entry.ID] = histories
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func CORS() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 允许 Origin 字段中的域发送请求
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Max-Age", "86400")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length，Apitoken")
		context.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Headers")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(200)
		} else {
			context.Next()
		}
	}
}
