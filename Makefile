# See LICENSE file for copyright and license details
.POSIX:

include config.mk

all: shrt

go.mod: ${SRC}
	${GO} mod tidy
	@touch go.mod

shrt: go.mod
	${GO} build -ldflags ${GO_LDFLAGS} ./cmd/shrt

clean:
	rm -f shrt

install: shrt
	install -Dm755 shrt ${DESTDIR}${PREFIX}/bin/shrt
	install -Dm644 shrt.1 ${DESTDIR}${MANPATH}/man1/shrt.1
	install -Dm644 shrt.1 ${DESTDIR}${MANPATH}/man5/shrtfile.5

uninstall:
	rm -f ${DESTDIR}${PREFIX}/bin/shrt
	rm -f ${DESTDIR}${MANPATH}/man1/shrt.1
	rm -f ${DESTDIR}${MANPATH}/man5/shrtfile.5

dist:
	rm -rf shrt-${VERSION}
	mkdir shrt-${VERSION}
	cp -r ${DIST_SRC} shrt-${VERSION}
	tar -cf - shrt-${VERSION} | gzip > shrt-${VERSION}.tar.gz

distclean:
	rm -rf shrt-${VERSION}
	rm -f shrt-${VERSION}.tar.gz

.PHONY: all clean install uninstall dist distclean
