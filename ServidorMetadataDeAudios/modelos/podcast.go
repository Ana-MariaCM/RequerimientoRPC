/**
 * @file podcast.go
 * @brief Modelo con los metadatos de un audio de tipo Podcast.
 */
package modelos

import "fmt"

/**
 * @brief Audio de tipo Podcast.
 *
 * El título del episodio se almacena en AudioBase::Titulo.
 */
type Podcast struct {
	AudioBase                     // Datos comunes (id, tipo, título del episodio y archivo).
	NombrePodcast          string ///< Título general del programa o serie.
	Anfitrion              string ///< Persona que dirige o presenta el contenido (host).
	Temporada              int    ///< Número de temporada dentro de la serie.
	Episodio               int    ///< Número de episodio dentro de la temporada.
	NotasDelShow           string ///< Descripción detallada y enlaces de referencia.
	ClasificacionContenido string ///< "Explícito" o "Para toda la familia".
	Duracion               string ///< Duración del episodio (mm:ss).
}

/**
 * @brief Construye la lista ordenada de metadatos del episodio.
 * @return Metadatos del podcast listos para ser presentados.
 */
func (thisP Podcast) ObtenerMetadatos() []Metadato {
	return []Metadato{
		{Nombre: "Nombre del podcast", Valor: thisP.NombrePodcast},
		{Nombre: "Título del episodio", Valor: thisP.Titulo},
		{Nombre: "Anfitrión (Host)", Valor: thisP.Anfitrion},
		{Nombre: "Temporada/Episodio", Valor: fmt.Sprintf("Temporada %d / Episodio %d", thisP.Temporada, thisP.Episodio)},
		{Nombre: "Notas del show", Valor: thisP.NotasDelShow},
		{Nombre: "Clasificación de contenido", Valor: thisP.ClasificacionContenido},
		{Nombre: "Duración", Valor: thisP.Duracion},
	}
}
