/**
 * @file clienteMetadatos.go
 * @brief Componente que registra, mediante REST, los metadatos de los audios nuevos
 *        en el servidor de metadatos.
 */
package componenteclientemetadatos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	dtos "almacenamiento/capaFachadaServices/DTOs"
)

/**
 * @brief Cliente REST del servidor de metadatos.
 */
type ClienteMetadatos struct {
	urlBase    string       ///< URL base del servidor de metadatos.
	clienteWeb *http.Client ///< Cliente HTTP reutilizable.
}

/**
 * @brief Crea el cliente del servidor de metadatos.
 * @param urlBase URL base (por ejemplo http://localhost:5001).
 * @return Puntero al nuevo cliente.
 */
func NuevoClienteMetadatos(urlBase string) *ClienteMetadatos {
	return &ClienteMetadatos{urlBase: urlBase, clienteWeb: &http.Client{Timeout: 10 * time.Second}}
}

/**
 * @brief Registra los metadatos de un audio (POST /audios del servidor de metadatos).
 * @param audio Metadatos del audio.
 * @return Datos del audio registrado (incluye el id asignado) o un error.
 */
func (thisC *ClienteMetadatos) RegistrarAudio(audio dtos.AudioRegistrarDTOOutput) (dtos.AudioRegistradoDTOInput, error) {
	var registrado dtos.AudioRegistradoDTOInput

	cuerpo, err := json.Marshal(audio)
	if err != nil {
		return registrado, err
	}

	fmt.Printf("[REST] Invocando POST %s/audios para registrar \"%s\"\n", thisC.urlBase, audio.Titulo)
	respuesta, err := thisC.clienteWeb.Post(thisC.urlBase+"/audios", "application/json", bytes.NewReader(cuerpo))
	if err != nil {
		return registrado, fmt.Errorf("no fue posible comunicarse con el servidor de metadatos (%s)", thisC.urlBase)
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusCreated && respuesta.StatusCode != http.StatusOK {
		var errorDTO dtos.RespuestaErrorDTOOutput
		if json.NewDecoder(respuesta.Body).Decode(&errorDTO) == nil && errorDTO.Mensaje != "" {
			return registrado, fmt.Errorf("%s", errorDTO.Mensaje)
		}
		return registrado, fmt.Errorf("el servidor de metadatos respondió con el código %d", respuesta.StatusCode)
	}

	if err := json.NewDecoder(respuesta.Body).Decode(&registrado); err != nil {
		return registrado, fmt.Errorf("respuesta inválida del servidor de metadatos: %v", err)
	}
	return registrado, nil
}
