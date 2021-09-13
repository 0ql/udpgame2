package server

import (
	"fmt"
	"log"
	"main/extrapolation/util"
	"math"
	"net"
	"os"
	"sync"
)

type Vec2 struct {
	X int32
	Y int32
}

func (p Vec2) Distance(p2 Vec2) float64 {
	first := math.Pow(float64(p2.X-p.X), 2)
	second := math.Pow(float64(p2.Y-p.Y), 2)
	return math.Sqrt(first + second)
}

type client struct {
	addr             net.Addr
	sid              uint16 // unique to the current session
	lastGameTickRecv uint32
	pos              Vec2
}

const (
	bufferSize = 100
)

var (
	confirmPack   = make([]byte, 1)
	lp            int
	err           error
	con           net.PacketConn
	connectionMap = map[string]client{}
	sids          = map[uint16]bool{}
	wg            *sync.WaitGroup
	ls            = log.New(os.Stderr, "\033[1;35mSERVER\033[0m ", log.Ltime)
	ticks         = uint32(0)
	incrementor   = uint16(0)
)

func init() {
	for i := uint16(0); i < 65535; i++ {
		sids[i] = false
	}
}

func removeClient(id string) {
	delete(connectionMap, id)
}

func confirmCon(addr net.Addr) {
	for { // TODO this is not smort
		if !sids[incrementor] {
			break
		}
		incrementor++
	}
	cl := client{
		addr: addr,
		sid:  incrementor,
		pos: Vec2{
			X: int32(incrementor * 50),
			Y: int32(incrementor * 50),
		},
	}
	incrementor++
	connectionMap[addr.String()] = cl
	_, err := con.WriteTo(confirmPack, addr)
	if err != nil {
		ls.Println(err)
	}
}

func sendInitialState(addr net.Addr) {
	packet := util.NewBuilder()
	packet.AddByte(1)
	packet.AddUint32(ticks)
	packet.AddUint32(uint32(connectionMap[addr.String()].pos.X))
	packet.AddUint32(uint32(connectionMap[addr.String()].pos.Y))
	_, err := con.WriteTo(packet.Buf, addr)
	if err != nil {
		ls.Println(err)
	}
}

func listenForPackets() {
	buf := make([]byte, bufferSize)

	defer wg.Done()
	decoder := util.NewDecoder(nil)
	for {
		_, addr, err := con.ReadFrom(buf)
		if err != nil {
			removeClient(addr.String())
		}

		decoder.SetPacket(buf)
		switch decoder.GetByte() {
		case 0:
			ls.Println("New Connection")
			confirmCon(addr)
			sendInitialState(addr)
			// case 1: // state packet
			// 	p := connectionMap[addr.String()]
			// 	p.lastGameTickRecv = decoder.GetUint32()
			// 	p.pos.X = int32(decoder.GetUint32())
			// 	p.pos.Y = int32(decoder.GetUint32())
		}
	}
}

func StartServer(localPort int, wgroup *sync.WaitGroup, hz int) {
	lp = localPort
	wg = wgroup
	con, err = net.ListenPacket("udp", fmt.Sprintf(":%d", lp))
	if err != nil {
		ls.Panic(err)
	}
	wg.Add(2)
	go listenForPackets()
	go clock()
	go gameLoop()
	go stateSender(hz)
}
