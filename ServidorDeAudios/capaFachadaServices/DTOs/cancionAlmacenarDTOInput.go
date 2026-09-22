package dtos

type CancionAlmacenarDTOInput struct {
	Titulo  string `json: "titulo"`
	Artista string `json: "artista"`
	Genero  string `json: genero`
}
