#!/bin/env bash
gox -os="windows darwin linux" -arch="amd64"
mv "steam-id-checker_darwin_amd64" "Steam-ID-Checker-MacOS"
mv "steam-id-checker_windows_amd64.exe" "Steam-ID-Checker-Windows.exe"
mv "steam-id-checker_linux_amd64" "Steam-ID-Checker-Linux"

