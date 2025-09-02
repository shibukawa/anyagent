# Go Language Specific Rules

## Code Standards
- Follow Go conventions: use gofmt, golint, and go vet
- Use meaningful package names (short, lowercase, no underscores)
- Follow Go naming conventions (CamelCase for exported, camelCase for unexported)
- Use go modules for dependency management

## Best Practices
- Handle errors explicitly, don't ignore them
- Use interfaces to define behavior
- Prefer composition over inheritance
- Keep functions small and focused
- Use context.Context for cancellation and timeouts

## Testing
- Write table-driven tests when applicable
- Use testify/assert for better test readability
- Write benchmarks for performance-critical code
- Use go test -race to detect race conditions

## Project Structure
- Follow standard Go project layout
- Separate main packages from library packages
- Use internal/ for private packages
- Keep vendor/ out of version control if using go modules

## Performance
- Profile before optimizing
- Use sync.Pool for object reuse
- Avoid premature optimization
- Be mindful of memory allocations in hot paths

## Concurrency
- Use goroutines and channels idiomatically
- Prefer select over polling
- Use sync package primitives when appropriate
- Always clean up goroutines to prevent leaks
