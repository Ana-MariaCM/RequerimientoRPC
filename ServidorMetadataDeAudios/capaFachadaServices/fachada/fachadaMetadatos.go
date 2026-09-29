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
	"strings"
	"sync"

	capaaccesoadatos "metadatos/capaAccesoADatos"
	dtos "metadatos/capaFachadaServices/DTOs"
	"metadatos/modelos"
)

/** @brief Error que indica que el tipo de audio solicitado no existe. */
var ErrTipoNoEncontrado = errors.New("tipo de audio no encontrado")

/** @brief Error que indica que el audio solicitado no existe. */
var ErrAudioNoEncontrado = errors.New("audio no encontrado")

/** @brief Error que indica que los metadatos no se pudieron guardar en el archivo JSON. */
var ErrPersistencia = errors.New("no se pudieron guardar los metadatos")

/**
 * @brief Fachada de servicios de metadatos de audios.
 */
type FachadaMetadatos struct {
	repositorioTipos       *capaaccesoadatos.RepositorioTiposAudio        ///< Acceso a los tipos de audio.
	repositorioAudios      *capaaccesoadatos.RepositorioAudios            ///< Acceso (en memoria) a los audios y sus metadatos.
	repositorioRegistrados *capaaccesoadatos.RepositorioAudiosRegistrados ///< Archivo JSON con los audios que registra el administrador.
	muRegistro             sync.Mutex                                     ///< Evita que dos registros simultáneos obtengan el mismo id.
}

/**
 * @brief Crea la fachada de metadatos y carga los audios registrados guardados en el archivo JSON.
 * @return Puntero a la nueva fachada.
 */
func NuevaFachadaMetadatos() *FachadaMetadatos {
	fmt.Println("Inicializando fachada de metadatos...")
	fachada := &FachadaMetadatos{
		repositorioTipos:       capaaccesoadatos.GetRepositorioTiposAudio(),
		repositorioAudios:      capaaccesoadatos.GetRepositorioAudios(),
		repositorioRegistrados: capaaccesoadatos.GetRepositorioAudiosRegistrados(),
	}
	fachada.cargarAudiosRegistrados()
	return fachada
}

/**
 * @brief Carga en memoria los audios que el administrador registró en ejecuciones anteriores.
 *
 * Cada registro del archivo JSON se convierte en su modelo (Musica, Podcast, ...).
 * Los registros inválidos se omiten con un aviso. Si el identificador guardado
 * ya lo usa otro audio (por ejemplo, porque se agregaron audios precargados),
 * se le asigna uno nuevo y el archivo se actualiza.
 */
func (thisF *FachadaMetadatos) cargarAudiosRegistrados() {
	registros, err := thisF.repositorioRegistrados.Cargar()
	if err != nil {
		fmt.Println("Advertencia:", err)
		return
	}

	cargados := []modelos.AudioRegistrado{}
	huboCambios := false
	for _, registro := range registros {
		if _, existe := thisF.repositorioTipos.BuscarTipoPorId(registro.IdTipo); !existe {
			fmt.Printf("Advertencia: se omite \"%s\" porque el tipo %d no existe\n", registro.Titulo, registro.IdTipo)
			huboCambios = true
			continue
		}
		if registro.Id <= 0 || thisF.repositorioAudios.IdOcupadoPorOtroAudio(registro.Id, registro.NombreArchivo) {
			registro.Id, _ = thisF.repositorioAudios.ObtenerIdParaArchivo(registro.NombreArchivo)
			huboCambios = true
		}

		base := modelos.AudioBase{Id: registro.Id, IdTipo: registro.IdTipo, Titulo: registro.Titulo, NombreArchivo: registro.NombreArchivo}
		audio, err := construirAudio(base, registro.Metadatos)
		if err != nil {
			fmt.Printf("Advertencia: se omite \"%s\": %v\n", registro.Titulo, err)
			huboCambios = true
			continue
		}
		thisF.repositorioAudios.GuardarAudio(audio)
		cargados = append(cargados, registro)
	}

	if huboCambios {
		if err := thisF.repositorioRegistrados.GuardarTodos(cargados); err != nil {
			fmt.Println("Advertencia:", err)
		}
	}
	fmt.Printf("Audios registrados por el administrador cargados desde %s: %d\n",
		thisF.repositorioRegistrados.ObtenerRutaArchivo(), len(cargados))
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
 * @brief Registra los metadatos de un audio nuevo almacenado por el administrador.
 *
 * Los metadatos se guardan primero en el archivo JSON (para que no se pierdan
 * al reiniciar el servidor) y luego en memoria, desde donde los consulta el
 * cliente. Si ya existe un audio con el mismo archivo mp3, sus metadatos se actualizan.
 * @param audioDTO Datos del audio recibidos del servidor de audios.
 * @return Detalle del audio registrado, true si es nuevo (false si se actualizó),
 *         o un error (ErrDatosInvalidos o ErrPersistencia).
 */
func (thisF *FachadaMetadatos) RegistrarAudio(audioDTO dtos.AudioRegistrarDTOInput) (dtos.AudioDetalleDTOOutput, bool, error) {
	tipo, existe := thisF.repositorioTipos.BuscarTipoPorId(audioDTO.IdTipo)
	if !existe {
		return dtos.AudioDetalleDTOOutput{}, false, fmt.Errorf("%w: el tipo %d no existe", ErrDatosInvalidos, audioDTO.IdTipo)
	}
	titulo := strings.TrimSpace(audioDTO.Titulo)
	nombreArchivo := strings.TrimSpace(audioDTO.NombreArchivo)
	if titulo == "" || nombreArchivo == "" {
		return dtos.AudioDetalleDTOOutput{}, false, fmt.Errorf("%w: el título y el nombre del archivo son obligatorios", ErrDatosInvalidos)
	}

	thisF.muRegistro.Lock()
	defer thisF.muRegistro.Unlock()

	idAudio, esNuevo := thisF.repositorioAudios.ObtenerIdParaArchivo(nombreArchivo)
	base := modelos.AudioBase{Id: idAudio, IdTipo: tipo.Id, Titulo: titulo, NombreArchivo: nombreArchivo}
	audio, err := construirAudio(base, audioDTO.Metadatos)
	if err != nil {
		return dtos.AudioDetalleDTOOutput{}, false, err
	}

	registro := modelos.AudioRegistrado{
		Id:            idAudio,
		IdTipo:        tipo.Id,
		Titulo:        titulo,
		NombreArchivo: nombreArchivo,
		Metadatos:     audioDTO.Metadatos,
	}
	if err := thisF.repositorioRegistrados.Guardar(registro); err != nil {
		return dtos.AudioDetalleDTOOutput{}, false, fmt.Errorf("%w: %v", ErrPersistencia, err)
	}

	thisF.repositorioAudios.GuardarAudio(audio)
	return convertirADetalleDTO(audio, tipo), esNuevo, nil
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
