# React Specific Rules

## Component Design
- Use functional components with hooks
- Keep components small and focused
- Follow single responsibility principle
- Use TypeScript for type safety
- Prefer composition over inheritance

## State Management
- Use useState for local component state
- Use useEffect appropriately with dependency arrays
- Consider useReducer for complex state logic
- Use context sparingly for shared state
- Implement proper state lifting patterns

## Performance Optimization
- Use React.memo for expensive components
- Implement useMemo and useCallback appropriately
- Avoid inline object/function creation in render
- Use React DevTools Profiler to identify bottlenecks
- Implement code splitting with React.lazy

## Best Practices
- Use meaningful component and prop names
- Implement proper error boundaries
- Use keys appropriately in lists
- Avoid mutating props or state directly
- Follow consistent file and folder naming

## Testing
- Write unit tests for components
- Use React Testing Library for user-centric tests
- Test component behavior, not implementation details
- Mock external dependencies appropriately
- Test accessibility features

## Development Workflow
- Use ESLint with React rules
- Configure Prettier for consistent formatting
- Use React DevTools for debugging
- Implement proper error handling
- Use Storybook for component documentation
