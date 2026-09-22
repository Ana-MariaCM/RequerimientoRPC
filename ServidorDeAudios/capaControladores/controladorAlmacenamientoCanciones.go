package capacontroladores

import (
	dtos "almacenamiento/capaFachadaServices/DTOs"
	capafachada "almacenamiento/capaFachadaServices/fachada"
	"fmt"
	"io"
	"net/http"
)

type ControladorAlmacenamientoCanciones struct {
	fachada *capafachada.FachadaAlmacenamiento
}

func NuevoControladorAlmacenamientoCanciones() *ControladorAlmacenamientoCanciones {
	return &ControladorAlmacenamientoCanciones{
		fachada: capafachada.NuevaFachadaAlmacenamiento(),
	}
}

func (thisC *ControladorAlmacenamientoCanciones) AlmacenarAudioCancion(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Almacenando cancion...\n")
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(50 << 20)
	file, _, err := r.FormFile("archivo")
	if err != nil {
		http.Error(w, "Error leyendo el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()
	data, _ := io.ReadAll(file)

	dto := dtos.CancionAlmacenarDTOInput{
		Titulo:  r.FormValue("titulo"),
		Artista: r.FormValue("artista"),
		Genero:  r.FormValue("genero"),
	}
	thisC.fachada.GuardarCancion(dto, data)
}
