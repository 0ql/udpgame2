package client

import (
	"fmt"
	"log"
	"main/extrapolation/util"
	"math"
	"net"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

type Vec struct {
	x int32
	y int32
}

type dude struct {
	sid   uint16
	state Vec
}

type client struct {
	lastServerTick uint32
	con            net.Conn
	state          Vec
	lc             *log.Logger
	others         map[uint16]dude
	playersOnline  uint32
}

var (
	clientSideGameTicks = uint32(0)
)

func (c client) sendConnectRequest() {
	c.lc.Println("Sending Con Request")
	packet := make([]byte, 1)
	_, err := c.con.Write(packet)
	if err != nil {
		c.lc.Panic(err)
	}
}

func (c client) clientListenForPackets() {
	buf := make([]byte, 100)

	decoder := util.NewDecoder(nil)
	for {
		n, err := c.con.Read(buf)
		if err != nil {
			log.Panic(err)
		}

		decoder.SetPacket(buf)
		switch decoder.GetByte() {
		case 0:
			c.lc.Println("Connection Confirmed")
		case 1:
			clientSideGameTicks = decoder.GetUint32()
			c.state.x = int32(decoder.GetUint32())
			c.state.y = int32(decoder.GetUint32())
			c.lc.Printf("Initial State Receved Pos: %v", c.state)
			go c.sendStatePackets(20)
		case 2: // myself not in state packet
			c.lastServerTick = decoder.GetUint32()
			c.playersOnline = decoder.GetUint32()
			playersInRange := math.Floor(float64(n-7) / 12)
			for i := float64(0); i < playersInRange; i++ {
				sid := decoder.GetUint16()
				c.others[sid] = dude{
					state: Vec{
						x: int32(decoder.GetUint32()),
						y: int32(decoder.GetUint32()),
					},
				}
			}
		}
	}
}

// accumulates all inputs and stores them in a buffer for evaluation
func inputManager(win pixelgl.Window) {
	// win.Pressed()
}

func NewClient(remoteAddress string, hz int) client {

	con, err := net.Dial("udp", remoteAddress)
	if err != nil {
		log.Panic(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:  "CLIENT",
		Bounds: pixel.R(0, 0, 400, 400),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Panic(err)
	}

	lc := log.New(os.Stderr, fmt.Sprintf("\033[1;34mCLIENT [%s]\033[0m ", con.LocalAddr().String()), log.Ltime)
	c := client{
		others: map[uint16]dude{},
		con:    con,
		lc:     lc,
	}

	go c.clientListenForPackets()
	go c.renderer(win)
	go c.ticker(hz)

	c.sendConnectRequest()

	return c
}
