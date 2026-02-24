package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBAggregateToolInput struct {
	DatabaseName   *string  `json:"database_name,omitempty" jsonschema:"Optional name of the database to find the document in"`
	CollectionName string   `json:"collection_name" jsonschema:"Name of the collection to find the document in"`
	Pipeline       []bson.M `json:"pipeline" jsonschema:"The aggregation pipeline to apply to the collection"`
	AllowDiskUse   *bool    `json:"allow_disk_use,omitempty" jsonschema:"Optional flag to allow disk use for the aggregation operation"`
	BatchSize      *int32   `json:"batch_size,omitempty" jsonschema:"Optional batch size for the aggregation operation"`
}

type MongoDBAggregateToolOutput struct {
	Result []bson.M `json:"result" jsonschema:"The result of the aggregation operation"`
}

type NewMongoDBAggregateTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBAggregateTool() *NewMongoDBAggregateTool {
	return &NewMongoDBAggregateTool{
		tool: t,
	}
}

func (t *NewMongoDBAggregateTool) name() string {
	return "[MongoDB] Aggregate Tool"
}

func (t *NewMongoDBAggregateTool) description() string {
	return "# Aggregate documents in MongoDB.\n\n" +
		"This tool can be used to perform aggregation operations on a MongoDB collection.\n\n"
}

func (t *NewMongoDBAggregateTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBAggregateToolInput,
) (
	*mcp.CallToolResult,
	MongoDBAggregateToolOutput,
	error,
) {
	defResponse := MongoDBAggregateToolOutput{
		Result: nil,
	}

	DB, err := t.tool.Database(input.DatabaseName)
	if err != nil {
		return nil, defResponse, err
	}

	collection := DB.Collection(input.CollectionName)

	opts := options.Aggregate()

	if input.AllowDiskUse != nil && *input.AllowDiskUse {
		opts.SetAllowDiskUse(true)
	}
	if input.BatchSize != nil && *input.BatchSize > 0 {
		opts.SetBatchSize(*input.BatchSize)
	}

	res, err := collection.Aggregate(ctx, input.Pipeline, opts)
	if err != nil {
		return nil, defResponse, err
	}

	var docs []bson.M
	if err := res.All(ctx, &docs); err != nil {
		return nil, defResponse, err
	}

	return nil, MongoDBAggregateToolOutput{
		Result: docs,
	}, nil
}

func (t *NewMongoDBAggregateTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
