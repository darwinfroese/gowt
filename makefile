.DEFAULT_GOAL: all

.PHONY: all
all: mux

.PHONY: mux
mux:
	go build mux/mux.go mux/utils.go mux/default_handlers.go mux/route.go mux/logger.go mux/type.go

