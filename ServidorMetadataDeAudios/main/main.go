/**
 * @file main.go
 * @brief Punto de entrada del servidor de metadatos de audios (servicios REST).
 *
 * Servicios publicados:
 *  - GET /tipos                    Lista los tipos de audio.
 *  - GET /tipos/{idTipo}/audios    Lista los audios de un tipo.
 *  - GET /audios/{idAudio}         Consulta los metadatos de un audio.
 *  - POST /audios                  Registra los metadatos de un audio nuevo
 *                                  (lo invoca el servidor de audios).
 */
package main

import (
	"fmt"
	"net/http"

	controlador "metadatos/capaControladores"
	"metadatos/configuracion"
)

/**
 * @brief Registra las rutas REST y pone a escuchar el servidor HTTP.
 */
func main() {
	ctrl := controlador.NuevoControladorMetadatos()

	http.HandleFunc("GET /tipos", ctrl.ListarTiposAudio)
	http.HandleFunc("GET /tipos/{idTipo}/audios", ctrl.ListarAudiosPorTipo)
	http.HandleFunc("GET /audios/{idAudio}", ctrl.ConsultarDetalleAudio)
	http.HandleFunc("POST /audios", ctrl.RegistrarAudio)

	puerto := configuracion.ObtenerPuerto()
	fmt.Printf("Servidor de metadatos de audios escuchando en el puerto %s...\n", puerto)
	if err := http.ListenAndServe(":"+puerto, nil); err != nil {
		fmt.Println("Error iniciando el servidor:", err)
	}
}
