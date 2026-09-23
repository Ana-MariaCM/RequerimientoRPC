package capacontroladores

import (
	capaFachada "servidor.local/grpc-servidor/capaFachada"
	pb "servidor.local/grpc-servidor/serviciosAudio"
)

type ControladorServidor struct {
	pb.UnimplementedAudioServiceServer
}

// Implementación del procedimiento remoto
func (s *ControladorServidor) AudioStream(req *pb.AudioRequest, stream pb.AudioService_AudioStreamServer) error {

	// Delegar la lógica de lectura y envío de chunks (fragmentos) a la fachada.
	return capaFachada.EnviarFragmentosAudio(req.Filename, stream)
}
