/**
 * @file menuAdministrador.go
 * @brief Vista por consola con el menú del administrador.
 */
package vistas

import (
	"fmt"

	capacontroladores "administrador/capaControladores"
	dtos "administrador/capaFachadaServices/DTOs"
	"administrador/utilidades"
)

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
 * @brief Opción 1: solicita el mp3, el tipo y los metadatos de un audio y lo envía al servidor de audios.
 *
 * Los metadatos (título, artista, autor, etc.) los digita el administrador. El
 * servidor de audios guarda el mp3 y registra esos metadatos en el servidor de
 * metadatos, por lo que el audio aparece en el cliente de inmediato.
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
		fmt.Printf("  %d. %s\n", indice+1, tipo.Nombre)
	}
	tipo := tiposAudio[utilidades.LeerOpcion("Seleccione el tipo: ", 1, len(tiposAudio))-1]
	audio.IdTipo = tipo.Id

	fmt.Printf("\n--- Metadatos de %s ---\n", tipo.Nombre)
	audio.Titulo = utilidades.LeerTextoObligatorio("  " + tipo.EtiquetaTitulo + ": ")
	audio.Metadatos = leerMetadatos(tipo)

	fmt.Println("\nEnviando audio al servidor de audios...")
	respuesta, err := controlador.AlmacenarAudio(audio)
	if err != nil {
		fmt.Println("No se pudo almacenar el audio:", err)
		return
	}
	fmt.Println(respuesta.Mensaje + ".")
	fmt.Printf("  Id: %d | Tipo: %s | Título: %s\n", respuesta.IdAudio, respuesta.Tipo, respuesta.Titulo)
	fmt.Printf("  Archivo: %s (%s)\n", respuesta.Archivo.NombreArchivo, utilidades.FormatearTamanio(respuesta.Archivo.TamanioBytes))
	fmt.Println("  El audio ya está disponible para el cliente.")
}

/**
 * @brief Solicita por consola los metadatos propios del tipo de audio (todos obligatorios).
 * @param tipo Tipo de audio seleccionado.
 * @return Metadatos ingresados (clave -> valor).
 */
func leerMetadatos(tipo TipoAudio) map[string]string {
	metadatos := map[string]string{}
	for _, campo := range tipo.Campos {
		mensaje := "  " + campo.Etiqueta + ": "
		if campo.Numerico {
			metadatos[campo.Clave] = utilidades.LeerNumeroObligatorio(mensaje)
		} else {
			metadatos[campo.Clave] = utilidades.LeerTextoObligatorio(mensaje)
		}
	}
	return metadatos
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
