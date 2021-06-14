package easyjsonprocess

type BuyReq struct {
	Username    string `json:"username"`
	ProductId   string `json:"productId"`
	PurchaseNum int    `json:"purchaseNum"`
}

type CancelBuyReq struct {
	Username string `json:"username"`
	OrderNum string `json:"orderNum"`
}

//func DecodeBuyReq(buyReq io.ReadCloser) (*BuyReq, error) {
//	b := new(BuyReq)
//	err := json.NewDecoder(buyReq).Decode(b)
//	if err != nil {
//		return new(BuyReq), err
//	}
//	return b, nil
//}
//
//func DecodeCancelBuyReq(cancelBuyReq io.ReadCloser) (*CancelBuyReq, error) {
//	b := new(CancelBuyReq)
//	err := json.NewDecoder(cancelBuyReq).Decode(b)
//	if err != nil {
//		return new(CancelBuyReq), err
//	}
//	return b, nil
//}
