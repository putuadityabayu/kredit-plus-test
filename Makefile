mock:
	@echo "mock repositories"
	@mockgen --source=./repository/user.go -destination=./mock/repository/user.go --package=mock_repository

test:
	@go test -v -cover -coverpkg=./internal/service/... -coverprofile=./coverage/cover.out ./...

coverage:
	@go tool cover -func=./coverage/cover.out
	@go tool cover -html=./coverage/cover.out  -o ./coverage/cover.html

.PHONY: mock test coverage