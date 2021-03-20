package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

// Init 初始化node
func Init(startTime string, machineTD int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineTD)
	return
}

// GenID 生成ID
func GenID() int64 {
	return node.Generate().Int64()
}
