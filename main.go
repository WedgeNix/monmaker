package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/WedgeNix/awsapi"
	"github.com/WedgeNix/monmaker/maker"
	"github.com/WedgeNix/warn"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	new, err := regexp.Compile(`\r\n`)
	if err != nil {
		log.Panic(err)
	}

	// start user inputs.
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter full brand name: ")
	brandStr, err := reader.ReadString('\n')
	if err != nil {
		log.Panic(err)
	}
	brand := new.ReplaceAllString(brandStr, "")
	fmt.Print("Start date for sales data (YYYY-MM-DD): ")
	date, err := reader.ReadString('\n')
	if err != nil {
		log.Panic(err)
	}
	from := dateValidate(date, reader)

	fmt.Println("     Start ShipStation call...")
	// Start ShipStation Call.
	now := time.Now()
	resp, err := maker.Ship(from, now, brand).GetOrdersShipments()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("     Makeing Mon file data...")
	// Make mon file data.
	monFileDate := maker.GetItems(*resp)

	fmt.Println("     Makeing Mon file...")
	// Make Mon File.
	f, err := os.Create(brand + ".json")
	if err != nil {
		log.Panic(err)
	}
	json.NewEncoder(f).Encode(monFileDate)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("     Check file for AWS...")
	// Check if file is good to send to AWS.
	path, err := filepath.Abs(".")
	if err != nil {
		log.Panic(err)
	}
	err = exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path+"/"+f.Name()).Run()
	if err != nil {
		log.Panic(err)
	}
	warn.Do("Check file for erros. Uplad to AWS:")
	fmt.Println("     Sending file to AWS...")
	// Send to aws
	aws, err := awsapi.New()
	if err != nil {
		log.Panic(err)
	}
	err = aws.SaveFile("hit-the-bananas/mon/", f)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("     Monitor file created!")

}

func dateValidate(readStr string, reader *bufio.Reader) time.Time {
	new, err := regexp.Compile(`\r\n`)
	if err != nil {
		log.Panic(err)
	}
	for {
		dateStr := new.ReplaceAllString(readStr, "")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			fmt.Print("Incorrect date format, try again: ")
			readStr, err = reader.ReadString('\n')
			if err != nil {
				log.Panic(err)
			}
			continue
		}
		return date
	}
}
