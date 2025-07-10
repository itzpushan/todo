# 📝 Todo Backend API

A fully-featured RESTful backend service to manage user accounts and their todo tasks.

Built with:
- Golang (1.24)
- Ent ORM for schema & queries
- PostgreSQL as the database
- Atlas for versioned migrations (DB schema changelog)
- JWT + bcrypt for secure authentication
- Docker + Docker Compose for containerization

---

## 📁 Project Structure

```
.
├── cmd/
│   ├── server/                 # Main API server entrypoint
│   └── seed/                   # Database seeder logic
├── internal/
│   ├── db/                     # DB connection logic
│   ├── user/                   # User routes & handlers
│   ├── todos/                  # Todo routes & handlers
│   └── middleware/             # JWT auth middleware
├── ent/                        # Ent schema & generated code
├── ent/migrate/                # Atlas migration files
├── .env                        # Environment configuration
├── Makefile
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## Setup

Add .env file and copy contents of .env.example to .env

Change the .env file according to the local environment. Not needed if running using docker-compose.

### Running Locally

```bash
make generate ./ent
make migrate
make seed
make dev
```

### Running with Docker

```bash
docker-compose up --build
```

---

# 📮 API Routes

## 🔐 Auth

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/user/signup` | Register new user |
| POST | `/api/v1/user/signin` | Login and get token |

## ✅ Todos (JWT required)

Don't forget to add `Bearer <token_id>` in Authorization field of Header

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/todos/` | Create new todo |
| GET | `/api/v1/todos/` | Get todos of first page (paginated, 5/page) |
| GET | `/api/v1/todos/?page=i` | Get all todos of i-th page (paginated, 5/page) |
| GET | `/api/v1/todos/recently-viewed` | Get recently viewed todos up to 10 |
| GET | `/api/v1/todos/{id}` | Get todo by ID |
| PUT | `/api/v1/todos/{id}` | Update todo by ID |
| DELETE | `/api/v1/todos/{id}` | Delete todo by ID |

---
