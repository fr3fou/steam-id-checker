# steam-id-checker
Simple tool that helps you check if a list of IDs are available on Steam.

## Features

- Supports piping
- Multithreaded
- Can check using scraping or API
- Has an interactive CLI

## Usage

1. Download [here](https://github.com/fr3fou/steam-id-checker/releases)
2. Open a terminal / command prompt in the same folder as the binary / .exe
3. Run it using `./Steam-ID-Checker-Linux` or `Steam-ID-Checker-Windows.exe` with the appropriate flags:

```
  -f string
    	path to the file which contains the IDs (default "example")
  -file string
    	path to the file which contains the IDs (default "example")
  -i	display an interactive prompt to check IDs - when using this mode, both taken and free IDs are printed
  -interactive
    	display an interactive prompt to check IDs - when using this mode, both taken and free IDs are printed
  -w int
    	amount of workers (default 10)
  -workers int
    	amount of workers (default 10)
```
