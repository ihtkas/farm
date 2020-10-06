package main

import (
	"github.com/golang/glog"
	"github.com/ihtkas/farm/account"
)

func main() {
	m := account.Manager{}
	err := m.Start()
	if err != nil {
		glog.Errorln(err)
		return
	}
	glog.Errorln("Done....")
}
