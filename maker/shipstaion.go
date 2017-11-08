package maker

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/WedgeNix/util"
)

// ShipControl conrtols api call to shipstaion
type ShipControl struct {
	from    time.Time
	to      time.Time
	keyword string
	*util.HTTPLogin
}

// Ship starts Shipstaion api
func Ship(from time.Time, to time.Time, keyword string) *ShipControl {
	return &ShipControl{
		from:    from,
		to:      to,
		keyword: keyword,
		HTTPLogin: &util.HTTPLogin{
			User: os.Getenv("SHIP_API_KEY"),
			Pass: os.Getenv("SHIP_API_SECRET"),
		},
	}

}

// GetOrdersShipments grabs an HTTP response of orders, filtering in those shipment.
func (c *ShipControl) GetOrdersShipments() (*Payload, error) {
	pg := 1

	pay := &Payload{}
	reqs, secs, err := c.getPage(pg, pay)
	if err != nil {
		return pay, util.Err(err)
	}
	for pay.Page < pay.Pages {
		pg++
		ords := pay.Orders

		pay = &Payload{}
		reqs, secs, err = c.getPage(pg, pay)
		if err != nil {
			return pay, util.Err(err)
		}
		if reqs < 1 {
			time.Sleep(time.Duration(secs) * time.Second)
		}

		pay.Orders = append(ords, pay.Orders...)
	}
	return pay, nil
}

func (c *ShipControl) getPage(page int, pay *Payload) (int, int, error) {
	shipURL := "https://ssapi.shipstation.com/"

	query := url.Values(map[string][]string{})
	query.Set(`page`, strconv.Itoa(page))
	query.Set(`createDateStart`, c.from.Format("2006-01-02 15:04:05"))
	query.Set(`createDateEnd`, c.to.Format("2006-01-02 15:04:05"))
	query.Set(`itemKeyword`, c.keyword)
	query.Set(`orderStatus`, `shipped`)
	query.Set(`sortBy`, `OrderDate`)
	query.Set(`sortDir`, `ASC`)
	query.Set(`pageSize`, `500`)

	resp, err := c.HTTPLogin.Get(shipURL + `orders?` + query.Encode())
	if err != nil {
		return 0, 0, util.Err(err)
	}

	err = json.NewDecoder(resp.Body).Decode(pay)
	if err != nil {
		return 0, 0, util.Err(err)
	}
	defer resp.Body.Close()

	remaining := resp.Header.Get("X-Rate-Limit-Remaining")
	reqs, err := strconv.Atoi(remaining)
	if err != nil {
		return 0, 0, util.Err(err)
	}
	reset := resp.Header.Get("X-Rate-Limit-Reset")
	secs, err := strconv.Atoi(reset)
	if err != nil {
		return reqs, 0, util.Err(err)
	}

	return reqs, secs, nil
}

// Payload is the first level of a ShipStation HTTP response body.
type Payload struct {
	Orders []Order
	Total  int
	Page   int
	Pages  int
}

// Order is the second level of a ShipStation HTTP response body.
type Order struct {
	OrderID                  int
	OrderNumber              string
	OrderKey                 string
	OrderDate                string
	CreateDate               string
	ModifyDate               string
	PaymentDate              string
	ShipByDate               string
	OrderStatus              string
	CustomerID               int
	CustomerUsername         string
	CustomerEmail            string
	BillTo                   interface{}
	ShipTo                   interface{}
	Items                    []Item
	OrderTotal               float32
	AmountPaid               float32
	TaxAmount                float32
	ShippingAmount           float32
	CustomerNotes            string
	InternalNotes            string
	Gift                     bool
	GiftMessage              string
	PaymentMethod            string
	RequestedShippingService string
	CarrierCode              string
	ServiceCode              string
	PackageCode              string
	Confirmation             string
	ShipDate                 string
	HoldUntilDate            string
	Weight                   interface{}
	Dimensions               interface{}
	InsuranceOptions         interface{}
	InternationalOptions     interface{}
	AdvancedOptions          AdvancedOptions
	TagIDs                   []int
	UserID                   string
	ExternallyFulfilled      bool
	ExternallyFulfilledBy    string
}

// Item is the third level of a ShipStation HTTP response body.
type Item struct {
	OrderItemID       int
	LineItemKey       string
	SKU               string
	Name              string
	ImageURL          string
	Weight            interface{}
	Quantity          int
	UnitPrice         float32
	TaxAmount         float32
	ShippingAmount    float32
	WarehouseLocation string
	Options           interface{}
	ProductID         int
	FulfillmentSKU    string
	Adjustment        bool
	UPC               string
	CreateDate        string
	ModifyDate        string
}

// AdvancedOptionsUpdate advanced options
type AdvancedOptionsUpdate struct {
	StoreID      int
	CustomField1 string
	CustomField2 string
	CustomField3 string
}

// AdvancedOptions holds the "needed" custom fields for post-email tagging.
type AdvancedOptions struct {
	WarehouseID       int
	NonMachinable     bool
	SaturdayDelivery  bool
	ContainsAlcohol   bool
	MergedOrSplit     bool
	MergedIDs         interface{}
	ParentID          interface{}
	StoreID           int
	CustomField1      string
	CustomField2      string
	CustomField3      string
	Source            string
	BillToParty       interface{}
	BillToAccount     interface{}
	BillToPostalCode  interface{}
	BillToCountryCode interface{}
}
