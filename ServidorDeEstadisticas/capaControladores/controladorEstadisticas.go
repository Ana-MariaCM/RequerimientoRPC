/**
 * @file controladorEstadisticas.go
 * @brief Controlador que procesa los mensajes consumidos de la cola de reproducciones.
 */
package capacontroladores

import (
	"encoding/json"
	"fmt"

	dtos "estadisticas/capaFachadaServices/DTOs"
	capafachada "estadisticas/capaFachadaServices/fachada"
	"estadisticas/vistas"
)

/**
 * @brief Controlador del servidor de estadísticas.
 */
type ControladorEstadisticas struct {
	fachada *capafachada.FachadaEstadisticas ///< Fachada que registra las reproducciones.
}

/**
 * @brief Crea el controlador de estadísticas con su fachada.
 * @return Puntero al nuevo controlador.
 */
func NuevoControladorEstadisticas() *ControladorEstadisticas {
	return &ControladorEstadisticas{fachada: capafachada.NuevaFachadaEstadisticas()}
}

/**
 * @brief Procesa el cuerpo JSON de un mensaje de reproducción.
 *
 * Deserializa el DTO, registra la reproducción mediante la fachada y actualiza
 * la vista con la reproducción y el resumen de estadísticas.
 * @param cuerpo Bytes del mensaje consumido de la cola.
 * @return Error si el mensaje no tiene un formato JSON válido.
 */
func (thisC *ControladorEstadisticas) ProcesarMensaje(cuerpo []byte) error {
	var reproduccionDTO dtos.ReproduccionDTOInput
	if err := json.Unmarshal(cuerpo, &reproduccionDTO); err != nil {
		return fmt.Errorf("mensaje con formato inválido: %v", err)
	}

	reproduccion, resumen := thisC.fachada.RegistrarReproduccion(reproduccionDTO)
	vistas.MostrarReproduccion(reproduccion)
	vistas.MostrarResumen(resumen)
	return nil
}
