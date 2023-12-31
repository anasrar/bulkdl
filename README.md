# bulkdl

Bulk download files.

## How to use

YAML file:

```yaml
files:
  - url: "https://example.com/file_0.txt"
  - url: "https://example.com/file_1.txt"
```

```bash
bulkdl file.yml
```

### Example

https://github.com/anasrar/.dotfiles/assets/38805204/1d638e80-a20b-41e4-b40d-46cfcee3dbe1

## YAML structure

```yaml
config: # global config
  method: "GET" # support: GET, POST, PUT, PATCH, DELETE, HEAD
  headers: # map[string]string
    apikey: "RaND0m"
  proxy: "http://127.0.0.1:9876"
  timeout: 10 # seconds
files:
  - url: "https://example.com/file_0.txt"
    filename: "change_file_name.txt" # change file name output

  - url: "https://example.com/file_1.txt" # filename will file_1.txt
    config: # override global config
      timeout: 5 # override config timeout
```

## Built with

- https://github.com/go-yaml/yaml
- https://github.com/go-zoox/fetch
- https://github.com/zenthangplus/goccm
- https://github.com/charmbracelet/bubbletea
- https://github.com/charmbracelet/lipgloss
