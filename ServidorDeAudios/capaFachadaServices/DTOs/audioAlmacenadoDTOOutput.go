/**
 * @file audioAlmacenadoDTOOutput.go
 * @brief DTOs de respuesta del servidor de audios (formato JSON).
 */
package dtos

/**
 * @brief Archivo mp3 almacenado en el servidor.
 */
type ArchivoAudioDTOOutput struct {
	NombreArchivo string `json:"nombreArchivo"` ///< Nombre del archivo mp3.
	TamanioBytes  int64  `json:"tamanioBytes"`  ///< Tamaño del archivo en bytes.
}

/**
 * @brief Respuesta de la operación de almacenamiento de un audio.
 */
type AudioAlmacenadoDTOOutput struct {
	Mensaje string                `json:"mensaje"` ///< Mensaje de confirmación.
	IdAudio int                   `json:"idAudio"` ///< Identificador asignado por el servidor de metadatos.
	Titulo  string                `json:"titulo"`  ///< Título del audio almacenado.
	Tipo    string                `json:"tipo"`    ///< Nombre del tipo de audio.
	Archivo ArchivoAudioDTOOutput `json:"archivo"` ///< Datos del archivo guardado.
}

/**
 * @brief Mensaje de error devuelto por los servicios REST.
 */
type RespuestaErrorDTOOutput struct {
	Codigo  int    `json:"codigo"`  ///< Código de estado HTTP.
	Mensaje string `json:"mensaje"` ///< Descripción del error.
}

/**
 * @brief Metadatos del audio que se envían al servidor de metadatos (POST /audios).
 */
type AudioRegistrarDTOOutput struct {
	IdTipo        int               `json:"idTipo"`        ///< Identificador del tipo de audio.
	Titulo        string            `json:"titulo"`        ///< Título del audio.
	NombreArchivo string            `json:"nombreArchivo"` ///< Archivo mp3 almacenado.
	Metadatos     map[string]string `json:"metadatos"`     ///< Metadatos propios del tipo de audio.
}

/**
 * @brief Respuesta del servidor de metadatos al registrar un audio.
 */
type AudioRegistradoDTOInput struct {
	Id         int    `json:"id"`         ///< Identificador asignado al audio.
	NombreTipo string `json:"nombreTipo"` ///< Nombre del tipo de audio.
	Titulo     string `json:"titulo"`     ///< Título del audio.
}
