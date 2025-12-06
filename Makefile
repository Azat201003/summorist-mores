all: run

ENVFILE=config.develop.env

include $(ENVFILE)
export $(shell sed '/^#/d; s/=.*//' $(ENVFILE))
test:
	go test ./tests/... -v

