
# Main Makefile for cryptcheck
#
# Copyright 2018 © by Ollivier Robert <roberto@keltia.net>
#

.PATH=	cmd/observatory:.
GOBIN=	${GOPATH}/bin

GO=		go
GSRCS=	cmd/observatory/main.go
SRCS=	mozilla.go mozilla_subr.go types.go utils.go

BIN=	observatory
EXE=	${BIN}.exe

OPTS=	-ldflags="-s -w" -v

all: ${BIN}

${BIN}: ${GSRCS} ${SRCS} ${USRCS}
	${GO} build ${OPTS} ./cmd/...

${EXE}: ${GSRCS} ${SRCS} ${USRCS}
	GOOS=windows ${GO} build ${OPTS} ./cmd/...

test: ${BIN}
	${GO} test .

windows: ${EXE}
	GOOS=windows ${GO} build ${OPTS} ./cmd/...

install: ${BIN}
	${GO} install ${OPTS} ./cmd/...

lint: ${GSRCS} ${SRCS} ${USRCS}
	gometalinter .

clean:
	${GO} clean .
	${GO} clean ./cmd/...
	-rm -f ${BIN} ${EXE}

push:
	git push --all
	git push --tags
