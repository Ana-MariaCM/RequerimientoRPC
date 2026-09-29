/**
 * @file audio.go
 * @brief Abstracciones comunes a todos los audios: interfaz Audio, AudioBase y Metadato.
 */
package modelos

/**
 * @brief Par etiqueta-valor que describe un metadato de un audio.
 *
 * Se utiliza para presentar de forma uniforme los metadatos de cualquier
 * tipo de audio, conservando el orden en que deben mostrarse.
 */
type Metadato struct {
	Nombre string ///< Etiqueta del metadato (por ejemplo "Autor").
	Valor  string ///< Valor del metadato (por ejemplo "Gabriel García Márquez").
}

/**
 * @brief Comportamiento común de todos los audios (música, podcast, audiolibro y ruido blanco).
 *
 * Cada tipo concreto implementa esta interfaz, de modo que la fachada y el
 * repositorio pueden tratarlos de forma polimórfica.
 */
type Audio interface {
	/** @brief Devuelve el identificador del audio. */
	ObtenerId() int
	/** @brief Devuelve el identificador del tipo de audio al que pertenece. */
	ObtenerIdTipo() int
	/** @brief Devuelve el título que se muestra en la lista de audios. */
	ObtenerTitulo() string
	/** @brief Devuelve el nombre del archivo mp3 asociado en el servidor de audios. */
	ObtenerNombreArchivo() string
	/** @brief Devuelve la lista ordenada de metadatos propios del tipo de audio. */
	ObtenerMetadatos() []Metadato
}

/**
 * @brief Datos comunes a cualquier audio.
 *
 * Los tipos concretos (Musica, Podcast, Audiolibro, RuidoBlanco) embeben esta
 * estructura para reutilizar sus campos y métodos.
 */
type AudioBase struct {
	Id            int    ///< Identificador único del audio.
	IdTipo        int    ///< Identificador del tipo de audio al que corresponde el audio.
	Titulo        string ///< Título del audio que se muestra en la lista.
	NombreArchivo string ///< Nombre del archivo mp3 que reproduce el servidor de streaming.
}

/**
 * @brief Devuelve el identificador del audio.
 * @return Identificador del audio.
 */
func (thisA AudioBase) ObtenerId() int {
	return thisA.Id
}

/**
 * @brief Devuelve el identificador del tipo de audio.
 * @return Identificador del tipo.
 */
func (thisA AudioBase) ObtenerIdTipo() int {
	return thisA.IdTipo
}

/**
 * @brief Devuelve el título del audio.
 * @return Título del audio.
 */
func (thisA AudioBase) ObtenerTitulo() string {
	return thisA.Titulo
}

/**
 * @brief Devuelve el nombre del archivo mp3 asociado al audio.
 * @return Nombre del archivo mp3.
 */
func (thisA AudioBase) ObtenerNombreArchivo() string {
	return thisA.NombreArchivo
}
