# ==============================================================================
# Script Name: run-tests.ps1
# Description: Recursively discovers Go modules and runs test suites,
#              consolidating coverage profiling data inside .github/ directory.
# ==============================================================================

# Ensures that the .github folder exists in the current pipeline environment.
New-Item -ItemType Directory -Force -Path ".github" | Out-Null

# Defines the absolute path to the unified coverage file.
$CoveragePath = Join-Path $env:GITHUB_WORKSPACE ".github/coverage.out"

# Cleans up previous residual executions, if any exist.
if (Test-Path $CoveragePath) { Remove-Item $CoveragePath }

# Recursively scans all folders for go.mod files.
Get-ChildItem -Recurse -Filter "go.mod" | ForEach-Object {
    $dir = $_.DirectoryName
    Write-Host "========================================="
    Write-Host "Running tests in module: $dir"
    Write-Host "========================================="

    # Navigate to the detected Go module.
    Set-Location $dir

    # Executes the native suite, saving the centralized coverage output.
    go test -v -race "-coverprofile=$CoveragePath" -covermode atomic ./...

    # Returns to the structural root of the global workspace.
    Set-Location $env:GITHUB_WORKSPACE
}
