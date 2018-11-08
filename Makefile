
build:
	CGO_ENABLED=0 go build --ldflags '${EXTLDFLAGS}' -o ./bin/moonphases github.com/m-okeefe/moonphases
container:
	docker build -t meganokeefe/moonphases:latest . 
