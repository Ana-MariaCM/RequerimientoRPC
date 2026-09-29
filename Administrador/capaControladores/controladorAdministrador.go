/**
 * @file controladorAdministrador.go
 * @brief Controlador que atiende las opciones del menú del administrador.
 */
package capacontroladores

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	dtos "administrador/capaFachadaServices/DTOs"
	"administrador/capaFachadaServices/fachada"
)

/**
 * @brief Controlador del administrador.
 */
type ControladorAdministrador struct {
	fachada *fachada.FachadaAudios ///< Fachada de acceso al servidor de audios.
}

/**
 * @brief Crea el controlador del administrador.
 * @param fachadaAudios Fachada de acceso al servidor de audios.
 * @return Puntero al nuevo controlador.
 */
func NuevoControladorAdministrador(fachadaAudios *fachada.FachadaAudios) *ControladorAdministrador {
	return &ControladorAdministrador{fachada: fachadaAudios}
}

/**
 * @brief Valida los datos ingresados y solicita el almacenamiento del audio.
 * @param audio Datos del audio a almacenar.
 * @return Respuesta del servidor o un error de validación o de comunicación.
 */
func (thisC *ControladorAdministrador) AlmacenarAudio(audio dtos.AudioAlmacenarDTO) (dtos.AudioAlmacenadoDTO, error) {
	rutaLocal, err := thisC.ValidarArchivoLocal(audio.RutaLocal)
	if err != nil {
		return dtos.AudioAlmacenadoDTO{}, err
	}
	audio.RutaLocal = rutaLocal

	return thisC.fachada.AlmacenarAudio(audio)
}

/**
 * @brief Verifica que la ruta indicada corresponda a un archivo .mp3 existente.
 * @param ruta Ruta escrita por el administrador.
 * @return Ruta normalizada o un error si el archivo no existe o no es .mp3.
 */
func (thisC *ControladorAdministrador) ValidarArchivoLocal(ruta string) (string, error) {
	ruta = normalizarRuta(ruta)

	informacion, err := os.Stat(ruta)
	if err != nil || informacion.IsDir() {
		return "", fmt.Errorf("el archivo %s no existe", ruta)
	}
	if !strings.EqualFold(filepath.Ext(ruta), ".mp3") {
		return "", fmt.Errorf("el archivo debe tener extensión .mp3")
	}
	return ruta, nil
}

/**
 * @brief Obtiene la lista de audios almacenados en el servidor.
 * @return Lista de archivos o un error.
 */
func (thisC *ControladorAdministrador) ListarAudios() ([]dtos.ArchivoAudioDTO, error) {
	return thisC.fachada.ListarAudios()
}

/**
 * @brief Limpia la ruta ingresada (comillas y ~ del directorio personal).
 * @param ruta Ruta escrita por el usuario.
 * @return Ruta normalizada.
 */
func normalizarRuta(ruta string) string {
	ruta = strings.Trim(strings.TrimSpace(ruta), "\"'")
	if strings.HasPrefix(ruta, "~/") {
		if directorioPersonal, err := os.UserHomeDir(); err == nil {
			ruta = filepath.Join(directorioPersonal, ruta[2:])
		}
	}
	return ruta
}
