# zui

A minimalistic terminal user interface for [Zotero](https://www.zotero.org/).

![demo](https://github.com/user-attachments/assets/e3913f93-cf5f-40d3-a41c-20d6883410e2)


## Install

```bash
go install github.com/camilo-zuluaga/zui@latest
```

Or download a prebuilt binary from [Releases](https://github.com/camilo-zuluaga/zui/releases).

Or build from source:

```bash
git clone https://github.com/camilo-zuluaga/zui.git
cd zui
go build -o zui .
```

## Setup

On first run, zui will prompt you for:

1. **API Key** — Create one at [zotero.org/settings/keys/new](https://www.zotero.org/settings/keys/new)
2. **User ID** — Found at [zotero.org/settings/security](https://www.zotero.org/settings/security#applications)

Credentials are stored securely in your system keyring.

## Keybindings

### Collections

| Key | Action |
|---|---|
| `enter` | Open collection |
| `s` | Search items |
| `/` | Filter collections |
| `ctrl+r` | Re-fetch collections |
| `q` | Quit |

### Items

| Key | Action |
|---|---|
| `enter` | Load item details |
| `tab` | Switch pane |
| `n` | Create/edit note |
| `r` | Open PDF |
| `b` | Copy bibliography |
| `/` | Filter items |
| `ctrl+r` | Re-fetch items |
| `esc` | Back |
| `q` | Quit |

## Config

There will be an optional config at `~/.config/zui/config.toml` _that needs to be created_.

#### Format & Style
The formats available for exporting are the following:
- `bibtex`, `biblatex`, `bookmarks`, `coins`, `csljson`, `csv`, `mods`, `refer`, `rdf_bibliontology`, `rdf_dc`, `rdf_zotero`, `ris`, `tei`, `wikipedia`

The styles are found at https://www.zotero.org/styles/

#### Max-items
By default the application will go up to maximum 200 items to fetch and show in the tui, you can modify this by editing this field in the config.

#### Note editor
By default, zui will use the built in note editor from bubbletea, you can modify to use your terminal editor like `nvim`, `nano`, etc.

```toml
# default config
format = "biblatex"
style = "apa"
max-items = 200
note-editor = ""
```

If no config file exists, defaults to `biblatex` format with `apa` style.

## Notes
- This is my first project with go, there is a lot that can be done better, and probably (highly) errors will come out.
- Main limitation is that zui uses the official API from zotero to fetch your items, this means it won't show any items that are not from zotero API's response.
- There is a possibility to try and implement an offline version since zotero uses an sqlite database, although this would require quite some work.
- [The main inspiration for this project](https://github.com/jbaiter/zotero-cli).

#### Built with

[bubbletea](https://github.com/charmbracelet/bubbletea), [bubbles](https://github.com/charmbracelet/bubbles), [lipgloss](https://github.com/charmbracelet/lipgloss)
