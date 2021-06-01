/*
使用了传说中的easyjson, 貌似速度并没有好很多
*/
package jsonStruct

type ReqBuy struct {
	UserId      string `json:"userId"`
	ProductId   string `json:"productId"`
	PurchaseNum int    `json:"purchaseNum"`
}
