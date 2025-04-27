package model

import "time"

type Customers struct {
	CustomerId string `gorm:"column:customer_id"`
	Name       string `gorm:"column:name"`
	Email      string `gorm:"column:email"`
	Address    string `gorm:"column:address"`
	// CreatedBy   string    `gorm:"column:created_by"`
	// CreatedDate time.Time `gorm:"column:created_date;type:datetime"`
	// UpdatedBy   string    `gorm:"column:updated_by"`
	// UpdatedDate time.Time `gorm:"column:updated_date;type:datetime"`
}

type Products struct {
	ProductId   string `gorm:"column:product_id"`
	ProductName string `gorm:"column:product_name"`
	Category    string `gorm:"column:category"`
	// CreatedBy   string    `gorm:"column:created_by"`
	// CreatedDate time.Time `gorm:"column:created_date;type:datetime"`
	// UpdatedBy   string    `gorm:"column:updated_by"`
	// UpdatedDate time.Time `gorm:"column:updated_date;type:datetime"`
}

type Orders struct {
	OrderId       string    `gorm:"column:order_id"`
	CustomerId    string    `gorm:"column:customer_id"`
	PaymentMethod string    `gorm:"column:payment_method"`
	DateofSale    time.Time `gorm:"column:date_of_sale;type:datetime"`
	Region        string    `gorm:"column:region"`
	// CreatedBy     string    `gorm:"column:created_by"`
	// CreatedDate   time.Time `gorm:"column:created_date;type:datetime"`
	// UpdatedBy     string    `gorm:"column:updated_by"`
	// UpdatedDate   time.Time `gorm:"column:updated_date;type:datetime"`
}

type OrderItems struct {
	OrderItemId  uint64  `gorm:"column:order_item_id;primaryKey;autoIncrement"`
	OrderId      string  `gorm:"column:order_id"`
	ProductId    string  `gorm:"column:product_id"`
	QuantitySold int     `gorm:"column:quantity_sold"`
	UnitPrice    float64 `gorm:"column:unit_price"`
	Discount     float64 `gorm:"column:discount"`
	ShippingCost float64 `gorm:"column:shipping_cost"`
	// CreatedBy    string    `gorm:"column:created_by"`
	// CreatedDate  time.Time `gorm:"column:created_date;type:datetime"`
	// UpdatedBy    string    `gorm:"column:updated_by"`
	// UpdatedDate  time.Time `gorm:"column:updated_date;type:datetime"`
}

type SalesProcessData struct {
	Customer   Customers
	Product    Products
	Orders     Orders
	OrderItems OrderItems
}

type SalesData struct {
	OrderID         string `csv:"Order ID"`
	ProductID       string `csv:"Product ID"`
	CustomerID      string `csv:"Customer ID"`
	ProductName     string `csv:"Product Name"`
	Category        string `csv:"Category"`
	Region          string `csv:"Region"`
	DateOfSale      string `csv:"Date of Sale"`
	QuantitySold    string `csv:"Quantity Sold"`
	UnitPrice       string `csv:"Unit Price"`
	Discount        string `csv:"Discount"`
	ShippingCost    string `csv:"Shipping Cost"`
	PaymentMethod   string `csv:"Payment Method"`
	CustomerName    string `csv:"Customer Name"`
	CustomerEmail   string `csv:"Customer Email"`
	CustomerAddress string `csv:"Customer Address"`
}
type RefreshDataResp struct {
	Status     string `json:"status"`
	StatusCode string `json:"statuscode"`
	Message    string `json:"message"`
}
type CommonRevenueResp struct {
	Status      string  `json:"status"`
	StatusCode  string  `json:"statuscode"`
	Message     string  `json:"message"`
	RevenueType string  `json:"revenuetype"`
	Revenue     float64 `json:"revenue"`
}

type SalesRefreshLog struct {
	Id          int       `gorm:"column:id;primarykey"`
	Status      string    `gorm:"column:status"`
	Type        string    `gorm:"column:type"`
	CreatedDate time.Time `gorm:"column:created_date"`
	UpdatedDate time.Time `gorm:"column:updated_date"`
}
