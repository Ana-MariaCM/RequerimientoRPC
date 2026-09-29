/**
 * @file reproduccion.go
 * @brief Modelos que representan una reproducción y las estadísticas acumuladas.
 */
package modelos

/**
 * @brief Reproducción de un audio registrada por el servidor de estadísticas.
 */
type Reproduccion struct {
	Numero           int    ///< Número consecutivo de la reproducción recibida.
	IdAudio          int    ///< Identificador del audio reproducido.
	Titulo           string ///< Título del audio.
	Tipo             string ///< Tipo de audio.
	NombreArchivo    string ///< Archivo mp3 transmitido.
	Usuario          string ///< Usuario que reprodujo el audio.
	DireccionCliente string ///< Dirección IP:puerto del cliente.
	FechaHora        string ///< Fecha y hora de la reproducción.
}

/**
 * @brief Cantidad de reproducciones asociadas a un nombre (tipo o audio).
 */
type Conteo struct {
	Nombre   string ///< Nombre del tipo o del audio.
	Cantidad int    ///< Número de reproducciones.
}

/**
 * @brief Resumen de las estadísticas acumuladas de reproducción.
 */
type ResumenEstadisticas struct {
	TotalReproducciones int      ///< Total de reproducciones registradas.
	PorTipo             []Conteo ///< Reproducciones agrupadas por tipo de audio.
	PorAudio            []Conteo ///< Reproducciones agrupadas por audio.
	PorUsuario          []Conteo ///< Reproducciones agrupadas por usuario.
}
