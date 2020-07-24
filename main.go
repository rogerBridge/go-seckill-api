package main

import (
	"log"
	"net/http"
)

func init() {
	//var wg sync.WaitGroup
	//// 给连接池多搞点conn
	//for i := 0; i < 20000; i++ {
	//	wg.Add(1)
	//	pushConnToPool(&wg)
	//}
	//wg.Wait()
	//log.Printf("conn 预热完成!\n")

	err := InitStore()
	if err != nil {
		log.Println(err)
		return
	}
}

//func pushConnToPool(wg *sync.WaitGroup) {
//	conn := pool.Get()
//	defer conn.Close()
//	wg.Done()
//}

func main() {
	mux := http.NewServeMux()
	//mux.HandleFunc("/album", showAlbum)
	//mux.HandleFunc("/like", addLike)
	//mux.HandleFunc("/createAlbum", createAlbum)
	mux.HandleFunc("/buy", buy)
	// "/cancelBuy" 这个接口只能由后台来调用
	mux.HandleFunc("/cancelBuy", cancelBuy)
	log.Println("Listening on 0.0.0.0:4000")
	err := http.ListenAndServe("0.0.0.0:4000", mux)
	if err != nil {
		log.Println(err)
		return
	}
	//u := new(User)
	//u.UserId = "leo2n"
	//u.orderGenerator("123456", 1)
}
