#!/usr/bin/env pwsh

Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

# Generate image and container names using the data in the "$PSScriptRoot/component.json" file
$component = Get-Content -Path "$PSScriptRoot/component.json" | ConvertFrom-Json

$protoImage = "$($component.registry)/$($component.name):$($component.version)-$($component.build)-protos"
$container = $component.name

# Remove build files
if (Test-Path "$PSScriptRoot/protos") {
    Remove-Item -Recurse -Force -Path "$PSScriptRoot/protos/*.go"
} else {
    New-Item -ItemType Directory -Force -Path "$PSScriptRoot/protos"
}

# Build docker image
docker build -f "$PSScriptRoot/docker/Dockerfile.proto" -t $protoImage .

# Create and copy compiled files, then destroy the container
docker create --name $container $protoImage
docker cp "$($container):/app/protos" "$PSScriptRoot/"
docker rm $container

# Verify that protos folder was indeed created after generating proto files
if (-not (Test-Path "$PSScriptRoot/protos")) {
    Write-Error "protos folder doesn't exist in root dir. Build failed. See logs above for more information."
}
