/**
 * @file audiolibro.go
 * @brief Modelo con los metadatos de un audio de tipo Audiolibro.
 */
package modelos

import "fmt"

/**
 * @brief Audio de tipo Audiolibro.
 *
 * El título del libro se almacena en AudioBase::Titulo.
 */
type Audiolibro struct {
	AudioBase             // Datos comunes (id, tipo, título del libro y archivo).
	Autor          string ///< Escritor original del texto.
	Narrador       string ///< Actor de voz que realiza la lectura.
	Editorial      string ///< Casa editora responsable de la producción.
	ISBN           string ///< Código internacional único de identificación del libro.
	Capitulo       int    ///< Marcador de posición (capítulo) dentro de la estructura narrativa.
	TotalCapitulos int    ///< Número total de capítulos del libro.
	Duracion       string ///< Duración del audio.
	Genero         string ///< Género literario.
}

/**
 * @brief Construye la lista ordenada de metadatos del audiolibro.
 * @return Metadatos del audiolibro listos para ser presentados.
 */
func (thisA Audiolibro) ObtenerMetadatos() []Metadato {
	return []Metadato{
		{Nombre: "Título del audio libro", Valor: thisA.Titulo},
		{Nombre: "Autor", Valor: thisA.Autor},
		{Nombre: "Narrador", Valor: thisA.Narrador},
		{Nombre: "Editorial", Valor: thisA.Editorial},
		{Nombre: "ISBN", Valor: thisA.ISBN},
		{Nombre: "Capítulo", Valor: fmt.Sprintf("%d de %d", thisA.Capitulo, thisA.TotalCapitulos)},
		{Nombre: "Duración", Valor: thisA.Duracion},
		{Nombre: "Género", Valor: thisA.Genero},
	}
}
