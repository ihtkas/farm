package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/google/gops/agent"
	"github.com/ihtkas/farm/seller"
)

func main() {
	runtime.SetBlockProfileRate(1)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	time.Sleep(time.Second * 5)
	m := seller.Manager{}
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
	go func() {
		err := m.Start()
		if err != nil {
			glog.Errorln(err)
			return
		}
		glog.Errorln("Done....")
	}()

	go load()
	l := &sync.Mutex{}
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	go getLock(l)
	select {}
}

// go:noinline
func getLock(l1 *sync.Mutex) {
	fmt.Println("before Acquired lock")
	// go:noinline
	l1.Lock()
	fmt.Println("Acquired lock")
	x := 100
	for i := 1; i <= 100; i++ {
		x = x * i / i
	}
	fmt.Println(x)
}

func load() {
	for {
		src := rand.NewSource(time.Now().UnixNano())
		r := rand.New(src)
		src = rand.NewSource(time.Now().UnixNano())
		r2 := rand.New(src)
		x := r.Intn(1000)
		y := r2.Intn(1000)
		arr := append(r.Perm(x), r.Perm(y)...)
		ch := make(chan bool)
		go mergeSort(arr, ch)
		<-ch
	}
}

func mergeSort(arr []int, done chan bool) {
	mergeSortRec(arr, 0, len(arr)-1)
	close(done)
}

func mergeSortRec(arr []int, s, e int) {
	if s < e {
		m := (s + e) / 2
		mergeSortRec(arr, s, m)
		mergeSortRec(arr, m+1, e)
		merge(arr, s, e)
	}
}

func merge(arr []int, s, e int) {
	clone := make([]int, e-s+1)
	ind := 0
	i := s
	m := (s + e) / 2
	j := m + 1
	// fmt.Println(s, e, m)
	for i <= m && j <= e {
		if arr[i] < arr[j] {
			clone[ind] = arr[i]
			i++
		} else {
			clone[ind] = arr[j]
			j++
		}
		ind++
	}
	for i <= m {
		clone[ind] = arr[i]
		ind++
		i++
	}

	for j <= e {
		clone[ind] = arr[j]
		ind++
		j++
	}
	for _, e := range clone {
		arr[s] = e
		s++
	}

}
