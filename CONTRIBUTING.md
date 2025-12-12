# Contributing to Movie App

Thank you for considering contributing to this project! To ensure code quality, consistency, and maintainability, please follow these guidelines.

## 🌿 Git Workflow

We use a simplified **Feature Branch Workflow**:

1.  **Main Branch (`main`)**: This branch contains production-ready code. **Do not push directly to main.**
2.  **Feature/Fix Branches**: Create a new branch for every task from the `main` branch.
    * **Features**: `feature/feature-name` (e.g., `feature/add-promo-module`)
    * **Bug Fixes**: `fix/bug-name` (e.g., `fix/email-sender-error`)
    * **Refactor**: `refactor/component-name` (e.g., `refactor/auth-middleware`)

### Commit Messages
We follow [Conventional Commits](https://www.conventionalcommits.org/) to keep our git history clean and readable:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code (white-space, formatting, etc)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `chore`: Maintenance tasks (e.g., updating dependencies)

**Example:** `feat: implement background worker for email reminders`

## 💻 Coding Standards

### Clean Architecture Layers
This project strictly follows **Clean Architecture**. Ensure your code resides in the correct layer:

1.  **Handler (`internal/delivery/http`)**:
    * Responsible for parsing HTTP requests (JSON binding).
    * Validates input using the Validator package.
    * Calls the UseCase layer.
    * Formats and returns HTTP responses.
    * **⛔ No business logic allowed here.**

2.  **UseCase (`internal/usecase`)**:
    * Contains purely business logic.
    * Orchestrates data flow between Repositories and external services (Mailer, etc.).
    * Agnostic of HTTP or Database implementation details.

3.  **Repository (`internal/repository`)**:
    * Handles direct Database interactions (GORM queries).
    * **⛔ No business logic allowed here.**

### Go Style Guide
- **Formatting**: Always run `go fmt ./...` before committing.
- **Naming**:
    - Use `camelCase` for local variables and unexported functions.
    - Use `PascalCase` for exported structs, interfaces, and functions.
- **Error Handling**: Always handle errors gracefully. Avoid using `_` to ignore errors unless absolutely necessary.

## 📄 Documentation (Swagger)

If you modify or add a new API Endpoint/Handler, you **MUST** update the Swagger annotations.

1.  Add `godoc` comments above the handler function.
2.  Run the generator command from the root directory:

```bash
swag init -g cmd/api/main.go --parseDependency --parseInternal
```

## 🧪 Testing
### Manual Testing
Currently, we rely on manual integration testing via Swagger UI or Postman.
1. Ensure the server runs without errors.
2. Test the happy paths (Success cases).
3. Test edge cases (Invalid input, Unauthorized access, etc.).

## 📝 Pull Request Process
1. Ensure your code builds locally: go build ./...
2. Update README.md if you introduce new environment variables or major architectural changes.
3. Submit a Pull Request (PR) to the main branch.
4. Provide a clear description of what you implemented or fixed.