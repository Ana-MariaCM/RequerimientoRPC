/**
 * @file audioDetalleDTOOutput.go
 * @brief DTOs con el detalle completo de un audio y sus metadatos.
 */
package dtos

/**
 * @brief Metadato (etiqueta y valor) enviado en formato JSON.
 */
type MetadatoDTOOutput struct {
	Nombre string `json:"nombre"` ///< Etiqueta del metadato.
	Valor  string `json:"valor"`  ///< Valor del metadato.
}

/**
 * @brief Detalle de un audio enviado al cliente (GET /audios/{idAudio}).
 */
type AudioDetalleDTOOutput struct {
	Id            int                 `json:"id"`            ///< Identificador del audio.
	IdTipo        int                 `json:"idTipo"`        ///< Identificador del tipo de audio.
	NombreTipo    string              `json:"nombreTipo"`    ///< Nombre del tipo de audio.
	Titulo        string              `json:"titulo"`        ///< Título del audio.
	NombreArchivo string              `json:"nombreArchivo"` ///< Archivo mp3 que se solicitará al servidor de streaming.
	Metadatos     []MetadatoDTOOutput `json:"metadatos"`     ///< Metadatos propios del tipo de audio, en orden de presentación.
}
