package all

import (
	"errors"
	"testing"
	"time"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

type testConn struct {
	database.Connection
	CountNode    int
	CountLink    int
	CountGlobals int
	CountPrune   int
	CountClose   int
}

func (c *testConn) InsertNode(node *runtime.Node) {
	c.CountNode++
}
func (c *testConn) InsertLink(link *runtime.Link, time time.Time) {
	c.CountLink++
}
func (c *testConn) InsertGlobals(stats *runtime.GlobalStats, time time.Time) {
	c.CountGlobals++
}
func (c *testConn) PruneNodes(time.Duration) {
	c.CountPrune++
}
func (c *testConn) Close() {
	c.CountClose++
}

func TestStart(t *testing.T) {
	assert := assert.New(t)

	globalConn := &testConn{}
	database.RegisterAdapter("a", func(config interface{}) (database.Connection, error) {
		return globalConn, nil
	})
	database.RegisterAdapter("b", func(config interface{}) (database.Connection, error) {
		return globalConn, nil
	})
	database.RegisterAdapter("c", func(config interface{}) (database.Connection, error) {
		return globalConn, nil
	})
	database.RegisterAdapter("d", func(config interface{}) (database.Connection, error) {
		return nil, nil
	})
	database.RegisterAdapter("e", func(config interface{}) (database.Connection, error) {
		return nil, errors.New("blub")
	})
	allConn, err := Connect(map[string][]interface{}{
		"a": []interface{}{"a1", "a2"},
		"b": nil,
		"c": []interface{}{"c1"},
		"d": []interface{}{"d0"}, // fetch continue command in Connect
	})
	assert.NoError(err)

	assert.Equal(0, globalConn.CountNode)
	allConn.InsertNode(nil)
	assert.Equal(3, globalConn.CountNode)

	assert.Equal(0, globalConn.CountLink)
	allConn.InsertLink(nil, time.Now())
	assert.Equal(3, globalConn.CountLink)

	assert.Equal(0, globalConn.CountGlobals)
	allConn.InsertGlobals(nil, time.Now())
	assert.Equal(3, globalConn.CountGlobals)

	assert.Equal(0, globalConn.CountPrune)
	allConn.PruneNodes(time.Second)
	assert.Equal(3, globalConn.CountPrune)

	assert.Equal(0, globalConn.CountClose)
	allConn.Close()
	assert.Equal(3, globalConn.CountClose)

	_, err = Connect(map[string][]interface{}{
		"e": []interface{}{"give me an error"},
	})
	assert.Error(err)
}
