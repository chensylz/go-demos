package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/1278651995/go-demos/kits/cron_ui/ui"
	"github.com/robfig/cron/v3"

	"github.com/1278651995/go-demos/kits/cron_ui/tasks"
)

type JobFunc func(ctx context.Context, config tasks.Config)

type Job struct {
	Spec string
	Func JobFunc
	Name string
}

// 需要执行的定时任务.
var jobs = []Job{
	// 每分钟.
	{Spec: "@every 1m", Func: tasks.HelloWorld, Name: "HelloWorld"},
}

var entriesName = map[cron.EntryID]ui.JobInfo{}

// SetupCronJob 初始化cronjob.
func SetupCronJob(config tasks.Config) *cron.Cron {
	c := cron.New()
	entriesName = make(map[cron.EntryID]ui.JobInfo)
	for _, job := range jobs {
		f := job.Func
		entryID, err := c.AddFunc(job.Spec, func() { f(context.Background(), config) })
		if err != nil {
			log.Fatalln(err)
		}
		entriesName[entryID] = ui.JobInfo{
			Name: job.Name,
			Spec: job.Spec,
		}
	}
	return c
}

func main() {
	cronjob := SetupCronJob(tasks.Config{})
	cronjob.Start()
	go func() {
		cronUI := ui.NewCronUI(cronjob, entriesName)
		cronUI.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx := cronjob.Stop()
	log.Println("Shutting down cron...")
	select {
	case <-time.After(10 * time.Second):
		log.Fatal("Cron forced to shutdown...")
	case <-ctx.Done():
		log.Println("Cron exiting...")
	}
}
