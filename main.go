package main

import (
	"runtime"
	// "fmt"
	// "time"
	"github.com/kopjenmbeng/evermos_online_store/cmd"
)

func main() {
	// fmt.Println(time.Now().Unix())
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Run()
}
