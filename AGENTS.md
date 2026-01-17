# Repository Guidelines

## Project Structure & Module Organization
 
- `README.md` captures the project vision, requirements, and UX goals.
- `HUMAN.md` contains human-facing notes.
 
### System Flow
**Launchd** (Watcher) -> **Filelauncher** (Go Binary/Router) -> **LLM Agent** (Python/OpenAI).
- **Launchd**: Triggers on file modification.
- **Go Binary**: Deduplicates events using Dolt history and executes rules.
- **Python Agent**: Generates content via OpenAI API.
 
- `actions/`: Contains the LLM agent scripts.

## Build, Test, and Development Commands

- No build or runtime commands are defined in this repository today.
- Once tooling exists, list the exact commands here with short explanations (for example `make dev` for local runs, `npm test` for tests).

## Coding Style & Naming Conventions

- Current content is Markdown-only. Keep headings descriptive and use short, direct paragraphs.
- Follow the existing top-level naming pattern for contributor-facing docs (uppercase filenames such as `README.md`, `HUMAN.md`, `AGENTS.md`).
- If you introduce code, add a formatter/linter and document its usage in this section.

## Testing Guidelines

## Testing Guidelines

- **End-to-End Verification**: The primary testing method is **observational**.
  1.  Create a test file (e.g., `notes/test.md`).
  2.  Observe the filesystem for the creation of expected output files (e.g., `test.medium.md`).
  3.  Verify the content of the generated files matches the prompts/rules.
- No automated test framework is defined yet. Future tests should automate this file creation/observation loop.

## Commit & Pull Request Guidelines

## Commit & Pull Request Guidelines

- This directory IS a Git repository.
- Use concise, imperative commit summaries (for example “Add file watcher prototype”) and include rationale in the body when needed.
- For pull requests, include a short description, a checklist of user-visible changes, and any required setup or validation steps.

## Security & Configuration Tips

- The project’s intent includes filesystem monitoring and user-configurable rules. Document any required permissions or OS-specific setup as soon as implementation begins.
- Keep configuration examples in the repository (for example `config.example.yaml`) and avoid committing secrets.

## Library & Model Documentation

- When integrating external libraries or APIs (like OpenAI, LangChain, etc.), **always research and document the currently available models/options**.
- Create a dedicated section or file (e.g., in `AGENTS.md` or a specific library doc) listing the valid values (e.g., `gpt-4o`, `gpt-4o-mini`, `o1-preview`) to help users configure the tool correctly.
- **Verify model availability** (e.g. via the provider's list endpoint) before implementation to avoid errors like `model_not_found`. Do not guess model names.
- Keep this list updated as new models or versions are released.
