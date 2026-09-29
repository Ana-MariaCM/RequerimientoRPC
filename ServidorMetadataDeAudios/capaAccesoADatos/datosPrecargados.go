/**
 * @file datosPrecargados.go
 * @brief Datos que se cargan en memoria al iniciar el servidor de metadatos.
 *
 * Se cargan los cuatro tipos de audio y cuatro audios por cada tipo (el
 * requerimiento exige al menos dos). Cada audio referencia, mediante
 * NombreArchivo, el archivo mp3 que se encuentra en el servidor de audios y que
 * transmite el servidor de streaming.
 *
 * @note Datos de ejemplo con fines académicos. Algunos valores (narradores,
 *       duraciones, ISBN, podcasts y proveedores de ruido blanco) son ilustrativos.
 */
package capaaccesoadatos

import "metadatos/modelos"

/** @brief Identificador del tipo Música. */
const IdTipoMusica = 1

/** @brief Identificador del tipo Podcasts. */
const IdTipoPodcasts = 2

/** @brief Identificador del tipo Audiolibros. */
const IdTipoAudiolibros = 3

/** @brief Identificador del tipo Ruido Blanco. */
const IdTipoRuidoBlanco = 4

/**
 * @brief Construye la lista de tipos de audio precargados.
 * @return Tipos de audio disponibles en el sistema.
 */
func cargarTiposAudio() []modelos.TipoAudio {
	return []modelos.TipoAudio{
		{Id: IdTipoMusica, Nombre: "Música"},
		{Id: IdTipoPodcasts, Nombre: "Podcasts"},
		{Id: IdTipoAudiolibros, Nombre: "Audiolibros"},
		{Id: IdTipoRuidoBlanco, Nombre: "Ruido Blanco"},
	}
}

/**
 * @brief Construye la lista de audios precargados (cuatro por cada tipo).
 * @return Audios de todos los tipos.
 */
func cargarAudios() []modelos.Audio {
	audios := []modelos.Audio{}
	audios = append(audios, cargarMusica()...)
	audios = append(audios, cargarPodcasts()...)
	audios = append(audios, cargarAudiolibros()...)
	audios = append(audios, cargarRuidoBlanco()...)
	return audios
}

/**
 * @brief Construye los audios de tipo Música.
 * @return Canciones precargadas.
 */
func cargarMusica() []modelos.Audio {
	return []modelos.Audio{
		modelos.Musica{
			AudioBase:        modelos.AudioBase{Id: 1, IdTipo: IdTipoMusica, Titulo: "Sonido Bestial", NombreArchivo: "Sonido_Bestial.mp3"},
			ArtistaPrincipal: "Ricardo Ray & Bobby Cruz", Album: "Sonido Bestial Plus", GeneroMusical: "Salsa",
			SelloDiscografico: "Ricardo Ray & Bobby Cruz", AnioLanzamiento: 2000, Duracion: "6:48",
		},
		modelos.Musica{
			AudioBase:        modelos.AudioBase{Id: 2, IdTipo: IdTipoMusica, Titulo: "Mis Ojos Lloran Por Ti", NombreArchivo: "Mi_Ojos_Lloran_Por_Ti.mp3"},
			ArtistaPrincipal: "Big Boy", Album: "Mis Ojos Lloran Por Ti", GeneroMusical: "Reggaeton",
			SelloDiscografico: "Big Boy", AnioLanzamiento: 2003, Duracion: "4:55",
		},
	}
}

/**
 * @brief Construye los audios de tipo Podcast.
 * @return Episodios de podcast precargados.
 */
func cargarPodcasts() []modelos.Audio {
	return []modelos.Audio{
		modelos.Podcast{
			AudioBase:     modelos.AudioBase{Id: 3, IdTipo: IdTipoPodcasts, Titulo: "Podcast Educativo, Tabaquismo", NombreArchivo: "Podcast_Educativo_Tabaquismo.mp3"},
			NombrePodcast: "Podcast educativo sobre el tabaquismo", Anfitrion: "José José", Temporada: 1, Episodio: 1,
			NotasDelShow:           "Por qué el consumo de tabaco es perjudicicial para la salud de todos.",
			ClasificacionContenido: "Para toda la familia", Duracion: "01:15",
		},
		modelos.Podcast{
			AudioBase:     modelos.AudioBase{Id: 4, IdTipo: IdTipoPodcasts, Titulo: "Trastornos del Aprendizaje", NombreArchivo: "Podcast_Trastornos_Del_Aprendizaje.mp3"},
			NombrePodcast: "Trastornos del Aprendizaje", Anfitrion: "Laura Muñoz", Temporada: 1, Episodio: 2,
			NotasDelShow:           "Identificar los diferentes trastornos del aprendizaje que nos pueden afectar a todos.",
			ClasificacionContenido: "Para toda la familia", Duracion: "03:00",
		},
	}
}

/**
 * @brief Construye los audios de tipo Audiolibro.
 * @return Audiolibros precargados.
 */
func cargarAudiolibros() []modelos.Audio {
	return []modelos.Audio{
		modelos.Audiolibro{
			AudioBase: modelos.AudioBase{Id: 5, IdTipo: IdTipoAudiolibros, Titulo: "El silencio de las sirenas", NombreArchivo: "Audiolibro_El_silencio_de_las_sirenas__Franz_Kafka.mp3"},
			Autor:     "Franz Kafka", Narrador: "Andrés Salazar", Editorial: "Salamandra", ISBN: "978-84-7888-445-2",
			Capitulo: 1, TotalCapitulos: 2, Duracion: "04:17", Genero: "Ficción",
		},
		modelos.Audiolibro{
			AudioBase: modelos.AudioBase{Id: 6, IdTipo: IdTipoAudiolibros, Titulo: "Beatriz, una palabra enorme", NombreArchivo: "Audiolibro_Beatriz_una_palabra_enorme__Mario_Benedetti.mp3"},
			Autor:     "Mario Benedetti", Narrador: "Ana María Castro", Editorial: "Penguin Random House Grupo Editorial", ISBN: "978-0-307-47472-8",
			Capitulo: 1, TotalCapitulos: 3, Duracion: "09:00", Genero: "Historia",
		},
	}
}

/**Historia
 * @brief Construye los audios de tipo Ruido Blanco.
 * @return Bucles de ruido precargados.
 */
func cargarRuidoBlanco() []modelos.Audio {
	return []modelos.Audio{
		modelos.RuidoBlanco{
			AudioBase:  modelos.AudioBase{Id: 7, IdTipo: IdTipoRuidoBlanco, Titulo: "Ventilador para concentrarse", NombreArchivo: "ruido_ventilador.mp3"},
			TipoSonido: "Ruido Blanco", FuenteAudio: "Ventilador", UsoSugerido: "Concentración",
			ProveedorContenido: "Estudio Enfoque", DuracionBucle: "00:20", FrecuenciaDominante: "Agudos",
		},
		modelos.RuidoBlanco{
			AudioBase:  modelos.AudioBase{Id: 8, IdTipo: IdTipoRuidoBlanco, Titulo: "Olas en la costa", NombreArchivo: "ruido_olas_costa.mp3"},
			TipoSonido: "Ruido Rosa", FuenteAudio: "Mar", UsoSugerido: "Meditación",
			ProveedorContenido: "Sonidos del Cauca", DuracionBucle: "00:40", FrecuenciaDominante: "Graves",
		},
	}
}
