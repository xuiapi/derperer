package speedtest

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/sourcegraph/conc"
	"tailscale.com/derp"
	"tailscale.com/derp/derphttp"
	"tailscale.com/types/key"
)

type BandWidthResult struct {
	TotalBytesSent Unit
	Bps            Unit
	Latency        time.Duration
}

func measure(c1, c2 *derphttp.Client, c2DstKey key.NodePublic, duration time.Duration) (*BandWidthResult, error) {
	packetSize := 64 * 1024
	var packetCount int
	var totalLatency time.Duration
	res := &BandWidthResult{}

	var wg conc.WaitGroup

	wg.Go(func() {
		t := time.After(duration)

		randBuf := make([]byte, packetSize)
		if _, err := rand.Read(randBuf); err != nil {
			panic(err)
		}
		for {
			select {
			case <-t:
				return
			default:
				// construct packet
				// marshal the timestamp into first 8 bytes
				binary.LittleEndian.PutUint64(randBuf, uint64(time.Now().UnixNano()))

				if err := c1.Send(c2DstKey, randBuf); err != nil {
					panic(err)
				}
			}
		}
	})

	wg.Go(func() {
		t := time.After(duration)
		start := time.Now()
		for {
			select {
			case <-t:
				elapsed := time.Since(start)
				res.Bps = Unit{float64(packetCount*packetSize*8) / elapsed.Seconds(), "bps"}
				res.TotalBytesSent = Unit{float64(packetCount * packetSize), "bytes"}
				res.Latency = totalLatency / time.Duration(packetCount) / 2
				return
			default:
				pkt, err := c2.Recv()
				if err != nil {
					panic(err)
				}
				p, ok := pkt.(derp.ReceivedPacket)
				if !ok {
					panic(fmt.Errorf("got %T, want ReceivedPacket", p))
				}
				// unmarshal the timestamp from first 8 bytes
				timestamp := int64(binary.LittleEndian.Uint64(p.Data))

				totalLatency += time.Since(time.Unix(0, timestamp))

				// if len(p.Data) != packetSize {
				// 	panic(fmt.Errorf("got %d bytes, want %d bytes", len(p.Data), packetSize))
				// }

				packetCount++
			}
		}
	})

	if err := wg.WaitAndRecover(); err != nil {
		return nil, err.AsError()
	}

	return res, nil
}
