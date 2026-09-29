/**
 * @file fachadaStreaming.go
 * @brief Fachada que invoca el procedimiento remoto AudioStream del servidor de streaming (gRPC).
 */
package fachada

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	dtos "cliente.local/cliente/capaFachadaServices/DTOs"
	pb "streaming.local/servidor-streaming/serviciosAudio"
)

/**
 * @brief Fachada de acceso remoto al servidor de streaming.
 */
type FachadaStreaming struct {
	clienteGRPC pb.AudioServiceClient ///< Stub del servicio gRPC AudioService.
	usuario     string                ///< Usuario que se reporta en cada reproducción.
}

/**
 * @brief Crea la fachada de streaming.
 * @param clienteGRPC Stub del servicio AudioService.
 * @param usuario Identificación del usuario del cliente.
 * @return Puntero a la nueva fachada.
 */
func NuevaFachadaStreaming(clienteGRPC pb.AudioServiceClient, usuario string) *FachadaStreaming {
	return &FachadaStreaming{clienteGRPC: clienteGRPC, usuario: usuario}
}

/**
 * @brief Invoca el procedimiento remoto AudioStream y recibe el primer fragmento.
 *
 * Recibir el primer fragmento de forma síncrona permite informar de inmediato
 * si el audio no está disponible en el servidor.
 * @param ctx Contexto que permite cancelar la reproducción.
 * @param audio Audio que se desea reproducir.
 * @return El flujo gRPC, el primer fragmento recibido, o un error.
 */
func (thisF *FachadaStreaming) SolicitarAudio(ctx context.Context, audio dtos.AudioDetalleDTO) (pb.AudioService_AudioStreamClient, *pb.AudioChunk, error) {
	peticion := &pb.AudioRequest{
		Filename: audio.NombreArchivo,
		IdAudio:  int32(audio.Id),
		Titulo:   audio.Titulo,
		Tipo:     audio.NombreTipo,
		Usuario:  thisF.usuario,
	}

	// Invocación del procedimiento remoto.
	stream, err := thisF.clienteGRPC.AudioStream(ctx, peticion)
	if err != nil {
		return nil, nil, traducirError(err)
	}

	primerFragmento, err := stream.Recv()
	if err != nil {
		return nil, nil, traducirError(err)
	}
	return stream, primerFragmento, nil
}

/**
 * @brief Convierte un error gRPC en un mensaje comprensible para el usuario.
 * @param err Error devuelto por gRPC.
 * @return Error con un mensaje descriptivo.
 */
func traducirError(err error) error {
	estado, _ := status.FromError(err)
	switch estado.Code() {
	case codes.NotFound:
		return fmt.Errorf("%s", estado.Message())
	case codes.Unavailable:
		return fmt.Errorf("no fue posible comunicarse con el servidor de streaming")
	default:
		return fmt.Errorf("error en el streaming: %s", estado.Message())
	}
}
