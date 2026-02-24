package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBFindOneAndUpdateToolInput struct {
	DatabaseName   *string `json:"database_name,omitempty" jsonschema:"Optional name of the database to find the document in"`
	CollectionName string  `json:"collection_name" jsonschema:"Name of the collection to find the document in"`
	Filter         bson.M  `json:"filter" jsonschema:"The filter to find the document with"`
	Update         bson.M  `json:"update" jsonschema:"The update to apply to the document"`
	Upsert         *bool   `json:"upsert,omitempty" jsonschema:"Optional whether to insert the document if it doesn't exist, defaults to false"`
}

type MongoDBFindOneAndUpdateToolOutput struct {
	Document bson.M `json:"document" jsonschema:"The document that was updated in the collection"`
}

type NewMongoDBFindOneAndUpdateTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBFindOneAndUpdateTool() *NewMongoDBFindOneAndUpdateTool {
	return &NewMongoDBFindOneAndUpdateTool{
		tool: t,
	}
}

func (t *NewMongoDBFindOneAndUpdateTool) name() string {
	return "[MongoDB] Find One and Update Tool"
}

func (t *NewMongoDBFindOneAndUpdateTool) description() string {
	return "# Find one document and update it in MongoDB.\n\n" +
		"This tool can be used to find one document in a MongoDB collection and update it.\n\n"
}

func (t *NewMongoDBFindOneAndUpdateTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBFindOneAndUpdateToolInput,
) (
	*mcp.CallToolResult,
	MongoDBFindOneAndUpdateToolOutput,
	error,
) {
	defResponse := MongoDBFindOneAndUpdateToolOutput{
		Document: bson.M{},
	}

	DB, err := t.tool.Database(input.DatabaseName)
	if err != nil {
		return nil, defResponse, err
	}

	collection := DB.Collection(input.CollectionName)

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if input.Upsert != nil && *input.Upsert {
		opts.SetUpsert(*input.Upsert)
	}

	var result bson.M
	err = collection.FindOneAndUpdate(ctx, input.Filter, input.Update, opts).
		Decode(&result)
	if err != nil {
		return nil, defResponse, err
	}

	return nil, MongoDBFindOneAndUpdateToolOutput{
		Document: result,
	}, nil
}

func (t *NewMongoDBFindOneAndUpdateTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
