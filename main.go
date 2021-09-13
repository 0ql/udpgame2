package main

import (
	"log"
	"main/interpolation/server"
	"os"
	"os/exec"
	"sync"

	"github.com/faiface/pixel/pixelgl"
)

func Start(args ...string) (p *os.Process, err error) {
	if args[0], err = exec.LookPath(args[0]); err == nil {
		var procAttr os.ProcAttr
		procAttr.Files = []*os.File{os.Stdin,
			os.Stdout, os.Stderr}
		p, err := os.StartProcess(args[0], args, &procAttr)
		if err == nil {
			return p, nil
		}
	}
	return nil, err
}

var wg sync.WaitGroup

func run() {
	log.SetPrefix("MAIN THREAD")
	server.StartServer(8080, &wg, 3)
	// client.NewClient("localhost:8080", 3)
	Start("./client")
	Start("./client")
	Start("./client")
	wg.Wait()
}

func main() {
	pixelgl.Run(run)
}

// func main() {
// 	go ping.CreateServer(":8080")
// 	ping.CreateClient("localhost:8080")
// }

// dl
// func main() {
// 	go dl.StartServer(":8080")
// 	go dl.StartClient("localhost:8080")

// 	dl.InitialiseWindow()
// }

/** first
func main() {
	wg.Add(1)
	defer wg.Done()
	go first.StartServer("localhost:8080")
	go first.StartClient("localhost:8080")
	wg.Wait()
}*/
