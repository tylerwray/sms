# Smurf Messaging Service (SMS)

We have a big problem! The smurfs are in quarantine and they can't communicate directly with eachother anymore 😱 We need to build them a messaging platform so they can keep in touch

## Start Server

1. Install go `1.16`.
2. Create a `.env` file in the root of the project that looks like [this example](#example.env)
3. Run `bin/server` to start the server quickly.
4. Go to `http://localhost:8080/conversations/1956/messages` in the browser to see a response.

> Run `bin/server` to build and run a production like server.

## Ingest updates

1. Run `bin/ingest my_file.csv` to run the ingest task on a CSV file.

## Run tests

```bash
go test ./...
```

This project adheres to the [golang-standard project layout](https://github.com/golang-standards/project-layout).

## Example .env

```bash
DATABASE_URL="postgresql://postgres:postgres@localhost:5432/podium_messenger?sslmode=disable"
HOST="localhost"
PORT="8080"
````

