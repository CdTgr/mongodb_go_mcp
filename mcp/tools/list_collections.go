package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoDBListCollectionsToolInput struct {
	DatabaseName *string `json:"database_name,omitempty" jsonschema:"Optional name of the database to list collections from"`
}

type MongoDBListCollectionsToolOutput struct {
	Collections []string `json:"collections" jsonschema:"The list of collections in the database"`
}

type MongoDBListCollectionsTool struct {
	tool *Tool
}

func (t *Tool) NewMongoDBListCollectionsTool() *MongoDBListCollectionsTool {
	return &MongoDBListCollectionsTool{
		tool: t,
	}
}

func (t *MongoDBListCollectionsTool) name() string {
	return "[MongoDB] List Collections Tool"
}

func (t *MongoDBListCollectionsTool) description() string {
	return "# List collections in MongoDB.\n\n" +
		"This tool can be used to list all collections in a MongoDB database.\n\n"
}

func (t *MongoDBListCollectionsTool) toolCall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MongoDBListCollectionsToolInput,
) (
	*mcp.CallToolResult,
	MongoDBListCollectionsToolOutput,
	error,
) {
	defResponse := MongoDBListCollectionsToolOutput{
		Collections: []string{},
	}
	DB, err := t.tool.Database(input.DatabaseName)
	if err != nil {
		return nil, defResponse, err
	}
	collections, err := DB.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, defResponse, err
	}

	output := MongoDBListCollectionsToolOutput{
		Collections: collections,
	}

	return nil, output, nil
}

func (t *MongoDBListCollectionsTool) AttachTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        t.name(),
		Description: t.description(),
	}, t.toolCall)
}
