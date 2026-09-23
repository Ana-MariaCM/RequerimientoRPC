package vistas

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"

	util "cliente.local/grpc-cliente/utilidades"
	pb "servidor.local/grpc-servidor/serviciosAudio"
)

// MostrarMenuPrincipal despliega el menu principal de audio mediante streaming gRPC.
// Devuelve false cuando el usuario decide salir de la aplicacion.
func MostrarMenuPrincipal(client pb.AudioServiceClient, ctxBase context.Context, readerInput *bufio.Reader) bool {

	fmt.Print("\n--- Menu de audio mediante streaming gRPC ---\n")
	fmt.Println("1. Reproducir un audio")
	fmt.Println("2. Salir")
	fmt.Print("Seleccione una opcion: ")

	opcion, _ := readerInput.ReadString('\n')
	opcion = strings.TrimSpace(opcion)

	switch opcion {
	case "1":
		reproducirAudio(client, ctxBase, readerInput)
		return true
	case "2":
		fmt.Println("Saliendo de la aplicacion...")
		return false
	default:
		fmt.Println("Opcion invalida, intente de nuevo.")
		return true
	}
}

// reproducirAudio solicita el titulo del audio, invoca el procedimiento remoto AudioStream
// y reproduce el audio recibido mediante streaming. El usuario puede abandonar la
// reproduccion en cualquier momento presionando Ctrl+C, regresando asi al menu principal
// (sin cerrar la aplicacion).
func reproducirAudio(client pb.AudioServiceClient, ctxBase context.Context, readerInput *bufio.Reader) {
	fmt.Print("Ingrese el titulo del audio: ")
	titulo, _ := readerInput.ReadString('\n')
	titulo = strings.TrimSpace(titulo)

	ctx, cancel := context.WithCancel(ctxBase)
	defer cancel()

	// Invocacion del procedimiento remoto (echo del lado del cliente)
	fmt.Printf("[RPC] Invocando AudioStream(filename=%s)...\n", titulo)
	stream, err := client.AudioStream(ctx, &pb.AudioRequest{Filename: titulo})
	if err != nil {
		log.Printf("Error invocando AudioStream: %v", err)
		return
	}

	fmt.Println("Recibiendo y reproduciendo audio en vivo...")
	fmt.Println("(Presione Ctrl+C en cualquier momento para detener y volver al menu)")

	reader, writer := io.Pipe()
	canalSincronizacion := make(chan struct{})
	canalFinRecepcion := make(chan struct{})

	// Captura Ctrl+C solo durante esta reproduccion, sin terminar el programa.
	senalInterrupcion := make(chan os.Signal, 1)
	signal.Notify(senalInterrupcion, os.Interrupt)
	defer signal.Stop(senalInterrupcion)

	// Arranca la goroutine de decodificacion y reproduccion de los fragmentos
	go util.DecodificarReproducir(reader, canalSincronizacion)

	// Arranca la recepcion de los fragmentos desde el servidor
	go func() {
		util.RecibirAudio(stream, writer, canalSincronizacion)
		close(canalFinRecepcion)
	}()

	select {
	case <-canalFinRecepcion:
		// La reproduccion finalizo de forma natural.
	case <-senalInterrupcion:
		fmt.Println("\nReproduccion cancelada por el usuario. Volviendo al menu...")
		// Al cancelar el contexto, stream.Recv() retorna con error y RecibirAudio
		// cierra el pipe y desbloquea la reproduccion por su cuenta.
		cancel()
		<-canalFinRecepcion
	}
}
