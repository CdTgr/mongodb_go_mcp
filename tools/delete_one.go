package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBDeleteOneToolInput struct {
	DatabaseName   *string `json:"database_name,omitempty" jsonschema:"Optional name of the database to find the document in"`
	CollectionName string  `json:"collection_name" jsonschema:"Name of the collection to find the document in"`
	Filter         bson.M  `json:"filter" jsonschema:"The filter to find the document with"`
}

type MongoDBDeleteOneToolOutput struct {
	Result *mongo.DeleteResult `json:"result" jsonschema:"The result of the delete operation"`
}

type NewMongoDBDeleteOneTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBDeleteOneTool() *NewMongoDBDeleteOneTool {
	return &NewMongoDBDeleteOneTool{
		tool: t,
	}
}

func (t *NewMongoDBDeleteOneTool) name() string {
	return "[MongoDB] Delete One Tool"
}

func (t *NewMongoDBDeleteOneTool) description() string {
	return "# Delete one document in MongoDB.\n\n" +
		"This tool can be used to delete one document in a MongoDB collection.\n\n"
}

func (t *NewMongoDBDeleteOneTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBDeleteOneToolInput,
) (
	*mcp.CallToolResult,
	MongoDBDeleteOneToolOutput,
	error,
) {
	defResponse := MongoDBDeleteOneToolOutput{
		Result: nil,
	}

	DB, err := t.tool.Database(input.DatabaseName)
	if err != nil {
		return nil, defResponse, err
	}

	collection := DB.Collection(input.CollectionName)

	opts := options.DeleteOne()

	res, err := collection.DeleteOne(ctx, input.Filter, opts)
	if err != nil {
		return nil, defResponse, err
	}

	return nil, MongoDBDeleteOneToolOutput{
		Result: res,
	}, nil
}

func (t *NewMongoDBDeleteOneTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
