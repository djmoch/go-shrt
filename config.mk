PREFIX := /usr/local
MANPATH := ${PREFIX}/man

SRC = cmd/shrt/main.go shrtfile.go
DIST_SRC = cmd shrtfile.go Makefile config.mk shrt.1 shrtfile.5 README LICENSE go.mod go.sum
VERSION := 0.1.0

GO_LDFLAGS := "-X main.version=${VERSION}"
GO := go
