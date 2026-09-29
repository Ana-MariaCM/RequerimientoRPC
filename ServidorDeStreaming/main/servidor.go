/**
 * @file servidor.go
 * @brief Punto de entrada del servidor de streaming (gRPC).
 */
package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	capacontroladores "streaming.local/servidor-streaming/capaControladores"
	"streaming.local/servidor-streaming/configuracion"
	pb "streaming.local/servidor-streaming/serviciosAudio"
)

/**
 * @brief Registra el servicio AudioService y pone a escuchar el servidor gRPC.
 */
func main() {
	puerto := configuracion.ObtenerPuerto()
	escucha, err := net.Listen("tcp", ":"+puerto)
	if err != nil {
		panic(err)
	}

	servidorGRPC := grpc.NewServer()
	pb.RegisterAudioServiceServer(servidorGRPC, capacontroladores.NuevoControladorStreaming())

	fmt.Printf("Servidor de streaming gRPC escuchando en :%s...\n", puerto)
	if err := servidorGRPC.Serve(escucha); err != nil {
		panic(err)
	}
}
