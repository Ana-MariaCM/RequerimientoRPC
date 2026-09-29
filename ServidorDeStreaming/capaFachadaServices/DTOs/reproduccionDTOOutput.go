/**
 * @file reproduccionDTOOutput.go
 * @brief DTO que encapsula la información de una reproducción que se publica en la cola.
 */
package dtos

/**
 * @brief Información de una reproducción enviada en formato JSON al servidor de estadísticas.
 */
type ReproduccionDTOOutput struct {
	IdAudio          int32  `json:"idAudio"`          ///< Identificador del audio reproducido.
	Titulo           string `json:"titulo"`           ///< Título del audio reproducido.
	Tipo             string `json:"tipo"`             ///< Tipo de audio (Música, Podcasts, ...).
	NombreArchivo    string `json:"nombreArchivo"`    ///< Archivo mp3 transmitido.
	Usuario          string `json:"usuario"`          ///< Usuario que solicitó la reproducción.
	DireccionCliente string `json:"direccionCliente"` ///< Dirección IP:puerto del cliente.
	FechaHora        string `json:"fechaHora"`        ///< Fecha y hora de inicio de la reproducción.
}
