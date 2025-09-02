# Docker Specific Rules

## Dockerfile Best Practices
- Use official base images when possible
- Minimize layer count by combining RUN commands
- Use multi-stage builds for smaller production images
- Set USER to run as non-root when possible
- Use .dockerignore to exclude unnecessary files

## Image Optimization
- Order instructions from least to most frequently changing
- Use COPY instead of ADD unless specifically needed
- Clean up package manager caches in the same RUN instruction
- Use specific version tags, avoid 'latest'
- Remove unnecessary packages and files

## Security
- Scan images for vulnerabilities regularly
- Use distroless or minimal base images
- Don't include secrets in images
- Use secrets management for sensitive data
- Set appropriate file permissions

## Development Workflow
- Use docker-compose for multi-service applications
- Mount source code as volumes for development
- Use bind mounts for configuration files
- Implement health checks for services
- Use networks for service isolation

## Production Considerations
- Use init system for proper signal handling
- Implement graceful shutdown
- Set resource limits (memory, CPU)
- Use logging drivers for centralized logging
- Consider using admission controllers

## Container Registry
- Use semantic versioning for image tags
- Implement image signing and verification
- Use private registries for proprietary code
- Clean up old/unused images regularly
