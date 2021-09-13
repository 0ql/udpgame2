package server

import (
	"main/extrapolation/util"
	"time"
)

var progress = false

func clock() {
	ticker := time.NewTicker(1000 / 20 * time.Millisecond)
	for {
		<-ticker.C
		progress = true
		ticks++
	}
}

func gameLoop() {
	for progress {
		progress = false
	}
}

func stateSender(hz int) {
	ticker := time.NewTicker(time.Duration(1000/hz) * time.Millisecond)
	builder := util.NewBuilder()
	var buf []byte
	for {
		<-ticker.C
		builder.Reset()
		builder.AddByte(2)
		builder.AddUint32(ticks)
		builder.AddUint32(uint32(int32(len(connectionMap))))
		buf = builder.Buf
		for _, client := range connectionMap {
			builder.Reset()
			builder.AddBytes(buf)
			for _, otherClient := range connectionMap {
				// check if other in range
				if client == otherClient {
					continue
				}
				if client.pos.Distance(otherClient.pos) < 100 {
					builder.AddUint16(otherClient.sid)
					builder.AddUint32(uint32(otherClient.pos.X))
					builder.AddUint32(uint32(otherClient.pos.Y))
				}
			}
			ls.Printf("Sending to %s: %v", client.addr.String(), builder.Buf)
			_, err := con.WriteTo(builder.Buf, client.addr)
			if err != nil {
				ls.Println(err)
			}
		}
	}
}
