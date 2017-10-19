package influxdb

import (
	"testing"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
)

func TestGlobalStats(t *testing.T) {
	stats := runtime.NewGlobalStats(createTestNodes())

	assert := assert.New(t)
	fields := GlobalStatsFields(stats)

	// check fields
	assert.EqualValues(3, fields["nodes"])
	conn := &Connection{
		points: make(chan *client.Point),
	}

	global := 0
	model := 0
	firmware := 0
	go func() {
		for p := range conn.points {
			switch p.Name() {
			case "global":
				global++
				break
			case "model":
				model++
				break
			default:
				firmware++
			}
		}
	}()
	conn.InsertGlobals(stats, time.Now())
	time.Sleep(time.Millisecond * 100)
	assert.Equal(1, global)
	assert.Equal(2, model)
	assert.Equal(1, firmware)
}

func createTestNodes() *runtime.Nodes {
	nodes := runtime.NewNodes(&runtime.Config{})

	nodeData := &runtime.Node{
		Online: true,
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 23,
			},
		},
		Nodeinfo: &data.NodeInfo{
			NodeID: "abcdef012345",
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
		},
	}
	nodeData.Nodeinfo.Software.Firmware.Release = "2016.1.6+entenhausen1"
	nodes.AddNode(nodeData)

	nodes.AddNode(&runtime.Node{
		Online: true,
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 2,
			},
		},
		Nodeinfo: &data.NodeInfo{
			NodeID: "112233445566",
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
		},
	})

	nodes.AddNode(&runtime.Node{
		Online: true,
		Nodeinfo: &data.NodeInfo{
			NodeID: "0xdeadbeef0x",
			VPN:    true,
			Hardware: data.Hardware{
				Model: "Xeon Multi-Core",
			},
		},
	})

	return nodes
}
