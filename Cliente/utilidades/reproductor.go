/**
 * @file reproductor.go
 * @brief Recepción de los fragmentos por streaming y reproducción del audio mientras llega.
 *
 * Los fragmentos recibidos se escriben en una tubería (io.Pipe); en paralelo, el
 * decodificador mp3 lee de la tubería y el altavoz reproduce el audio, de modo
 * que la reproducción inicia sin descargar el archivo completo.
 */
package utilidades

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"

	pb "streaming.local/servidor-streaming/serviciosAudio"
)

/** @brief Error usado para cerrar la tubería cuando el usuario abandona la reproducción. */
var errReproduccionDetenida = errors.New("reproducción detenida por el usuario")

/**
 * @brief Reproductor de audio mp3 alimentado por streaming.
 */
type Reproductor struct {
	lector     *io.PipeReader ///< Extremo de lectura de la tubería (lo consume el decodificador).
	escritor   *io.PipeWriter ///< Extremo de escritura de la tubería (lo alimentan los fragmentos).
	finalizado chan struct{}  ///< Se cierra cuando la reproducción termina o se detiene.
	unaVez     sync.Once      ///< Garantiza que el canal finalizado se cierre una sola vez.
	mu         sync.Mutex     ///< Protege el estado detenido durante la inicialización del altavoz.
	detenido   bool           ///< Indica si el usuario abandonó la reproducción.
	muError    sync.Mutex     ///< Protege el campo err.
	err        error          ///< Primer error ocurrido durante la recepción o la reproducción.
}

/**
 * @brief Crea un reproductor con su tubería de datos.
 * @return Puntero al nuevo reproductor.
 */
func NuevoReproductor() *Reproductor {
	lector, escritor := io.Pipe()
	return &Reproductor{lector: lector, escritor: escritor, finalizado: make(chan struct{})}
}

/**
 * @brief Decodifica el mp3 a medida que llegan los fragmentos y lo envía al altavoz.
 *
 * Se debe ejecutar en una goroutine: bloquea hasta recibir la cabecera del mp3.
 */
func (thisR *Reproductor) Reproducir() {
	streamer, formato, err := mp3.Decode(thisR.lector)
	if err != nil {
		thisR.finalizarConError(fmt.Errorf("no se pudo decodificar el audio: %v", err))
		return
	}

	thisR.mu.Lock()
	if thisR.detenido {
		thisR.mu.Unlock()
		streamer.Close()
		return
	}

	errAltavoz := speaker.Init(formato.SampleRate, formato.SampleRate.N(time.Second/10))
	if errAltavoz == nil {
		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			if errStreamer := streamer.Err(); errStreamer != nil && !errors.Is(errStreamer, errReproduccionDetenida) {
				thisR.registrarError(errStreamer)
			}
			thisR.marcarFinalizado()
		})))
	}
	thisR.mu.Unlock()

	if errAltavoz != nil {
		streamer.Close()
		thisR.finalizarConError(fmt.Errorf("no se pudo inicializar el dispositivo de audio: %v", errAltavoz))
	}
}

/**
 * @brief Detiene la reproducción en cualquier momento.
 *
 * Primero cierra la tubería para desbloquear al decodificador y luego retira
 * el audio del altavoz.
 */
func (thisR *Reproductor) Detener() {
	thisR.mu.Lock()
	thisR.detenido = true
	thisR.mu.Unlock()

	thisR.lector.CloseWithError(errReproduccionDetenida)
	speaker.Clear()
	thisR.marcarFinalizado()
}

/**
 * @brief Devuelve un canal que se cierra cuando la reproducción termina.
 * @return Canal de sólo lectura.
 */
func (thisR *Reproductor) Finalizado() <-chan struct{} {
	return thisR.finalizado
}

/**
 * @brief Devuelve el error que interrumpió la reproducción (si lo hubo).
 * @return Error ocurrido o nil.
 */
func (thisR *Reproductor) Error() error {
	thisR.muError.Lock()
	defer thisR.muError.Unlock()
	return thisR.err
}

/**
 * @brief Registra el primer error ocurrido.
 * @param err Error a registrar.
 */
func (thisR *Reproductor) registrarError(err error) {
	thisR.muError.Lock()
	defer thisR.muError.Unlock()
	if thisR.err == nil {
		thisR.err = err
	}
}

/**
 * @brief Registra un error, cierra la tubería y marca la reproducción como finalizada.
 * @param err Error ocurrido.
 */
func (thisR *Reproductor) finalizarConError(err error) {
	thisR.mu.Lock()
	detenido := thisR.detenido
	thisR.mu.Unlock()
	if !detenido {
		thisR.registrarError(err)
	}
	thisR.lector.CloseWithError(err)
	thisR.marcarFinalizado()
}

/**
 * @brief Cierra (una sola vez) el canal que indica el fin de la reproducción.
 */
func (thisR *Reproductor) marcarFinalizado() {
	thisR.unaVez.Do(func() { close(thisR.finalizado) })
}

/**
 * @brief Recibe los fragmentos del servidor de streaming y los entrega al reproductor.
 *
 * Termina cuando se recibe el audio completo, cuando ocurre un error o cuando
 * el usuario abandona la reproducción.
 * @param stream Flujo gRPC del procedimiento remoto AudioStream.
 * @param primerFragmento Fragmento recibido al invocar el procedimiento remoto.
 * @param reproductor Reproductor que consume los fragmentos.
 */
func RecibirAudio(stream pb.AudioService_AudioStreamClient, primerFragmento *pb.AudioChunk, reproductor *Reproductor) {
	fragmento := primerFragmento
	for {
		if _, err := reproductor.escritor.Write(fragmento.Data); err != nil {
			// El reproductor cerró la tubería (usuario salió o error de audio).
			return
		}

		siguiente, err := stream.Recv()
		if err == io.EOF {
			reproductor.escritor.Close()
			return
		}
		if err != nil {
			if stream.Context().Err() == nil {
				reproductor.registrarError(fmt.Errorf("se interrumpió la recepción del audio"))
			}
			reproductor.escritor.CloseWithError(err)
			return
		}
		fragmento = siguiente
	}
}
