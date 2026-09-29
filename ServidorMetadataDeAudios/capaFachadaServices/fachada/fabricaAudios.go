/**
 * @file fabricaAudios.go
 * @brief Construye el modelo concreto (Musica, Podcast, Audiolibro o RuidoBlanco)
 *        a partir de los metadatos recibidos en un DTO.
 */
package fachada

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	capaaccesodatos "metadatos/capaAccesoADatos"
	"metadatos/modelos"
)

/** @brief Error que indica que los datos recibidos para registrar un audio no son válidos. */
var ErrDatosInvalidos = errors.New("datos del audio inválidos")

/**
 * @brief Construye el modelo del audio según su tipo.
 * @param base Datos comunes del audio (id, tipo, título y archivo).
 * @param metadatos Metadatos propios del tipo (clave -> valor).
 * @return El modelo construido, o ErrDatosInvalidos si un valor numérico no es válido.
 */
func construirAudio(base modelos.AudioBase, metadatos map[string]string) (modelos.Audio, error) {
	lector := lectorMetadatos{valores: metadatos}

	var audio modelos.Audio
	switch base.IdTipo {
	case capaaccesodatos.IdTipoMusica:
		audio = modelos.Musica{
			AudioBase:         base,
			ArtistaPrincipal:  lector.texto("artistaPrincipal"),
			Album:             lector.texto("album"),
			GeneroMusical:     lector.texto("generoMusical"),
			SelloDiscografico: lector.texto("selloDiscografico"),
			AnioLanzamiento:   lector.entero("anioLanzamiento"),
			Duracion:          lector.texto("duracion"),
		}
	case capaaccesodatos.IdTipoPodcasts:
		audio = modelos.Podcast{
			AudioBase:              base,
			NombrePodcast:          lector.texto("nombrePodcast"),
			Anfitrion:              lector.texto("anfitrion"),
			Temporada:              lector.entero("temporada"),
			Episodio:               lector.entero("episodio"),
			NotasDelShow:           lector.texto("notasDelShow"),
			ClasificacionContenido: lector.texto("clasificacionContenido"),
			Duracion:               lector.texto("duracion"),
		}
	case capaaccesodatos.IdTipoAudiolibros:
		audio = modelos.Audiolibro{
			AudioBase:      base,
			Autor:          lector.texto("autor"),
			Narrador:       lector.texto("narrador"),
			Editorial:      lector.texto("editorial"),
			ISBN:           lector.texto("isbn"),
			Capitulo:       lector.entero("capitulo"),
			TotalCapitulos: lector.entero("totalCapitulos"),
			Duracion:       lector.texto("duracion"),
			Genero:         lector.texto("genero"),
		}
	case capaaccesodatos.IdTipoRuidoBlanco:
		audio = modelos.RuidoBlanco{
			AudioBase:           base,
			TipoSonido:          lector.texto("tipoSonido"),
			FuenteAudio:         lector.texto("fuenteAudio"),
			UsoSugerido:         lector.texto("usoSugerido"),
			ProveedorContenido:  lector.texto("proveedorContenido"),
			DuracionBucle:       lector.texto("duracionBucle"),
			FrecuenciaDominante: lector.texto("frecuenciaDominante"),
		}
	default:
		return nil, fmt.Errorf("%w: el tipo %d no existe", ErrDatosInvalidos, base.IdTipo)
	}

	if lector.err != nil {
		return nil, lector.err
	}
	return audio, nil
}

/**
 * @brief Lee valores de un mapa de metadatos registrando el primer error de conversión.
 */
type lectorMetadatos struct {
	valores map[string]string ///< Metadatos recibidos (clave -> valor).
	err     error             ///< Primer error encontrado al convertir un valor.
}

/**
 * @brief Obtiene un metadato de texto.
 * @param clave Nombre del metadato.
 * @return Valor sin espacios al inicio ni al final (vacío si no existe).
 */
func (thisL *lectorMetadatos) texto(clave string) string {
	return strings.TrimSpace(thisL.valores[clave])
}

/**
 * @brief Obtiene un metadato numérico entero.
 * @param clave Nombre del metadato.
 * @return Valor entero (0 si está vacío o no es válido; en ese último caso se registra el error).
 */
func (thisL *lectorMetadatos) entero(clave string) int {
	valor := thisL.texto(clave)
	if valor == "" {
		return 0
	}
	numero, err := strconv.Atoi(valor)
	if err != nil && thisL.err == nil {
		thisL.err = fmt.Errorf("%w: el metadato %s debe ser un número entero", ErrDatosInvalidos, clave)
	}
	return numero
}
