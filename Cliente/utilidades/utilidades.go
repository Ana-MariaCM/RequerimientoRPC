package utilidades

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"

	pb "servidor.local/grpc-servidor/serviciosAudio"
)

// DecodificarReproducir decodifica el flujo mp3 que llega por el pipe y lo reproduce
// mediante el parlante del sistema. Si la reproduccion es abandonada por el usuario
// (pipe cerrado antes de tiempo) o el mp3 no puede decodificarse, se cierra el canal
// de sincronizacion para no dejar al llamador bloqueado indefinidamente.
func DecodificarReproducir(reader io.Reader, canalSincronizacion chan struct{}) {
	streamer, format, err := mp3.Decode(io.NopCloser(reader))
	if err != nil {
		log.Printf("Error decodificando MP3 (posible cancelacion): %v", err)
		close(canalSincronizacion)
		return
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/2))

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(canalSincronizacion)
	})))
}

// RecibirAudio invoca de forma continua stream.Recv() para recibir los fragmentos
// (chunks) de audio enviados por el servidor de streaming y los escribe en el pipe
// que alimenta al decodificador/reproductor. Si la reproduccion es cancelada por el
// usuario (contexto cancelado) el ciclo se detiene sin terminar la aplicacion.
func RecibirAudio(
	stream pb.AudioService_AudioStreamClient,
	writer *io.PipeWriter,
	canalSincronizacion chan struct{}) {
	noFragmento := 0
	for {
		fragmento, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("\nCanción recibida completa.")
			writer.Close()
			break
		}
		if err != nil {
			log.Printf("Recepcion detenida (%v)", err)
			writer.CloseWithError(err)
			break
		}
		noFragmento++
		fmt.Printf("\nFragmento #%d recibido (%d bytes) reproduciendo ...", noFragmento, len(fragmento.Data))

		if _, err := writer.Write(fragmento.Data); err != nil {
			log.Printf("Error escribiendo en pipe: %v", err)
			break
		}
	}
	// Esperar hasta que termine la reproducción (o se detecte la cancelación)
	<-canalSincronizacion
	fmt.Println("Reproducción finalizada.")
}
