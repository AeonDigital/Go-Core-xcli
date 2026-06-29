GitHub CI/CD & Automation Infrastructure
================================================================

&nbsp;

> This directory contains the immutable, data-driven automation, Continuous Integration (CI), Continuous Deployment (CD), and release management infrastructure for this repository.


&nbsp;
&nbsp;


________________________________________________________________________________

## Directory Structure

```text
.github/
  ├── release/
  │     └── config.yaml             # Static cross-platform compilation configuration (GoReleaser)
  │
  ├── scripts/
  │     └── ci-cd/
  │           ├── 00-trigger-release.sh     # Local utility script to trigger releases via terminal
  │           ├── 01-load-metadata.sh       # Dynamic metadata extractor (reads go.mod and config.txt)
  │           ├── 02-run-tests.ps1          # Batch multi-module test suite executor (PowerShell)
  │           ├── 03-pre-build-verify.sh    # Pre-flight compilation validation across submodules
  │           ├── 04-check-entrypoints.sh   # Data evaluator to enable/disable GoReleaser engine
  │           └── 05-manage-version-tag.sh  # Automated semantic versioning and Git tagging engine
  │
  ├── workflows/
  │     └── ci-cd.yaml              # Static unified GitHub Actions pipeline (Orchestrator)
  │
  ├── ACTIONS.md                    # This technical documentation guide
  ├── config.txt                    # Project-specific pipeline configuration (e.g. coverage thresholds)
  └── entrypoints.txt               # Data file containing paths to main.go targets (project-specific)
```


&nbsp;
&nbsp;


________________________________________________________________________________

## Pipeline Architecture & Lifecycle

The `workflows/ci-cd.yaml` file orchestrates the complete software development lifecycle automatically through three sequential stages (Jobs):


&nbsp;


### 1. Test Suite & Security Job (`test`)

*   **Trigger:** Executes on any `push` or `pull_request` targeting the `main` branch.
*   **Core Responsibilities:**
    *   Loads dynamic metadata using `01-load-metadata.sh` to extract the Go version and Module name.
    *   Provisions the Go environment utilizing the auto-detected version from `go.mod`.
    *   **Vulnerability Scanning (`govulncheck`):** Statically analyzes the codebase on Linux to detect known security vulnerabilities within upstream dependencies.
    *   **Multi-Platform Matrix Testing:** Runs the full unit test suite via `02-run-tests.ps1` concurrently across three distinct operating systems: **Linux (Ubuntu), Windows, and macOS**, centralizing the coverage profile inside `.github/coverage.out`.
    *   **Coverage Badge & Gate:** Evaluates the `COVERAGE_THRESHOLD` from `config.txt` and automatically renders a dynamic SVG badge to the `badges` branch on Linux.


&nbsp;


### 2. Compilation Integrity Check Job (`build-check`)

*   **Trigger:** Executes on `push` events to the `main` branch, **if and only if** all checks in the preceding `test` job complete successfully.
*   **Core Responsibilities:**
    *   Executes `03-pre-build-verify.sh` to perform an integrity verification compilation (`go build`) across all local modules recursively. This acts as a critical safety gate before any deployment or version marking occurs.


&nbsp;


### 3. Release & Artifact Generation Job (`release`)

*   **Trigger:** Executes safely last, **if and only if** the preceding `build-check` verification passes with a 100% success rate.
*   **Core Responsibilities:**
    *   **Semantic Tagging:** Executes `05-manage-version-tag.sh` to calculate the next Semantic Versioning (SemVer) target based on commit history and publishes the new immutable version tag (`vX.Y.Z`) upstream.
    *   **Conditional Artifact Compilation:** Executes `04-check-entrypoints.sh` to inspect `.github/entrypoints.txt`. If valid active paths are found, it triggers GoReleaser using the static configuration file (`-f .github/release/config.yaml`) to cross-compile binary executables and publish a formal **GitHub Release** populated with distribution artifacts. If the file is empty or contains only comments, it safely skips binary compilation, treating the repository strictly as a library.


&nbsp;
&nbsp;


________________________________________________________________________________

## Semantic Versioning Strategy

The automated versioning system parses the message payload of the **latest commit** to compute the next release iteration:

*   **Repository Bootstrapping:** If no prior Git tags are registered in the repository, the pipeline initializes the version baseline automatically at `v0.0.1`.
*   **Feature Increment (`Minor`):** Commits prefixed with `feat:` (e.g., `feat: add logging system`) increment the minor version component (`v0.1.0`).
*   **Patch / Maintenance Increment (`Patch`):** Standard commits or those prefixed with `fix:`, `chore:`, or `docs:` increment the lower patch component (`v0.0.2`).
*   **Breaking Changes (`Major`):** Any commit containing the phrase `BREAKING CHANGE` within its body or footer increments the major version component (`v1.0.0`).
*   **Manual Override Configuration:** Commits prefixed explicitly with `release: vX.Y.Z` bypass the semantic heuristic engine and force the application of the designated version string.


&nbsp;
&nbsp;


________________________________________________________________________________

## Execution: How to Trigger a Release

To prevent tag pollution and overhead during continuous development, you can push multiple iterative commits to `main` without creating a version. Once you are satisfied with the state of the codebase, trigger a formal release using one of the following methods:


&nbsp;


### Option A: Local Terminal Script (Automated Empty Commit)

Execute the bundled shell script to generate an infrastructure-only empty commit that signals the deployment engine:

1.  Grant executable permissions (first-time setup only):
    ```bash
    chmod +x .github/scripts/ci-cd/00-trigger-release.sh
    ```
2.  To trigger a default automated semantic patch increment (`+1 patch`):
    ```bash
    .github/scripts/ci-cd/00-trigger-release.sh
    ```
3.  To enforce a strict manual target version:
    ```bash
    .github/scripts/ci-cd/00-trigger-release.sh v1.0.0
    ```


&nbsp;


### Option B: GitHub UI Manual Dispatch (`workflow_dispatch`)

1.  Navigate to the **Actions** tab of your repository on GitHub.
2.  Under the left-hand workflows sidebar, select **CI/CD Pipeline**.
3.  Click the grey **Run workflow** dropdown component located on the right side of the interface.
4.  Target the `main` branch and click the green **Run workflow** button to initiate deployment without generating additional commit logs.
