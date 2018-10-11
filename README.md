Configuration/Data.go
  - defaultConfiguration: mettre la configuration par défaut
  - struct Data: définir les configurations spécifiques d'accès aux ressources, ainsi que les paramètrages spécifiques

Configuration/Application.go
  - current.Name: Mettre le nom de l'application que l'on souhaite voir apparaître dans les infos
  - current.Version: Numéroter la version

main.go
  - Mettre le nom de l'application pour le middleware Prometheus
  
Nom du binaire
docker/Dockerfile.openshift:
  - COPY --chown=1001:runner binaryName .
  - RUN chmod +x binaryName
  - CMD ["./binaryName"]
Makefile
  - APP_NAME = binaryName
  - OC_REPO = $(REGISTRY)/$(OC_PROJECT)/binaryName:develop
  - DEV_IMAGE = binaryName
  - DEV_CONTAINER = binaryName-run
.dockerignore
  - !binaryName

binaryName = projectName ??


Cas 2:
Makefile
  - PROJECT_NAME = name
  - OC_REPO = $(REGISTRY)/$(OC_PROJECT)/$(PROJECT_NAME):develop
  - DEV_IMAGE = $(PROJECT_NAME)
  - DEV_CONTAINER = $(PROJECT_NAME)-run
