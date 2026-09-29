/**
 * @file respuestaErrorDTOOutput.go
 * @brief DTO para informar errores al cliente en formato JSON.
 */
package dtos

/**
 * @brief Mensaje de error devuelto por los servicios REST.
 */
type RespuestaErrorDTOOutput struct {
	Codigo  int    `json:"codigo"`  ///< Código de estado HTTP.
	Mensaje string `json:"mensaje"` ///< Descripción del error.
}
