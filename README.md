# go-boba

This project contains a Go Bubble Tea application and `vhs` demos.

## About

This application is a terminal-based UI that allows users to select one or more items from a scrollable list.

- Use the `up`/`down` arrow keys (or `k`/`j`) to navigate.
- Use the `-`/`=` keys to expand or shrink the selection size.
- Press `enter` or `space` to confirm your selection.

Upon confirmation, the selected items are printed to standard output as a JSON array. This allows the output to be piped to other command-line tools.

## Updating

To update all Go dependencies, run this custom script:

```bash
./goupdate.sh
```

## Testing and Linting

The `golint.sh` script handles linting and formatting. Run it after a successful `go test`.

```bash
./golint.sh
```
