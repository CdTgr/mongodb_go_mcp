# MCP Server for MongoDB in Golang

This project implements a MongoDB-compatible server in Golang, allowing clients to connect and interact with it as if it were a standard MongoDB server. The server supports basic CRUD operations, authentication, and other MongoDB features.

## Features

The following mongodb features are supported:

- Aggregate
- CountDocuments
- DeleteMany
- DeleteOne
- FindOneAndDelete
- FindOneAndReplace
- FindOneAndUpdate
- FindOne
- Find
- InsertMany
- InsertOne
- UpdateMany
- UpdateOne
- ListCollections

## Configurations

The server can be configured to run with the following environment variables:

```sh
DB_URL=
DB_NAME=
READ_ONLY=
```

| Variable | Description | Required | Default |
| --- | --- | --- | --- |
| `DB_URL` | The connection string for the MongoDB instance to connect to. | Yes | None |
| `DB_NAME` | The name of the MongoDB database to use. If not provided, the server will require the database name to be specified in each query. | No | None |
| `READ_ONLY` | If set to "true" or "1", the server will operate in read-only mode, disallowing any write operations. | No | false |


## Usage

To run the server, use the following command:

```sh
go run main.go
```

## Testing with MCP

Install the MCP Inspector using the following command:

```sh
npm install -g @modelcontextprotocol/inspector
```

Now, set your environment variables and run the server:

```sh
mcp-inspector go run .
```

## Installing and running the project for production

Install the project with the following command:

```sh
go install github.com/CdTgr/mongodb_go_mcp@latest
```

Now, you can freely connect your favourite MCP client to the server with the following execution command:

```sh
mongodb_go_mcp
```

> Note: Make sure to set the required environment variables before running the server.
