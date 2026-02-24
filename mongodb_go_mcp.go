package main

import (
	"context"
	"log"

	"github.com/CdTgr/mongodb_go_mcp/mcp/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "MongoDB MCP", Version: "v1.0.0"}, nil)

	// Create a tool and add it to the server.
	coreTools := tools.NewTool()
	coreTools.NewMongoDBListCollectionsTool().AttachTool(server)
	coreTools.NewMongoDBCountDocumentsTool().AttachTool(server)
	coreTools.NewMongoDBFindOneTool().AttachTool(server)
	coreTools.NewMongoDBFindTool().AttachTool(server)
	if !coreTools.ReadOnly {
		// Insert tools
		coreTools.NewMongoDBInsertOneTool().AttachTool(server)
		coreTools.NewMongoDBInsertManyTool().AttachTool(server)
		// Update tools
		coreTools.NewMongoDBFindOneAndUpdateTool().AttachTool(server)
		coreTools.NewMongoDBFindOneAndReplaceTool().AttachTool(server)
		coreTools.NewMongoDBUpdateOneTool().AttachTool(server)
		coreTools.NewMongoDBUpdateManyTool().AttachTool(server)
		// Delete tools
		coreTools.NewMongoDBDeleteOneTool().AttachTool(server)
		coreTools.NewMongoDBDeleteManyTool().AttachTool(server)
		coreTools.NewMongoDBFindOneAndDeleteTool().AttachTool(server)
	}

	if coreTools.AllowAggregates {
		// Aggregate tool
		coreTools.NewMongoDBAggregateTool().AttachTool(server)
	}

	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
