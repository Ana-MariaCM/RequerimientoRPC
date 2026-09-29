/**
 * @file controladorMetadatos.go
 * @brief Controlador REST que atiende las peticiones del cliente sobre tipos y audios.
 *
 * Cada servicio imprime un eco por pantalla (printf) indicando el servicio REST
 * invocado, tal como lo exige el requerimiento.
 */
package capacontroladores

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	dtos "metadatos/capaFachadaServices/DTOs"
	capafachada "metadatos/capaFachadaServices/fachada"
)

/**
 * @brief Controlador de los servicios REST de metadatos.
 */
type ControladorMetadatos struct {
	fachada *capafachada.FachadaMetadatos ///< Fachada que resuelve la lógica de los servicios.
}

/**
 * @brief Crea un controlador de metadatos con su fachada.
 * @return Puntero al nuevo controlador.
 */
func NuevoControladorMetadatos() *ControladorMetadatos {
	return &ControladorMetadatos{fachada: capafachada.NuevaFachadaMetadatos()}
}

/**
 * @brief Servicio REST GET /tipos: lista los tipos de audio registrados.
 * @param w Escritor de la respuesta HTTP.
 * @param r Petición HTTP recibida.
 */
func (thisC *ControladorMetadatos) ListarTiposAudio(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[REST] GET /tipos invocado por %s\n", r.RemoteAddr)

	tipos := thisC.fachada.ListarTiposAudio()

	fmt.Printf("[REST] GET /tipos -> se enviaron %d tipos de audio\n", len(tipos))
	responderJSON(w, http.StatusOK, tipos)
}

/**
 * @brief Servicio REST GET /tipos/{idTipo}/audios: lista los audios de un tipo.
 * @param w Escritor de la respuesta HTTP.
 * @param r Petición HTTP recibida (contiene el parámetro de ruta idTipo).
 */
func (thisC *ControladorMetadatos) ListarAudiosPorTipo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[REST] GET /tipos/%s/audios invocado por %s\n", r.PathValue("idTipo"), r.RemoteAddr)

	idTipo, err := strconv.Atoi(r.PathValue("idTipo"))
	if err != nil {
		fmt.Println("[REST] Identificador de tipo inválido")
		responderError(w, http.StatusBadRequest, "El identificador del tipo debe ser un número entero")
		return
	}

	audios, err := thisC.fachada.ListarAudiosPorTipo(idTipo)
	if errors.Is(err, capafachada.ErrTipoNoEncontrado) {
		fmt.Printf("[REST] El tipo %d no existe\n", idTipo)
		responderError(w, http.StatusNotFound, err.Error())
		return
	}

	fmt.Printf("[REST] GET /tipos/%d/audios -> se enviaron %d audios\n", idTipo, len(audios))
	responderJSON(w, http.StatusOK, audios)
}

/**
 * @brief Servicio REST GET /audios/{idAudio}: consulta los metadatos de un audio.
 * @param w Escritor de la respuesta HTTP.
 * @param r Petición HTTP recibida (contiene el parámetro de ruta idAudio).
 */
func (thisC *ControladorMetadatos) ConsultarDetalleAudio(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[REST] GET /audios/%s invocado por %s\n", r.PathValue("idAudio"), r.RemoteAddr)

	idAudio, err := strconv.Atoi(r.PathValue("idAudio"))
	if err != nil {
		fmt.Println("[REST] Identificador de audio inválido")
		responderError(w, http.StatusBadRequest, "El identificador del audio debe ser un número entero")
		return
	}

	detalle, err := thisC.fachada.ConsultarDetalleAudio(idAudio)
	if errors.Is(err, capafachada.ErrAudioNoEncontrado) {
		fmt.Printf("[REST] El audio %d no existe\n", idAudio)
		responderError(w, http.StatusNotFound, err.Error())
		return
	}

	fmt.Printf("[REST] GET /audios/%d -> se envió el detalle de \"%s\"\n", idAudio, detalle.Titulo)
	responderJSON(w, http.StatusOK, detalle)
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
