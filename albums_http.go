package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

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
	http.Error(w, err.Error(), statusCode)
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
