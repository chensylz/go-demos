package main

/**
https://segmentfault.com/a/1190000016412013
 */

import (
	"log"
	"net/http"
	_ "net/http/pprof" // focus on here
	"time"
)

var testData []string

func AddData(str string) string {
	data := []byte(str)
	sData := string(data)
	testData = append(testData, sData)

	return sData
}

func main() {
	go func() {
		for {
			log.Println(AddData("https://github.com/EDDYCJY"))
			time.Sleep(1 * time.Second)
		}
	}()

	http.ListenAndServe("0.0.0.0:6060", nil) // http://127.0.0.1:6060/debug/pprof/
}
