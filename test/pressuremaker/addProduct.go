package pressuremaker

import (
	"encoding/json"
	"go-seckill/internal/db/shop_orm"
	"log"
	"net/url"
	"strconv"

	"github.com/valyala/fasthttp"
)

type Product shop_orm.Good

func AddProducts() error {
	productList := make([]*Product, 0, 8)
	for i := 0; i < 6; i++ {
		productList = append(productList, &Product{
			ProductCategory: "Mobile Phone",
			ProductName:     "Xiaomi-" + strconv.Itoa(i+1),
			Inventory:       200,
			Price:           1000 * (i + 1) * 1e2,
		})
	}
	var err error = nil
	for _, v := range productList {
		// log.Println(v.ProductCategory, v.ProductName, v.Inventory, v.Price)
		err = v.CreateGood()
		if err != nil {
			break
		}
	}
	return err
}

func (p *Product) CreateGood() error {
	client := FastHttpClient

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	url := &url.URL{
		Scheme: "http",
		Host:   "go-seckill:4000",
		Path:   "/api/v0/admin/goodCreate",
	}
	reqBody, err := json.Marshal(p)
	if err != nil {
		return err
	}
	req.Header.SetMethod(fasthttp.MethodPost)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJncm91cCI6ImFkbWluIiwidXNlcm5hbWUiOiJhZG1pbiIsImV4cCI6MTYzMTg3Mjk1Mn0.GK1IA3T4bVvGFwmUeqrRhaxTSI7TK5rBrFvnvEGwUVY"
	req.Header.Set("Authorization", token)
	req.SetBody(reqBody)
	req.SetRequestURI(url.String())
	log.Println("request url is: ", url.String())

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)
	err = client.Do(req, res)
	log.Println(string(res.Body()))
	if err != nil {
		return err
	}
	return nil
}
