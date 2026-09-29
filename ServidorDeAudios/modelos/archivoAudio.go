/**
 * @file archivoAudio.go
 * @brief Modelo que representa un archivo mp3 almacenado en el servidor de audios.
 */
package modelos

/**
 * @brief Archivo de audio almacenado en disco.
 */
type ArchivoAudio struct {
	NombreArchivo string ///< Nombre del archivo mp3.
	TamanioBytes  int64  ///< Tamaño del archivo en bytes.
}
