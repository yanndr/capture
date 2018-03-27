VERSION=`cat version`
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.build=${BUILD}"

build:
	go install ${LDFLAGS} 

linux:
	GOOS=linux go build ${LDFLAGS} 

windows:
	GOOS=windows go build ${LDFLAGS} 

install:
	go install ${LDFLAGS} 