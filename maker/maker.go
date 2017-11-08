package maker

import (
	"log"
	"strings"
	"time"

	"github.com/OuttaLineNomad/skuvault"
	"github.com/OuttaLineNomad/skuvault/products"

	"github.com/WedgeNix/awsapi/file"
	"github.com/WedgeNix/awsapi/types"
)

var (
	data *file.BananasMon
)

type skuDate map[string]time.Time

// GetItems makes the mon file data.
func GetItems(pay Payload) *file.BananasMon {
	data = &file.BananasMon{
		SKUs: types.SKUs{},
	}
	prodDate := getSKUAge(pay)
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
				daysOld := max(days, 1)
				monSKU.Days += daysOld
				monSKU.ProbationPeriod = min(8100/daysOld, 90)
				monSKU.LastUTC = t
			} else {
				diff := t.Sub(prodDate[itm.SKU])
				days := int(diff.Hours()/24 + 0.5)
				daysOld := max(days, 1)
				monSKU.ProbationPeriod = min(8100/daysOld, 90)
				monSKU.LastUTC = t
				monSKU.Days = 1
			}

			monSKU.Sold += itm.Quantity
			monSKU.Pending = time.Time{}

			addAvgWait(ord)

			data.SKUs[itm.SKU] = monSKU

		}
	}
	cleanData(data, prodDate)
	return data
}

func cleanData(data *file.BananasMon, prodDate skuDate) {

	for sku, monSKU := range data.SKUs {
		t := time.Now()
		diff := t.Sub(monSKU.LastUTC)
		days := int(diff.Hours()/24 + 0.5)
		daysOld := max(days, 1)
		expired := daysOld > monSKU.ProbationPeriod
		if expired {
			delete(data.SKUs, sku)
		}
	}
}

func getSKUAge(pay Payload) skuDate {
	skuD := map[string]time.Time{}
	SKUs := []string{}
	for _, ord := range pay.Orders {
		for _, itm := range ord.Items {
			SKUs = append(SKUs, itm.SKU)
		}
	}
	sv := skuvault.NewEnvCredSession()
	prodPay := &products.GetProducts{
		PageSize:    10000,
		ProductSKUs: SKUs,
	}

	resp := sv.Products.GetProducts(prodPay)
	for _, prod := range resp.Products {
		skuD[prod.Sku] = prod.CreatedDateUtc
	}
	return skuD
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
	days := shipDate.Sub(poDate).Hours()/24 + .5
	if days == 0 {
		return
	}
	if data.AvgWait == 0 {
		data.AvgWait = 5
	}
	data.OrdSKUCnt++
	rez := 1.0 / data.OrdSKUCnt
	data.AvgWait = data.AvgWait*(1-rez) + days*rez
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
