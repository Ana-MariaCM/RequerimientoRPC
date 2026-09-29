/**
 * @file repositorioAudios.go
 * @brief Repositorio en memoria de los metadatos de los audios.
 */
package capaaccesoadatos

import (
	"sync"

	"metadatos/modelos"
)

/**
 * @brief Repositorio que almacena en memoria los audios y sus metadatos.
 *
 * Se implementa como singleton para que todas las capas compartan la misma
 * instancia de los datos.
 */
type RepositorioAudios struct {
	mu     sync.RWMutex    ///< Protege el acceso concurrente a la lista de audios.
	audios []modelos.Audio ///< Audios almacenados (música, podcasts, audiolibros y ruido blanco).
}

var (
	instanciaRepositorioAudios *RepositorioAudios ///< Única instancia del repositorio de audios.
	onceRepositorioAudios      sync.Once          ///< Garantiza que la instancia se cree una sola vez.
)

/**
 * @brief Obtiene la única instancia del repositorio de audios (patrón singleton).
 *
 * En la primera invocación se cargan los audios precargados.
 * @return Puntero al repositorio de audios.
 */
func GetRepositorioAudios() *RepositorioAudios {
	onceRepositorioAudios.Do(func() {
		instanciaRepositorioAudios = &RepositorioAudios{audios: cargarAudios()}
	})
	return instanciaRepositorioAudios
}

/**
 * @brief Lista los audios que pertenecen a un tipo.
 * @param idTipo Identificador del tipo de audio.
 * @return Audios del tipo indicado (lista vacía si no hay).
 */
func (thisR *RepositorioAudios) ListarAudiosPorTipo(idTipo int) []modelos.Audio {
	thisR.mu.RLock()
	defer thisR.mu.RUnlock()

	resultado := []modelos.Audio{}
	for _, audio := range thisR.audios {
		if audio.ObtenerIdTipo() == idTipo {
			resultado = append(resultado, audio)
		}
	}
	return resultado
}

/**
 * @brief Busca un audio por su identificador.
 * @param idAudio Identificador del audio.
 * @return El audio encontrado y true, o nil y false si no existe.
 */
func (thisR *RepositorioAudios) BuscarAudioPorId(idAudio int) (modelos.Audio, bool) {
	thisR.mu.RLock()
	defer thisR.mu.RUnlock()

	for _, audio := range thisR.audios {
		if audio.ObtenerId() == idAudio {
			return audio, true
		}
	}
	return nil, false
}
