# go-template

[![codecov](https://codecov.io/gh/Ingenieria-de-Software-2-Gupo-14/user-api/graph/badge.svg?token=E1LZ14J5QR)](https://codecov.io/gh/Ingenieria-de-Software-2-Gupo-14/user-api)

## Descripción

Api enfocada en manejo de data relacionada a usuario de una aplicacion 

Para correr esta aplciacion es necesario:

- Goland 1.24
- Docker
- DockerCompose
- Posgresql


## Guía de Uso

### ENV

```bash
PORT = "8080"
HOST = "0.0.0.0"
ENVIRONMENT = "development"

JWT_SECRET = "jwt-secret"

DATABASE_URL="postgres://postgres:postgres123@postgres:5432/postgres?sslmode=disable"

DD_CLIENT_TYPE=""
DD_SITE="us5.datadoghq.com"
DD_API_KEY=""
DD_AGENT_HOST=""
DD_AGENT_STATSD_PORT=""
EMAIL_API_KEY = ""
CHAT_GPT_KEY = ""
FCM_PROJECT_ID = ""
FIREBASE_SERVICE_ACCOUNT = ""
```

- PORT: Puerto en el que se ejecutará el servidor.
- HOST: Host en el que se ejecutará el servidor.
- ENVIRONMENT: Entorno en el que se ejecutará el servidor.
- JWT_SECRET: Secret para generar y validar tokens JWT.
- DATABASE_URL: URL de la base de datos.
  - Corriendo con docker-compose: postgres://postgres:postgres123@postgres:5432/postgres?sslmode=disable
- DD_CLIENT_TYPE: Tipo de cliente de Datadog.
  - "api": Usa la API de Datadog. Requiere DD_API_KEY.
  - "statsd","agent": Usa el agente de Datadog. Requiere DD_AGENT_HOST y DD_AGENT_STATSD_PORT.
  - Cualquier otro valor: Usa el cliente mock que no hace nada.
- DD_SITE: Sitio de Datadog.
- DD_API_KEY: API Key de Datadog.
- DD_AGENT_HOST: Host del agente de Datadog.
- DD_AGENT_STATSD_PORT: Puerto del agente de Datadog.
- EMAIL_API_KEY: API Key para el servicio de emails
- CHAT_GPT_KEY: API Key de ChatGPT
- FCM_PROJECT_ID: id del projecto en Firebase
- FIREBASE_SERVICE_ACCOUNT: secrets necesarios para el uso del sistema de messaging de firebase

### Correr local

Para correr el proyecto local:
```bash
make run
```

Para correr los tests:
```bash
make test
```

### Construir la Imagen Docker
Ejecuta el siguiente comando para construir la imagen Docker:
```bash
make docker-build
```

### Correr el Servicio
Para correr el servicio con Docker Compose:
```bash
make docker-run
```

Para detener y limpiar los contenedores:
```bash
make clean
```
