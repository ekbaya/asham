# ASHAM - ARSO Platform

## Prerequisites

- Docker and Docker Compose
- Go 1.19 or later

## Setup Instructions

### 1. Start the Database Services

First, start the PostgreSQL and Redis services using Docker Compose:

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- Redis cache on port 6379

### 2. Verify Database Connection

Check if the services are running:

```bash
docker-compose ps
```

You should see both `asham-postgres` and `asham-redis` containers running.

### 3. Run the Application

Once the database is running, start the Go application:

```bash
go run cmd/main.go
```

The application will:
1. Connect to the PostgreSQL database
2. Run migrations automatically (create all tables)
3. Seed the database with initial data
4. Start the web server on port 8080

### 4. Verify Migrations

The migrations should run automatically when the application starts. You can verify by checking the logs for:
- "Created X stages"
- "Created role: ROLE_NAME"
- "Created permission: PERMISSION_NAME"

### Troubleshooting

#### Database Connection Issues

If you get database connection errors:

1. Ensure Docker containers are running:
   ```bash
   docker-compose ps
   ```

2. Check container logs:
   ```bash
   docker-compose logs postgres
   ```

3. Test database connection:
   ```bash
   docker exec -it asham-postgres psql -U postgres -d asham -c "\dt"
   ```

#### Migration Issues

If migrations don't run:

1. Check if the database exists and is accessible
2. Verify the .env file has correct database credentials
3. Look for error messages in the application logs

### Stopping Services

To stop all services:

```bash
docker-compose down
```

To stop and remove volumes (this will delete all data):

```bash
docker-compose down -v
```

## Environment Variables

Make sure your `.env` file contains the correct database configuration:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=Eliasbaya@1223
DB_NAME=asham
SERVER_PORT=8080
```