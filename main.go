package main

import (
	"challenging/service"
	"fmt"
)

func main() {
	err := service.CalculateStatus("./data/data.json")
	if err != nil {
		fmt.Println(err)
	}
}
