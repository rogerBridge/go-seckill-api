package main

import (
	"log"

	"github.com/valyala/fasthttp"
)

func main() {
	fs := fasthttp.FSHandler("./static", 0)
	log.Println("fasthttp file server running on port 3000 :)")
	log.Println(fasthttp.ListenAndServe(":3000", fs))

	// fs := http.FileServer(http.Dir("./static"))
	// http.Handle("/", fs)

	// log.Println("Listening on :3000...")
	// log.Println(http.ListenAndServe(":3000", nil))
	// err := http.ListenAndServe(":3000", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
