# go-boba

This project contains a Go Bubble Tea application and a `vhs` demo.

## Running the Demo

To run the demo, you'll need to have `vhs` installed.

### `vhs` Installation

1.  **Install `ttyd` and `ffmpeg`:**

    ```bash
    sudo apt-get update && sudo apt-get install -y ttyd ffmpeg
    ```

2.  **Install `vhs`:**

    ```bash
    go install github.com/charmbracelet/vhs@latest
    ```

3.  **Add `go/bin` to your `PATH`:**

    ```bash
    export PATH=$PATH:~/go/bin
    ```

### Generating the Demo

Once `vhs` is installed, you can generate the demo by running:

```bash
vhs demo.tape
```

