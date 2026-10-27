.PHONY: test

VIEW_PORT ?= 8080

test:
	go test -tags="unit" ./...

view:
	go run ./_tools/visualizer -addr localhost:$(VIEW_PORT) -timeout 5s