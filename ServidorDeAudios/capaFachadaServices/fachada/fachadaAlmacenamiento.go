package fachada

import (
	capaAccesoADatos "almacenamiento/capaAccesoADatos"
	dtos "almacenamiento/capaFachadaServices/DTOs"
	componenteConexionCola "almacenamiento/componenteConexionCola"
	"fmt"
)

type FachadaAlmacenamiento struct {
	repo         *capaAccesoADatos.RepositorioCanciones
	conexionCola *componenteConexionCola.RabbitPublisher
}

func NuevaFachadaAlmacenamiento() *FachadaAlmacenamiento {
	fmt.Println("Inicializando fachada de almacenamiento...")

	repo := capaAccesoADatos.GetRepositorioCanciones()

	conexionCola, err := componenteConexionCola.NewRabbitPublisher()
	if err != nil {
		fmt.Println("Error al conectar con RabbitMQ: ", err)
		conexionCola = nil
	}

	return &FachadaAlmacenamiento{
		repo:         repo,
		conexionCola: conexionCola,
	}
}

func (thisF *FachadaAlmacenamiento) GuardarCancion(objCancion dtos.CancionAlmacenarDTOInput, data []byte) error {
	thisF.conexionCola.PublicarNotificacion(componenteConexionCola.NotificacionCancion{
		Titulo:  objCancion.Titulo,
		Artista: objCancion.Artista,
		Genero:  objCancion.Genero,
		Mensaje: "Nueva cancion almacenada: " + objCancion.Titulo + " de " + objCancion.Artista,
	})

	return thisF.repo.GuardarCancion(objCancion.Titulo, objCancion.Genero, objCancion.Artista, data)
}
