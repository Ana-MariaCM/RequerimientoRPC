/**
 * @file musica.go
 * @brief Modelo con los metadatos de un audio de tipo Música.
 */
package modelos

import "strconv"

/**
 * @brief Audio de tipo Música.
 *
 * El título de la canción se almacena en AudioBase::Titulo.
 */
type Musica struct {
	AudioBase                // Datos comunes (id, tipo, título de la canción y archivo).
	ArtistaPrincipal  string ///< Intérprete o banda líder del proyecto.
	Album             string ///< Disco, EP o sencillo que agrupa la obra.
	GeneroMusical     string ///< Categoría rítmica o cultural (Pop, Rock, Jazz).
	SelloDiscografico string ///< Empresa dueña de los derechos fonográficos.
	AnioLanzamiento   int    ///< Año oficial de publicación.
	Duracion          string ///< Duración de la canción (mm:ss).
}

/**
 * @brief Construye la lista ordenada de metadatos de la canción.
 * @return Metadatos de la canción listos para ser presentados.
 */
func (thisM Musica) ObtenerMetadatos() []Metadato {
	return []Metadato{
		{Nombre: "Título de la canción", Valor: thisM.Titulo},
		{Nombre: "Artista principal", Valor: thisM.ArtistaPrincipal},
		{Nombre: "Álbum", Valor: thisM.Album},
		{Nombre: "Género musical", Valor: thisM.GeneroMusical},
		{Nombre: "Sello discográfico", Valor: thisM.SelloDiscografico},
		{Nombre: "Año de lanzamiento", Valor: strconv.Itoa(thisM.AnioLanzamiento)},
		{Nombre: "Duración", Valor: thisM.Duracion},
	}
}
