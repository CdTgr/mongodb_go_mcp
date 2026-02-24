package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBDeleteManyToolInput struct {
	DatabaseName   *string `json:"database_name,omitempty" jsonschema:"Optional name of the database to find the document in"`
	CollectionName string  `json:"collection_name" jsonschema:"Name of the collection to find the document in"`
	Filter         bson.M  `json:"filter" jsonschema:"The filter to find the document with"`
}

type MongoDBDeleteManyToolOutput struct {
	Result *mongo.DeleteResult `json:"result" jsonschema:"The result of the delete operation"`
}

type NewMongoDBDeleteManyTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBDeleteManyTool() *NewMongoDBDeleteManyTool {
	return &NewMongoDBDeleteManyTool{
		tool: t,
	}
}

func (t *NewMongoDBDeleteManyTool) name() string {
	return "[MongoDB] Delete Many Tool"
}

func (t *NewMongoDBDeleteManyTool) description() string {
	return "# Delete many documents in MongoDB.\n\n" +
		"This tool can be used to delete many documents in a MongoDB collection.\n\n"
}

func (t *NewMongoDBDeleteManyTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBDeleteManyToolInput,
) (
	*mcp.CallToolResult,
	MongoDBDeleteManyToolOutput,
	error,
) {
	defResponse := MongoDBDeleteManyToolOutput{
		Result: nil,
	}

	DB, err := t.tool.Database(input.DatabaseName)
	if err != nil {
		return nil, defResponse, err
	}

	collection := DB.Collection(input.CollectionName)

	opts := options.DeleteMany()

	res, err := collection.DeleteMany(ctx, input.Filter, opts)
	if err != nil {
		return nil, defResponse, err
	}

	return nil, MongoDBDeleteManyToolOutput{
		Result: res,
	}, nil
}

func (t *NewMongoDBDeleteManyTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
