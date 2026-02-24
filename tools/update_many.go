package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBUpdateManyToolInput struct {
	DatabaseName   *string `json:"database_name,omitempty" jsonschema:"Optional name of the database to find the document in"`
	CollectionName string  `json:"collection_name" jsonschema:"Name of the collection to find the document in"`
	Filter         bson.M  `json:"filter" jsonschema:"The filter to find the document with"`
	Update         bson.M  `json:"update" jsonschema:"The update to apply to the document"`
	Upsert         *bool   `json:"upsert,omitempty" jsonschema:"Optional whether to insert the document if it doesn't exist, defaults to false"`
}

type MongoDBUpdateManyToolOutput struct {
	Result *mongo.UpdateResult `json:"result" jsonschema:"The result of the update operation"`
}

type NewMongoDBUpdateManyTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBUpdateManyTool() *NewMongoDBUpdateManyTool {
	return &NewMongoDBUpdateManyTool{
		tool: t,
	}
}

func (t *NewMongoDBUpdateManyTool) name() string {
	return "[MongoDB] Update Many Tool"
}

func (t *NewMongoDBUpdateManyTool) description() string {
	return "# Update many documents in MongoDB.\n\n" +
		"This tool can be used to update many documents in a MongoDB collection.\n\n"
}

func (t *NewMongoDBUpdateManyTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBUpdateManyToolInput,
) (
	*mcp.CallToolResult,
	MongoDBUpdateManyToolOutput,
	error,
) {
	defResponse := MongoDBUpdateManyToolOutput{
		Result: nil,
	}

	DB, err := t.tool.Database(input.DatabaseName)
	if err != nil {
		return nil, defResponse, err
	}

	collection := DB.Collection(input.CollectionName)

	opts := options.UpdateMany()
	if input.Upsert != nil && *input.Upsert {
		opts.SetUpsert(*input.Upsert)
	}

	res, err := collection.UpdateMany(ctx, input.Filter, input.Update, opts)
	if err != nil {
		return nil, defResponse, err
	}

	return nil, MongoDBUpdateManyToolOutput{
		Result: res,
	}, nil
}

func (t *NewMongoDBUpdateManyTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
