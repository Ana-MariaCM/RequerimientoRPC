/**
 * @file controladorAlmacenamientoAudios.go
 * @brief Controlador REST que atiende las peticiones del administrador.
 */
package capacontroladores

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	dtos "almacenamiento/capaFachadaServices/DTOs"
	capafachada "almacenamiento/capaFachadaServices/fachada"
	"almacenamiento/configuracion"
)

/**
 * @brief Controlador de los servicios REST de almacenamiento de audios.
 */
type ControladorAlmacenamientoAudios struct {
	fachada *capafachada.FachadaAlmacenamiento ///< Fachada que valida y almacena los audios.
}

/**
 * @brief Crea el controlador de almacenamiento con su fachada.
 * @return Puntero al nuevo controlador.
 */
func NuevoControladorAlmacenamientoAudios() *ControladorAlmacenamientoAudios {
	return &ControladorAlmacenamientoAudios{fachada: capafachada.NuevaFachadaAlmacenamiento()}
}

/**
 * @brief Servicio REST POST /audios/almacenamiento: almacena un nuevo audio mp3.
 *
 * Recibe un formulario multipart con los campos "archivo" (mp3), "titulo",
 * "tipo" y "nombreArchivo".
 * @param w Escritor de la respuesta HTTP.
 * @param r Petición HTTP recibida.
 */
func (thisC *ControladorAlmacenamientoAudios) AlmacenarAudio(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[REST] POST /audios/almacenamiento invocado por %s\n", r.RemoteAddr)

	r.Body = http.MaxBytesReader(w, r.Body, configuracion.TamanioMaximoAudio)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		fmt.Println("[REST] Formulario inválido:", err)
		responderError(w, http.StatusBadRequest, "El formulario es inválido o supera el tamaño máximo permitido")
		return
	}

	archivo, _, err := r.FormFile("archivo")
	if err != nil {
		fmt.Println("[REST] No se recibió el campo \"archivo\"")
		responderError(w, http.StatusBadRequest, "Error leyendo el archivo")
		return
	}
	defer archivo.Close()

	datos, err := io.ReadAll(archivo)
	if err != nil {
		responderError(w, http.StatusBadRequest, "Error leyendo el archivo")
		return
	}

	audioDTO := dtos.AudioAlmacenarDTOInput{
		Titulo:        r.FormValue("titulo"),
		Tipo:          r.FormValue("tipo"),
		NombreArchivo: r.FormValue("nombreArchivo"),
	}
	fmt.Printf("[REST] Almacenando audio \"%s\" [%s] (%d bytes)\n", audioDTO.Titulo, audioDTO.Tipo, len(datos))

	respuesta, err := thisC.fachada.GuardarAudio(audioDTO, datos)
	if errors.Is(err, capafachada.ErrAudioInvalido) {
		fmt.Println("[REST]", err)
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		fmt.Println("[REST]", err)
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Printf("[REST] POST /audios/almacenamiento -> audio almacenado como %s\n", respuesta.Archivo.NombreArchivo)
	responderJSON(w, http.StatusCreated, respuesta)
}

/**
 * @brief Servicio REST GET /audios: lista los audios almacenados.
 * @param w Escritor de la respuesta HTTP.
 * @param r Petición HTTP recibida.
 */
func (thisC *ControladorAlmacenamientoAudios) ListarAudios(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[REST] GET /audios invocado por %s\n", r.RemoteAddr)

	audios, err := thisC.fachada.ListarAudios()
	if err != nil {
		fmt.Println("[REST]", err)
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Printf("[REST] GET /audios -> se enviaron %d audios\n", len(audios))
	responderJSON(w, http.StatusOK, audios)
}

/**
 * @brief Serializa un objeto a JSON y lo escribe en la respuesta HTTP.
 * @param w Escritor de la respuesta HTTP.
 * @param codigo Código de estado HTTP.
 * @param cuerpo Objeto (DTO) a serializar.
 */
func responderJSON(w http.ResponseWriter, codigo int, cuerpo any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codigo)
	if err := json.NewEncoder(w).Encode(cuerpo); err != nil {
		fmt.Println("[REST] Error serializando la respuesta:", err)
	}
}

/**
 * @brief Escribe una respuesta de error en formato JSON.
 * @param w Escritor de la respuesta HTTP.
 * @param codigo Código de estado HTTP.
 * @param mensaje Descripción del error.
 */
func responderError(w http.ResponseWriter, codigo int, mensaje string) {
	responderJSON(w, codigo, dtos.RespuestaErrorDTOOutput{Codigo: codigo, Mensaje: mensaje})
}
