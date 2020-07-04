package main

import (
	"log"
	"net/http"
	"sync"
)

//
//import (
//	"github.com/gomodule/redigo/redis"
//	"log"
//	"strconv"
//)
//
//func main() {
//	conn, err := redis.Dial("tcp", "localhost:6379", redis.DialPassword("hello"))
//	if err!=nil {
//		log.Fatalln(err)
//	}
//	defer conn.Close()
//
//	//_, err = conn.Do("del", "album:2")
//	//if err!=nil {
//	//	log.Fatalln(err)
//	//}
//	//log.Println("album:2 removed!")
//
//	_, err = conn.Do("hmset", "album:2", "title", "yuan", "artist", "Jimi", "price", 4.95, "likes", 8)
//	if err!=nil {
//		log.Fatalln(err)
//	}
//	log.Println("album:2 added!")
//
//	reply, err := redis.StringMap(conn.Do("hgetall", "album:2"))
//	if err!=nil {
//		log.Fatalln(err)
//	}
//
//	album, err := populateAlbum(reply)
//	if err!=nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v", album)
//
//	////_, err := conn.Do("")
//	//title, err := redis.String(conn.Do("hget", "album:2", "title"))
//	//if err!=nil {
//	//	log.Fatalln(err)
//	//}
//	//log.Printf("title is %s", title)
//	//
//	//artist, err := redis.String(conn.Do("hget", "album:2", "artist"))
//	//if err!=nil {
//	//	log.Fatalln(err)
//	//}
//	//log.Printf("artist is %s", artist)
//	//
//	//price, err := redis.Float64(conn.Do("hget", "album:2", "price"))
//	//if err!=nil {
//	//	log.Fatalln(err)
//	//}
//	//log.Printf("price is %.2f", price)
//	//
//	//likes, err := redis.Int(conn.Do("hget", "album:2", "likes"))
//	//if err!=nil {
//	//	log.Fatalln(err)
//	//}
//	//log.Printf("likes is %d", likes)
//}
//
//func populateAlbum(reply map[string]string) (*Album, error) {
//	var err error
//	album := new(Album)
//	album.Title = reply["title"]
//	album.Artist = reply["artist"]
//	album.Price, err = strconv.ParseFloat(reply["price"], 64)
//	if err!=nil {
//		return nil, err
//	}
//	album.Likes, err = strconv.Atoi(reply["likes"])
//	if err!=nil {
//		return nil, err
//	}
//	return album, nil
//}
func init() {
	//// 初始化album
	//err := InitAlbumData()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	// 初始化productStore
	var wg sync.WaitGroup
	// 给连接池多搞点conn
	for i := 0; i < 20000; i++ {
		wg.Add(1)
		pushConnToPool(&wg)
	}
	wg.Wait()
	log.Printf("conn 预热完成!\n")

	err := InitStore()
	if err != nil {
		log.Println(err)
		return
	}
}

func pushConnToPool(wg *sync.WaitGroup) {
	conn := pool.Get()
	defer conn.Close()
	wg.Done()
}

func main() {
	mux := http.NewServeMux()
	//mux.HandleFunc("/album", showAlbum)
	//mux.HandleFunc("/like", addLike)
	//mux.HandleFunc("/createAlbum", createAlbum)
	mux.HandleFunc("/buy", buy)
	// "/cancelBuy" 这个接口只能由后台来调用
	mux.HandleFunc("/cancelBuy", cancelBuy)
	log.Println("Listening on 127.0.0.1:4000")
	err := http.ListenAndServe("127.0.0.1:4000", mux)
	if err != nil {
		log.Println(err)
		return
	}
	//u := new(User)
	//u.UserId = "leo2n"
	//u.orderGenerator("123456", 1)
}
