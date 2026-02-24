package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBFindToolInput struct {
	DatabaseName   *string `json:"database_name,omitempty" jsonschema:"Optional name of the database to find the document in"`
	CollectionName string  `json:"collection_name" jsonschema:"Name of the collection to find the document in"`
	Filter         bson.M  `json:"filter" jsonschema:"The filter to find the document with"`
	Skip           *int64  `json:"skip,omitempty" jsonschema:"Optional number of documents to skip"`
	Limit          *int64  `json:"limit,omitempty" jsonschema:"Optional maximum number of documents to return, defaults to 10"`
}

type MongoDBFindToolOutput struct {
	Documents []bson.M `json:"documents" jsonschema:"The documents found in the collection"`
	HasMore   bool     `json:"has_more" jsonschema:"Whether there are more documents to find"`
	Total     int64    `json:"total" jsonschema:"The total number of documents that match the filter"`
}

type NewMongoDBFindTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBFindTool() *NewMongoDBFindTool {
	return &NewMongoDBFindTool{
		tool: t,
	}
}

func (t *NewMongoDBFindTool) name() string {
	return "[MongoDB] Find Tool"
}

func (t *NewMongoDBFindTool) description() string {
	return "# Find documents in MongoDB.\n\n" +
		"This tool can be used to find multiple documents in a MongoDB collection.\n\n"
}

func (t *NewMongoDBFindTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBFindToolInput,
) (
	*mcp.CallToolResult,
	MongoDBFindToolOutput,
	error,
) {
	defResponse := MongoDBFindToolOutput{
		Documents: []bson.M{},
		HasMore:   false,
		Total:     0,
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

	filterOptions := options.Find().SetLimit(limit).SetSkip(skip)

	total, err := collection.CountDocuments(ctx, input.Filter)
	if err != nil {
		return nil, defResponse, err
	}

	var results []bson.M = []bson.M{}
	cursor, err := collection.Find(ctx, input.Filter, filterOptions)
	if err != nil {
		return nil, defResponse, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result bson.M
		err = cursor.Decode(&result)
		if err != nil {
			return nil, defResponse, err
		}
		results = append(results, result)
	}

	output := MongoDBFindToolOutput{
		Documents: results,
		HasMore:   int64(len(results))+skip < total,
		Total:     total,
	}

	return nil, output, nil
}

func (t *NewMongoDBFindTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
