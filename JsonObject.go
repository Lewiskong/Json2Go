type JsonObject struct {
	Result	bool	`json:"result"`
	Orders	[]OrdersItem	`json:"orders"`
}
type OrdersItem struct {
	Amount	float64	`json:"amount"`
	Avg_price	float64	`json:"avg_price"`
	Create_date	float64	`json:"create_date"`
	Order_id	float64	`json:"order_id"`
	Price	float64	`json:"price"`
	Type	string	`json:"type"`
	Deal_amount	float64	`json:"deal_amount"`
	Orders_id	float64	`json:"orders_id"`
	Status	float64	`json:"status"`
	Symbol	string	`json:"symbol"`
}
