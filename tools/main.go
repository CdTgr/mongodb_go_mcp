package tools

import (
	"fmt"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Tool struct {
	ReadOnly         bool
	connectionString string
	database         string
	client           *mongo.Client
}

func NewTool() *Tool {
	tool := &Tool{
		ReadOnly: false,
	}

	tool.validateArgs()

	return tool
}

func (t *Tool) validateArgs() {
	dbName := os.Getenv("DB_NAME")
	dbURL := os.Getenv("DB_URL")
	ReadOnly := strings.ToLower(strings.TrimSpace(os.Getenv("READ_ONLY")))

	if ReadOnly == "true" || ReadOnly == "1" {
		t.ReadOnly = true
	}

	if dbURL == "" {
		log.Fatal("missing required environment variable: DB_URL")
	}

	t.database = dbName

	t.connectionString = dbURL

	t.connect()
}

func (t *Tool) connect() {
	client, err := mongo.Connect(
		options.Client().
			ApplyURI(t.connectionString).
			SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)),
	)

	if err != nil {
		log.Fatalf("Error connecting to database: %s", err.Error())
	}

	t.client = client
}

func (t *Tool) Database(database *string) (*mongo.Database, error) {
	if database != nil && *database != "" {
		return t.client.Database(*database), nil
	}

	if t.database != "" {
		return t.client.Database(t.database), nil
	}

	return nil, fmt.Errorf("Database selection is missing to execute the query")
}
