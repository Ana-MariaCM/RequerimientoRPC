/**
 * @file controladorAudios.go
 * @brief Controlador que atiende las consultas de tipos, listas y detalles de audios.
 */
package capacontroladores

import (
	dtos "cliente.local/cliente/capaFachadaServices/DTOs"
	"cliente.local/cliente/capaFachadaServices/fachada"
)

/**
 * @brief Controlador de consultas de metadatos de audios.
 */
type ControladorAudios struct {
	fachada *fachada.FachadaMetadatos ///< Fachada de acceso al servidor de metadatos.
}

/**
 * @brief Crea el controlador de audios.
 * @param fachadaMetadatos Fachada de acceso al servidor de metadatos.
 * @return Puntero al nuevo controlador.
 */
func NuevoControladorAudios(fachadaMetadatos *fachada.FachadaMetadatos) *ControladorAudios {
	return &ControladorAudios{fachada: fachadaMetadatos}
}

/**
 * @brief Obtiene los tipos de audio registrados.
 * @return Lista de tipos de audio o un error.
 */
func (thisC *ControladorAudios) ObtenerTiposAudio() ([]dtos.TipoAudioDTO, error) {
	return thisC.fachada.ListarTiposAudio()
}

/**
 * @brief Obtiene la lista de audios de un tipo.
 * @param idTipo Identificador del tipo de audio.
 * @return Lista de audios o un error.
 */
func (thisC *ControladorAudios) ObtenerAudiosPorTipo(idTipo int) ([]dtos.AudioResumenDTO, error) {
	return thisC.fachada.ListarAudiosPorTipo(idTipo)
}

/**
 * @brief Obtiene el detalle (metadatos) de un audio.
 * @param idAudio Identificador del audio.
 * @return Detalle del audio o un error.
 */
func (thisC *ControladorAudios) ObtenerDetalleAudio(idAudio int) (dtos.AudioDetalleDTO, error) {
	return thisC.fachada.ConsultarDetalleAudio(idAudio)
}
