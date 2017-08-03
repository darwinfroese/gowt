.DEFAULT_GOAL: all

.PHONY: all
all: mux

.PHONY: mux
mux:
	go build mux/mux.go mux/utils.go mux/muxHandlers.go mux/route.go mux/muxLogger.go mux/type.go

