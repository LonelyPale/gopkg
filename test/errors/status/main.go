package main

import (
	"fmt"

	"github.com/lonelypale/gopkg/errors/status"
)

func main() {
	e1 := status.Error("test error")
	fmt.Println(e1)
}
