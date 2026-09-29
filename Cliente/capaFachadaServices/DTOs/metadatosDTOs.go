/**
 * @file metadatosDTOs.go
 * @brief DTOs que encapsulan los datos recibidos del servidor de metadatos en formato JSON.
 */
package dtos

/**
 * @brief Tipo de audio (Música, Podcasts, Audiolibros, Ruido Blanco).
 */
type TipoAudioDTO struct {
	Id     int    `json:"id"`     ///< Identificador del tipo de audio.
	Nombre string `json:"nombre"` ///< Nombre del tipo de audio.
}

/**
 * @brief Resumen de un audio usado en la lista de audios de un tipo.
 */
type AudioResumenDTO struct {
	Id     int    `json:"id"`     ///< Identificador del audio.
	IdTipo int    `json:"idTipo"` ///< Identificador del tipo de audio.
	Titulo string `json:"titulo"` ///< Título del audio.
}

/**
 * @brief Metadato (etiqueta y valor) de un audio.
 */
type MetadatoDTO struct {
	Nombre string `json:"nombre"` ///< Etiqueta del metadato.
	Valor  string `json:"valor"`  ///< Valor del metadato.
}

/**
 * @brief Detalle completo de un audio.
 */
type AudioDetalleDTO struct {
	Id            int           `json:"id"`            ///< Identificador del audio.
	IdTipo        int           `json:"idTipo"`        ///< Identificador del tipo de audio.
	NombreTipo    string        `json:"nombreTipo"`    ///< Nombre del tipo de audio.
	Titulo        string        `json:"titulo"`        ///< Título del audio.
	NombreArchivo string        `json:"nombreArchivo"` ///< Archivo mp3 que se solicitará al servidor de streaming.
	Metadatos     []MetadatoDTO `json:"metadatos"`     ///< Metadatos del audio en orden de presentación.
}

/**
 * @brief Mensaje de error devuelto por el servidor de metadatos.
 */
type RespuestaErrorDTO struct {
	Codigo  int    `json:"codigo"`  ///< Código de estado HTTP.
	Mensaje string `json:"mensaje"` ///< Descripción del error.
}
