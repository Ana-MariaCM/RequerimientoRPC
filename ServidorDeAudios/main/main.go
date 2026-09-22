package main

import (
	controlador "almacenamiento/capaControladores"
	"fmt"
	"net/http"
)

func main() {
	ctrl := controlador.NuevoControladorAlmacenamientoCanciones()

	http.HandleFunc("/canciones/almacenamiento", ctrl.AlmacenarAudioCancion)

	fmt.Println("Servicio de almacenamiento escuchando en el puerto 5000...")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		fmt.Println("Error iniciando el servidor:", err)
	}
}
