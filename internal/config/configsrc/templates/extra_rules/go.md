# Go Language Specific Rules

## Code Standards
- Follow Go conventions: use gofmt, golint, and go vet
- Use meaningful package names (short, lowercase, no underscores)
- Follow Go naming conventions (CamelCase for exported, camelCase for unexported)
- Use go modules for dependency management

## Best Practices
- Handle errors explicitly, don't ignore them
- Use interfaces to define behavior
- Keep functions small and focused
- Prefer composition over inheritance
- Use context for cancellation and timeouts

## Testing
- Write table-driven tests
- Use testify or Go's testing package
- Mock interfaces for dependencies
- Aim for high coverage on critical code paths

