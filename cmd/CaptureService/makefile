VERSION=`cat version`

LDFLAGS=-ldflags "-X main.version=${VERSION}"

build:
	go install ${LDFLAGS} 

linux:
	GOOS=linux go build ${LDFLAGS} 

windows:
	GOOS=windows go build ${LDFLAGS} 

install:
	go install ${LDFLAGS} 