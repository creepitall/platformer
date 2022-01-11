LOCAL_BIN=$(CURDIR)/bin
PROJECT_NAME=platformer

export GO111MODULE=on
GOENV:=GOPRIVATE="github.com/*" GO111MODULE=on

RUN_ARGS:=
ifneq (,$(wildcard cfg/values.yaml))
    RUN_ARGS=--local-config=cfg/values.yaml
endif

.PHONY: run
run:
	$(GOENV) go run cmd/main.go $(RUN_ARGS)