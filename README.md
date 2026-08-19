# Sharing Vision — Backend API

REST API backend untuk aplikasi Post Article. Dibangun dengan Golang (Gin) + GORM + MySQL.

## Tech Stack
- **Language:** Golang 1.25
- **Framework:** Gin
- **ORM:** GORM v2
- **Database:** MySQL 8
- **Hosting:** Railway

## Endpoints

| Method | URL | Description |
|---|---|---|
| GET | `/health` | Health check |
| POST | `/article/` | Create article |
| GET | `/article/:limit/:offset` | List articles with pagination |
| GET | `/article/detail/:id` | Get article by ID |
| PUT | `/article/:id` | Update article |
| DELETE | `/article/:id` | Delete article |

## Validation Rules
- `title`: required, min 20 characters
- `content`: required, min 200 characters
- `category`: required, min 3 characters
- `status`: required, one of: `publish` | `draft` | `thrash`

## Local Development

### Prerequisites
- Go 1.21+
- MySQL 8 running locally
- Create database: `CREATE DATABASE article;`

### Setup
```bash
# Clone & enter directory
cd backend

# Copy env file
cp .env.example .env
# Edit .env with your MySQL credentials

# Run (auto-migrates on start)
go run main.go
```

Server starts at `http://localhost:8080`

### API Testing
Import `postman_collection.json` into Postman.  
Set `{{base_url}}` to `http://localhost:8080`.

## Deployment (Railway)

1. Push this repo to GitHub
2. New Project on Railway → Deploy from GitHub
3. Add MySQL Plugin
4. Set environment variables:
   ```
   DB_HOST=${{MYSQLHOST}}
   DB_PORT=${{MYSQLPORT}}
   DB_USER=${{MYSQLUSER}}
   DB_PASSWORD=${{MYSQLPASSWORD}}
   DB_NAME=${{MYSQLDATABASE}}
   PORT=8080
   ```
5. Generate domain in Railway dashboard
6. Update CORS in `routes/routes.go` with Vercel frontend URL
