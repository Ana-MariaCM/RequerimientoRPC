/**
 * @file main.go
 * @brief Punto de entrada del administrador, que almacena audios en el servidor de audios mediante REST.
 */
package main

import (
	capacontroladores "administrador/capaControladores"
	"administrador/capaFachadaServices/fachada"
	"administrador/configuracion"
	"administrador/vistas"
)

/**
 * @brief Construye las capas del administrador y muestra el menú principal.
 */
func main() {
	fachadaAudios := fachada.NuevaFachadaAudios(configuracion.ObtenerURLServidorAudios())
	controlador := capacontroladores.NuevoControladorAdministrador(fachadaAudios)
	vistas.MostrarMenuPrincipal(controlador)
}
