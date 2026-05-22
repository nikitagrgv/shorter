package id

import "github.com/bwmarrin/snowflake"

type SnowflakeIdGenerator struct {
	node *snowflake.Node
}

func NewSnowflakeIdGenerator(id int64) (*SnowflakeIdGenerator, error) {
	node, err := snowflake.NewNode(id)
	if err != nil {
		return nil, err
	}
	return &SnowflakeIdGenerator{node}, nil
}

func (s *SnowflakeIdGenerator) Generate() int64 {
	return s.node.Generate().Int64()
}
