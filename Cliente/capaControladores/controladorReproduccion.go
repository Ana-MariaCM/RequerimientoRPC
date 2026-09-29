/**
 * @file controladorReproduccion.go
 * @brief Controlador que inicia y detiene la reproducción de audios mediante streaming.
 */
package capacontroladores

import (
	"context"

	dtos "cliente.local/cliente/capaFachadaServices/DTOs"
	"cliente.local/cliente/capaFachadaServices/fachada"
	"cliente.local/cliente/utilidades"
)

/**
 * @brief Controlador de reproducción de audios.
 */
type ControladorReproduccion struct {
	fachada *fachada.FachadaStreaming ///< Fachada de acceso al servidor de streaming.
}

/**
 * @brief Reproducción en curso que el usuario puede abandonar en cualquier momento.
 */
type SesionReproduccion struct {
	cancelar    context.CancelFunc      ///< Cancela la llamada gRPC en curso.
	reproductor *utilidades.Reproductor ///< Reproductor que decodifica y reproduce el audio.
}

/**
 * @brief Crea el controlador de reproducción.
 * @param fachadaStreaming Fachada de acceso al servidor de streaming.
 * @return Puntero al nuevo controlador.
 */
func NuevoControladorReproduccion(fachadaStreaming *fachada.FachadaStreaming) *ControladorReproduccion {
	return &ControladorReproduccion{fachada: fachadaStreaming}
}

/**
 * @brief Solicita el audio al servidor de streaming e inicia su reproducción.
 *
 * La recepción de fragmentos y la reproducción se ejecutan en goroutines, de
 * modo que la vista puede seguir atendiendo al usuario.
 * @param audio Audio a reproducir.
 * @return La sesión de reproducción o un error si el audio no pudo solicitarse.
 */
func (thisC *ControladorReproduccion) IniciarReproduccion(audio dtos.AudioDetalleDTO) (*SesionReproduccion, error) {
	ctx, cancelar := context.WithCancel(context.Background())

	stream, primerFragmento, err := thisC.fachada.SolicitarAudio(ctx, audio)
	if err != nil {
		cancelar()
		return nil, err
	}

	reproductor := utilidades.NuevoReproductor()
	go reproductor.Reproducir()
	go func() {
		utilidades.RecibirAudio(stream, primerFragmento, reproductor)
		<-reproductor.Finalizado()
		cancelar()
	}()

	return &SesionReproduccion{cancelar: cancelar, reproductor: reproductor}, nil
}

/**
 * @brief Abandona la reproducción: cancela el streaming y detiene el audio.
 */
func (thisS *SesionReproduccion) Detener() {
	thisS.cancelar()
	thisS.reproductor.Detener()
}

/**
 * @brief Devuelve un canal que se cierra cuando la reproducción termina.
 * @return Canal de sólo lectura.
 */
func (thisS *SesionReproduccion) Finalizada() <-chan struct{} {
	return thisS.reproductor.Finalizado()
}

/**
 * @brief Devuelve el error que interrumpió la reproducción, si lo hubo.
 * @return Error ocurrido o nil.
 */
func (thisS *SesionReproduccion) Error() error {
	return thisS.reproductor.Error()
}
