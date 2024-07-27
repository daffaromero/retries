package main

func main() {
	httpServer := NewHTTPServer("localhost:9000")
	go httpServer.Run()

	gRPCServer := NewgRPCServer("localhost:8086")
	gRPCServer.Run()
}
