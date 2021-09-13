package main

import (
	"log"
	"net/http"
)

func main() {
	// fs := fasthttp.FSHandler("./static", 0)
	// log.Println("fasthttp file server running on port 3000 :)")
	// log.Println(fasthttp.ListenAndServe(":3000", fs))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	// http.Handle("/static/", http.StripPrefix("/static", fs))

	log.Println("Listening on :3000...")
	log.Fatal((http.ListenAndServe(":3000", nil)))
}
