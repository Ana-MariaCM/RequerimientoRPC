/**
 * @file menuAdministrador.go
 * @brief Vista por consola con el menú del administrador.
 */
package vistas

import (
	"fmt"
	"path/filepath"

	capacontroladores "administrador/capaControladores"
	dtos "administrador/capaFachadaServices/DTOs"
	"administrador/utilidades"
)

/** @brief Tipos de audio que el administrador puede seleccionar. */
var tiposAudio = []string{"Música", "Podcasts", "Audiolibros", "Ruido Blanco"}

/**
 * @brief Muestra el menú principal del administrador hasta que elija salir.
 * @param controlador Controlador que atiende las opciones del menú.
 */
func MostrarMenuPrincipal(controlador *capacontroladores.ControladorAdministrador) {
	for {
		fmt.Println("\n==================================")
		fmt.Println("  Administrador - Servidor de audios")
		fmt.Println("==================================")
		fmt.Println("  1. Almacenar audio")
		fmt.Println("  2. Ver audios almacenados")
		fmt.Println("  3. Salir")

		switch utilidades.LeerOpcion("Seleccione una opción: ", 1, 3) {
		case 1:
			opcionAlmacenarAudio(controlador)
		case 2:
			opcionVerAudiosAlmacenados(controlador)
		case 3:
			fmt.Println("Hasta pronto.")
			return
		}
	}
}

/**
 * @brief Opción 1: solicita los datos de un audio y lo envía al servidor de audios.
 * @param controlador Controlador del administrador.
 */
func opcionAlmacenarAudio(controlador *capacontroladores.ControladorAdministrador) {
	fmt.Println("\n--- Almacenar audio ---")
	audio := dtos.AudioAlmacenarDTO{}
	rutaLocal, err := controlador.ValidarArchivoLocal(utilidades.LeerTexto("Ruta del archivo mp3: "))
	if err != nil {
		fmt.Println("No se puede almacenar el audio:", err)
		return
	}
	audio.RutaLocal = rutaLocal

	fmt.Println("Tipo de audio:")
	for indice, tipo := range tiposAudio {
		fmt.Printf("  %d. %s\n", indice+1, tipo)
	}
	audio.Tipo = tiposAudio[utilidades.LeerOpcion("Seleccione el tipo: ", 1, len(tiposAudio))-1]
	audio.Titulo = utilidades.LeerTexto("Título del audio: ")
	audio.NombreArchivo = utilidades.LeerTexto(fmt.Sprintf("Nombre en el servidor (Enter = %s): ", filepath.Base(audio.RutaLocal)))

	respuesta, err := controlador.AlmacenarAudio(audio)
	if err != nil {
		fmt.Println("No se pudo almacenar el audio:", err)
		return
	}
	fmt.Printf("%s: %s (%s)\n", respuesta.Mensaje, respuesta.Archivo.NombreArchivo,
		utilidades.FormatearTamanio(respuesta.Archivo.TamanioBytes))
}

/**
 * @brief Opción 2: muestra los audios almacenados en el servidor de audios.
 * @param controlador Controlador del administrador.
 */
func opcionVerAudiosAlmacenados(controlador *capacontroladores.ControladorAdministrador) {
	audios, err := controlador.ListarAudios()
	if err != nil {
		fmt.Println("No se pudo consultar los audios:", err)
		return
	}

	fmt.Printf("\n--- Audios almacenados (%d) ---\n", len(audios))
	for indice, audio := range audios {
		fmt.Printf("  %2d. %-45s %10s\n", indice+1, audio.NombreArchivo, utilidades.FormatearTamanio(audio.TamanioBytes))
	}
}
