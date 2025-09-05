# ASHAM Production Deployment Guide

This guide will help you deploy the ASHAM application on a remote server using Docker.

## Prerequisites

- Docker and Docker Compose installed on your server
- Git installed on your server
- At least 2GB RAM and 10GB disk space
- Open ports: 8080 (application), 5432 (PostgreSQL), 6379 (Redis)

## Deployment Steps

### 1. Setup GitHub Actions (One-time setup)

**Configure GitHub Secrets:**

Add these secrets to your GitHub repository (Settings → Secrets and variables → Actions):

```
DOCKERHUB_USERNAME: your-dockerhub-username
DOCKERHUB_TOKEN: your-dockerhub-access-token
```

**Update Docker Compose:**

Edit `docker-compose.prod.yml` and replace `your-dockerhub-username` with your actual Docker Hub username:

```yaml
image: your-actual-username/asham:latest
```

**Note:** Docker images are automatically built and pushed to Docker Hub when you push to the main branch.

### 2. Deploy on Remote Server

**On your production server:**

```bash
git clone <your-repository-url>
cd asham
```

### 3. Configure Environment Variables

The project already includes a `.env` file with production configuration. Review and update the values as needed:

```bash
# Edit the environment file with your actual values
nano .env
```

**Important**: Verify and update these values in your `.env` file:
- `DB_PASSWORD`: Use a strong, unique password
- `AZURE_*`: Your Azure AD application credentials
- `GOOGLE_CLIENT_TOKEN`: Your Google API token
- `EMAIL_*`: Your SMTP email configuration

### 4. Update Docker Compose Configuration

Edit `docker-compose.prod.yml` and update the image name and password:

```bash
nano docker-compose.prod.yml
```

- Replace `your-dockerhub-username/asham:latest` with your actual Docker Hub image
- Replace `your_secure_password_here` with the same password you set in the `.env` file

### 5. Deploy the Application

```bash
# Pull the latest image
docker-compose -f docker-compose.prod.yml pull

# Start all services
docker-compose -f docker-compose.prod.yml up -d

# Check if all services are running
docker-compose -f docker-compose.prod.yml ps

# View logs to ensure everything started correctly
docker-compose -f docker-compose.prod.yml logs -f
```

### 5. Verify Deployment

```bash
# Check application health
curl http://localhost:8080/health

# Check database connection
docker exec asham-postgres-prod psql -U postgres -d asham -c "\dt"

# Check Redis connection
docker exec asham-redis-prod redis-cli ping
```

## Management Commands

### Start Services
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### Stop Services
```bash
docker-compose -f docker-compose.prod.yml down
```

### View Logs
```bash
# All services
docker-compose -f docker-compose.prod.yml logs -f

# Specific service
docker-compose -f docker-compose.prod.yml logs -f app
```

### Update Application
```bash
# Push changes to main branch (triggers automatic Docker build)
git push origin main

# On production server: pull and restart
docker-compose -f docker-compose.prod.yml pull app
docker-compose -f docker-compose.prod.yml up -d app
```

### Backup Database
```bash
docker exec asham-postgres-prod pg_dump -U postgres asham > backup_$(date +%Y%m%d_%H%M%S).sql
```

### Restore Database
```bash
docker exec -i asham-postgres-prod psql -U postgres asham < backup_file.sql
```

## Security Considerations

1. **Firewall**: Configure your server firewall to only allow necessary ports
2. **SSL/TLS**: Use a reverse proxy (nginx) with SSL certificates
3. **Environment Variables**: Never commit real credentials to version control
4. **Database**: Use strong passwords and consider restricting database access
5. **Updates**: Regularly update Docker images and dependencies

## Troubleshooting

### Application Won't Start
```bash
# Check logs
docker-compose -f docker-compose.prod.yml logs app

# Check if database is ready
docker-compose -f docker-compose.prod.yml logs postgres
```

### Database Connection Issues
```bash
# Test database connectivity
docker exec asham-app ping postgres

# Check database status
docker exec asham-postgres-prod pg_isready -U postgres
```

### Performance Issues
```bash
# Monitor resource usage
docker stats

# Check disk space
df -h

# Check memory usage
free -h
```

## Production Recommendations

1. **Reverse Proxy**: Use nginx or Traefik for SSL termination and load balancing
2. **Monitoring**: Set up monitoring with Prometheus and Grafana
3. **Backups**: Implement automated database backups
4. **Logging**: Configure centralized logging with ELK stack or similar
5. **Scaling**: Consider using Docker Swarm or Kubernetes for scaling

## Support

For issues or questions, please check the application logs and refer to the main README.md file.