/**
 * @file audioRegistrado.go
 * @brief Modelo con los datos de un audio registrado por el administrador,
 *        tal como se guarda en el archivo JSON.
 */
package modelos

/**
 * @brief Audio registrado por el administrador (registro persistente).
 *
 * Se guarda en datos/audiosRegistrados.json para que los audios nuevos no se
 * pierdan al reiniciar el servidor de metadatos. Al iniciar, cada registro se
 * convierte de nuevo en su modelo concreto (Musica, Podcast, Audiolibro o RuidoBlanco).
 */
type AudioRegistrado struct {
	Id            int               `json:"id"`            ///< Identificador del audio.
	IdTipo        int               `json:"idTipo"`        ///< Identificador del tipo de audio.
	Titulo        string            `json:"titulo"`        ///< Título del audio.
	NombreArchivo string            `json:"nombreArchivo"` ///< Archivo mp3 en el servidor de audios.
	Metadatos     map[string]string `json:"metadatos"`     ///< Metadatos propios del tipo (clave -> valor).
}
