/**
 * @file fachadaAlmacenamiento.go
 * @brief Fachada que valida y almacena los audios enviados por el administrador.
 */
package fachada

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	capaaccesodatos "almacenamiento/capaAccesoADatos"
	dtos "almacenamiento/capaFachadaServices/DTOs"
)

/** @brief Error que indica que los datos recibidos no corresponden a un audio mp3 válido. */
var ErrAudioInvalido = errors.New("audio inválido")

/**
 * @brief Fachada del servidor de audios.
 */
type FachadaAlmacenamiento struct {
	repositorio *capaaccesodatos.RepositorioAudios ///< Acceso a los archivos en disco.
}

/**
 * @brief Crea la fachada de almacenamiento.
 * @return Puntero a la nueva fachada.
 */
func NuevaFachadaAlmacenamiento() *FachadaAlmacenamiento {
	fmt.Println("Inicializando fachada de almacenamiento...")
	return &FachadaAlmacenamiento{repositorio: capaaccesodatos.GetRepositorioAudios()}
}

/**
 * @brief Valida y almacena un audio mp3.
 * @param audio Datos descriptivos del audio (título, tipo y nombre de archivo).
 * @param datos Contenido binario del archivo mp3.
 * @return DTO con la confirmación, o un error (ErrAudioInvalido si los datos no son válidos).
 */
func (thisF *FachadaAlmacenamiento) GuardarAudio(audio dtos.AudioAlmacenarDTOInput, datos []byte) (dtos.AudioAlmacenadoDTOOutput, error) {
	if len(datos) == 0 {
		return dtos.AudioAlmacenadoDTOOutput{}, fmt.Errorf("%w: el archivo está vacío", ErrAudioInvalido)
	}
	if !esMP3(datos) {
		return dtos.AudioAlmacenadoDTOOutput{}, fmt.Errorf("%w: el archivo no tiene formato mp3", ErrAudioInvalido)
	}

	nombreArchivo := normalizarNombreArchivo(audio.NombreArchivo, audio.Titulo, audio.Tipo)
	if nombreArchivo == "" {
		return dtos.AudioAlmacenadoDTOOutput{}, fmt.Errorf("%w: se requiere un título o un nombre de archivo", ErrAudioInvalido)
	}

	archivo, err := thisF.repositorio.GuardarAudio(nombreArchivo, datos)
	if err != nil {
		return dtos.AudioAlmacenadoDTOOutput{}, err
	}

	return dtos.AudioAlmacenadoDTOOutput{
		Mensaje: "Audio almacenado correctamente",
		Titulo:  audio.Titulo,
		Tipo:    audio.Tipo,
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
 * @brief Construye un nombre de archivo seguro terminado en .mp3.
 *
 * Si no se indica un nombre de archivo, se construye a partir del título y del
 * tipo (titulo_tipo.mp3). Se eliminan rutas y caracteres no permitidos.
 * @param nombreArchivo Nombre sugerido por el administrador (puede estar vacío).
 * @param titulo Título del audio.
 * @param tipo Tipo del audio.
 * @return Nombre de archivo normalizado, o cadena vacía si no hay datos suficientes.
 */
func normalizarNombreArchivo(nombreArchivo string, titulo string, tipo string) string {
	nombre := strings.TrimSpace(filepath.Base(strings.ReplaceAll(nombreArchivo, "\\", "/")))
	if nombre == "" || nombre == "." || nombre == "/" {
		if strings.TrimSpace(titulo) == "" {
			return ""
		}
		nombre = strings.TrimSpace(titulo)
		if strings.TrimSpace(tipo) != "" {
			nombre += "_" + strings.TrimSpace(tipo)
		}
	}

	nombre = strings.TrimSuffix(nombre, filepath.Ext(nombre))
	nombre = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '.' {
			return r
		}
		if unicode.IsSpace(r) {
			return '_'
		}
		return -1
	}, nombre)
	nombre = strings.Trim(nombre, ".")

	if nombre == "" {
		return ""
	}
	return nombre + ".mp3"
}
