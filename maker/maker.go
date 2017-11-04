package maker

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/WedgeNix/awsapi/file"
	"github.com/WedgeNix/awsapi/types"
)

var (
	data *file.BananasMon
)

func GetItems(pay Payload) *file.BananasMon {
	data = &file.BananasMon{
		SKUs: types.SKUs{},
	}
	for _, ord := range pay.Orders {
		for _, itm := range ord.Items {
			monSKU := data.SKUs[itm.SKU]
			t, err := time.Parse("2006-01-02T15:04:05", ord.OrderDate)
			if err != nil {
				log.Panic(err)
			}
			if !monSKU.LastUTC.IsZero() {
				diff := t.Sub(monSKU.LastUTC)
				days := int(diff.Hours()/24 + 0.5)
				monSKU.Days += days
			}

			if err != nil {
				log.Panic(err)
			}

			monSKU.Sold += itm.Quantity
			monSKU.LastUTC = t
			monSKU.Pending = time.Time{}

			addAvgWait(ord)

			data.SKUs[itm.SKU] = monSKU
		}
	}
	return data
}

func addAvgWait(ord Order) {
	if len(ord.AdvancedOptions.CustomField1) == 0 {
		return
	}
	cusF1 := strings.Split(ord.AdvancedOptions.CustomField1, ",")
	poDateStr := strings.Split(cusF1[0], "-")[1]

	poDate, err := time.Parse("20060102", poDateStr)
	if err != nil {
		log.Panic(err)
	}
	shipDate, err := time.Parse("2006-01-02", ord.ShipDate)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(poDate.Sub(shipDate))

}
