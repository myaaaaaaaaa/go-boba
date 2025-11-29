# go-boba

This project contains a Go Bubble Tea application and a `vhs` demo.

## Running the Demo

If you would like to run the demo, you'll need to install dependencies first.

1.  **Install Dependencies:**

    ```bash
    sudo apt-get update && sudo apt-get install -y ttyd ffmpeg
    ```

2.  **Generating the Demo:**

    ```bash
    ./vhs.sh demo.tape
    ```

## Testing and Linting

To ensure code quality, run this custom script after any successful `go test`s:

```bash
./golint.sh
```
