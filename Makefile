PACKAGE_URL=https://ohse.de/uwe/software/streamedit/streamedit.html
PACKAGE_VERSION=0.1
PACKAGE_BUGREPORT=uwe@ohse.de
PACKAGE_ISSUES=https://github.com/UweOhse/streamedit/issues
PACKAGE_NAME=streamedit
PREFIX=/usr/local
BINDIR=$(PREFIX)/bin
MAN1DIR=$(PREFIX)/share/man/man1/

BUILD=`git rev-parse HEAD`

GO_LDFLAGS=-ldflags "-X main.Build=${BUILD}"
GOCMD=go
GOBUILD=$(GOCMD) build $(go_ldflags)

all: streamedit streamedit.md README.md streamedit.1

streamedit: streamedit.go version.go
	$(GOBUILD) -o $@ $^

version.go: version.stamp

version.stamp: version.in Makefile
	@echo recreating version.go if needed
	@sed -e 's%@PACKAGE_VERSION@%$(PACKAGE_VERSION)%' \
		-e 's%@PACKAGE_URL@%$(PACKAGE_URL)%' \
		-e 's%@PACKAGE_BUGREPORT@%$(PACKAGE_BUGREPORT)%' \
		-e 's%@PACKAGE_NAME@%$(PACKAGE_NAME)%' \
		-e 's%@PACKAGE_ISSUES@%$(PACKAGE_ISSUES)%' \
		$< >version.t
	@cmp -s version.t version.go || mv version.t version.go
	@rm -f version.t

cover: cover.out

streamedit.md: streamedit.md.in Makefile
	@sed -e 's%@PACKAGE_VERSION@%$(PACKAGE_VERSION)%' \
		-e 's%@PACKAGE_URL@%$(PACKAGE_URL)%' \
		-e 's%@PACKAGE_BUGREPORT@%$(PACKAGE_BUGREPORT)%' \
		-e 's%@PACKAGE_NAME@%$(PACKAGE_NAME)%' \
		-e 's%@PACKAGE_ISSUES@%$(PACKAGE_ISSUES)%' \
		$< >$@.t
	mv $@.t $@

install-man: streamedit.1
	install -c -d $(MAN1DIR)
	install -c $< $(MAN1DIR)
install-bin: streamedit
	install -c -d $(BINDIR)
	install -c $< $(BINDIR)

streamedit.1: streamedit.md
	pandoc --standalone --from markdown-smart-definition_lists+escaped_line_breaks --to man $< -o $@.t
	mv $@.t $@

install: install-man install-bin
	

cover.out: streamedit-covering test
	COVER=1 sh test.sh >/dev/null
	go tool cover -func cover.out
	rm cover1.out

streamedit-covering: streamedit.go version.go
	go test -coverpkg="./..." -c -tags testrunmain -o $@ .

coverupload: cover.html
	scp cover.html uwe@ohse.de:oldweb/uwe/misc/streamedit.cover.html

cover.html: cover.out
	go tool cover --html=cover.out -o cover.html

test check: streamedit test.sh
	sh test.sh >test.out
	diff test.expect test.out

clean:
	rm -f streamedit streamedit-covering test.out cover1.out cover.out cover.html \
		test.tmp version.go
