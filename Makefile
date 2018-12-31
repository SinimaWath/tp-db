
SWAGGER = github.com/go-swagger/go-swagger
BINDATA = github.com/jteeuwen/go-bindata
GENERATE_PATH = ./internal/restapi
CMD_PATH = ./internal/cmd/forum-server
EXE_NAME = forum-server

run:
	$(EXE_NAME) --scheme=http --port=5001 --host=0.0.0.0 --database=postgres://postgres:admin@localhost/forum?sslmode=disable
install:
	go install $(CMD_PATH)

generate:
	go generate -x $(GENERATE_PATH)



start: build install

init:
	dep init

build: generate
	dep ensure

download-generators:
	rm -rf vendor/github.com/go-swagger/go-swagger/
	rm -rf vendor/github.com/jteeuwen/go-bindata/
	mkdir -p vendor/github.com/go-swagger/
	git clone https://github.com/go-swagger/go-swagger vendor/github.com/go-swagger/go-swagger/
	mkdir -p vendor/github.com/jteeuwen/go-bindata/
	git clone https://github.com/jteeuwen/go-bindata vendor/github.com/jteeuwen/go-bindata/

docker-start:
	docker build -t tp-db -f Dockerfile .
	./deploy/runDocker.sh


docker-stop:
	./deploy/stopDocker.sh

docker-run:
	./deploy/runDocker.sh

clear:
	rm -rf internal/cmd/ internal/restapi/operations/ internal/models internal/modules/assets/
	rm $(shell ls -p internal/restapi/ | grep -v -E "^configure" | grep -v / | sed 's/ /\n/g' | sed 's|^|internal/restapi/|')