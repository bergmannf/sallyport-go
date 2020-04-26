GO="go"

all: server client

clean:
	rm -f ./client
	rm -f ./server

client: clean
	${GO} build ./src/client

server: clean
	${GO} build ./src/server
