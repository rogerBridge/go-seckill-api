package jsonStruct

type ReqBuy struct {
	UserId      string `json:"userId"`
	ProductId   string `json:"productId"`
	PurchaseNum int    `json:"purchaseNum"`
}
