package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBInsertOneToolInput struct {
	DatabaseName   *string `json:"database_name,omitempty" jsonschema:"Optional name of the database to insert the document in"`
	CollectionName string  `json:"collection_name" jsonschema:"Name of the collection to insert the document in"`
	Document       bson.M  `json:"document" jsonschema:"The document to insert into the collection"`
}

type MongoDBInsertOneToolOutput struct {
	Result *mongo.InsertOneResult `json:"result" jsonschema:"The result of the insert operation"`
}

type NewMongoDBInsertOneTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBInsertOneTool() *NewMongoDBInsertOneTool {
	return &NewMongoDBInsertOneTool{
		tool: t,
	}
}

func (t *NewMongoDBInsertOneTool) name() string {
	return "[MongoDB] Insert One Tool"
}

func (t *NewMongoDBInsertOneTool) description() string {
	return "# Insert one document into a MongoDB collection.\n\n" +
		"This tool can be used to insert one document into a MongoDB collection.\n\n"
}

func (t *NewMongoDBInsertOneTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBInsertOneToolInput,
) (
	*mcp.CallToolResult,
	MongoDBInsertOneToolOutput,
	error,
) {
	defResponse := MongoDBInsertOneToolOutput{
		Result: nil,
	}

	DB, err := t.tool.Database(input.DatabaseName)
	if err != nil {
		return nil, defResponse, err
	}

	collection := DB.Collection(input.CollectionName)

	opts := options.InsertOne()

	res, err := collection.InsertOne(ctx, input.Document, opts)
	if err != nil {
		return nil, defResponse, err
	}

	return nil, MongoDBInsertOneToolOutput{
		Result: res,
	}, nil
}

func (t *NewMongoDBInsertOneTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
