package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	capacontroladores "servidor.local/grpc-servidor/capaControladores"
	pb "servidor.local/grpc-servidor/serviciosAudio"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAudioServiceServer(grpcServer, &capacontroladores.ControladorServidor{})

	fmt.Println("Servidor gRPC escuchando en :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
