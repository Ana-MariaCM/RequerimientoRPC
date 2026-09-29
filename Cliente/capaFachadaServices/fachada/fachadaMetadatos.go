/**
 * @file fachadaMetadatos.go
 * @brief Fachada que consume los servicios REST del servidor de metadatos.
 */
package fachada

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	dtos "cliente.local/cliente/capaFachadaServices/DTOs"
)

/**
 * @brief Fachada de acceso remoto al servidor de metadatos.
 */
type FachadaMetadatos struct {
	urlBase    string       ///< URL base del servidor de metadatos.
	clienteWeb *http.Client ///< Cliente HTTP reutilizable.
}

/**
 * @brief Crea la fachada de metadatos.
 * @param urlBase URL base del servidor (por ejemplo http://localhost:5001).
 * @return Puntero a la nueva fachada.
 */
func NuevaFachadaMetadatos(urlBase string) *FachadaMetadatos {
	return &FachadaMetadatos{urlBase: urlBase, clienteWeb: &http.Client{Timeout: 10 * time.Second}}
}

/**
 * @brief Consulta los tipos de audio (GET /tipos).
 * @return Lista de tipos de audio o un error.
 */
func (thisF *FachadaMetadatos) ListarTiposAudio() ([]dtos.TipoAudioDTO, error) {
	var tipos []dtos.TipoAudioDTO
	err := thisF.consultar("/tipos", &tipos)
	return tipos, err
}

/**
 * @brief Consulta los audios de un tipo (GET /tipos/{idTipo}/audios).
 * @param idTipo Identificador del tipo de audio.
 * @return Lista de audios del tipo o un error.
 */
func (thisF *FachadaMetadatos) ListarAudiosPorTipo(idTipo int) ([]dtos.AudioResumenDTO, error) {
	var audios []dtos.AudioResumenDTO
	err := thisF.consultar(fmt.Sprintf("/tipos/%d/audios", idTipo), &audios)
	return audios, err
}

/**
 * @brief Consulta el detalle de un audio (GET /audios/{idAudio}).
 * @param idAudio Identificador del audio.
 * @return Detalle del audio o un error.
 */
func (thisF *FachadaMetadatos) ConsultarDetalleAudio(idAudio int) (dtos.AudioDetalleDTO, error) {
	var detalle dtos.AudioDetalleDTO
	err := thisF.consultar(fmt.Sprintf("/audios/%d", idAudio), &detalle)
	return detalle, err
}

/**
 * @brief Realiza una petición GET y deserializa la respuesta JSON.
 * @param ruta Ruta del servicio REST (por ejemplo "/tipos").
 * @param destino Puntero donde se deserializa la respuesta.
 * @return Error de comunicación, el mensaje de error del servidor, o nil.
 */
func (thisF *FachadaMetadatos) consultar(ruta string, destino any) error {
	respuesta, err := thisF.clienteWeb.Get(thisF.urlBase + ruta)
	if err != nil {
		return fmt.Errorf("no fue posible comunicarse con el servidor de metadatos (%s)", thisF.urlBase)
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		var errorDTO dtos.RespuestaErrorDTO
		if json.NewDecoder(respuesta.Body).Decode(&errorDTO) == nil && errorDTO.Mensaje != "" {
			return fmt.Errorf("%s", errorDTO.Mensaje)
		}
		return fmt.Errorf("el servidor de metadatos respondió con el código %d", respuesta.StatusCode)
	}

	if err := json.NewDecoder(respuesta.Body).Decode(destino); err != nil {
		return fmt.Errorf("respuesta inválida del servidor de metadatos: %v", err)
	}
	return nil
}
