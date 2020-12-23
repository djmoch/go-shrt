# See LICENSE file for copyright and license details
PREFIX := /usr/local
MANPATH := ${PREFIX}/man

SRC = cmd/shrt/*.go *.go
DIST_SRC = cmd *.go Makefile config.mk shrt.1 shrtfile.5 README LICENSE go.mod go.sum
VERSION := 0.2.3-dev0

GO_LDFLAGS := "-X main.version=${VERSION}"
GO := go
