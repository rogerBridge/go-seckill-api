package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
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
	err := InitAlbumData()
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/album", showAlbum)
	mux.HandleFunc("/like", addLike)
	mux.HandleFunc("/createAlbum", createAlbum)
	log.Println("Listening on port 4000")
	err := http.ListenAndServe("127.0.0.1:4000", mux)
	if err != nil {
		log.Println(err)
		return
	}
}

// 通过http post方法对redis进行操作
func addLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	id := r.PostFormValue("id") // 只会读取postForm里面的数值, 不会读取url里面的键值对
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	// 验证请求传入的id是否合法, 只能是int数据
	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	// 通过请求传入的id, 增加对应的数值
	err := IncrementLikes(id)
	if err == ErrNoAlbum {
		http.Error(w, http.StatusText(404), 404)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// 303 跳转到其他的可用的网页
	http.Redirect(w, r, "/album?id="+id, 303)
}

// 使用http post方法在redis里面写入album相关的数据
func createAlbum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), 405)
		return // 在外部调用这看来, 这个步骤已经可以停止了
	}
	// 需要验证数据的格式是否正确
	// 一般常用两种传参方式: application/x-www-form-urlencoded and application/json
	// x-www-form-urlencoded
	//err := r.ParseForm()
	//if err!=nil {
	//	errorHandle(w, err, 500)
	//	return
	//}
	//title := r.PostForm["title"][0]
	//if title == "" {
	//	errorHandle(w, errors.New("title参数有误!"), 400)
	//	return
	//}
	//artist := r.PostForm["artist"][0]
	//// Q: 正则表达式会更好??
	//if artist == "" {
	//	errorHandle(w, errors.New("artist参数有误!"), 400)
	//	return
	//}
	//price := r.PostForm["price"][0]
	//_, err = strconv.ParseFloat(r.PostForm["price"][0], 64)
	//if err!=nil {
	//	errorHandle(w, err, 400)
	//	return
	//}
	//likes := r.PostForm["likes"][0]
	//_, err = strconv.Atoi(r.PostForm["likes"][0])
	//if err!=nil {
	//	errorHandle(w, err, 400)
	//	return
	//}
	//albumName := r.PostForm["albumName"][0]
	//argList := []string{
	//	albumName, "title", title, "artist", artist, "price", price, "likes", likes,
	//}
	//// 构造一个Album arg list
	//err = createAlbumInRedis(argList)
	//if err!=nil {
	//	errorHandle(w, err, 500)
	//	return
	//}
	// application/json 类型的
	ap := new(Album) // albumPointer
	err := json.NewDecoder(r.Body).Decode(ap)
	if err != nil {
		errorHandle(w, err, 500)
		return
	}
	// 对上传的数据的判断
	albumName := ap.AlbumName
	if albumName == "" {
		errorHandle(w, errors.New("albumName参数有误!"), 400)
		return
	}
	title := ap.Title
	if title == "" {
		errorHandle(w, errors.New("title参数有误!"), 400)
		return
	}
	artist := ap.Artist
	// Q: 正则表达式会更好??
	if artist == "" {
		errorHandle(w, errors.New("artist参数有误!"), 400)
		return
	}
	price := strconv.FormatFloat(ap.Price, 'f', -1, 64)
	likes := strconv.Itoa(ap.Likes)
	argList := []string{
		albumName, "title", ap.Title, "artist", ap.Artist, "price", price, "likes", likes,
	}
	err = createAlbumInRedis(argList)
	if err != nil {
		errorHandle(w, err, 500)
		return
	}
}

// 错误信息的处理
func errorHandle(w http.ResponseWriter, err error, statusCode int) {
	log.Println(err)
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func showAlbum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), 405)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), 400)
		return
	}

	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, errors.New("argument type is bad").Error(), 400)
		return
	}

	bk, err := FindAlbum(id)
	if err == ErrNoAlbum {
		log.Println(ErrNoAlbum.Error())
		http.Error(w, http.StatusText(404), 404)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	log.Printf("%v", bk)
	s, err := json.Marshal(bk)
	if err != nil {
		log.Fatalln(err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(s))
	if err != nil {
		log.Println(err)
		return
	}
}
