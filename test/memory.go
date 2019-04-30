package main

import "pkg.deepin.io/gir/gio-2.0"
import "time"
import "runtime"
import "syscall"
import "fmt"
import "os"

func RSSinMB() int {
	var r syscall.Rusage
	syscall.Getrusage(syscall.RUSAGE_SELF, &r)
	return int(r.Maxrss) / 1024
}

func TestFunc() {
	gio.AppInfoGetAll()
	gio.DesktopAppInfoSearch("d")
}

func main() {
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
		const limit = 30
		increasedInMB := getIncreased()
		if increasedInMB > limit {
			fmt.Printf("increased: %dMB\n", increasedInMB)
			panic(fmt.Errorf("increased beyond %dMB, detect a memory leak", limit))
		}
		runtime.GC()
	}
}
