/**
 * @file main.go
 * @brief Punto de entrada del servidor de audios (servicios REST para el administrador).
 *
 * Servicios publicados:
 *  - POST /audios/almacenamiento   Almacena un nuevo audio mp3 (multipart/form-data) y
 *                                  registra sus metadatos en el servidor de metadatos.
 *  - GET  /audios                  Lista los audios almacenados.
 */
package main

import (
	"fmt"
	"net/http"

	controlador "almacenamiento/capaControladores"
	"almacenamiento/configuracion"
)

/**
 * @brief Registra las rutas REST y pone a escuchar el servidor HTTP.
 */
func main() {
	ctrl := controlador.NuevoControladorAlmacenamientoAudios()

	http.HandleFunc("POST /audios/almacenamiento", ctrl.AlmacenarAudio)
	http.HandleFunc("GET /audios", ctrl.ListarAudios)

	puerto := configuracion.ObtenerPuerto()
	fmt.Printf("Servidor de audios escuchando en el puerto %s (carpeta: %s)...\n", puerto, configuracion.ObtenerRutaAudios())
	if err := http.ListenAndServe(":"+puerto, nil); err != nil {
		fmt.Println("Error iniciando el servidor:", err)
	}
}
