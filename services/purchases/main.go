package main

func main() {
	gRPCServer := NewgRPCServer("localhost:8086")
	gRPCServer.Run()
}