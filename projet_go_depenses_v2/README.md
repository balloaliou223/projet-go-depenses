# Учёт личных трат — étapes 1 à 4

Application pédagogique en Go : serveur HTTP, templates HTML, SQLite, CRUD, catégories et filtres.

## Prérequis
Go 1.22 ou plus récent. Une connexion Internet est nécessaire la première fois pour télécharger le pilote SQLite.

## Lancement
À la racine du dossier :
```bash
go mod tidy
go run .
```
Puis ouvrir http://localhost:8080. Le fichier `expenses.db` et les catégories par défaut sont créés automatiquement.

## Routes
- `GET /`, `GET /about`, `GET /ping`
- `GET /expenses` : liste et filtres `?category=2&from=2025-01-01&to=2025-01-31`
- `GET /expenses/new`, `POST /expenses` : ajout
- `GET /expenses/{id}/edit`, `POST /expenses/{id}/edit` : modification
- `POST /expenses/{id}/delete` : suppression

Le SQL est dans `internal/repository`. Migration : `internal/repository/migrations/001_create_tables.sql`.
