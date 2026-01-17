# Abritary process launcher

Objective:

Monitor a path in the file system for creation of files. 
So when a file is created or updated this process will run. 

The first implementation is a folder with md files. 
This files are then processed by this process.
And create other files from them. 
Update the file in place. 
According with a predeterminate rules. 

Example:

Create articles from this note in. 

- Medium
- Linkedin
- X

Architecture:

The system follows a lightweight, event-driven architecture:

```mermaid
graph TD
    A[Launchd (macOS)] -- Watches 'notes/' --> B(Filelauncher Binary / Go)
    B -- Matches Rules --> C{Check History}
    C -- New/Updated --> D[Python Agent]
    C -- No Change --> E[Skip]
    D -- API Call --> F[OpenAI]
    F -- Generate --> G[Output Files]
```

### Execution Flow
1.  **Launchd**: macOS native daemon monitors the configured paths (e.g., `notes/`). When a file changes, it triggers the `filelauncher` binary.
2.  **Filelauncher (Go)**:
    -   Loads `config.yaml`.
    -   Checks **Dolt** history to verify if the file state has changed (deduplication).
    -   Executes the configured action.
3.  **Python Agent**:
    -   The Go binary spawns a Python process (using the local `.venv`).
    -   Runs `actions/llm_agent.py` with input/output arguments.
    -   Connects to OpenAI API using `gpt-4.1` (or configured model).
    -   Writes generated content (e.g., `.medium.md`) back to the filesystem.
User Experience:

- Users need to have an easy way to configure paths.
- It needs to be mac os native
- It can be installed with brew
- Users can create custom rules per file path / type etc.
- Everything that can be infered as a sensible default needs to be infered as a sensible defautl. 

