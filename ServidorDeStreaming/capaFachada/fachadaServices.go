package fachada

import (
	"io"
	"log"
	"os"

	capaaccesodatos "servidor.local/grpc-servidor/capaAccesoDatos"
	pb "servidor.local/grpc-servidor/serviciosAudio"
)

// tamanoFragmento define el tamano en bytes de cada fragmento (chunk) enviado por streaming.
const tamanoFragmento = 32 * 1024 // 32KB por chunk

// abrirArchivo actua como fachada para obtener el archivo de la cancion usando la capa de acceso a datos.
// Recibe el titulo (nombre del audio) y devuelve el *os.File abierto o un error.
func abrirArchivo(titulo string) (*os.File, error) {
	log.Printf("GetAudioFile llamado con titulo=%s", titulo)
	return capaaccesodatos.AbrirArchivo(titulo)
}

// EnviarFragmentosAudio lee el archivo de la cancion en chunks y los envia al stream gRPC.
// Esta funcion encapsula la logica de lectura y envio para mantener el servidor limpio.
func EnviarFragmentosAudio(titulo string, stream pb.AudioService_AudioStreamServer) error {
	// Eco obligatorio: se notifica en el servidor cada llamado a un procedimiento remoto.
	log.Printf("[RPC] Llamado a AudioStream con titulo=%s", titulo)

	file, err := abrirArchivo(titulo)
	if err != nil {
		log.Printf("Error abriendo el archivo de audio '%s': %v", titulo, err)
		return err
	}
	defer file.Close()

	buf := make([]byte, tamanoFragmento)
	chunkNum := 0

	for {
		// n es la cantidad de bytes que devolvio file.Read(buf).
		// en buf se guarda el fragmento leido
		n, err := file.Read(buf)

		if n > 0 {
			chunkNum++
			log.Printf("[RPC] Enviando fragmento #%d (%d bytes) de '%s'", chunkNum, n, titulo)

			if errEnvio := stream.Send(&pb.AudioChunk{Data: buf[:n]}); errEnvio != nil {
				log.Printf("Error enviando fragmento #%d: %v", chunkNum, errEnvio)
				return errEnvio
			}
		}

		if err == io.EOF {
			log.Printf("[RPC] Envio de '%s' finalizado. Total de fragmentos: %d", titulo, chunkNum)
			break
		}
		if err != nil {
			log.Printf("Error leyendo el archivo de audio '%s': %v", titulo, err)
			return err
		}
	}

	return nil
}
