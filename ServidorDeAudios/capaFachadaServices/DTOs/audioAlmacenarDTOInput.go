/**
 * @file audioAlmacenarDTOInput.go
 * @brief DTO con los datos que envía el administrador para almacenar un audio.
 */
package dtos

/**
 * @brief Datos del audio a almacenar (campos del formulario multipart).
 */
type AudioAlmacenarDTOInput struct {
	Titulo        string `json:"titulo"`        ///< Título del audio.
	Tipo          string `json:"tipo"`          ///< Tipo de audio (Música, Podcasts, Audiolibros, Ruido Blanco).
	NombreArchivo string `json:"nombreArchivo"` ///< Nombre con el que se guardará el mp3.
}
