package main

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
)

var ErrNoAlbum = errors.New("no album found")

// Define a custom struct to hold Album data.
type Album struct {
	AlbumName string  `json:"albumName"`
	Title     string  `json:"title"`
	Artist    string  `json:"artist"`
	Price     float64 `json:"price"`
	Likes     int     `json:"likes"`
}

// 初始化演示demo所需要的数据
func InitAlbumData() error {
	conn := pool.Get()
	defer conn.Close()
	// First of all, flushdb
	conn.Send("flushdb")

	// 然后, 创建album相关数据
	err := conn.Send("hmset", "album:1", "title", "Electric Ladyland", "artist", "Jimi Hungry", "price", "4.95", "likes", "8")
	if err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	err = conn.Send("hmset", "album:2", "title", "Back In Black", "artist", "AC/DC", "price", "5.95", "likes", "3")
	if err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	err = conn.Send("hmset", "album:3", "title", "Rumours", "artist", "Fleetwood Mac", "price", "7.95", "likes", "12")
	if err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	err = conn.Send("hmset", "album:4", "title", "Nevermind", "artist", "Nirvana", "price", "5.95", "likes", "8")
	if err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	// After, create orderSet 相关数据
	err = conn.Send("zadd", "likes", "8", "1", "3", "2", "12", "3", "8", "4")
	if err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	return nil
}

// 解析argList并且写入redis数据库中
func createAlbumInRedis(argList []string) error {
	conn := pool.Get()
	defer conn.Close()

	// 转化 []string -> []interface{}
	argList1 := make([]interface{}, len(argList))
	for i, v := range argList {
		argList1[i] = v
	}
	err := conn.Send("hmset", argList1...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 根据传入的id来查找album的相关信息
func FindAlbum(id string) (*Album, error) {
	conn := pool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("hgetall", "album:"+id))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, ErrNoAlbum // 语句执行过程正确, 但是... 没有值返回的时候, 就会触发这个错误!
	}

	albumPointer := new(Album)
	err = redis.ScanStruct(values, albumPointer)
	if err != nil {
		return nil, err
	}
	return albumPointer, nil
}

// 增加hash: album:id 和 sorted_set likes
func IncrementLikes(id string) error {
	conn := pool.Get()
	defer conn.Close()

	// 首先, 检查特定的album:id 是否存在于redis中
	exists, err := redis.Int(conn.Do("exists", "album:"+id))
	if err != nil {
		return err
	} else if exists == 0 {
		return ErrNoAlbum
	}

	// 准备开始执行事务
	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	// hincrby 增加 hash table 特定key的values的数值
	err = conn.Send("hincrby", "album:"+id, "likes", 1)
	if err != nil {
		return err
	}
	// zincrby 增加set特定value的score
	err = conn.Send("zincrby", "likes", 1, id)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}
