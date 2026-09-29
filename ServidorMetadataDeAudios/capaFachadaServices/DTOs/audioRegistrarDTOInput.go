/**
 * @file audioRegistrarDTOInput.go
 * @brief DTO con los metadatos de un audio nuevo que envía el servidor de audios (JSON).
 */
package dtos

/**
 * @brief Datos para registrar (o actualizar) los metadatos de un audio (POST /audios).
 *
 * Las claves de Metadatos dependen del tipo de audio:
 *  - Música: artistaPrincipal, album, generoMusical, selloDiscografico, anioLanzamiento, duracion.
 *  - Podcasts: nombrePodcast, anfitrion, temporada, episodio, notasDelShow, clasificacionContenido, duracion.
 *  - Audiolibros: autor, narrador, editorial, isbn, capitulo, totalCapitulos, duracion, genero.
 *  - Ruido Blanco: tipoSonido, fuenteAudio, usoSugerido, proveedorContenido, duracionBucle, frecuenciaDominante.
 */
type AudioRegistrarDTOInput struct {
	IdTipo        int               `json:"idTipo"`        ///< Identificador del tipo de audio.
	Titulo        string            `json:"titulo"`        ///< Título del audio.
	NombreArchivo string            `json:"nombreArchivo"` ///< Archivo mp3 almacenado en el servidor de audios.
	Metadatos     map[string]string `json:"metadatos"`     ///< Metadatos propios del tipo de audio.
}
