package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBCountDocumentsToolInput struct {
	DatabaseName   *string `json:"database_name,omitempty" jsonschema:"Optional name of the database to find the document in"`
	CollectionName string  `json:"collection_name" jsonschema:"Name of the collection to find the document in"`
	Filter         bson.M  `json:"filter" jsonschema:"The filter to find the document with"`
	Skip           *int64  `json:"skip,omitempty" jsonschema:"Optional number of documents to skip"`
	Limit          *int64  `json:"limit,omitempty" jsonschema:"Optional maximum number of documents to return, defaults to 10"`
}

type MongoDBCountDocumentsToolOutput struct {
	Count int64 `json:"count" jsonschema:"The number of documents that match the filter"`
}

type NewMongoDBCountDocumentsTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBCountDocumentsTool() *NewMongoDBCountDocumentsTool {
	return &NewMongoDBCountDocumentsTool{
		tool: t,
	}
}

func (t *NewMongoDBCountDocumentsTool) name() string {
	return "[MongoDB] Count Documents Tool"
}

func (t *NewMongoDBCountDocumentsTool) description() string {
	return "# Count documents in MongoDB.\n\n" +
		"This tool can be used to count the number of documents in a MongoDB collection that match a given filter.\n\n"
}

func (t *NewMongoDBCountDocumentsTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBCountDocumentsToolInput,
) (
	*mcp.CallToolResult,
	MongoDBCountDocumentsToolOutput,
	error,
) {
	defResponse := MongoDBCountDocumentsToolOutput{
		Count: 0,
	}

	DB, err := t.tool.Database(input.DatabaseName)
	if err != nil {
		return nil, defResponse, err
	}

	collection := DB.Collection(input.CollectionName)

	var limit int64 = 10
	if input.Limit != nil && *input.Limit > 0 {
		limit = *input.Limit
	}
	var skip int64 = 0
	if input.Skip != nil && *input.Skip > 0 {
		skip = *input.Skip
	}

	filterOptions := options.Count().SetLimit(limit).SetSkip(skip)

	total, err := collection.CountDocuments(ctx, input.Filter, filterOptions)
	if err != nil {
		return nil, defResponse, err
	}

	return nil, MongoDBCountDocumentsToolOutput{
		Count: total,
	}, nil
}

func (t *NewMongoDBCountDocumentsTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
