/**
 * @file audioAlmacenarDTOInput.go
 * @brief DTO con los datos que envía el administrador para almacenar un audio.
 */
package dtos

/**
 * @brief Datos del audio a almacenar (campos del formulario multipart).
 */
type AudioAlmacenarDTOInput struct {
	IdTipo    int               `json:"idTipo"`    ///< Identificador del tipo de audio.
	Titulo    string            `json:"titulo"`    ///< Título del audio que verá el cliente.
	Metadatos map[string]string `json:"metadatos"` ///< Metadatos propios del tipo de audio (clave -> valor).
}
