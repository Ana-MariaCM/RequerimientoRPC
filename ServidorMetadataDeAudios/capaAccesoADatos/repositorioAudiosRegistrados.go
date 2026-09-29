/**
 * @file repositorioAudiosRegistrados.go
 * @brief Repositorio que guarda en un archivo JSON los audios registrados por el administrador.
 */
package capaaccesoadatos

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"metadatos/configuracion"
	"metadatos/modelos"
)

/**
 * @brief Repositorio persistente (archivo JSON) de los audios registrados.
 *
 * Los audios precargados siguen en datosPrecargados.go; en este archivo solo se
 * guardan los que agrega el administrador. Se implementa como singleton.
 */
type RepositorioAudiosRegistrados struct {
	mu          sync.Mutex ///< Serializa las lecturas y escrituras del archivo.
	rutaArchivo string     ///< Ruta del archivo JSON.
}

var (
	instanciaRepositorioRegistrados *RepositorioAudiosRegistrados ///< Única instancia del repositorio.
	onceRepositorioRegistrados      sync.Once                     ///< Garantiza que la instancia se cree una sola vez.
)

/**
 * @brief Obtiene la única instancia del repositorio de audios registrados (patrón singleton).
 * @return Puntero al repositorio.
 */
func GetRepositorioAudiosRegistrados() *RepositorioAudiosRegistrados {
	onceRepositorioRegistrados.Do(func() {
		instanciaRepositorioRegistrados = &RepositorioAudiosRegistrados{rutaArchivo: configuracion.ObtenerRutaAudiosRegistrados()}
	})
	return instanciaRepositorioRegistrados
}

/**
 * @brief Devuelve la ruta del archivo JSON.
 * @return Ruta del archivo.
 */
func (thisR *RepositorioAudiosRegistrados) ObtenerRutaArchivo() string {
	return thisR.rutaArchivo
}

/**
 * @brief Lee todos los audios registrados del archivo JSON.
 * @return Lista de registros (vacía si el archivo aún no existe) o un error si el archivo está dañado.
 */
func (thisR *RepositorioAudiosRegistrados) Cargar() ([]modelos.AudioRegistrado, error) {
	thisR.mu.Lock()
	defer thisR.mu.Unlock()
	return thisR.leerArchivo()
}

/**
 * @brief Guarda un audio registrado en el archivo JSON.
 *
 * Si ya existe un registro con el mismo archivo mp3 se reemplaza. El archivo
 * se escribe primero en un temporal y luego se renombra, para no dañarlo si
 * el proceso se interrumpe a mitad de la escritura.
 * @param registro Audio a guardar.
 * @return Error si no se pudo escribir el archivo.
 */
func (thisR *RepositorioAudiosRegistrados) Guardar(registro modelos.AudioRegistrado) error {
	thisR.mu.Lock()
	defer thisR.mu.Unlock()

	registros, err := thisR.leerArchivo()
	if err != nil {
		return err
	}

	reemplazado := false
	for indice, existente := range registros {
		if existente.NombreArchivo == registro.NombreArchivo {
			registros[indice] = registro
			reemplazado = true
		}
	}
	if !reemplazado {
		registros = append(registros, registro)
	}
	return thisR.escribirArchivo(registros)
}

/**
 * @brief Reemplaza el contenido del archivo JSON con la lista indicada.
 * @param registros Registros a guardar.
 * @return Error si no se pudo escribir el archivo.
 */
func (thisR *RepositorioAudiosRegistrados) GuardarTodos(registros []modelos.AudioRegistrado) error {
	thisR.mu.Lock()
	defer thisR.mu.Unlock()
	return thisR.escribirArchivo(registros)
}

/**
 * @brief Lee y deserializa el archivo JSON (sin bloquear el mutex).
 * @return Registros leídos o un error.
 */
func (thisR *RepositorioAudiosRegistrados) leerArchivo() ([]modelos.AudioRegistrado, error) {
	contenido, err := os.ReadFile(thisR.rutaArchivo)
	if errors.Is(err, os.ErrNotExist) {
		return []modelos.AudioRegistrado{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer %s: %v", thisR.rutaArchivo, err)
	}

	registros := []modelos.AudioRegistrado{}
	if len(contenido) == 0 {
		return registros, nil
	}
	if err := json.Unmarshal(contenido, &registros); err != nil {
		return nil, fmt.Errorf("el archivo %s no tiene un JSON válido: %v", thisR.rutaArchivo, err)
	}
	return registros, nil
}

/**
 * @brief Serializa y escribe el archivo JSON de forma segura (temporal + renombrar).
 * @param registros Registros a escribir.
 * @return Error si no se pudo escribir el archivo.
 */
func (thisR *RepositorioAudiosRegistrados) escribirArchivo(registros []modelos.AudioRegistrado) error {
	if err := os.MkdirAll(filepath.Dir(thisR.rutaArchivo), os.ModePerm); err != nil {
		return fmt.Errorf("no se pudo crear la carpeta de datos: %v", err)
	}

	contenido, err := json.MarshalIndent(registros, "", "  ")
	if err != nil {
		return err
	}

	temporal := thisR.rutaArchivo + ".tmp"
	if err := os.WriteFile(temporal, contenido, 0644); err != nil {
		return fmt.Errorf("no se pudo escribir %s: %v", temporal, err)
	}
	if err := os.Rename(temporal, thisR.rutaArchivo); err != nil {
		return fmt.Errorf("no se pudo guardar %s: %v", thisR.rutaArchivo, err)
	}
	fmt.Printf("[JSON] %d audio(s) registrado(s) guardado(s) en %s\n", len(registros), thisR.rutaArchivo)
	return nil
}
