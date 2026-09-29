/**
 * @file fachadaAudios.go
 * @brief Fachada que consume los servicios REST del servidor de audios.
 */
package fachada

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	dtos "administrador/capaFachadaServices/DTOs"
)

/**
 * @brief Fachada de acceso remoto al servidor de audios.
 */
type FachadaAudios struct {
	urlBase    string       ///< URL base del servidor de audios.
	clienteWeb *http.Client ///< Cliente HTTP reutilizable.
}

/**
 * @brief Crea la fachada de audios.
 * @param urlBase URL base del servidor de audios (por ejemplo http://localhost:5000).
 * @return Puntero a la nueva fachada.
 */
func NuevaFachadaAudios(urlBase string) *FachadaAudios {
	return &FachadaAudios{urlBase: urlBase, clienteWeb: &http.Client{Timeout: 2 * time.Minute}}
}

/**
 * @brief Envía un archivo mp3 al servidor de audios (POST /audios/almacenamiento).
 * @param audio Datos del audio y ruta local del archivo.
 * @return Respuesta del servidor o un error.
 */
func (thisF *FachadaAudios) AlmacenarAudio(audio dtos.AudioAlmacenarDTO) (dtos.AudioAlmacenadoDTO, error) {
	var respuesta dtos.AudioAlmacenadoDTO

	archivo, err := os.Open(audio.RutaLocal)
	if err != nil {
		return respuesta, fmt.Errorf("no se pudo abrir el archivo local: %v", err)
	}
	defer archivo.Close()

	cuerpo := &bytes.Buffer{}
	formulario := multipart.NewWriter(cuerpo)
	campoArchivo, err := formulario.CreateFormFile("archivo", filepath.Base(audio.RutaLocal))
	if err != nil {
		return respuesta, err
	}
	if _, err := io.Copy(campoArchivo, archivo); err != nil {
		return respuesta, fmt.Errorf("no se pudo leer el archivo local: %v", err)
	}
	formulario.WriteField("titulo", audio.Titulo)
	formulario.WriteField("tipo", audio.Tipo)
	formulario.WriteField("nombreArchivo", audio.NombreArchivo)
	formulario.Close()

	peticion, err := http.NewRequest(http.MethodPost, thisF.urlBase+"/audios/almacenamiento", cuerpo)
	if err != nil {
		return respuesta, err
	}
	peticion.Header.Set("Content-Type", formulario.FormDataContentType())

	err = thisF.ejecutar(peticion, http.StatusCreated, &respuesta)
	return respuesta, err
}

/**
 * @brief Consulta los audios almacenados en el servidor (GET /audios).
 * @return Lista de archivos o un error.
 */
func (thisF *FachadaAudios) ListarAudios() ([]dtos.ArchivoAudioDTO, error) {
	var respuesta []dtos.ArchivoAudioDTO

	peticion, err := http.NewRequest(http.MethodGet, thisF.urlBase+"/audios", nil)
	if err != nil {
		return nil, err
	}
	err = thisF.ejecutar(peticion, http.StatusOK, &respuesta)
	return respuesta, err
}

/**
 * @brief Ejecuta una petición HTTP y deserializa la respuesta JSON.
 * @param peticion Petición a ejecutar.
 * @param codigoEsperado Código HTTP que indica éxito.
 * @param destino Puntero donde se deserializa la respuesta.
 * @return Error de conexión, el mensaje de error del servidor, o nil.
 */
func (thisF *FachadaAudios) ejecutar(peticion *http.Request, codigoEsperado int, destino any) error {
	respuestaHTTP, err := thisF.clienteWeb.Do(peticion)
	if err != nil {
		return fmt.Errorf("no fue posible comunicarse con el servidor de audios (%s)", thisF.urlBase)
	}
	defer respuestaHTTP.Body.Close()

	if respuestaHTTP.StatusCode != codigoEsperado {
		var errorDTO dtos.RespuestaErrorDTO
		if json.NewDecoder(respuestaHTTP.Body).Decode(&errorDTO) == nil && errorDTO.Mensaje != "" {
			return fmt.Errorf("%s", errorDTO.Mensaje)
		}
		return fmt.Errorf("el servidor respondió con el código %d", respuestaHTTP.StatusCode)
	}

	if err := json.NewDecoder(respuestaHTTP.Body).Decode(destino); err != nil {
		return fmt.Errorf("respuesta inválida del servidor: %v", err)
	}
	return nil
}
