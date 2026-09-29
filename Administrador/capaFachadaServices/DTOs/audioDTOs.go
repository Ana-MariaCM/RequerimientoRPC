/**
 * @file audioDTOs.go
 * @brief DTOs intercambiados con el servidor de audios en formato JSON.
 */
package dtos

/**
 * @brief Datos del audio que el administrador desea almacenar.
 */
type AudioAlmacenarDTO struct {
	RutaLocal string            ///< Ruta del archivo mp3 en el equipo del administrador.
	IdTipo    int               ///< Identificador del tipo de audio.
	Titulo    string            ///< Título del audio que verá el cliente.
	Metadatos map[string]string ///< Metadatos propios del tipo de audio (clave -> valor).
}

/**
 * @brief Archivo mp3 almacenado en el servidor.
 */
type ArchivoAudioDTO struct {
	NombreArchivo string `json:"nombreArchivo"` ///< Nombre del archivo mp3.
	TamanioBytes  int64  `json:"tamanioBytes"`  ///< Tamaño del archivo en bytes.
}

/**
 * @brief Respuesta del servidor al almacenar un audio.
 */
type AudioAlmacenadoDTO struct {
	Mensaje string          `json:"mensaje"` ///< Mensaje de confirmación.
	IdAudio int             `json:"idAudio"` ///< Identificador asignado por el servidor de metadatos.
	Titulo  string          `json:"titulo"`  ///< Título del audio almacenado.
	Tipo    string          `json:"tipo"`    ///< Nombre del tipo de audio.
	Archivo ArchivoAudioDTO `json:"archivo"` ///< Datos del archivo guardado.
}

/**
 * @brief Mensaje de error devuelto por el servidor.
 */
type RespuestaErrorDTO struct {
	Codigo  int    `json:"codigo"`  ///< Código de estado HTTP.
	Mensaje string `json:"mensaje"` ///< Descripción del error.
}
