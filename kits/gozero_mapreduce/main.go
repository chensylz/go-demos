package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tal-tech/go-zero/core/mr"
)

func main() {
	urls := []string{
		"https://www.baidu.com",
		"https://kezaihui.com",
		"https://zaihuiba.com",
	}
	client := resty.New()
	rr := mr.Map(func(source chan<- interface{}) {
		for _, url := range urls {
			source <- url
		}
	}, func(item interface{}, writer mr.Writer) {
		url := item.(string)
		resp, err := client.R().Get(url)
		if err == nil {
			writer.Write(resp.Body())
		}
	})
	for content := range rr {
		fmt.Println(content)
	}
}
