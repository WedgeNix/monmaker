package main

import (
	"fmt"
	"log"
	"time"

	"github.com/WedgeNix/monmaker/maker"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	now := time.Now()
	from := now.AddDate(0, -6, 0)

	resp, err := maker.Ship(from, now, "MyPakage").GetOrdersShipments()
	if err != nil {
		log.Panic(err)
	}
	// fmt.Println(resp.Orders)
	fmt.Println(len(resp.Orders))
	maker.GetItems(*resp)
	// b, _ := json.MarshalIndent(maker.GetItems(*resp), "", "    ")
	// fmt.Println(string(b))
}
