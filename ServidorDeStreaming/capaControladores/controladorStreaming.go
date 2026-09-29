/**
 * @file controladorStreaming.go
 * @brief Controlador que implementa el procedimiento remoto AudioStream del servicio gRPC.
 */
package capacontroladores

import (
	"fmt"

	"google.golang.org/grpc/peer"

	capafachada "streaming.local/servidor-streaming/capaFachadaServices/fachada"
	pb "streaming.local/servidor-streaming/serviciosAudio"
)

/**
 * @brief Controlador del servidor gRPC de streaming.
 */
type ControladorStreaming struct {
	pb.UnimplementedAudioServiceServer                               // Implementación por defecto requerida por gRPC.
	fachada                            *capafachada.FachadaStreaming ///< Fachada que realiza el streaming.
}

/**
 * @brief Crea el controlador de streaming con su fachada.
 * @return Puntero al nuevo controlador.
 */
func NuevoControladorStreaming() *ControladorStreaming {
	return &ControladorStreaming{fachada: capafachada.NuevaFachadaStreaming()}
}

/**
 * @brief Implementación del procedimiento remoto AudioStream.
 *
 * Imprime el eco de la invocación y delega a la fachada la lectura y el envío
 * de los fragmentos del audio.
 * @param peticion Petición con el archivo y los datos del audio a reproducir.
 * @param stream Flujo de salida por el que se envían los fragmentos.
 * @return Error gRPC si el audio no pudo transmitirse.
 */
func (thisC *ControladorStreaming) AudioStream(peticion *pb.AudioRequest, stream pb.AudioService_AudioStreamServer) error {
	direccionCliente := "desconocida"
	if datosPeer, ok := peer.FromContext(stream.Context()); ok {
		direccionCliente = datosPeer.Addr.String()
	}

	fmt.Printf("\n[gRPC] AudioStream invocado por %s (usuario: %s)\n", direccionCliente, peticion.Usuario)
	fmt.Printf("[gRPC] Audio solicitado: \"%s\" [%s] -> archivo %s\n", peticion.Titulo, peticion.Tipo, peticion.Filename)

	return thisC.fachada.TransmitirAudio(peticion, direccionCliente, stream)
}
