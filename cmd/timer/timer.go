//+build timer

package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	timer(context.Background(), time.Hour, "alert.mp3")
}

func timer(ctx context.Context, interval time.Duration, alertmp3 string) {
	// tick := time.NewTicker(interval)
	// ping := time.NewTicker(time.Minute)
	// // glog.Errorln("Next alert at ", time.Now().Add(interval).Format("3:04:05 PM"))
	// // time.Sleep(time.Second * 10)
	// // func() {
	// // 	f, err := os.Open(alertmp3)
	// // 	if err != nil {
	// // 		log.Fatal(err)
	// // 	}
	// // 	defer f.Close()
	// // 	streamer, format, err := mp3.Decode(f)
	// // 	if err != nil {
	// // 		log.Fatal(err)
	// // 	}
	// // 	defer streamer.Close()

	// // 	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	// // 	speaker.Play(streamer)
	// // 	glog.Errorln("Next alert at ", time.Now().Add(interval).Format("3:04:05 PM"))
	// // }()
	// for {

	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	case <-ping.C:
	// 		fmt.Print(".")
	// 	case <-tick.C:
	// 		func() {
	// 			f, err := os.Open(alertmp3)
	// 			if err != nil {
	// 				log.Fatal(err)
	// 			}
	// 			defer f.Close()
	// 			streamer, format, err := mp3.Decode(f)
	// 			if err != nil {
	// 				log.Fatal(err)
	// 			}
	// 			defer streamer.Close()

	// 			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	// 			speaker.Play(streamer)
	// 			done := make(chan bool)
	// 			speaker.Play(beep.Seq(streamer, beep.Callback(func() {
	// 				done <- true
	// 			})))
	// 			speaker.Close()
	// 			<-done
	// 			glog.Errorln("Next alert at ", time.Now().Add(interval).Format("3:04:05 PM"))
	// 		}()
	// 	}
	// }

	buf := make(chan bool, 100)

	for i := 0; i < 10000; i++ {
		buf <- true
		go func() {

			out, err := exec.Command("attack.sh", strconv.FormatInt(fmt.Printf("%06d\n", i))).Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("The date is %s\n", out)
		}(i)

	}
}
