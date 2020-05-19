include version.mk
GOARCH ?= amd64
GO_TEST_FLAGS ?= -race
SERVICES = server finder
REGISTRY = docker.pkg.github.com/kind84/gospiga

# set pkgs to all packages
PKGS = ./...

# verbose mode
ifdef VERBOSE
	GO_TEST_FLAGS += -v
	GO_BUILD_FLAGS += -v
endif

default: build

go-generate:
	go generate $(PKGS)

build: go-generate
	go build $(GO_BUILD_FLAGS) -ldflags "all=$(GO_LDFLAGS)" $(PKGS)

test: go-generate
	go test $(GO_TEST_FLAGS) -ldflags "all=$(GO_LDFLAGS)" $(PKGS)

docker: docker-server docker-finder

docker-server: build-dependencies
	docker build -t gospiga/server server

docker-finder: build-dependencies
	docker build -t gospiga/finder finder

docker-dev: docker-server-dev docker-finder-dev docker-redis-dev

docker-server-dev: build-dependencies
	docker build -t gospiga/server-dev -f ./server/dev.Dockerfile server

docker-finder-dev: build-dependencies
	docker build -t gospiga/finder-dev -f ./finder/dev.Dockerfile finder

build-dependencies:
	docker build -t dependencies -f ./dependencies.Dockerfile .

docker-redis-dev:
	docker build -t redisearch-dev -f ./redisearch.Dockerfile .

docker-build: build-dependencies
	docker-compose build

docker-run: docker-build
	docker-compose up

release: docker
	for service in $(SERVICES); do \
		docker tag gospiga/$$service $(REGISTRY)/$$service:$(DOCKER_TAG); \
		docker push $(REGISTRY)/$$service:$(DOCKER_TAG); \
	done

release-dev: docker-dev
	for service in $(SERVICES); do \
		docker tag gospiga/$$service-dev $(REGISTRY)/$$service-dev:$(DOCKER_TAG); \
		docker push $(REGISTRY)/$$service-dev:$(DOCKER_TAG); \
	done; \
	docker tag gospiga/redisearch-dev $(REGISTRY)/redisearch-dev; \
	docker push $(REGISTRY)/redisearch-dev; \
	dgraph $(DGRAPH_TAG); \
	docker tag gospiga/dgraph-dev:$(DGRAPH_TAG) $(REGISTRY)/dgraph-dev:$(DGRAPH_TAG); \
	docker push $(REGISTRY)/dgraph-dev:$(DGRAPH_TAG)

