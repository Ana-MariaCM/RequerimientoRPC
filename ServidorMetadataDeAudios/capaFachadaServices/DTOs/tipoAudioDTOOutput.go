/**
 * @file tipoAudioDTOOutput.go
 * @brief DTO que encapsula un tipo de audio para enviarlo en formato JSON.
 */
package dtos

/**
 * @brief Tipo de audio enviado al cliente (GET /tipos).
 */
type TipoAudioDTOOutput struct {
	Id     int    `json:"id"`     ///< Identificador del tipo de audio.
	Nombre string `json:"nombre"` ///< Nombre del tipo de audio.
}
