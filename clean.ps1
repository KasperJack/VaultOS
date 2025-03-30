# Define the file where software paths are stored
$pathFile = "./software_paths.json"

# Ensure the file exists
if (-not (Test-Path $pathFile)) {
    New-Item -ItemType File -Path $pathFile -Force | Out-Null
    Set-Content $pathFile "{}" # Initialize with an empty JSON object
}

# Load the paths from the file
function Load-Paths {
    $json = Get-Content $pathFile -Raw
    return ConvertFrom-Json $json
}

# Save the paths to the file
function Save-Paths ($paths) {
    $json = $paths | ConvertTo-Json -Depth 10
    Set-Content $pathFile $json
}

# Add a path for a specific software
function Add-Path ($software, $path) {
    $paths = Load-Paths
    if (-not $paths.$software) {
        $paths.$software = @()
    }
    if ($paths.$software -notcontains $path) {
        $paths.$software += $path
        Write-Host "Added path '$path' for software '$software'."
    } else {
        Write-Host "Path '$path' already exists for software '$software'."
    }
    Save-Paths $paths
}

# Remove all paths for a specific software or all software
function Remove-Paths ($software = $null) {
    $paths = Load-Paths
    if ($software) {
        if ($paths.$software) {
            foreach ($path in $paths.$software) {
                if (Test-Path $path) {
                    Remove-Item -Recurse -Force $path
                    Write-Host "Removed path '$path'."
                } else {
                    Write-Host "Path '$path' does not exist."
                }
            }
        } else {
            Write-Host "No paths found for software '$software'."
        }
    } else {
        foreach ($software in $paths.PSObject.Properties.Name) {
            foreach ($path in $paths.$software) {
                if (Test-Path $path) {
                    Remove-Item -Recurse -Force $path
                    Write-Host "Removed path '$path'."
                } else {
                    Write-Host "Path '$path' does not exist."
                }
            }
        }
    }
}

# Main logic
param (
    [string]$Command,
    [string]$Software,
    [string]$Path
)

switch ($Command) {
    "add" {
        if (-not $Software -or -not $Path) {
            Write-Host "Usage: tool add <software> <path>"
            break
        }
        Add-Path $Software $Path
    }
    "remove" {
        if ($Software) {
            Remove-Paths $Software
        } else {
            Write-Host "Are you sure you want to remove all paths? (y/n)"
            $confirmation = Read-Host
            if ($confirmation -eq "y") {
                Remove-Paths
            } else {
                Write-Host "Operation canceled."
            }
        }
    }
    default {
        Write-Host "Unknown command. Usage:"
        Write-Host "  tool add <software> <path>"
        Write-Host "  tool remove <software>"
        Write-Host "  tool remove (to remove all paths)"
    }
}