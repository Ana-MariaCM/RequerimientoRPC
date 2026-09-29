/**
 * @file audioResumenDTOOutput.go
 * @brief DTO con los datos mínimos de un audio para construir la lista de audios de un tipo.
 */
package dtos

/**
 * @brief Resumen de un audio enviado al cliente (GET /tipos/{idTipo}/audios).
 */
type AudioResumenDTOOutput struct {
	Id     int    `json:"id"`     ///< Identificador del audio.
	IdTipo int    `json:"idTipo"` ///< Identificador del tipo de audio.
	Titulo string `json:"titulo"` ///< Título del audio.
}
