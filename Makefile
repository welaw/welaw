-include .env
.DEFAULT_GOAL := build
DB_URL := "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):5432/$(DB_NAME)?sslmode=$(DB_SSL_MODE)"

api:
	@echo "==> Making API ..."
	protoc api/v1/*.proto --go_out=plugins=grpc:.

build:
	@echo "==> Building ..."
	go build -v ./...

clean:
	@echo "==> Cleaning ..."
	rm -rf $(REPOS_DIR)

image:
	@echo "==> Building Docker image ..."
	docker build -t welaw .

install:
	@echo "==> Installing ..."
	go install -v ./...

start: install
	@echo "==> Installing and starting ..."
	welaw

new: reset start

reset: reset-db clean

reset-db:
	@echo "==> Resetting database ..."
	-dropdb -h $(DB_HOST) -U $(DB_USERNAME) $(DB_NAME)
	createdb -h $(DB_HOST) -U $(DB_USERNAME) $(DB_NAME)
	migrate -path _migrations -database $(DB_URL) up

test:
	@echo "==> Executing tests ..."
	go test -v ./...

# Dependencies

install-libgit2:
	@echo "==> Installing libgit2 ..."
	wget https://github.com/libgit2/libgit2/archive/v0.25.0.tar.gz
	tar xzf v0.25.0.tar.gz
	cd libgit2-0.25.0/; \
		cmake .; \
		cmake --build . --target install; \
		ldconfig;
	rm -rf libgit2-0.25.0/
	rm -f v0.25.0.tar.gz

install-migrate:
	@echo "==> Installing migrate ..."
	go get -u -d github.com/mattes/migrate/cli
	go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli

install-protoc:
	@echo "==> Installing protobuf ..."
	go get github.com/golang/protobuf/...

install-go-deps:
	@echo "==> Installing Go dependencies ..."
	go get github.com/dgrijalva/jwt-go
	go get github.com/go-kit/kit/auth/jwt
	go get github.com/go-kit/kit/endpoint
	go get github.com/go-kit/kit/log
	go get github.com/go-kit/kit/metrics
	go get github.com/go-kit/kit/metrics/prometheus
	go get github.com/go-kit/kit/tracing/opentracing
	go get github.com/go-kit/kit/transport/grpc
	go get github.com/go-kit/kit/transport/http
	go get github.com/go-kit/kit/circuitbreaker
	go get github.com/go-kit/kit/ratelimit
	go get github.com/golang/protobuf/...
	go get github.com/google/uuid
	go get github.com/gorilla/mux
	go get github.com/gorilla/securecookie
	go get github.com/joho/godotenv
	go get github.com/opentracing/opentracing-go
	go get github.com/prometheus/client_golang/prometheus
	go get github.com/rs/cors
	go get golang.org/x/net/context
	go get golang.org/x/oauth2
	go get golang.org/x/oauth2/google
	go get google.golang.org/grpc
	go get google.golang.org/grpc/metadata
	-go get gopkg.in/libgit2/git2go.v25
	go get github.com/sony/gobreaker
	go get golang.org/x/time/rate
	go get cloud.google.com/go/storage
	go get github.com/aws/aws-sdk-go/aws
	go get github.com/aws/aws-sdk-go/aws/credentials
	go get github.com/aws/aws-sdk-go/aws/session
	go get github.com/aws/aws-sdk-go/service/s3
	go get github.com/aws/aws-sdk-go/service/s3/s3manager
	go get github.com/pierrre/archivefile/zip
	go get github.com/lib/pq

install-deps: install-libgit2 install-migrate install-protoc install-go-deps
	@echo "==> Installing dependencies ..."

install-test-deps: install-go-deps
	@echo "==> Installing testing dependencies ..."
	go get github.com/stretchr/testify/...

.PHONY: api clean reset-db install install-migrate test
