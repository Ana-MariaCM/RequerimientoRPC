/**
 * @file fachadaMetadatos.go
 * @brief Fachada que expone los servicios de consulta de metadatos a la capa de controladores.
 *
 * Coordina los repositorios de la capa de acceso a datos y transforma los
 * modelos en DTOs que viajarán en formato JSON.
 */
package fachada

import (
	"errors"
	"fmt"

	capaaccesoadatos "metadatos/capaAccesoADatos"
	dtos "metadatos/capaFachadaServices/DTOs"
	"metadatos/modelos"
)

/** @brief Error que indica que el tipo de audio solicitado no existe. */
var ErrTipoNoEncontrado = errors.New("tipo de audio no encontrado")

/** @brief Error que indica que el audio solicitado no existe. */
var ErrAudioNoEncontrado = errors.New("audio no encontrado")

/**
 * @brief Fachada de servicios de metadatos de audios.
 */
type FachadaMetadatos struct {
	repositorioTipos  *capaaccesoadatos.RepositorioTiposAudio ///< Acceso a los tipos de audio.
	repositorioAudios *capaaccesoadatos.RepositorioAudios     ///< Acceso a los audios y sus metadatos.
}

/**
 * @brief Crea la fachada de metadatos obteniendo las instancias de los repositorios.
 * @return Puntero a la nueva fachada.
 */
func NuevaFachadaMetadatos() *FachadaMetadatos {
	fmt.Println("Inicializando fachada de metadatos...")
	return &FachadaMetadatos{
		repositorioTipos:  capaaccesoadatos.GetRepositorioTiposAudio(),
		repositorioAudios: capaaccesoadatos.GetRepositorioAudios(),
	}
}

/**
 * @brief Lista los tipos de audio registrados.
 * @return Lista de DTOs de tipos de audio.
 */
func (thisF *FachadaMetadatos) ListarTiposAudio() []dtos.TipoAudioDTOOutput {
	tipos := thisF.repositorioTipos.ListarTipos()

	respuesta := make([]dtos.TipoAudioDTOOutput, 0, len(tipos))
	for _, tipo := range tipos {
		respuesta = append(respuesta, dtos.TipoAudioDTOOutput{Id: tipo.Id, Nombre: tipo.Nombre})
	}
	return respuesta
}

/**
 * @brief Lista los audios de un tipo determinado.
 * @param idTipo Identificador del tipo de audio.
 * @return Lista de resúmenes de audio, o ErrTipoNoEncontrado si el tipo no existe.
 */
func (thisF *FachadaMetadatos) ListarAudiosPorTipo(idTipo int) ([]dtos.AudioResumenDTOOutput, error) {
	if _, existe := thisF.repositorioTipos.BuscarTipoPorId(idTipo); !existe {
		return nil, ErrTipoNoEncontrado
	}

	audios := thisF.repositorioAudios.ListarAudiosPorTipo(idTipo)
	respuesta := make([]dtos.AudioResumenDTOOutput, 0, len(audios))
	for _, audio := range audios {
		respuesta = append(respuesta, dtos.AudioResumenDTOOutput{
			Id:     audio.ObtenerId(),
			IdTipo: audio.ObtenerIdTipo(),
			Titulo: audio.ObtenerTitulo(),
		})
	}
	return respuesta, nil
}

/**
 * @brief Consulta el detalle completo (metadatos) de un audio.
 * @param idAudio Identificador del audio.
 * @return DTO con el detalle del audio, o ErrAudioNoEncontrado si no existe.
 */
func (thisF *FachadaMetadatos) ConsultarDetalleAudio(idAudio int) (dtos.AudioDetalleDTOOutput, error) {
	audio, existe := thisF.repositorioAudios.BuscarAudioPorId(idAudio)
	if !existe {
		return dtos.AudioDetalleDTOOutput{}, ErrAudioNoEncontrado
	}

	tipo, _ := thisF.repositorioTipos.BuscarTipoPorId(audio.ObtenerIdTipo())
	return convertirADetalleDTO(audio, tipo), nil
}

/**
 * @brief Convierte un modelo de audio en su DTO de detalle.
 * @param audio Modelo del audio.
 * @param tipo Tipo de audio al que pertenece.
 * @return DTO de detalle listo para serializar en JSON.
 */
func convertirADetalleDTO(audio modelos.Audio, tipo modelos.TipoAudio) dtos.AudioDetalleDTOOutput {
	metadatos := audio.ObtenerMetadatos()
	metadatosDTO := make([]dtos.MetadatoDTOOutput, 0, len(metadatos))
	for _, metadato := range metadatos {
		metadatosDTO = append(metadatosDTO, dtos.MetadatoDTOOutput{Nombre: metadato.Nombre, Valor: metadato.Valor})
	}

	return dtos.AudioDetalleDTOOutput{
		Id:            audio.ObtenerId(),
		IdTipo:        audio.ObtenerIdTipo(),
		NombreTipo:    tipo.Nombre,
		Titulo:        audio.ObtenerTitulo(),
		NombreArchivo: audio.ObtenerNombreArchivo(),
		Metadatos:     metadatosDTO,
	}
}
