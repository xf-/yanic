package database

import (
	"testing"
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

type testConn struct {
	Connection
	CountClose int
	CountPrune int
}

func (c *testConn) Close() {
	c.CountClose++
}
func (c *testConn) PruneNodes(time.Duration) {
	c.CountPrune++
}

func TestStart(t *testing.T) {
	assert := assert.New(t)

	conn := &testConn{}
	config := &runtime.Config{
		Database: struct {
			DeleteInterval runtime.Duration `toml:"delete_interval"`
			DeleteAfter    runtime.Duration `toml:"delete_after"`
			Connection     map[string][]interface{}
		}{
			DeleteInterval: runtime.Duration{Duration: time.Millisecond * 10},
		},
	}
	assert.Nil(quit)

	Start(conn, config)
	assert.NotNil(quit)

	assert.Equal(0, conn.CountPrune)
	time.Sleep(time.Millisecond * 12)
	assert.Equal(1, conn.CountPrune)

	assert.Equal(0, conn.CountClose)
	Close(conn)
	assert.NotNil(quit)
	assert.Equal(1, conn.CountClose)

	time.Sleep(time.Millisecond * 12) // to reach timer.Stop() line

}
