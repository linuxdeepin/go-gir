package main

import "gir/gio-2.0"
import "time"
import "runtime"
import "syscall"
import "fmt"
import "os"

func RSSinMB() int64 {
	var r syscall.Rusage
	syscall.Getrusage(syscall.RUSAGE_SELF, &r)
	return r.Maxrss / int64(1024)
}

func TestFunc() {
	gio.AppInfoGetAll()
	//TODO: fix memory leak of DesktopAppInfoSearch("d")
}

func main() {
	time.AfterFunc(time.Second*60, func() {
		fmt.Printf("Memory Used: %dMB\n", RSSinMB())
		os.Exit(0)
	})
	for {
		TestFunc()
		if RSSinMB() > 20 {
			fmt.Println("V:", RSSinMB())
			panic("RSS beyond 15MB, detect a memory leak!")
		}
		runtime.GC()
	}
}
