/**
 * @file reproduccionDTOInput.go
 * @brief DTO que encapsula la información de una reproducción recibida desde la cola (JSON).
 */
package dtos

/**
 * @brief Información de una reproducción publicada por el servidor de streaming.
 */
type ReproduccionDTOInput struct {
	IdAudio          int    `json:"idAudio"`          ///< Identificador del audio reproducido.
	Titulo           string `json:"titulo"`           ///< Título del audio reproducido.
	Tipo             string `json:"tipo"`             ///< Tipo de audio.
	NombreArchivo    string `json:"nombreArchivo"`    ///< Archivo mp3 transmitido.
	Usuario          string `json:"usuario"`          ///< Usuario que solicitó la reproducción.
	DireccionCliente string `json:"direccionCliente"` ///< Dirección IP:puerto del cliente.
	FechaHora        string `json:"fechaHora"`        ///< Fecha y hora de inicio de la reproducción.
}
