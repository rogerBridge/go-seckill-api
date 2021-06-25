/*
使用了传说中的easyjson, 貌似速度并没有好很多
*/
package jsonStruct

type ReqBuy struct {
	ProductId   int `json:"productID"`
	PurchaseNum int `json:"purchaseNum"`
}
