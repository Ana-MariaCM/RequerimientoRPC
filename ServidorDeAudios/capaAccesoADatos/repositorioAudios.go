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
 * @brief Guarda un archivo mp3 nuevo en la carpeta de audios sin reemplazar ninguno existente.
 *
 * Si ya existe un archivo con el nombre base, se agrega un sufijo numérico
 * (por ejemplo "musica_mi_cancion_2.mp3").
 * @param nombreBase Nombre deseado para el archivo (terminado en .mp3).
 * @param datos Contenido del archivo mp3.
 * @return El archivo almacenado (con el nombre definitivo) o un error si no se pudo escribir en disco.
 */
func (thisR *RepositorioAudios) GuardarAudio(nombreBase string, datos []byte) (modelos.ArchivoAudio, error) {
	thisR.mu.Lock()
	defer thisR.mu.Unlock()

	if err := os.MkdirAll(thisR.rutaAudios, os.ModePerm); err != nil {
		return modelos.ArchivoAudio{}, fmt.Errorf("error creando la carpeta de audios: %v", err)
	}

	nombreArchivo := nombreBase
	extension := filepath.Ext(nombreBase)
	for consecutivo := 2; ; consecutivo++ {
		if _, err := os.Stat(filepath.Join(thisR.rutaAudios, nombreArchivo)); os.IsNotExist(err) {
			break
		}
		nombreArchivo = fmt.Sprintf("%s_%d%s", strings.TrimSuffix(nombreBase, extension), consecutivo, extension)
	}

	ruta := filepath.Join(thisR.rutaAudios, nombreArchivo)
	if err := os.WriteFile(ruta, datos, 0644); err != nil {
		return modelos.ArchivoAudio{}, fmt.Errorf("error al guardar archivo: %v", err)
	}

	fmt.Println("Audio guardado en:", ruta)
	return modelos.ArchivoAudio{NombreArchivo: nombreArchivo, TamanioBytes: int64(len(datos))}, nil
}

/**
 * @brief Elimina un archivo mp3 de la carpeta de audios.
 *
 * Se usa para deshacer el almacenamiento cuando los metadatos no se pudieron registrar.
 * @param nombreArchivo Nombre del archivo a eliminar.
 */
func (thisR *RepositorioAudios) EliminarAudio(nombreArchivo string) {
	thisR.mu.Lock()
	defer thisR.mu.Unlock()

	ruta := filepath.Join(thisR.rutaAudios, filepath.Base(nombreArchivo))
	if err := os.Remove(ruta); err == nil {
		fmt.Println("Audio eliminado:", ruta)
	}
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
