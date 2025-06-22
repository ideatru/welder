# Contributing to Welder

Thank you for your interest in contributing to Welder! We welcome contributions from the community and are excited to collaborate with you.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Guidelines](#contributing-guidelines)
- [Testing](#testing)
- [Documentation](#documentation)
- [Submitting Contributions](#submitting-contributions)
- [Release Process](#release-process)

## Code of Conduct

By participating in this project, you agree to abide by our code of conduct. Please treat all community members with respect and create a welcoming environment for everyone.

## Getting Started

### Prerequisites

- Go 1.24.0 or later
- Git
- Make (optional, for convenience commands)

### Development Setup

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/welder.git
   cd welder
   ```

3. **Add the upstream remote**:
   ```bash
   git remote add upstream https://github.com/ideatru/welder.git
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

5. **Verify the setup** by running tests:
   ```bash
   go test ./...
   ```

## Contributing Guidelines

### Types of Contributions

We welcome various types of contributions:

- **Bug Reports**: Help us identify and fix issues
- **Feature Requests**: Suggest new functionality
- **Code Contributions**: Fix bugs, implement features, or improve performance
- **Documentation**: Improve existing docs or add new documentation
- **Examples**: Add usage examples or tutorials

### Before You Start

1. **Check existing issues** to avoid duplicate work
2. **Create an issue** for significant changes to discuss the approach
3. **Keep changes focused** - one pull request per feature/fix
4. **Follow the coding standards** outlined below

### Coding Standards

#### Go Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions small and focused
- Use Go modules for dependency management

#### Code Organization

```
welder/
├── ether/          # Ethereum-specific implementations
├── types/          # Core type definitions
├── internal/       # Internal utilities (not exported)
├── examples/       # Usage examples
├── docs/           # Documentation and assets
└── *.go           # Main library files
```

#### Naming Conventions

- **Packages**: lowercase, single word when possible
- **Functions**: CamelCase for exported, camelCase for unexported
- **Variables**: camelCase
- **Constants**: CamelCase for exported, camelCase for unexported
- **Types**: CamelCase for exported, camelCase for unexported

#### Error Handling

- Always handle errors explicitly
- Use meaningful error messages
- Wrap errors with context when appropriate:
  ```go
  if err != nil {
      return fmt.Errorf("failed to serialize schema: %w", err)
  }
  ```

### Schema and Type Safety

When working with schemas and types:

- Ensure all new types support the full range of Ethereum ABI types
- Add comprehensive validation for input data
- Maintain backward compatibility when possible
- Document any breaking changes clearly

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific test
go test -run TestSpecificFunction ./package
```

### Writing Tests

- Write tests for all new functionality
- Include edge cases and error conditions
- Use table-driven tests for multiple scenarios:
  ```go
  func TestWeld(t *testing.T) {
      tests := []struct {
          name     string
          schema   types.Elements
          payload  []byte
          expected interface{}
          wantErr  bool
      }{
          // test cases
      }

      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              // test implementation
          })
      }
  }
  ```

### Test Coverage

- Aim for at least 80% test coverage for new code
- Critical paths should have 100% coverage
- Include integration tests for complete workflows

## Documentation

### Code Documentation

- Document all exported functions, types, and packages
- Use Go doc conventions:
  ```go
  // Serialize converts a Welder schema into Ethereum ABI arguments.
  // It validates the schema structure and returns an Arguments object
  // that can be used for encoding and decoding data.
  func Serialize(schema types.Elements) (*Arguments, error) {
      // implementation
  }
  ```

### README Updates

- Update README.md if you add new features
- Include usage examples for new functionality
- Keep the quick start guide current

### Examples

- Add examples in the `examples/` directory
- Include comments explaining the code
- Test all examples to ensure they work

## Submitting Contributions

### Pull Request Process

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the guidelines above

3. **Test thoroughly**:
   ```bash
   go test ./...
   go vet ./...
   ```

4. **Commit with descriptive messages**:
   ```bash
   git commit -m "feat: add support for nested array schemas

   - Implement nested array handling in schema serialization
   - Add validation for deeply nested structures
   - Include comprehensive tests for edge cases"
   ```

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a pull request** on GitHub

### Commit Message Format

Use conventional commits format:

- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `test:` for adding or modifying tests
- `refactor:` for code refactoring
- `perf:` for performance improvements
- `chore:` for maintenance tasks

### Pull Request Guidelines

- **Provide a clear description** of what your PR does
- **Reference related issues** using keywords like "Fixes #123"
- **Include tests** for new functionality
- **Update documentation** as needed
- **Keep PRs focused** - avoid mixing unrelated changes
- **Respond to feedback** promptly and professionally

### PR Template

```markdown
## Description
Brief description of the changes.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Other (please describe):

## Testing
- [ ] Tests pass locally
- [ ] Added tests for new functionality
- [ ] Updated existing tests as needed

## Documentation
- [ ] Updated README if needed
- [ ] Added/updated code comments
- [ ] Added examples if applicable

## Checklist
- [ ] My code follows the project's coding standards
- [ ] I have performed a self-review of my code
- [ ] My changes generate no new warnings
- [ ] New and existing tests pass
```

## Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

1. Update version in relevant files
2. Update CHANGELOG.md
3. Ensure all tests pass
4. Create and push version tag
5. Create GitHub release with release notes

## Getting Help

- **Questions**: Open a discussion on GitHub
- **Bugs**: Create an issue with reproduction steps
- **Feature Ideas**: Open an issue to discuss before implementing
- **Chat**: Join our community discussions

## Recognition

Contributors will be recognized in our README.md and release notes. We appreciate all contributions, big and small!

## Development Roadmap

Check our [roadmap in the README](README.md#roadmap) to see where the project is heading and how you can contribute to upcoming features:

- Cross-chain support (Q3 2025)
- TypeScript library (Q4 2025)
- Enhanced data processing (In Progress)

---

Thank you for contributing to Welder! Your efforts help make Ethereum contract interactions easier for developers worldwide. 🚀