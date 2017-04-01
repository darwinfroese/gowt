.DEFAULT_GOAL: all

.PHONY: all
all: mux

.PHONY: mux
mux:
	go build mux/mux.go mux/muxUtils.go