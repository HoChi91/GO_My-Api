# Utilise une image de Go comme base
FROM golang:1.23.2

# Définir le dossier de travail
WORKDIR /app

# Copier les fichiers du projet et installer les dépendances
COPY . .
RUN go mod download

# Construire l'application
RUN go build -o main .

# Vérifier si le binaire 'main' est présent dans /app
RUN ls -l /app

# Rendre le fichier binaire exécutable
RUN chmod +x /app/main

# Exposer le port de l'API
EXPOSE 4000

# Lancer l'application
CMD ["go", "run", "main.go"]