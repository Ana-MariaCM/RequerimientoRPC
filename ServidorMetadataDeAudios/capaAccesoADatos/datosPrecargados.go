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
			AudioBase:        modelos.AudioBase{Id: 1, IdTipo: IdTipoMusica, Titulo: "Bohemian Rhapsody", NombreArchivo: "musica_bohemian_rhapsody.mp3"},
			ArtistaPrincipal: "Queen", Album: "A Night at the Opera", GeneroMusical: "Rock",
			SelloDiscografico: "EMI", AnioLanzamiento: 1975, Duracion: "5:55",
		},
		modelos.Musica{
			AudioBase:        modelos.AudioBase{Id: 2, IdTipo: IdTipoMusica, Titulo: "Take Five", NombreArchivo: "musica_take_five.mp3"},
			ArtistaPrincipal: "The Dave Brubeck Quartet", Album: "Time Out", GeneroMusical: "Jazz",
			SelloDiscografico: "Columbia Records", AnioLanzamiento: 1959, Duracion: "5:24",
		},
		modelos.Musica{
			AudioBase:        modelos.AudioBase{Id: 3, IdTipo: IdTipoMusica, Titulo: "Billie Jean", NombreArchivo: "musica_billie_jean.mp3"},
			ArtistaPrincipal: "Michael Jackson", Album: "Thriller", GeneroMusical: "Pop",
			SelloDiscografico: "Epic Records", AnioLanzamiento: 1982, Duracion: "4:54",
		},
		modelos.Musica{
			AudioBase:        modelos.AudioBase{Id: 4, IdTipo: IdTipoMusica, Titulo: "Hotel California", NombreArchivo: "musica_hotel_california.mp3"},
			ArtistaPrincipal: "Eagles", Album: "Hotel California", GeneroMusical: "Rock",
			SelloDiscografico: "Asylum Records", AnioLanzamiento: 1976, Duracion: "6:30",
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
			AudioBase:     modelos.AudioBase{Id: 5, IdTipo: IdTipoPodcasts, Titulo: "¿Qué es una llamada a procedimiento remoto?", NombreArchivo: "podcast_sd_t1_e1.mp3"},
			NombrePodcast: "Sistemas Distribuidos al Día", Anfitrion: "Laura Muñoz", Temporada: 1, Episodio: 1,
			NotasDelShow:           "Introducción al modelo RPC: stubs, serialización y transparencia. Referencia: https://grpc.io/docs/what-is-grpc/introduction/",
			ClasificacionContenido: "Para toda la familia", Duracion: "18:30",
		},
		modelos.Podcast{
			AudioBase:     modelos.AudioBase{Id: 6, IdTipo: IdTipoPodcasts, Titulo: "Streaming con gRPC en Go", NombreArchivo: "podcast_sd_t1_e2.mp3"},
			NombrePodcast: "Sistemas Distribuidos al Día", Anfitrion: "Laura Muñoz", Temporada: 1, Episodio: 2,
			NotasDelShow:           "Server streaming, fragmentación de archivos y control de flujo. Referencia: https://grpc.io/docs/languages/go/basics/",
			ClasificacionContenido: "Para toda la familia", Duracion: "22:10",
		},
		modelos.Podcast{
			AudioBase:     modelos.AudioBase{Id: 7, IdTipo: IdTipoPodcasts, Titulo: "El faro abandonado", NombreArchivo: "podcast_historias_t2_e5.mp3"},
			NombrePodcast: "Historias de Medianoche", Anfitrion: "Carlos Rivera", Temporada: 2, Episodio: 5,
			NotasDelShow:           "Relato de suspenso con escenas intensas y lenguaje fuerte.",
			ClasificacionContenido: "Explícito", Duracion: "35:45",
		},
		modelos.Podcast{
			AudioBase:     modelos.AudioBase{Id: 8, IdTipo: IdTipoPodcasts, Titulo: "¿Por qué el cielo es azul?", NombreArchivo: "podcast_ciencia_t3_e12.mp3"},
			NombrePodcast: "Ciencia en Cápsulas", Anfitrion: "Mariana Gómez", Temporada: 3, Episodio: 12,
			NotasDelShow:           "La dispersión de Rayleigh explicada de forma sencilla. Referencia: https://es.wikipedia.org/wiki/Dispersión_de_Rayleigh",
			ClasificacionContenido: "Para toda la familia", Duracion: "12:05",
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
			AudioBase: modelos.AudioBase{Id: 9, IdTipo: IdTipoAudiolibros, Titulo: "Harry Potter y la Piedra Filosofal", NombreArchivo: "audiolibro_harry_potter.mp3"},
			Autor:     "J. K. Rowling", Narrador: "Andrés Salazar", Editorial: "Salamandra", ISBN: "978-84-7888-445-2",
			Capitulo: 1, TotalCapitulos: 17, Duracion: "8:48:00", Genero: "Fantasía",
		},
		modelos.Audiolibro{
			AudioBase: modelos.AudioBase{Id: 10, IdTipo: IdTipoAudiolibros, Titulo: "Cien años de soledad", NombreArchivo: "audiolibro_cien_anios_de_soledad.mp3"},
			Autor:     "Gabriel García Márquez", Narrador: "Gustavo Bonfigli", Editorial: "Penguin Random House Grupo Editorial", ISBN: "978-0-307-47472-8",
			Capitulo: 1, TotalCapitulos: 20, Duracion: "4:42", Genero: "Ficción",
		},
		modelos.Audiolibro{
			AudioBase: modelos.AudioBase{Id: 11, IdTipo: IdTipoAudiolibros, Titulo: "El Señor de los Anillos: La Comunidad del Anillo", NombreArchivo: "audiolibro_senor_de_los_anillos.mp3"},
			Autor:     "J. R. R. Tolkien", Narrador: "Julián Restrepo", Editorial: "Minotauro", ISBN: "978-84-450-0006-9",
			Capitulo: 1, TotalCapitulos: 22, Duracion: "19:07:00", Genero: "Fantasía épica",
		},
		modelos.Audiolibro{
			AudioBase: modelos.AudioBase{Id: 12, IdTipo: IdTipoAudiolibros, Titulo: "Orgullo y Prejuicio", NombreArchivo: "audiolibro_orgullo_y_prejuicio.mp3"},
			Autor:     "Jane Austen", Narrador: "Lucía Fernández", Editorial: "Penguin Clásicos", ISBN: "978-84-9105-092-6",
			Capitulo: 1, TotalCapitulos: 61, Duracion: "11:35:00", Genero: "Novela romántica",
		},
	}
}

/**
 * @brief Construye los audios de tipo Ruido Blanco.
 * @return Bucles de ruido precargados.
 */
func cargarRuidoBlanco() []modelos.Audio {
	return []modelos.Audio{
		modelos.RuidoBlanco{
			AudioBase:  modelos.AudioBase{Id: 13, IdTipo: IdTipoRuidoBlanco, Titulo: "Lluvia suave para dormir", NombreArchivo: "ruido_lluvia_suave.mp3"},
			TipoSonido: "Ruido Rosa", FuenteAudio: "Lluvia", UsoSugerido: "Dormir",
			ProveedorContenido: "Sonidos del Cauca", DuracionBucle: "00:30", FrecuenciaDominante: "Graves",
		},
		modelos.RuidoBlanco{
			AudioBase:  modelos.AudioBase{Id: 14, IdTipo: IdTipoRuidoBlanco, Titulo: "Ventilador para concentrarse", NombreArchivo: "ruido_ventilador.mp3"},
			TipoSonido: "Ruido Blanco", FuenteAudio: "Ventilador", UsoSugerido: "Concentración",
			ProveedorContenido: "Estudio Enfoque", DuracionBucle: "00:20", FrecuenciaDominante: "Agudos",
		},
		modelos.RuidoBlanco{
			AudioBase:  modelos.AudioBase{Id: 15, IdTipo: IdTipoRuidoBlanco, Titulo: "Bosque profundo", NombreArchivo: "ruido_bosque_profundo.mp3"},
			TipoSonido: "Ruido Marrón", FuenteAudio: "Bosque", UsoSugerido: "Meditación",
			ProveedorContenido: "Naturaleza Viva", DuracionBucle: "00:45", FrecuenciaDominante: "Graves",
		},
		modelos.RuidoBlanco{
			AudioBase:  modelos.AudioBase{Id: 16, IdTipo: IdTipoRuidoBlanco, Titulo: "Olas en la costa", NombreArchivo: "ruido_olas_costa.mp3"},
			TipoSonido: "Ruido Rosa", FuenteAudio: "Mar", UsoSugerido: "Meditación",
			ProveedorContenido: "Sonidos del Cauca", DuracionBucle: "00:40", FrecuenciaDominante: "Graves",
		},
	}
}
