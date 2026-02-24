package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBFindOneAndReplaceToolInput struct {
	DatabaseName   *string `json:"database_name,omitempty" jsonschema:"Optional name of the database to find the document in"`
	CollectionName string  `json:"collection_name" jsonschema:"Name of the collection to find the document in"`
	Filter         bson.M  `json:"filter" jsonschema:"The filter to find the document with"`
	Replacement    bson.M  `json:"replacement" jsonschema:"The document to replace the existing document with"`
	Upsert         *bool   `json:"upsert,omitempty" jsonschema:"Optional whether to insert the document if it doesn't exist, defaults to false"`
}

type MongoDBFindOneAndReplaceToolOutput struct {
	Document bson.M `json:"document" jsonschema:"The document that was replaced in the collection"`
}

type NewMongoDBFindOneAndReplaceTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBFindOneAndReplaceTool() *NewMongoDBFindOneAndReplaceTool {
	return &NewMongoDBFindOneAndReplaceTool{
		tool: t,
	}
}

func (t *NewMongoDBFindOneAndReplaceTool) name() string {
	return "[MongoDB] Find One and Replace Tool"
}

func (t *NewMongoDBFindOneAndReplaceTool) description() string {
	return "# Find one document and replace it in MongoDB.\n\n" +
		"This tool can be used to find one document in a MongoDB collection and replace it.\n\n"
}

func (t *NewMongoDBFindOneAndReplaceTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBFindOneAndReplaceToolInput,
) (
	*mcp.CallToolResult,
	MongoDBFindOneAndReplaceToolOutput,
	error,
) {
	defResponse := MongoDBFindOneAndReplaceToolOutput{
		Document: bson.M{},
	}

	DB, err := t.tool.Database(input.DatabaseName)
	if err != nil {
		return nil, defResponse, err
	}

	collection := DB.Collection(input.CollectionName)

	opts := options.FindOneAndReplace().SetReturnDocument(options.After)
	if input.Upsert != nil && *input.Upsert {
		opts.SetUpsert(*input.Upsert)
	}

	var result bson.M
	err = collection.FindOneAndReplace(ctx, input.Filter, input.Replacement, opts).
		Decode(&result)
	if err != nil {
		return nil, defResponse, err
	}

	return nil, MongoDBFindOneAndReplaceToolOutput{
		Document: result,
	}, nil
}

func (t *NewMongoDBFindOneAndReplaceTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
