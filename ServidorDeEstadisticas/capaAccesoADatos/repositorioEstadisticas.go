/**
 * @file repositorioEstadisticas.go
 * @brief Repositorio en memoria de las reproducciones recibidas.
 */
package capaaccesodatos

import (
	"sort"
	"sync"

	"estadisticas/modelos"
)

/**
 * @brief Repositorio que almacena las reproducciones y calcula los conteos.
 *
 * Se implementa como singleton.
 */
type RepositorioEstadisticas struct {
	mu             sync.Mutex             ///< Protege el acceso concurrente a los datos.
	reproducciones []modelos.Reproduccion ///< Reproducciones almacenadas.
}

var (
	instancia *RepositorioEstadisticas ///< Única instancia del repositorio.
	once      sync.Once                ///< Garantiza que la instancia se cree una sola vez.
)

/**
 * @brief Obtiene la única instancia del repositorio (patrón singleton).
 * @return Puntero al repositorio de estadísticas.
 */
func GetRepositorioEstadisticas() *RepositorioEstadisticas {
	once.Do(func() {
		instancia = &RepositorioEstadisticas{}
	})
	return instancia
}

/**
 * @brief Almacena una nueva reproducción asignándole un número consecutivo.
 * @param reproduccion Reproducción a almacenar.
 * @return La reproducción almacenada con su número asignado.
 */
func (thisR *RepositorioEstadisticas) GuardarReproduccion(reproduccion modelos.Reproduccion) modelos.Reproduccion {
	thisR.mu.Lock()
	defer thisR.mu.Unlock()

	reproduccion.Numero = len(thisR.reproducciones) + 1
	thisR.reproducciones = append(thisR.reproducciones, reproduccion)
	return reproduccion
}

/**
 * @brief Calcula el resumen de estadísticas con los datos almacenados.
 * @return Total y conteos por tipo, por audio y por usuario, ordenados de mayor a menor.
 */
func (thisR *RepositorioEstadisticas) CalcularResumen() modelos.ResumenEstadisticas {
	thisR.mu.Lock()
	defer thisR.mu.Unlock()

	porTipo := map[string]int{}
	porAudio := map[string]int{}
	porUsuario := map[string]int{}
	for _, reproduccion := range thisR.reproducciones {
		porTipo[reproduccion.Tipo]++
		porAudio[reproduccion.Titulo]++
		porUsuario[reproduccion.Usuario]++
	}

	return modelos.ResumenEstadisticas{
		TotalReproducciones: len(thisR.reproducciones),
		PorTipo:             ordenarConteos(porTipo),
		PorAudio:            ordenarConteos(porAudio),
		PorUsuario:          ordenarConteos(porUsuario),
	}
}

/**
 * @brief Convierte un mapa de conteos en una lista ordenada de mayor a menor.
 * @param conteos Mapa nombre -> cantidad.
 * @return Lista de conteos ordenada por cantidad (y por nombre en caso de empate).
 */
func ordenarConteos(conteos map[string]int) []modelos.Conteo {
	lista := make([]modelos.Conteo, 0, len(conteos))
	for nombre, cantidad := range conteos {
		lista = append(lista, modelos.Conteo{Nombre: nombre, Cantidad: cantidad})
	}
	sort.Slice(lista, func(i, j int) bool {
		if lista[i].Cantidad != lista[j].Cantidad {
			return lista[i].Cantidad > lista[j].Cantidad
		}
		return lista[i].Nombre < lista[j].Nombre
	})
	return lista
}
