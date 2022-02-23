# topi

Terminal OpenAPI documentation viewer 🐐

## About

topi is the documentation viewer for OpenAPI v3 definitions in the terminal.

> topi is still under development... 🐐

## Installation

`$ go install github.com/lusingander/topi@latest`

## Usage

`$ topi <path>`

> `path` can be local file path or remote URL.

### Keybindings

#### Common

|Key|Description|
|-|-|
|<kbd>Backspace</kbd>|back to perv page|
|<kbd>Ctrl+c</kbd>|quit|

#### List page

|Key|Description|
|-|-|
|<kbd>j</kbd>|cursor down|
|<kbd>k</kbd>|cursor up|
|<kbd>f</kbd> <kbd>l</kbd>|next page|
|<kbd>b</kbd> <kbd>h</kbd>|prev page|
|<kbd>g</kbd>|go to start|
|<kbd>G</kbd>|go to end|
|<kbd>/</kbd>|Enter filtering mode|
|<kbd>Enter</kbd>|(default) select item, (filtering) apply filter|
|<kbd>Esc</kbd>|(filtering) cancel filter, (filter applied) remove filter|

#### Detail page

|Key|Description|
|-|-|
|<kbd>j</kbd>|page down one line|
|<kbd>k</kbd>|page up one line|
|<kbd>f</kbd>|page down|
|<kbd>b</kbd>|page up|
|<kbd>d</kbd>|half page down|
|<kbd>u</kbd>|half page up|
|<kbd>Tab</kbd>|select link|
|<kbd>x</kbd>|open selecting link|

## License

MIT
