package capaaccesodatos

import (
	"os"
	"path/filepath"
	"strings"
)

// directorioAudiosPorDefecto es la ruta del almacenamiento fisico de audios que
// gestiona el ServidorDeAudios (ver diagrama: "Audio.mp3" es un almacenamiento
// compartido, no propiedad del servidor de streaming). Al ejecutar el binario
// desde ServidorDeStreaming/, esta ruta relativa apunta a la carpeta hermana
// ServidorDeAudios/audios.
const directorioAudiosPorDefecto = "../ServidorDeAudios/audios"

// directorioAudios permite sobreescribir la ruta compartida mediante la variable
// de entorno AUDIOS_DIR (util si el ServidorDeAudios corre en otra ruta o maquina).
func directorioAudios() string {
	if ruta := os.Getenv("AUDIOS_DIR"); ruta != "" {
		return ruta
	}
	return directorioAudiosPorDefecto
}

// AbrirArchivo abre el archivo del audio dado su titulo y lo devuelve como *os.File.
// Si el titulo no trae la extension .mp3, se le agrega automaticamente.
func AbrirArchivo(titulo string) (*os.File, error) {
	nombreArchivo := titulo
	if !strings.HasSuffix(strings.ToLower(nombreArchivo), ".mp3") {
		nombreArchivo = nombreArchivo + ".mp3"
	}

	rutaArchivo := filepath.Join(directorioAudios(), nombreArchivo)

	file, err := os.Open(rutaArchivo)
	if err != nil {
		return nil, err
	}

	return file, nil
}
