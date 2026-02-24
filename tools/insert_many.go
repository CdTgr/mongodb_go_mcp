package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBInsertManyToolInput struct {
	DatabaseName   *string  `json:"database_name,omitempty" jsonschema:"Optional name of the database to insert the documents in"`
	CollectionName string   `json:"collection_name" jsonschema:"Name of the collection to insert the documents in"`
	Documents      []bson.M `json:"documents" jsonschema:"The documents to insert into the collection"`
}

type MongoDBInsertManyToolOutput struct {
	Result *mongo.InsertManyResult `json:"result" jsonschema:"The result of the insert operation"`
}

type NewMongoDBInsertManyTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBInsertManyTool() *NewMongoDBInsertManyTool {
	return &NewMongoDBInsertManyTool{
		tool: t,
	}
}

func (t *NewMongoDBInsertManyTool) name() string {
	return "[MongoDB] Insert Many Tool"
}

func (t *NewMongoDBInsertManyTool) description() string {
	return "# Insert many documents into a MongoDB collection.\n\n" +
		"This tool can be used to insert many documents into a MongoDB collection.\n\n"
}

func (t *NewMongoDBInsertManyTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBInsertManyToolInput,
) (
	*mcp.CallToolResult,
	MongoDBInsertManyToolOutput,
	error,
) {
	defResponse := MongoDBInsertManyToolOutput{
		Result: nil,
	}

	DB, err := t.tool.Database(input.DatabaseName)
	if err != nil {
		return nil, defResponse, err
	}

	collection := DB.Collection(input.CollectionName)

	opts := options.InsertMany()

	res, err := collection.InsertMany(ctx, input.Documents, opts)
	if err != nil {
		return nil, defResponse, err
	}

	return nil, MongoDBInsertManyToolOutput{
		Result: res,
	}, nil
}

func (t *NewMongoDBInsertManyTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
