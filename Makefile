
.PHONY: build
build:
	cd cmd/server && GOOS=linux go build
	cd cmd/analysis/footprint && go build

.PHONY: run
run: build
	docker run --rm -it -p 1234:1234 -v `pwd`:/src alpine /src/cmd/server/server

