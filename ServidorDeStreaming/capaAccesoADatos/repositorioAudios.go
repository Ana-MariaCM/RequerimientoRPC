/**
 * @file repositorioAudios.go
 * @brief Acceso a los archivos mp3 que se transmiten por streaming.
 */
package capaaccesodatos

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"streaming.local/servidor-streaming/configuracion"
)

/**
 * @brief Repositorio que abre los archivos mp3 almacenados por el servidor de audios.
 *
 * Se implementa como singleton.
 */
type RepositorioAudios struct {
	rutaAudios string ///< Carpeta que contiene los archivos mp3.
}

var (
	instancia *RepositorioAudios ///< Única instancia del repositorio.
	once      sync.Once          ///< Garantiza que la instancia se cree una sola vez.
)

/**
 * @brief Obtiene la única instancia del repositorio de audios (patrón singleton).
 * @return Puntero al repositorio.
 */
func GetRepositorioAudios() *RepositorioAudios {
	once.Do(func() {
		instancia = &RepositorioAudios{rutaAudios: configuracion.ObtenerRutaAudios()}
	})
	return instancia
}

/**
 * @brief Devuelve la carpeta desde la que se leen los audios.
 * @return Ruta de la carpeta de audios.
 */
func (thisR *RepositorioAudios) ObtenerRutaAudios() string {
	return thisR.rutaAudios
}

/**
 * @brief Abre un archivo de audio a partir de su nombre.
 *
 * Solo se usa el nombre base del archivo para impedir el acceso a rutas
 * fuera de la carpeta de audios.
 * @param nombreArchivo Nombre del archivo mp3 (por ejemplo "musica_take_five.mp3").
 * @return El archivo abierto o un error si no existe.
 */
func (thisR *RepositorioAudios) AbrirArchivo(nombreArchivo string) (*os.File, error) {
	ruta := filepath.Join(thisR.rutaAudios, filepath.Base(nombreArchivo))

	archivo, err := os.Open(ruta)
	if err != nil {
		fmt.Println("Error abriendo el audio:", ruta)
		return nil, fmt.Errorf("no se pudo abrir el archivo %s: %w", nombreArchivo, err)
	}
	fmt.Println("Audio abierto:", ruta)
	return archivo, nil
}
