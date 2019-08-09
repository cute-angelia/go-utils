package snowflake

import (
	"fmt"
	snowflakeLib "github.com/bwmarrin/snowflake"
)

func NewSnowId(nodeId int64) (snowflakeLib.ID, error) {
	// Create a new Node with a Node number of 1
	node, err := snowflakeLib.NewNode(nodeId)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	// Generate a snowflake ID.
	id := node.Generate()

	return id, nil
}
