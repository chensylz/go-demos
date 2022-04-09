package tasks

import (
	"context"
	"fmt"
)

func HelloWorld(ctx context.Context, _ Config) {
	fmt.Println("hello world")
}
