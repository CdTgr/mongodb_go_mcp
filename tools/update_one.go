package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBUpdateOneToolInput struct {
	DatabaseName   *string `json:"database_name,omitempty" jsonschema:"Optional name of the database to find the document in"`
	CollectionName string  `json:"collection_name" jsonschema:"Name of the collection to find the document in"`
	Filter         bson.M  `json:"filter" jsonschema:"The filter to find the document with"`
	Update         bson.M  `json:"update" jsonschema:"The update to apply to the document"`
	Upsert         *bool   `json:"upsert,omitempty" jsonschema:"Optional whether to insert the document if it doesn't exist, defaults to false"`
}

type MongoDBUpdateOneToolOutput struct {
	Result *mongo.UpdateResult `json:"result" jsonschema:"The result of the update operation"`
}

type NewMongoDBUpdateOneTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBUpdateOneTool() *NewMongoDBUpdateOneTool {
	return &NewMongoDBUpdateOneTool{
		tool: t,
	}
}

func (t *NewMongoDBUpdateOneTool) name() string {
	return "[MongoDB] Update One Tool"
}

func (t *NewMongoDBUpdateOneTool) description() string {
	return "# Update one document in MongoDB.\n\n" +
		"This tool can be used to update one document in a MongoDB collection.\n\n"
}

func (t *NewMongoDBUpdateOneTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBUpdateOneToolInput,
) (
	*mcp.CallToolResult,
	MongoDBUpdateOneToolOutput,
	error,
) {
	defResponse := MongoDBUpdateOneToolOutput{
		Result: nil,
	}

	DB, err := t.tool.Database(input.DatabaseName)
	if err != nil {
		return nil, defResponse, err
	}

	collection := DB.Collection(input.CollectionName)

	opts := options.UpdateOne()
	if input.Upsert != nil && *input.Upsert {
		opts.SetUpsert(*input.Upsert)
	}

	res, err := collection.UpdateOne(ctx, input.Filter, input.Update, opts)
	if err != nil {
		return nil, defResponse, err
	}

	return nil, MongoDBUpdateOneToolOutput{
		Result: res,
	}, nil
}

func (t *NewMongoDBUpdateOneTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
