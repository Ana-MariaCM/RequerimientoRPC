/**
 * @file fachadaStreaming.go
 * @brief Fachada que encapsula la lógica de lectura, envío por streaming y
 *        notificación asíncrona de las reproducciones.
 */
package fachada

import (
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	capaaccesodatos "streaming.local/servidor-streaming/capaAccesoADatos"
	dtos "streaming.local/servidor-streaming/capaFachadaServices/DTOs"
	componenteconexioncola "streaming.local/servidor-streaming/componenteConexionCola"
	"streaming.local/servidor-streaming/configuracion"
	pb "streaming.local/servidor-streaming/serviciosAudio"
)

/**
 * @brief Fachada del servidor de streaming.
 */
type FachadaStreaming struct {
	repositorio  *capaaccesodatos.RepositorioAudios       ///< Acceso a los archivos mp3.
	conexionCola *componenteconexioncola.RabbitPublicador ///< Publicador de reproducciones en la cola.
}

/**
 * @brief Crea la fachada de streaming y la conexión con la cola de mensajes.
 *
 * Si RabbitMQ no está disponible el servidor sigue funcionando y el publicador
 * intentará reconectarse en la siguiente reproducción.
 * @return Puntero a la nueva fachada.
 */
func NuevaFachadaStreaming() *FachadaStreaming {
	fmt.Println("Inicializando fachada de streaming...")

	repositorio := capaaccesodatos.GetRepositorioAudios()
	fmt.Println("Carpeta de audios:", repositorio.ObtenerRutaAudios())

	conexionCola, err := componenteconexioncola.NuevoRabbitPublicador(
		configuracion.ObtenerURLRabbitMQ(), configuracion.NombreColaReproducciones)
	if err != nil {
		fmt.Println("Advertencia:", err, "(se reintentará al publicar)")
	} else {
		fmt.Println("Conectado a RabbitMQ, cola:", configuracion.NombreColaReproducciones)
	}

	return &FachadaStreaming{repositorio: repositorio, conexionCola: conexionCola}
}

/**
 * @brief Transmite un archivo de audio al cliente fragmento a fragmento.
 *
 * Antes de iniciar el envío publica, de forma asíncrona, la información de la
 * reproducción en la cola de mensajes. Si el cliente abandona la reproducción
 * el envío se detiene.
 * @param peticion Petición gRPC con el archivo y los datos del audio.
 * @param direccionCliente Dirección IP:puerto del cliente.
 * @param stream Flujo gRPC por el que se envían los fragmentos.
 * @return nil si el audio se envió completo, o un error gRPC (NotFound, Canceled, Internal).
 */
func (thisF *FachadaStreaming) TransmitirAudio(peticion *pb.AudioRequest, direccionCliente string, stream pb.AudioService_AudioStreamServer) error {
	archivo, err := thisF.repositorio.AbrirArchivo(peticion.Filename)
	if err != nil {
		fmt.Printf("[gRPC] No se encontró el audio \"%s\"\n", peticion.Filename)
		return status.Errorf(codes.NotFound, "el audio %s no está disponible en el servidor", peticion.Filename)
	}
	defer archivo.Close()

	// Comunicación asíncrona: la publicación no bloquea el streaming.
	go thisF.publicarReproduccion(peticion, direccionCliente)

	return thisF.enviarFragmentos(archivo, stream)
}

/**
 * @brief Lee el archivo en fragmentos de 32 KB y los envía por el stream gRPC.
 * @param lector Archivo de audio abierto.
 * @param stream Flujo gRPC por el que se envían los fragmentos.
 * @return nil si se envió completo, o un error gRPC.
 */
func (thisF *FachadaStreaming) enviarFragmentos(lector io.Reader, stream pb.AudioService_AudioStreamServer) error {
	buffer := make([]byte, configuracion.TamanioFragmento)
	numeroFragmento := int32(0)

	for {
		n, err := lector.Read(buffer)
		if n > 0 {
			numeroFragmento++
			fragmento := &pb.AudioChunk{Data: buffer[:n], NumeroFragmento: numeroFragmento}
			if errEnvio := stream.Send(fragmento); errEnvio != nil {
				if stream.Context().Err() != nil {
					fmt.Printf("[gRPC] El cliente abandonó la reproducción tras %d fragmentos\n", numeroFragmento-1)
					return status.Error(codes.Canceled, "reproducción cancelada por el cliente")
				}
				return status.Errorf(codes.Internal, "error enviando el fragmento #%d: %v", numeroFragmento, errEnvio)
			}
			fmt.Printf("[gRPC] Fragmento #%d enviado (%d bytes)\n", numeroFragmento, n)
		}
		if err == io.EOF {
			fmt.Printf("[gRPC] Audio enviado completo (%d fragmentos)\n", numeroFragmento)
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "error leyendo el archivo: %v", err)
		}
	}
}

/**
 * @brief Publica en la cola la información de una nueva reproducción.
 * @param peticion Petición gRPC con los datos del audio.
 * @param direccionCliente Dirección IP:puerto del cliente.
 */
func (thisF *FachadaStreaming) publicarReproduccion(peticion *pb.AudioRequest, direccionCliente string) {
	if thisF.conexionCola == nil {
		return
	}

	reproduccion := dtos.ReproduccionDTOOutput{
		IdAudio:          peticion.IdAudio,
		Titulo:           peticion.Titulo,
		Tipo:             peticion.Tipo,
		NombreArchivo:    peticion.Filename,
		Usuario:          peticion.Usuario,
		DireccionCliente: direccionCliente,
		FechaHora:        time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := thisF.conexionCola.Publicar(reproduccion); err != nil {
		fmt.Println("[Cola] No se pudo publicar la reproducción:", err)
	}
}
