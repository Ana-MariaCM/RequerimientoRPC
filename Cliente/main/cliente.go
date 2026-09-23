package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"cliente.local/grpc-cliente/vistas"
	pb "servidor.local/grpc-servidor/serviciosAudio"
)

// direccionServidor es la direccion del servidor de streaming gRPC.
const direccionServidor = "localhost:50051"

func main() {
	client, conn := conectarServidorStreaming(direccionServidor)
	defer conn.Close()

	ejecutarAplicacion(client)
}

// conectarServidorStreaming establece la conexion gRPC con el servidor de streaming
// y devuelve el cliente listo para invocar el procedimiento remoto AudioStream.
func conectarServidorStreaming(direccion string) (pb.AudioServiceClient, *grpc.ClientConn) {
	conn, err := grpc.NewClient(direccion, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("No fue posible conectar con el servidor de streaming (%s): %v\n", direccion, err)
		os.Exit(1)
	}

	client := pb.NewAudioServiceClient(conn)
	return client, conn
}

// ejecutarAplicacion mantiene el ciclo del menu principal hasta que el usuario decide salir.
func ejecutarAplicacion(client pb.AudioServiceClient) {
	readerInput := bufio.NewReader(os.Stdin)
	ctx := context.Background()

	continuar := true
	for continuar {
		continuar = vistas.MostrarMenuPrincipal(client, ctx, readerInput)
	}
}
