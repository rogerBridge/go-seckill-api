package main
//
//import (
//	"log"
//	"sync"
//	"time"
//)
//
//func play(timeStatistics chan float64, wg *sync.WaitGroup) {
//	t0 := time.Now()
//	time.Sleep(time.Second)
//	t1 := time.Since(t0).Seconds()
//	timeStatistics <- t1
//	wg.Done()
//}
//
//func main() {
//	var wg sync.WaitGroup
//	t := make(chan float64, 10)
//	for i:=0; i<10; i++ {
//		wg.Add(1)
//		go play(t, &wg)
//	}
//	wg.Wait()
//	close(t)
//	for v := range t {
//		log.Println(v)
//	}
//}