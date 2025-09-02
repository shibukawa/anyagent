# TypeScript Specific Rules

## Code Standards
- Use strict TypeScript configuration (strict: true)
- Prefer interfaces for object shapes
- Use type annotations for function parameters and return types
- Follow naming conventions: PascalCase for types/interfaces, camelCase for variables/functions

## Best Practices
- Enable noImplicitAny and noImplicitReturns
- Use readonly for immutable data
- Prefer const assertions for literal types
- Use discriminated unions for type safety
- Avoid any type, use unknown when needed

## Modern TypeScript Features
- Use optional chaining (?.) and nullish coalescing (??)
- Utilize template literal types
- Use mapped types and conditional types appropriately
- Leverage utility types (Partial, Pick, Omit, etc.)

## Project Structure
- Separate types into .d.ts files or dedicated type files
- Use barrel exports (index.ts) for cleaner imports
- Follow consistent import/export patterns
- Use path mapping for cleaner imports

## Testing
- Use type-safe testing frameworks
- Test type definitions with type tests
- Use mock types for testing interfaces
- Ensure test files have proper TypeScript configuration

## Build and Tooling
- Use ESLint with TypeScript rules
- Configure prettier for consistent formatting
- Use appropriate tsconfig.json for different environments
- Consider using ts-node for development
