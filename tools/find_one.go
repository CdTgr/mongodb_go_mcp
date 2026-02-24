package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoDBFindOneToolInput struct {
	DatabaseName   *string `json:"database_name,omitempty" jsonschema:"Optional name of the database to find the document in"`
	CollectionName string  `json:"collection_name" jsonschema:"Name of the collection to find the document in"`
	Filter         bson.M  `json:"filter" jsonschema:"The filter to find the document with"`
}

type MongoDBFindOneToolOutput struct {
	Document bson.M `json:"document" jsonschema:"The document found in the collection"`
}

type NewMongoDBFindOneTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBFindOneTool() *NewMongoDBFindOneTool {
	return &NewMongoDBFindOneTool{
		tool: t,
	}
}

func (t *NewMongoDBFindOneTool) name() string {
	return "[MongoDB] Find One Tool"
}

func (t *NewMongoDBFindOneTool) description() string {
	return "# Find one document in MongoDB.\n\n" +
		"This tool can be used to find one document in a MongoDB collection.\n\n"
}

func (t *NewMongoDBFindOneTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBFindOneToolInput,
) (
	*mcp.CallToolResult,
	MongoDBFindOneToolOutput,
	error,
) {
	defResponse := MongoDBFindOneToolOutput{
		Document: bson.M{},
	}

	DB, err := t.tool.Database(input.DatabaseName)
	if err != nil {
		return nil, defResponse, err
	}

	collection := DB.Collection(input.CollectionName)

	var result bson.M
	err = collection.FindOne(ctx, input.Filter).Decode(&result)
	if err != nil {
		return nil, defResponse, err
	}

	output := MongoDBFindOneToolOutput{
		Document: result,
	}

	return nil, output, nil
}

func (t *NewMongoDBFindOneTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
