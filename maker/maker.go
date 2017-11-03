package maker

import (
	"log"
	"time"

	"github.com/WedgeNix/awsapi/file"
	"github.com/WedgeNix/awsapi/types"
)

func GetItems(pay Payload) file.BananasMon {
	data := file.BananasMon{
		SKUs: types.SKUs{},
	}
	for _, ord := range pay.Orders {
		for _, itm := range ord.Items {
			monSKU := data.SKUs[itm.SKU]
			t, err := time.Parse("2006-01-02T15:04:05", ord.OrderDate)
			if !monSKU.LastUTC.IsZero() {
				diff := monSKU.LastUTC.Sub(t)
				days := int(diff.Hours()/24 + 0.5)
				monSKU.Days += days
			}
			if err != nil {
				log.Panic(err)
			}

			monSKU.Sold += itm.Quantity
			monSKU.LastUTC = t
			monSKU.Pending = time.Time{}

			data.SKUs[itm.SKU] = monSKU
		}
	}
	return data
}
