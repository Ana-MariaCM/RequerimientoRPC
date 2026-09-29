# Spotify distribuido con RPC (Go)

Sistema distribuido para consultar y reproducir audios (Música, Podcasts,
Audiolibros y Ruido Blanco) usando **gRPC**, **REST** y una **cola de mensajes
(RabbitMQ)**, escrito en Go para Linux.

```
 Cliente ──REST──► ServidorMetadataDeAudios           (tipos, listas y metadatos)
 Cliente ──gRPC──► ServidorDeStreaming ──cola──► ServidorDeEstadisticas
                          │  (lee Audio.mp3)
 Administrador ──REST──► ServidorDeAudios              (almacena Audio.mp3)
```

## Carpetas

| Carpeta | Rol | Comunicación | Puerto por defecto |
|---|---|---|---|
| `Cliente` | Menú de consola: ver tipos, lista, detalle y reproducir | REST + gRPC | — |
| `Administrador` | Menú de consola para almacenar audios mp3 | REST | — |
| `ServidorMetadataDeAudios` | Tipos de audio y metadatos (4 audios por tipo) | REST | 5001 |
| `ServidorDeStreaming` | Reproduce un audio por streaming y publica la reproducción en la cola | gRPC + RabbitMQ | 50051 |
| `ServidorDeAudios` | Almacena los mp3 enviados por el administrador (`audios/`) | REST | 5000 |
| `ServidorDeEstadisticas` | Consume la cola, almacena e imprime las estadísticas | RabbitMQ | — |

Cada módulo usa la misma organización por capas:

- `main/` punto de entrada.
- `capaControladores/` controladores (MVC).
- `capaFachadaServices/fachada/` fachada con la lógica de los servicios.
- `capaFachadaServices/DTOs/` DTOs que viajan en JSON (los DTOs protobuffer están en `ServidorDeStreaming/serviciosAudio`).
- `capaAccesoADatos/` repositorios (patrón singleton).
- `modelos/` modelos (MVC) y `vistas/` vistas de consola (MVC).
- `componenteConexionCola/` conexión con RabbitMQ.
- `configuracion/` parámetros que se pueden cambiar con variables de entorno.

## Requisitos (Ubuntu/Debian)

```bash
# Go 1.24.5 o superior
sudo apt install build-essential pkg-config libasound2-dev   # para compilar el Cliente (cgo + audio ALSA)
sudo apt install rabbitmq-server             # o un broker RabbitMQ accesible en la red
```

## Compilar y ejecutar

Abrir una terminal por componente. Cada componente se ejecuta **desde su propia carpeta**.

```bash
# 0. RabbitMQ en ejecución (usuario guest/guest en localhost por defecto)
sudo systemctl start rabbitmq-server

# 1. Servidor de metadatos (REST :5001)
cd ServidorMetadataDeAudios && go run ./main

# 2. Servidor de audios (REST :5000)
cd ServidorDeAudios && go run ./main

# 3. Servidor de estadísticas (consumidor de la cola)
cd ServidorDeEstadisticas && go run ./main

# 4. Servidor de streaming (gRPC :50051)
cd ServidorDeStreaming && go run ./main

# 5. Cliente
cd Cliente && go run ./main

# 6. Administrador
cd Administrador && go run ./main
```

La primera vez, `go run` descarga las dependencias indicadas en `go.mod`/`go.sum`
(`go mod tidy` también las descarga). Para generar ejecutables: `go build -o servidor ./main`.

### Configuración (variables de entorno)

| Variable | Componente | Valor por defecto |
|---|---|---|
| `PUERTO_METADATOS` | ServidorMetadataDeAudios | `5001` |
| `PUERTO_AUDIOS`, `RUTA_AUDIOS` | ServidorDeAudios | `5000`, `audios` |
| `PUERTO_STREAMING`, `RUTA_AUDIOS` | ServidorDeStreaming | `50051`, `../ServidorDeAudios/audios` |
| `RABBITMQ_URL` | Streaming y Estadísticas | `amqp://guest:guest@localhost:5672/` |
| `URL_METADATOS`, `DIRECCION_STREAMING` | Cliente | `http://localhost:5001`, `localhost:50051` |
| `URL_SERVIDOR_AUDIOS` | Administrador | `http://localhost:5000` |

Ejemplo en máquinas distintas:

```bash
RABBITMQ_URL=amqp://admin:1234@10.0.0.5:5672/ go run ./main
URL_METADATOS=http://10.0.0.2:5001 DIRECCION_STREAMING=10.0.0.3:50051 go run ./main
```

El servidor de streaming lee los mp3 de la carpeta donde el servidor de audios
los almacena (el recurso *Audio.mp3* del diagrama). Si se ejecutan en máquinas
distintas, `RUTA_AUDIOS` debe apuntar a esa carpeta compartida (por ejemplo, NFS).

## Servicios

**REST – ServidorMetadataDeAudios**

- `GET /tipos` → lista de tipos (`id`, `nombre`).
- `GET /tipos/{idTipo}/audios` → audios de un tipo.
- `GET /audios/{idAudio}` → detalle y metadatos de un audio.

**REST – ServidorDeAudios**

- `POST /audios/almacenamiento` (multipart: `archivo`, `titulo`, `tipo`, `nombreArchivo`) → almacena el mp3.
- `GET /audios` → audios almacenados.

**gRPC – ServidorDeStreaming** (`audio.proto`)

- `rpc AudioStream(AudioRequest) returns (stream AudioChunk)` → fragmentos de 32 KB.
- Regenerar los stubs: `protoc --go_out=. --go-grpc_out=. audio.proto`.

**Cola – RabbitMQ**: cola durable `reproducciones_audios`; cada reproducción se
publica en JSON desde el servidor de streaming y la consume el servidor de estadísticas.

## Audios de ejemplo

`ServidorDeAudios/audios/` trae 16 mp3 cortos generados para el laboratorio
(pistas instrumentales sintetizadas, voces sintéticas y ruido). Para usar una
canción real, se sube con el Administrador indicando el mismo nombre de archivo
que tiene en los metadatos (por ejemplo `musica_bohemian_rhapsody.mp3`).

## Documentación (Doxygen)

```bash
doxygen Doxyfile      # genera documentacion/html/index.html
```
