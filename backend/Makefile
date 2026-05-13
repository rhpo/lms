.PHONY: run build seed test test-auth test-admin test-teacher test-student test-company test-upload test-notation reset-db

run:
	go run ./cmd/server

build:
	go build -o pfe-server ./cmd/server

seed:
	go run ./cmd/seed

test:
	go test ./tests/... -v

test-auth:
	go test ./tests/... -v -run TestAuth

test-admin:
	go test ./tests/... -v -run TestAdmin

test-teacher:
	go test ./tests/... -v -run TestTeacher

test-student:
	go test ./tests/... -v -run TestStudent

test-company:
	go test ./tests/... -v -run TestCompany

test-upload:
	go test ./tests/... -v -run TestUpload

test-notation:
	go test ./tests/... -v -run TestNotation

reset-db:
	del pfe.db