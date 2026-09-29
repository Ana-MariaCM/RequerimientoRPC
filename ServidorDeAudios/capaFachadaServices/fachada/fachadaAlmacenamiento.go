/**
 * @file fachadaAlmacenamiento.go
 * @brief Fachada que valida y almacena los audios enviados por el administrador.
 */
package fachada

import (
	"errors"
	"fmt"
	"strings"

	capaaccesodatos "almacenamiento/capaAccesoADatos"
	dtos "almacenamiento/capaFachadaServices/DTOs"
	componenteclientemetadatos "almacenamiento/componenteClienteMetadatos"
	"almacenamiento/configuracion"
)

/** @brief Error que indica que los datos recibidos no corresponden a un audio mp3 válido. */
var ErrAudioInvalido = errors.New("audio inválido")

/** @brief Error que indica que el servidor de metadatos no registró el audio (el mp3 se elimina). */
var ErrRegistroMetadatos = errors.New("no se registraron los metadatos")

/** @brief Prefijo del nombre de archivo según el identificador del tipo de audio. */
var prefijosPorTipo = map[int]string{1: "musica", 2: "podcast", 3: "audiolibro", 4: "ruido"}

/**
 * @brief Fachada del servidor de audios.
 */
type FachadaAlmacenamiento struct {
	repositorio      *capaaccesodatos.RepositorioAudios           ///< Acceso a los archivos en disco.
	clienteMetadatos *componenteclientemetadatos.ClienteMetadatos ///< Registro de metadatos en el servidor de metadatos.
}

/**
 * @brief Crea la fachada de almacenamiento.
 * @return Puntero a la nueva fachada.
 */
func NuevaFachadaAlmacenamiento() *FachadaAlmacenamiento {
	fmt.Println("Inicializando fachada de almacenamiento...")
	return &FachadaAlmacenamiento{
		repositorio:      capaaccesodatos.GetRepositorioAudios(),
		clienteMetadatos: componenteclientemetadatos.NuevoClienteMetadatos(configuracion.ObtenerURLMetadatos()),
	}
}

/**
 * @brief Valida y almacena un audio mp3 junto con los metadatos digitados por el administrador.
 *
 * Pasos:
 *  1. Valida el título, el tipo y que el archivo sea un mp3.
 *  2. Guarda el mp3 en la carpeta de audios (compartida con el servidor de
 *     streaming) con un nombre generado a partir del tipo y del título.
 *  3. Registra los metadatos en el servidor de metadatos (REST POST /audios),
 *     de modo que el audio aparece en el cliente en su siguiente consulta.
 *  4. Si el registro falla, elimina el mp3 para no dejar audios sin metadatos.
 * @param audio Tipo, título y metadatos del audio.
 * @param datos Contenido binario del archivo mp3.
 * @return DTO con la confirmación, o un error (ErrAudioInvalido o ErrRegistroMetadatos).
 */
func (thisF *FachadaAlmacenamiento) GuardarAudio(audio dtos.AudioAlmacenarDTOInput, datos []byte) (dtos.AudioAlmacenadoDTOOutput, error) {
	titulo := strings.TrimSpace(audio.Titulo)
	prefijo, tipoValido := prefijosPorTipo[audio.IdTipo]
	if titulo == "" || !tipoValido {
		return dtos.AudioAlmacenadoDTOOutput{}, fmt.Errorf("%w: el título y un tipo de audio válido son obligatorios", ErrAudioInvalido)
	}
	if len(datos) == 0 {
		return dtos.AudioAlmacenadoDTOOutput{}, fmt.Errorf("%w: el archivo está vacío", ErrAudioInvalido)
	}
	if !esMP3(datos) {
		return dtos.AudioAlmacenadoDTOOutput{}, fmt.Errorf("%w: el archivo no tiene formato mp3", ErrAudioInvalido)
	}

	archivo, err := thisF.repositorio.GuardarAudio(generarNombreArchivo(prefijo, titulo), datos)
	if err != nil {
		return dtos.AudioAlmacenadoDTOOutput{}, err
	}

	registrado, err := thisF.clienteMetadatos.RegistrarAudio(dtos.AudioRegistrarDTOOutput{
		IdTipo:        audio.IdTipo,
		Titulo:        titulo,
		NombreArchivo: archivo.NombreArchivo,
		Metadatos:     audio.Metadatos,
	})
	if err != nil {
		thisF.repositorio.EliminarAudio(archivo.NombreArchivo)
		return dtos.AudioAlmacenadoDTOOutput{}, fmt.Errorf("%w: %v", ErrRegistroMetadatos, err)
	}
	fmt.Printf("[REST] Metadatos registrados: audio %d en %s\n", registrado.Id, registrado.NombreTipo)

	return dtos.AudioAlmacenadoDTOOutput{
		Mensaje: "Audio almacenado y registrado correctamente",
		IdAudio: registrado.Id,
		Titulo:  registrado.Titulo,
		Tipo:    registrado.NombreTipo,
		Archivo: dtos.ArchivoAudioDTOOutput{NombreArchivo: archivo.NombreArchivo, TamanioBytes: archivo.TamanioBytes},
	}, nil
}

/**
 * @brief Lista los audios almacenados en el servidor.
 * @return Lista de archivos mp3 o un error si no se pudo leer la carpeta.
 */
func (thisF *FachadaAlmacenamiento) ListarAudios() ([]dtos.ArchivoAudioDTOOutput, error) {
	archivos, err := thisF.repositorio.ListarAudios()
	if err != nil {
		return nil, err
	}

	respuesta := make([]dtos.ArchivoAudioDTOOutput, 0, len(archivos))
	for _, archivo := range archivos {
		respuesta = append(respuesta, dtos.ArchivoAudioDTOOutput{NombreArchivo: archivo.NombreArchivo, TamanioBytes: archivo.TamanioBytes})
	}
	return respuesta, nil
}

/**
 * @brief Verifica la firma de un archivo mp3 (etiqueta ID3 o cabecera de trama MPEG).
 * @param datos Contenido del archivo.
 * @return true si los datos parecen un mp3.
 */
func esMP3(datos []byte) bool {
	if len(datos) >= 3 && string(datos[:3]) == "ID3" {
		return true
	}
	return len(datos) >= 2 && datos[0] == 0xFF && datos[1]&0xE0 == 0xE0
}

/**
 * @brief Genera un nombre de archivo seguro a partir del tipo y del título.
 *
 * Se pasan las letras a minúsculas, se quitan tildes y se reemplaza cualquier
 * otro carácter por "_" (por ejemplo "Música" + "La Bicicleta" ->
 * "musica_la_bicicleta.mp3"). Así el nombre funciona en cualquier sistema operativo.
 * @param prefijo Prefijo del tipo de audio (musica, podcast, audiolibro o ruido).
 * @param titulo Título del audio.
 * @return Nombre del archivo terminado en .mp3.
 */
func generarNombreArchivo(prefijo string, titulo string) string {
	sinTildes := strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n").
		Replace(strings.ToLower(titulo))

	var nombre strings.Builder
	for _, caracter := range sinTildes {
		if (caracter >= 'a' && caracter <= 'z') || (caracter >= '0' && caracter <= '9') {
			nombre.WriteRune(caracter)
		} else if nombre.Len() > 0 && !strings.HasSuffix(nombre.String(), "_") {
			nombre.WriteRune('_')
		}
	}

	base := strings.Trim(nombre.String(), "_")
	if base == "" {
		base = "audio"
	}
	return prefijo + "_" + base + ".mp3"
}
