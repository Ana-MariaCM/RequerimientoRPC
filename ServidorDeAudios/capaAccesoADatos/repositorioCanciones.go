package capaAccesoADatos

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type RepositorioCanciones struct {
	mu sync.Mutex
}

var (
	instancia *RepositorioCanciones
	once      sync.Once
)

func GetRepositorioCanciones() *RepositorioCanciones {
	once.Do(func() {
		instancia = &RepositorioCanciones{}
	})
	return instancia
}

func (r *RepositorioCanciones) GuardarCancion(titulo string, genero string, artista string, data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	os.MkdirAll("audios", os.ModePerm)

	fileName := fmt.Sprintf("%s_%s_%s.mp3", titulo, genero, artista)
	filePath := filepath.Join("audios", fileName)

	err := os.WriteFile(filePath, data, 0644)

	if err != nil {
		return fmt.Errorf("error al guardar archivo: %v", err)
	}

	return nil
}
