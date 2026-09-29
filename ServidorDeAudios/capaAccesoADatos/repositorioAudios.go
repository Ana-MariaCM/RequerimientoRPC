/**
 * @file repositorioAudios.go
 * @brief Repositorio que almacena y lista los archivos mp3 en disco.
 */
package capaaccesodatos

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"almacenamiento/configuracion"
	"almacenamiento/modelos"
)

/**
 * @brief Repositorio de archivos de audio (patrón singleton).
 */
type RepositorioAudios struct {
	mu         sync.Mutex ///< Serializa las escrituras en disco.
	rutaAudios string     ///< Carpeta donde se guardan los mp3.
}

var (
	instancia *RepositorioAudios ///< Única instancia del repositorio.
	once      sync.Once          ///< Garantiza que la instancia se cree una sola vez.
)

/**
 * @brief Obtiene la única instancia del repositorio de audios.
 * @return Puntero al repositorio.
 */
func GetRepositorioAudios() *RepositorioAudios {
	once.Do(func() {
		instancia = &RepositorioAudios{rutaAudios: configuracion.ObtenerRutaAudios()}
	})
	return instancia
}

/**
 * @brief Guarda un archivo mp3 en la carpeta de audios.
 *
 * Si ya existe un archivo con el mismo nombre se reemplaza.
 * @param nombreArchivo Nombre con el que se almacenará el archivo.
 * @param datos Contenido del archivo mp3.
 * @return El archivo almacenado o un error si no se pudo escribir en disco.
 */
func (thisR *RepositorioAudios) GuardarAudio(nombreArchivo string, datos []byte) (modelos.ArchivoAudio, error) {
	thisR.mu.Lock()
	defer thisR.mu.Unlock()

	if err := os.MkdirAll(thisR.rutaAudios, os.ModePerm); err != nil {
		return modelos.ArchivoAudio{}, fmt.Errorf("error creando la carpeta de audios: %v", err)
	}

	ruta := filepath.Join(thisR.rutaAudios, nombreArchivo)
	if err := os.WriteFile(ruta, datos, 0644); err != nil {
		return modelos.ArchivoAudio{}, fmt.Errorf("error al guardar archivo: %v", err)
	}

	fmt.Println("Audio guardado en:", ruta)
	return modelos.ArchivoAudio{NombreArchivo: nombreArchivo, TamanioBytes: int64(len(datos))}, nil
}

/**
 * @brief Lista los archivos mp3 almacenados.
 * @return Archivos mp3 ordenados por nombre, o un error si no se pudo leer la carpeta.
 */
func (thisR *RepositorioAudios) ListarAudios() ([]modelos.ArchivoAudio, error) {
	entradas, err := os.ReadDir(thisR.rutaAudios)
	if err != nil {
		if os.IsNotExist(err) {
			return []modelos.ArchivoAudio{}, nil
		}
		return nil, fmt.Errorf("error leyendo la carpeta de audios: %v", err)
	}

	archivos := []modelos.ArchivoAudio{}
	for _, entrada := range entradas {
		if entrada.IsDir() || !strings.HasSuffix(strings.ToLower(entrada.Name()), ".mp3") {
			continue
		}
		informacion, err := entrada.Info()
		if err != nil {
			continue
		}
		archivos = append(archivos, modelos.ArchivoAudio{NombreArchivo: entrada.Name(), TamanioBytes: informacion.Size()})
	}

	sort.Slice(archivos, func(i, j int) bool { return archivos[i].NombreArchivo < archivos[j].NombreArchivo })
	return archivos, nil
}
