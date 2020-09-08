clean:
	rm -rf vendor/

vendor: clean
	go mod vendor

gen: clean
	cd graph && go generate

ingest: vendor
	go run cmd/ingestion/main.go

run: ingest
	go run cmd/server/main.go
