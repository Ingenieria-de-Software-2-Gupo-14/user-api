# go-template

## Descripción

go-template es un proyecto de ejemplo de go que utiliza Docker y Docker Compose para crear un entorno de desarrollo y despliegue.


## Guía de Uso

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
