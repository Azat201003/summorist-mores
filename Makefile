all: test

MORES_FILE_PREFIX=../../storage/
MORES_FILE_SUFFIX=.txt

export MORES_FILE_PREFIX MORES_FILE_SUFFIX

test:
	@go test ./tests/... -v

