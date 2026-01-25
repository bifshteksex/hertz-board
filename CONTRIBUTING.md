# Contributing to HertzBoard

Thank you for considering contributing to HertzBoard! We welcome contributions from everyone.

## Code of Conduct

Please be respectful and constructive in all interactions.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/bifshteksex/hertz-board/issues)
2. If not, create a new issue using the bug report template
3. Include:
   - Clear description of the bug
   - Steps to reproduce
   - Expected vs actual behavior
   - Screenshots if applicable
   - Environment details (OS, browser, version)

### Suggesting Features

1. Check if the feature has already been suggested
2. Create a new issue using the feature request template
3. Describe:
   - The problem you're trying to solve
   - Your proposed solution
   - Alternative solutions you've considered
   - Use cases

### Pull Requests

1. Fork the repository
2. Create a new branch from `develop`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Make your changes
4. Write or update tests
5. Ensure all tests pass
6. Follow the code style guidelines
7. Commit your changes following [Conventional Commits](https://www.conventionalcommits.org/):
   ```bash
   git commit -m "feat: add new feature"
   ```
8. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
9. Create a Pull Request to the `develop` branch

## Development Setup

See [Development Setup Guide](docs/development/setup.md) for detailed instructions.

Quick start:
```bash
git clone https://github.com/yourusername/hertzboard.git
cd hertzboard
make init
```

### Commit Convention Setup

This project enforces [Conventional Commits](https://www.conventionalcommits.org/) at multiple levels:

**Local validation (Husky + Commitlint):**
- Automatically installed when you run `npm install` in the root directory
- Validates commit messages before they are created
- Prevents non-compliant commits on your local machine

**GitHub Actions:**
- Validates all commits in Pull Requests
- Checks both individual commits and PR title
- PRs with invalid commit messages will fail CI checks

**To test your commit message format:**
```bash
# This will validate your message before committing
git commit -m "feat(canvas): add new drawing tool"
```

If your commit message doesn't follow the convention, you'll see an error like:
```
⧗   input: invalid commit message
✖   subject may not be empty [subject-empty]
✖   type may not be empty [type-empty]
```

## Code Style

### Backend (Go)

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Write tests for new code
- Keep functions small and focused

### Frontend (TypeScript/Svelte)

- Use TypeScript for type safety
- Follow the Prettier configuration
- Use meaningful variable names
- Write unit tests for components
- Keep components focused on a single responsibility

## Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types (strictly enforced):**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `perf`: Performance improvements
- `ci`: CI/CD changes
- `build`: Build system changes
- `revert`: Revert previous commit

**Rules:**
- Type must be lowercase
- Subject cannot be empty
- Subject cannot end with a period
- Header max length: 100 characters
- Body and footer must have blank line before them

**Valid examples:**
```
feat(canvas): add shape rotation feature
fix(auth): resolve JWT token expiration issue
docs(api): update REST API documentation
perf(renderer): optimize canvas rendering performance
ci(actions): add commit message validation
```

**Invalid examples:**
```
Feature: add rotation        ❌ Type must be from allowed list
feat(Canvas): Add rotation   ❌ Type and scope must be lowercase
feat: add rotation.          ❌ Subject cannot end with period
feat:add rotation            ❌ Missing space after colon
add rotation                 ❌ Missing type
```

## Testing

### Backend Tests

```bash
make backend-test
```

### Frontend Tests

```bash
make frontend-test
```

### All Tests

```bash
make test
```

## Pull Request Checklist

- [ ] Code follows the project's style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Tests added/updated
- [ ] All tests passing
- [ ] No merge conflicts

## Review Process

1. At least one maintainer will review your PR
2. Address any feedback or requested changes
3. Once approved, a maintainer will merge your PR

## Questions?

Feel free to ask questions by:
- Opening an issue
- Commenting on an existing issue or PR
- Reaching out to the maintainers

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
