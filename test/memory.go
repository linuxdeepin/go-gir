package main

import (
	"pkg.deepin.io/gir/gio-2.0"
	"runtime"
	"strings"
)
import "time"
import "syscall"
import "fmt"
import "os"

func RSSinMB() int {
	var r syscall.Rusage
	err := syscall.Getrusage(syscall.RUSAGE_SELF, &r)
	if err != nil {
		panic(err)
	}
	return int(r.Maxrss) / 1024
}

func TestFunc() {
	gio.AppInfoGetAll()
	gio.DesktopAppInfoSearch("d")
	runtime.GC()
}

func main() {
	var limit = 30
	if strings.HasPrefix(runtime.GOARCH, "arm") {
		// GOARCH is arm or arm64
		limit = 200
	}
	fmt.Printf("limit: %dMB\n", limit)

	TestFunc()
	baseRSSinMB := RSSinMB()
	fmt.Printf("base: %dMB\n", baseRSSinMB)
	getIncreased := func() int {
		return  RSSinMB() - baseRSSinMB
	}

	time.AfterFunc(time.Second*60, func() {
		fmt.Printf("increased: %dMB\n", getIncreased())
		os.Exit(0)
	})
	for {
		TestFunc()
		increasedInMB := getIncreased()
		if increasedInMB > limit {
			fmt.Printf("increased: %dMB\n", increasedInMB)
			panic(fmt.Errorf("increased beyond %dMB, detect a memory leak", limit))
		}
	}
}
