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
	Titulo  string                `json:"titulo"`  ///< Título del audio almacenado.
	Tipo    string                `json:"tipo"`    ///< Tipo del audio almacenado.
	Archivo ArchivoAudioDTOOutput `json:"archivo"` ///< Datos del archivo guardado.
}

/**
 * @brief Mensaje de error devuelto por los servicios REST.
 */
type RespuestaErrorDTOOutput struct {
	Codigo  int    `json:"codigo"`  ///< Código de estado HTTP.
	Mensaje string `json:"mensaje"` ///< Descripción del error.
}
