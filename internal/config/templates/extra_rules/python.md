# Python Specific Rules

## Code Standards
- Follow PEP 8 style guide
- Use type hints for function parameters and return values
- Follow naming conventions: snake_case for variables/functions, PascalCase for classes
- Use docstrings for modules, classes, and functions
- Maximum line length of 88 characters (Black formatter standard)

## Best Practices
- Use virtual environments (venv, conda, pipenv)
- Pin dependency versions in requirements.txt
- Use f-strings for string formatting
- Prefer list/dict comprehensions when readable
- Use context managers (with statements) for resource management

## Modern Python Features
- Use dataclasses for simple data containers
- Utilize pathlib for file system operations
- Use Enum for constants
- Leverage asyncio for asynchronous programming
- Use match statements (Python 3.10+) for pattern matching

## Testing
- Use pytest as the testing framework
- Write unit tests with good coverage
- Use fixtures for test setup
- Mock external dependencies
- Test error conditions and edge cases

## Project Structure
- Follow standard Python project layout
- Use __init__.py files appropriately
- Separate tests from source code
- Use setup.py or pyproject.toml for packaging
- Include requirements.txt and requirements-dev.txt

## Performance and Quality
- Use profiling tools to identify bottlenecks
- Use linting tools (flake8, pylint)
- Use code formatters (black, autopep8)
- Use type checkers (mypy, pyright)
- Consider using pre-commit hooks
