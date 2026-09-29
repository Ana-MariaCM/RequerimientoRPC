/**
 * @file cliente.go
 * @brief Punto de entrada del cliente: consulta metadatos por REST y reproduce audios por gRPC.
 */
package main

import (
	"fmt"
	"os"
	"os/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	capacontroladores "cliente.local/cliente/capaControladores"
	"cliente.local/cliente/capaFachadaServices/fachada"
	"cliente.local/cliente/configuracion"
	"cliente.local/cliente/vistas"
	pb "streaming.local/servidor-streaming/serviciosAudio"
)

/**
 * @brief Crea la conexión gRPC, construye las capas del cliente y muestra el menú principal.
 */
func main() {
	conexion, err := grpc.NewClient(configuracion.ObtenerDireccionStreaming(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("No fue posible crear la conexión con el servidor de streaming:", err)
		os.Exit(1)
	}
	defer conexion.Close()

	fachadaMetadatos := fachada.NuevaFachadaMetadatos(configuracion.ObtenerURLMetadatos())
	fachadaStreaming := fachada.NuevaFachadaStreaming(pb.NewAudioServiceClient(conexion), obtenerUsuario())

	controladorAudios := capacontroladores.NuevoControladorAudios(fachadaMetadatos)
	controladorReproduccion := capacontroladores.NuevoControladorReproduccion(fachadaStreaming)

	vistas.NuevoNavegadorVistas(controladorAudios, controladorReproduccion).MostrarMenuPrincipal()
}

/**
 * @brief Identifica al usuario como usuario\@equipo para las estadísticas.
 * @return Identificación del usuario.
 */
func obtenerUsuario() string {
	nombreUsuario := "anonimo"
	if usuarioActual, err := user.Current(); err == nil {
		nombreUsuario = usuarioActual.Username
	}
	equipo, err := os.Hostname()
	if err != nil {
		return nombreUsuario
	}
	return nombreUsuario + "@" + equipo
}
