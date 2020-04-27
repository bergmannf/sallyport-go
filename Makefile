GO="go"

all: server client

clean: clean-server clean-client

clean-server:
	rm -rf ./server

clean-client:
	rm -rf ./client

client: clean-client
	${GO} build ./src/client

server: clean-server
	${GO} build ./src/server
