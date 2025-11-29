# go-boba

This project contains a Go Bubble Tea application and `vhs` demos.

## About

This application is a terminal-based UI that allows users to select one or more items from a scrollable list.

- Use the `up`/`down` arrow keys (or `k`/`j`) to navigate.
- Use the `-`/`=` keys to expand or shrink the selection size.
- Press `enter` or `space` to confirm your selection.

Upon confirmation, the selected items are printed to standard output as a JSON array. This allows the output to be piped to other command-line tools.

## Updating

To update all Go dependencies, run this script:

```bash
./goupdate.sh
```

## Testing and Linting

To ensure code quality, run this custom script after any successful `go test`s:

```bash
./golint.sh
```


<!-- Everything below this line is auto-generated. Manual changes will be clobbered. DO NOT EDIT. -->
## Demos

#### demo.gif
![demo.gif](demo.gif)

