/**
 * @file ruidoBlanco.go
 * @brief Modelo con los metadatos de un audio de tipo Ruido Blanco.
 */
package modelos

/**
 * @brief Audio de tipo Ruido Blanco.
 *
 * El nombre descriptivo del sonido se almacena en AudioBase::Titulo.
 */
type RuidoBlanco struct {
	AudioBase                  // Datos comunes (id, tipo, título y archivo).
	TipoSonido          string ///< Variante de frecuencia (Ruido Blanco, Marrón, Rosa).
	FuenteAudio         string ///< Origen del sonido (Lluvia, Ventilador, Bosque).
	UsoSugerido         string ///< Intención del audio (Dormir, Concentración, Meditación).
	ProveedorContenido  string ///< Canal o marca que genera el bucle sonoro.
	DuracionBucle       string ///< Tiempo que tarda el audio en repetirse (loop).
	FrecuenciaDominante string ///< Indica si el sonido tiende a tonos Graves o Agudos.
}

/**
 * @brief Construye la lista ordenada de metadatos del ruido blanco.
 * @return Metadatos del ruido blanco listos para ser presentados.
 */
func (thisR RuidoBlanco) ObtenerMetadatos() []Metadato {
	return []Metadato{
		{Nombre: "Nombre", Valor: thisR.Titulo},
		{Nombre: "Tipo de sonido", Valor: thisR.TipoSonido},
		{Nombre: "Fuente del audio", Valor: thisR.FuenteAudio},
		{Nombre: "Uso sugerido", Valor: thisR.UsoSugerido},
		{Nombre: "Proveedor de contenido", Valor: thisR.ProveedorContenido},
		{Nombre: "Duración del bucle", Valor: thisR.DuracionBucle},
		{Nombre: "Frecuencia dominante", Valor: thisR.FrecuenciaDominante},
	}
}
