all: golem

golem: main.go net.go mc/handshake.go mc/io.go mc/login.go mc/protocol.go mc/status.go
	go fmt 
	go vet && go build

clean:
	rm golem
