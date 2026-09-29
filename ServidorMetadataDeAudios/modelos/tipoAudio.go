/**
 * @file tipoAudio.go
 * @brief Modelo que representa un tipo de audio (Música, Podcasts, Audiolibros, Ruido Blanco).
 */
package modelos

/**
 * @brief Tipo de audio almacenado en el servidor de metadatos.
 *
 * Según el requerimiento, por cada tipo se almacena un identificador entero
 * y un nombre de tipo string.
 */
type TipoAudio struct {
	Id     int    ///< Identificador único del tipo de audio.
	Nombre string ///< Nombre del tipo de audio (por ejemplo "Música").
}
