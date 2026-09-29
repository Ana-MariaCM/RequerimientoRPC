/**
 * @file formularioMetadatos.go
 * @brief Definición de los tipos de audio y de los metadatos que el administrador
 *        debe ingresar para cada uno.
 */
package vistas

/**
 * @brief Campo de metadato que se solicita por consola.
 */
type CampoMetadato struct {
	Clave    string ///< Clave con la que se envía al servidor (por ejemplo "artistaPrincipal").
	Etiqueta string ///< Texto que se muestra al administrador.
	Numerico bool   ///< Indica si el valor debe ser un número entero.
}

/**
 * @brief Tipo de audio con los metadatos que lo describen.
 */
type TipoAudio struct {
	Id             int             ///< Identificador del tipo en el servidor de metadatos.
	Nombre         string          ///< Nombre del tipo de audio.
	EtiquetaTitulo string          ///< Texto con el que se solicita el título según el tipo.
	Campos         []CampoMetadato ///< Metadatos que se solicitan para este tipo.
}

/** @brief Tipos de audio y sus metadatos (mismos identificadores del servidor de metadatos). */
var tiposAudio = []TipoAudio{
	{Id: 1, Nombre: "Música", EtiquetaTitulo: "Título de la canción", Campos: []CampoMetadato{
		{Clave: "artistaPrincipal", Etiqueta: "Artista principal"},
		{Clave: "album", Etiqueta: "Álbum"},
		{Clave: "generoMusical", Etiqueta: "Género musical (Pop, Rock, Jazz...)"},
		{Clave: "selloDiscografico", Etiqueta: "Sello discográfico"},
		{Clave: "anioLanzamiento", Etiqueta: "Año de lanzamiento", Numerico: true},
		{Clave: "duracion", Etiqueta: "Duración (mm:ss)"},
	}},
	{Id: 2, Nombre: "Podcasts", EtiquetaTitulo: "Título del episodio", Campos: []CampoMetadato{
		{Clave: "nombrePodcast", Etiqueta: "Nombre del podcast"},
		{Clave: "anfitrion", Etiqueta: "Anfitrión (host)"},
		{Clave: "temporada", Etiqueta: "Número de temporada", Numerico: true},
		{Clave: "episodio", Etiqueta: "Número de episodio", Numerico: true},
		{Clave: "notasDelShow", Etiqueta: "Notas del show"},
		{Clave: "clasificacionContenido", Etiqueta: "Clasificación (Explícito / Para toda la familia)"},
		{Clave: "duracion", Etiqueta: "Duración (mm:ss)"},
	}},
	{Id: 3, Nombre: "Audiolibros", EtiquetaTitulo: "Título del libro", Campos: []CampoMetadato{
		{Clave: "autor", Etiqueta: "Autor"},
		{Clave: "narrador", Etiqueta: "Narrador"},
		{Clave: "editorial", Etiqueta: "Editorial"},
		{Clave: "isbn", Etiqueta: "ISBN"},
		{Clave: "capitulo", Etiqueta: "Capítulo", Numerico: true},
		{Clave: "totalCapitulos", Etiqueta: "Total de capítulos", Numerico: true},
		{Clave: "duracion", Etiqueta: "Duración"},
		{Clave: "genero", Etiqueta: "Género"},
	}},
	{Id: 4, Nombre: "Ruido Blanco", EtiquetaTitulo: "Nombre del sonido", Campos: []CampoMetadato{
		{Clave: "tipoSonido", Etiqueta: "Tipo de sonido (Ruido Blanco, Marrón, Rosa)"},
		{Clave: "fuenteAudio", Etiqueta: "Fuente del audio (Lluvia, Ventilador, Bosque...)"},
		{Clave: "usoSugerido", Etiqueta: "Uso sugerido (Dormir, Concentración, Meditación)"},
		{Clave: "proveedorContenido", Etiqueta: "Proveedor de contenido"},
		{Clave: "duracionBucle", Etiqueta: "Duración del bucle (mm:ss)"},
		{Clave: "frecuenciaDominante", Etiqueta: "Frecuencia dominante (Graves / Agudos)"},
	}},
}
