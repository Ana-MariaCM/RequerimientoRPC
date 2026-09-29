/**
 * @file repositorioTiposAudio.go
 * @brief Repositorio en memoria de los tipos de audio.
 */
package capaaccesoadatos

import (
	"sync"

	"metadatos/modelos"
)

/**
 * @brief Repositorio que almacena en memoria los tipos de audio.
 *
 * Se implementa como singleton para que todas las capas compartan la misma
 * instancia de los datos.
 */
type RepositorioTiposAudio struct {
	mu    sync.RWMutex        ///< Protege el acceso concurrente a la lista de tipos.
	tipos []modelos.TipoAudio ///< Tipos de audio almacenados.
}

var (
	instanciaRepositorioTipos *RepositorioTiposAudio ///< Única instancia del repositorio de tipos.
	onceRepositorioTipos      sync.Once              ///< Garantiza que la instancia se cree una sola vez.
)

/**
 * @brief Obtiene la única instancia del repositorio de tipos (patrón singleton).
 *
 * En la primera invocación se cargan los tipos precargados.
 * @return Puntero al repositorio de tipos de audio.
 */
func GetRepositorioTiposAudio() *RepositorioTiposAudio {
	onceRepositorioTipos.Do(func() {
		instanciaRepositorioTipos = &RepositorioTiposAudio{tipos: cargarTiposAudio()}
	})
	return instanciaRepositorioTipos
}

/**
 * @brief Lista todos los tipos de audio registrados.
 * @return Copia de la lista de tipos de audio.
 */
func (thisR *RepositorioTiposAudio) ListarTipos() []modelos.TipoAudio {
	thisR.mu.RLock()
	defer thisR.mu.RUnlock()

	copia := make([]modelos.TipoAudio, len(thisR.tipos))
	copy(copia, thisR.tipos)
	return copia
}

/**
 * @brief Busca un tipo de audio por su identificador.
 * @param idTipo Identificador del tipo buscado.
 * @return El tipo encontrado y true, o un tipo vacío y false si no existe.
 */
func (thisR *RepositorioTiposAudio) BuscarTipoPorId(idTipo int) (modelos.TipoAudio, bool) {
	thisR.mu.RLock()
	defer thisR.mu.RUnlock()

	for _, tipo := range thisR.tipos {
		if tipo.Id == idTipo {
			return tipo, true
		}
	}
	return modelos.TipoAudio{}, false
}
