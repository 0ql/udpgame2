package client

import (
	"main/interpolation/util"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func (c client) ticker(hz int) {
	pause := time.NewTicker(time.Duration(1000 / hz * int(time.Millisecond)))
	for {
		<-pause.C
	}
}

func (c client) renderer(win *pixelgl.Window) {
	defer win.Destroy()

	circle := imdraw.New(nil)
	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)

		for _, other := range c.others {
			circle.Color = colornames.Coral
			circle.Push(pixel.V(float64(other.state.x), float64(other.state.y)))
		}

		circle.Color = colornames.Blueviolet
		circle.Push(pixel.V(float64(c.state.x), float64(c.state.y)))
		circle.Circle(25, 0)
		circle.Draw(win)

		win.Update()
	}
}

func (c client) buildStatePacket() []byte {
	packet := util.NewBuilder()
	packet.AddByte(1)
	packet.AddUint32(clientSideGameTicks)
	packet.AddUint32(uint32(c.state.x))
	packet.AddUint32(uint32(c.state.y))
	return packet.Buf
}

func (c client) sendStatePackets(hz int) {
	c.lc.Printf("Sending State Packets @%dhz", hz)
	p := c.buildStatePacket()
	duration := time.Duration(1000 / hz * int(time.Millisecond))
	for {
		time.Sleep(duration)
		c.con.Write(p)
	}
}
