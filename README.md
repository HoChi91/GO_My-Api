# API de Gestion de Bibliothèque  
Une API RESTful pour gérer des livres, des achats de ces derniers, des auteurs et des utilisateurs dans une bibliothèque.  

## Installation  

### Prérequis  
- Go 1.21+  
- Docker (optionnel pour la base de données)  

### Étapes  
1. Clonez le dépôt :  
   ```bash
   git clone git@rendu-git.etna-alternance.net:module-9787/activity-52605/group-1044751.git
   cd votre-repo
   ```

2. Installez les dépendances :
    ```bash
    go mod tidy
    ````

3. Configurez le fichier ```.env``` :
    ```env
    PORT        # Port de l'API
    DB_HOST     # Nom du service db, utilisé pour la connexion à MariaDB
    DB_PORT     # Port de la db
    DB_USER     # L'utilisateur de la base de données MariaDB
    DB_PASSWORD # Le mot de passe de l'utilisateur API
    DB_NAME     # Le nom de la base de données à utiliser
    JWT_SECRET  # signature JWT
    ````

#### Sans Docker 

5. Lancez l'API :

    ```
    go run main.go
    ```


#### Avec Docker 

5. Build l'API :
    ```bash
    docker-compose up --build
    ```
 

4. Lancez l'API :

    Sur l'application ```Docker Desktop```, lancez l'application ```my-api```. 
    (qui contient les container de l'api (myapi-api-1) et de la db (myapi-db-1))

6. Accédez a l'API via :

    http://localhost:8080



## Structure du projet

```main.go``` : Point d’entrée de l’application.

```routes/``` : Contient les définitions des routes.

```controllers/``` : Gère la logique métier et la réponse des requêtes.

```models/``` : Définit les structures de données (structs) et les schémas.

```middlewares/``` : Middlewares pour les fonctionnalités comme l'authentification ou la gestion des erreurs.

```configs/``` : Fichiers de configuration, par exemple pour la base de données.

```services/``` : Logique réutilisable pour communiquer avec des bases de données ou des services externes.

## Utilisation  

**_Remarque : Les payload JSON serve uniquement d'exemple et doivent etre modifié par les données souhaitées._**

### Endpoints principaux  

#### Pour la gestion des utilisateurs

##### Routes non sécurisé

1. Pour l'incsription :

    **POST** `api/auth/signup`

    JSON payload :
    ```json
    {
        "FirstName":"test",
        "LastName":"admin",
        "Email":"test@test.com",
        "Password":"test123"
    }
    ```
    Réponse :
    ```json 
    {
        "message": "Utilisateur créer avec succès !",
        "details": {
            "ID": 0,
            "FirstName":"test",
            "LastName":"admin",
            "Email":"test@test.com",
            "Password":"test123",
            "Role": "admin"
            }
    }
    ```
    _Remarque : Le 1er utilisateur inscrit sera automatiquement un admin, les suivants des utilisateurs lambda. Pour modifier leur role, merci d'utiliser la route pour modifier un utilisateur._

2. Pour la connexion : 

    **POST** `api/auth/login`

    JSON payload : 
    ```json
    {
        "Email":"test@test.com",
        "Password":"test123"
    }
    ```
    Réponse : 
    ```json
    {
        "details": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyQHRlc3QuY29tIiwiZXhwIjoxNzMzMDA1NjMxLCJpZCI6IjQiLCJyb2xlIjoidXNlciJ9.qdypP-dpxWXK0npOuYbTWuN03PMlh_EAcCW7jGvebJ8",
        "message": "Connexion réussi !"
    }
    ```
     _Remarque : A la connexion, un token JWT d'une durée de 24h est envoyé dans les cookies._

##### Routes accessible à tout les utilisateurs connectés

3. Pour récuperer un utilisateur par ID :

    **GET** `api/users/:id`

    Réponse : 
    ```json 
    {
        "details": {
            "ID": 8,
            "FirstName": "test4",
            "LastName": "User",
            "Email": "test4@test.com",
            "Password": "",
            "Role": "user"
        },
        "message": "Utilisateur"
    }
    ```
    _Remarque : L'ID doit etre remplacé par celui souhaité, Ex : ```api/users/2```. Si l'utilisateur connecté n'est pas un admin, il peut récuperer seulement les informations de son ID._

##### Routes accessible uniquement aux admins 

4. Pour Récuperer tout les utilisateurs :

    **GET** `api/users`

    Réponse : 
    ```json
    {
        "details": [
            {
                "ID": 1,
                "FirstName": "",
                "LastName": "",
                "Email": "test@admin.com",
                "Password": "",
                "Role": ""
            },
            {
                "ID": 4,
                "FirstName": "",
                "LastName": "",
                "Email": "test2@test.com",
                "Password": "",
                "Role": ""
            },
            {
                "ID": 9,
                "FirstName": "",
                "LastName": "",
                "Email": "test2IU@test.com",
                "Password": "",
                "Role": ""
            },
            {
                "ID": 6,
                "FirstName": "",
                "LastName": "",
                "Email": "test3@test.com",
                "Password": "",
                "Role": ""
            },
            {
                "ID": 8,
                "FirstName": "",
                "LastName": "",
                "Email": "test4@test.com",
                "Password": "",
                "Role": ""
            }
        ],
        "message": "Utilisateurs"
    }
    ```
    _Remarque : Seul l'ID et l'Email sont renvoyés. La liste est non exhaustive et sert uniquement d'exemple._

5. Pour modifier un utilisateur :

    **PUT** `api/users/:id`

    JSON payload : 
    ```json 
    {
        "FirstName":"test3",
        "LastName":"test",
        "Email":"test3jj@test.com",
        "Password":"test123"
    } 
    ```

    Réponse :
    ```json
    {
        "details": {
            "ID": 0,
            "FirstName": "test3",
            "LastName": "test",
            "Email": "test3jj@test.com",
            "Password": "test123",
            "Role": ""
        },
        "message": "Utilisateur mis a jour"
    }
    ```
    _Remarque : L'ID doit etre remplacé par celui souhaité, Ex : ```api/users/3```._

6. Pour supprimer un utilisateur :

    **DELETE** 'api/users/:id`
    
    Réponse : 
    ```json
    {
        "details": null,
        "message": "Utilisateur supprimé"
    }
    ```
     _Remarque : L'ID doit etre remplacé par celui souhaité, Ex : ```api/users/2```._







#### Pour la gestion des livres

##### Routes accessible à tout les utilisateurs connectés

1. Récupérer tous les livres :

    **GET** `api/books`  

    Réponse :  
    ```json
    {
        "details": [
            {
                "ID": 3,
                "Name": "Karl Marx",
                "Birth_Date": "",
                "Description": ""
            },
            {
                "ID": 4,
                "Name": "Freud",
                "Birth_Date": "",
                "Description": ""
            }
        ],
        "message": "Auteurs"
    }
    ```
    _Remarque : La liste est non exhaustive et sert uniquement d'exemple._

2. Récuperer un livre spécifique par son ID :

    **GET** `api/books/:id`

    Réponse : 
    ```json
    [
        {
        "id": 3,
        "title": "La République",
        "author": "Platon",
        "description": "..."
        }
    ]
    ```
    _Remarque : L'ID doit etre remplacé par celui souhaité. Ex : ```api/books/3```._

##### Routes accessible uniquement aux admins

3. Pour ajouter un livre :

    **POST** `api/books`

    JSON payload:
    ```json
    {
        "title": "Le Moi et le Ça",
        "author": "Freud",
        "description": "..."
    }
    ```

    Réponse : 
    ```json
    {"message": "Livre ajouté avec succès"}
    ```

4. Modifier un livre :

    **PUT** `api/books/:id`

    JSON payload:
    ```json
    {
        "title": "Le Moi et le Ça",
        "author": "Freud",
        "description": "..."
    }
    ```

    Réponse :
    ```json
    {"message": "Livre mis à jour avec succès"}
    ```
    _Remarque : L'ID doit etre remplacé par celui souhaité. Ex : ```api/books/3```._

5. Supprimer un livre : 

    **DELETE** `api/books/:id`

    Réponse :
    ```json
    {"message": "Livre supprimé avec succès"}
    ```
    _Remarque : L'ID doit etre remplacé par celui souhaité. Ex : ```api/books/3```._

#### Pour la gestion des auteurs

##### Routes accessible à tout les utilisateurs connectés

1. Récupérer tous les auteurs :

    **GET** `api/authors`  

    Réponse :  
    ```json
    [
    {
        "id": 1,
        "name": "Victor Hugo",
        "birth_date": "1802-02-26",
        "description": "..."
    }
    {
        "id": 2,
        "name": "Molière",
        "birth_date": "1673-02-17",
        "description": "..."
    }
    {
        "id": 3,
        "name": "Platon",
        "birth_date": "-0428-00-00",
        "description": "..."
    }
    ]
    ```
    _Remarque : La liste est non exhaustive et sert uniquement d'exemple._

2. Récuperer un auteur spécifique par son ID :

    **GET** `api/authors/:id`

    Réponse : 
    ```json
    [
        {
        "id": 2,
        "name": "Molière",
        "birth_date": "1673-02-17",
        "description": "..."
        }
    ]
    ```
    _Remarque : L'ID doit etre remplacé par celui souhaité. Ex : ```api/authors/2```._

##### Routes accessible uniquement aux admins

3. Pour ajouter un auteur :

    **POST** `api/authors`

    JSON payload:
    ```json
    {
        "name": "Pierre de Marivaux",
        "birth_date": "1688-02-04",
        "description": "..."
    }
    ```

    Réponse : 
    ```json
    {"message": "Auteur ajouté avec succès"}
    ```

4. Modifier un auteur :

    **PUT** `api/authors/:id`

    JSON payload:
    ```json
    {
        "name": "Molière",
        "birth_date": "1673-02-17",
        "description": "..."
    }
    ```

    Réponse :
    ```json
    {"message": "Auteur mis à jour avec succès"}
    ```
    _Remarque : L'ID doit etre remplacé par celui souhaité. Ex : ```api/authors/2```._

5. Supprimer un auteur : 

    **DELETE** `api/authors/:id`

    Réponse :
    ```json
    {"message": "Auteur supprimé avec succès"}
    ```
    _Remarque : L'ID doit etre remplacé par celui souhaité. Ex : ```api/authors/2```._


#### Pour la gestion des achats 

##### Routes accessible à tout les utilisateurs connectés

1. Acheter un livre (sans ajout au panier):

    **POST** `api/purchase/order`

    JSON payload : 
    ```json
    {
        "BookID": 1,
        "Quantity":3
    }

    Réponse : 
    ```json
    {
        "message": "Achat effectué avec succès"
    }
    ```

2. Ajouter un livre au panier :

    **POST** `api/purchase/cart`

    JSON payload : 
    ```json
    {
        "BookID": 1,
        "Quantity":3
    }
    ```

    Réponse : 
    ```json
    {
        "message": "Article ajouté au panier"
    }
    ```

3. Supprimer le panier :

    **POST** `api/cart/remove`

    JSON payload :
    ```json
    {
        "BookID" : 1
    }
    ```

    Réponse : 
    ```json 
    {
        "message": "Article retiré du panier avec succès", "cart": votre_panier
    }
    ```
4. Finaliser l'achat du panier :

    **POST** `api/cart/finalize`

5. Pour avoir l'historique d'achat par ID : 

    **GET** `api/history/:id`

    Réponse : 
    ```json
    {
        "Id": 3,
        "User": 2,
        "Book": 4,
        "Quantity": 4,
        "Total_price": 20,
        "Payment_timestamp": "2006-01-02 15:04:05",
    }

##### Routes accessible uniquement aux admins

6. Pour avoir tout les historiques d'achat :

    **GET** `api/history`

    Réponse : 
    ```json
    [
        {
            "Id": 1,
            "User": 2,
            "Book": 3,
            "Quantity": 1,
            "Total_price": 20,
            "Payment_timestamp": "2006-02-02 17:56:34",
        }
        {
            "Id": 2,
            "User": 1,
            "Book": 2,
            "Quantity": 6,
            "Total_price": 100,
            "Payment_timestamp": "2006-05-02 15:04:05",
        }
        {
            "Id": 3,
            "User": 2,
            "Book": 9,
            "Quantity": 4,
            "Total_price": 70,
            "Payment_timestamp": "2024-10-02 04:04:10",
        }
        {
            "Id": 4,
            "User": 3,
            "Book": 4,
            "Quantity": 4,
            "Total_price": 34,
            "Payment_timestamp": "2024-11-02 16:45:24",
        }
    ]
    ```
    _Remarque : La liste est non exhaustive et sert uniquement d'exemple._

## Tests

### Pour faire les tests avec Postman

1. Ouvrir Postman 

2. Configurer la requête :

    * Type de requête : Sélectionnez par exemple POST (ou GET,PUT...) dans le menu déroulant à gauche de la barre d’adresse.
    * Entrer l'URL. Ex : http://localhost:8080/books.

3. Ajouter les en-têtes (si nécessaire)

4. Définir le corps de la requête :

    * Cliquez sur l'onglet **Body**
    * Sélectionnez l’option raw (dans les choix entre form-data, x-www-form-urlencoded, raw, etc.).
    * À droite, choisissez le type JSON dans le menu déroulant.
    * Ajoutez votre payload JSON dans la zone de texte :
    
        Exemple :

        ```json
        {
            "Email":"test@test.com",
            "Password":"test123"
        }
        ```

5. Envoyer la requête :

    * Cliquez sur Send.
    * La réponse de l’API s’affichera dans la section inférieure, avec le code HTTP (par exemple, ```200 OK``` si tout s'est bien passé).


## Info Utile 

### Documentation 
* [Go](https://go.dev/doc/) 
* [Gin](https://pkg.go.dev/github.com/gin-gonic/gin) 
* [Go MySQL Driver](https://pkg.go.dev/github.com/go-sql-driver/mysql) 
* [Package sql](https://pkg.go.dev/database/sql)

### Pour déployer sur Azure
1. Installer Azure cli (si pas déjà fait): ```https://learn.microsoft.com/en-us/cli/azure/install-azure-cli```

2. Se connecter a ton environnement azure cli :
    ```Az login``` (Choisis ton abonnement)

3. Connecte toi a ton Registre de conteneur (ACR, créé au prealable sur Azure) :
    ```Az acr login```

4. Construire l’image docker en local en arm64 :
    ```docker buildx build --platform linux/amd64 -t myapi-api:latest .```

5. Tague l’image docker que tu viens de build pour le lier a ton air :
    ```docker tag <Id de l’image que tu as créé juste au dessus> acrmyapidevwesteurope.azurecr.io/myapi-api:latest```

6. Ensuite tu push sur l’ACR :
    ```docker push acrmyapidevwesteurope.azurecr.io/myapi-api:latest```
