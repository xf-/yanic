package all

import (
	"log"
	"time"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Connection struct {
	database.Connection
	list []database.Connection
}

func Connect(allConnection map[string]interface{}) (database.Connection, error) {
	var list []database.Connection
	for dbType, conn := range database.Adapters {
		configForType := allConnection[dbType]
		if configForType == nil {
			log.Printf("the output type '%s' has no configuration\n", dbType)
			continue
		}
		dbConfigs, ok := configForType.([]map[string]interface{})
		if !ok {
			log.Panicf("the output type '%s' has the wrong format\n", dbType)
		}

		for _, config := range dbConfigs {
			if c, ok := config["enable"].(bool); ok && !c {
				continue
			}
			connected, err := conn(config)
			if err != nil {
				return nil, err
			}
			if connected == nil {
				continue
			}
			list = append(list, connected)
		}
	}
	return &Connection{list: list}, nil
}

func (conn *Connection) InsertNode(node *runtime.Node) {
	for _, item := range conn.list {
		item.InsertNode(node)
	}
}

func (conn *Connection) InsertLink(link *runtime.Link, time time.Time) {
	for _, item := range conn.list {
		item.InsertLink(link, time)
	}
}

func (conn *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time) {
	for _, item := range conn.list {
		item.InsertGlobals(stats, time)
	}
}

func (conn *Connection) PruneNodes(deleteAfter time.Duration) {
	for _, item := range conn.list {
		item.PruneNodes(deleteAfter)
	}
}

func (conn *Connection) Close() {
	for _, item := range conn.list {
		item.Close()
	}
}
