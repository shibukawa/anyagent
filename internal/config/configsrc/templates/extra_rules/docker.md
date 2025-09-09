# Docker Specific Rules

## Image Best Practices
- Use minimal base images (alpine, distroless)
- Pin exact versions for determinism
- Avoid running as root
- Reduce layers and clean up caches

## Build
- Use multi-stage builds
- Cache dependencies effectively
- Keep contexts small with .dockerignore

## Runtime
- Expose only necessary ports
- Use healthchecks
- Set resource limits where supported

