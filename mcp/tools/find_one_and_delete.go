package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBFindOneAndDeleteToolInput struct {
	DatabaseName   *string `json:"database_name,omitempty" jsonschema:"Optional name of the database to find the document in"`
	CollectionName string  `json:"collection_name" jsonschema:"Name of the collection to find the document in"`
	Filter         bson.M  `json:"filter" jsonschema:"The filter to find the document with"`
}

type MongoDBFindOneAndDeleteToolOutput struct {
	Document bson.M `json:"document" jsonschema:"The document that was deleted in the collection"`
}

type NewMongoDBFindOneAndDeleteTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBFindOneAndDeleteTool() *NewMongoDBFindOneAndDeleteTool {
	return &NewMongoDBFindOneAndDeleteTool{
		tool: t,
	}
}

func (t *NewMongoDBFindOneAndDeleteTool) name() string {
	return "[MongoDB] Find One and Delete Tool"
}

func (t *NewMongoDBFindOneAndDeleteTool) description() string {
	return "# Find one document and delete it in MongoDB.\n\n" +
		"This tool can be used to find one document in a MongoDB collection and delete it.\n\n"
}

func (t *NewMongoDBFindOneAndDeleteTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBFindOneAndDeleteToolInput,
) (
	*mcp.CallToolResult,
	MongoDBFindOneAndDeleteToolOutput,
	error,
) {
	defResponse := MongoDBFindOneAndDeleteToolOutput{
		Document: bson.M{},
	}

	DB, err := t.tool.Database(input.DatabaseName)
	if err != nil {
		return nil, defResponse, err
	}

	collection := DB.Collection(input.CollectionName)

	opts := options.FindOneAndDelete()
	var result bson.M
	err = collection.FindOneAndDelete(ctx, input.Filter, opts).
		Decode(&result)
	if err != nil {
		return nil, defResponse, err
	}

	return nil, MongoDBFindOneAndDeleteToolOutput{
		Document: result,
	}, nil
}

func (t *NewMongoDBFindOneAndDeleteTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
