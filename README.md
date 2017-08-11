# Json2Go
A command line tool convert Json 2 Go struct  

## **Installation**
```
go get github.com/lewiskong/Json2Go
```
Assuming that `GOPATH/bin` is in your `PATH` , you can now use `Json2Go` directly.

## **Usage**
```
usage:
    Json2Go [file]
    Json2Go [file] [target]
    Json2Go:
        -src	:	source file
        -dest	:	dest file
        -f		:	overwrite dest file. Default true
        -name	:	The name of the go struct
```

## **Example**
```
    {
    "result": true,
    "orders": [
        {
            "amount": 0.1,
            "avg_price": 0,
            "create_date": 1418008467000,
            "deal_amount": 0,
            "order_id": 10000591,
            "orders_id": 10000591,
            "price": 500,
            "status": 0,
            "symbol": "btc_cny",
            "type": "sell"
        },
        {
            "amount": 0.2,
            "avg_price": 0,
            "create_date": 1417417957000,
            "deal_amount": 0,
            "order_id": 10000724,
            "orders_id": 10000724,
            "price": 0.1,
            "status": 0,
            "symbol": "btc_cny",
            "type": "buy",
            "fuck" : "what"
        }
    ]
}
```

**After convert =>**

```
type JsonObject struct {
	Result	bool	`json:"result"`
	Orders	[]OrdersItem	`json:"orders"`
}
type OrdersItem struct {
	Fuck	string	`json:"fuck"`
	Create_date	float64	`json:"create_date"`
	Symbol	string	`json:"symbol"`
	Type	string	`json:"type"`
	Amount	float64	`json:"amount"`
	Deal_amount	float64	`json:"deal_amount"`
	Order_id	float64	`json:"order_id"`
	Avg_price	float64	`json:"avg_price"`
	Orders_id	float64	`json:"orders_id"`
	Price	float64	`json:"price"`
	Status	float64	`json:"status"`
}

```

## **Develop**
```
git clone https://github.com/ChimeraCoder/gojson.git

go build
```

## **License**
Lewiskong/Json2Go is licensed under the MIT License.
A short and simple permissive license with conditions only requiring preservation of copyright and license notices. Licensed works, modifications, and larger works may be distributed under different terms and without source code.

