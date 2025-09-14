# Play with UTF-8

A versatile Go application for encoding and decoding UTF-8 characters, designed to work as both a command-line tool and a web service. This project serves as a hands-on exploration of UTF-8 character encoding and HTTP server fundamentals in Go.

---

## Features

- **Dual-Mode Operation:** Run the application as an interactive command-line interface or as a self-contained web server.
- **Custom UTF-8 Implementation:** The core logic for encoding and decoding UTF-8 is implemented from scratch, without external libraries.
- **Single-Binary Deployment:** All necessary static assets (the HTML file for the web interface) are embedded directly into the Go binary.

---

## Prerequisites

- **Go:** Version 1.16 or newer.

---

## Project Structure

```
/main.go          # Contains the main.go file, the entry point for the application
/cli/             # Handles the interactive command-line interface logic
/services/        # Contains the core Utf8Encode and Utf8Decode functions
/controllers      # Manages the web server, including API routes and controllers
templates/        # Stores the index.html file that is embedded into the final binary
```

---

## Usage

The application can be run in two distinct modes:

### 1. Interactive CLI Mode

This is the default mode. The application will start an interactive session where you can enter commands directly.

To run the application in CLI mode, simply execute:

```sh
go run main.go
```

Once the prompt (`>`) appears, you can use the following commands:

- **Encode a string:**
    ```
    > --encode "Hello, 世界!"
    ```
- **Decode a hexadecimal byte sequence:**
    ```
    > --decode "48,65,6c,6c,6f,2c,20,e4,b8,96,e7,95,8c,21"
    ```

---

### 2. Web Server Mode

This mode starts an HTTP server that exposes a single API endpoint and serves a simple web page.

To run the application in web server mode, use the `--serve` flag:

```sh
go run main.go --serve
```

The server will start on [http://localhost:8080](http://localhost:8080). You can interact with it using a web browser or a tool like `curl`.

---

## API Endpoint

The server exposes a single `POST` endpoint for both encoding and decoding.

- **URL:** `http://localhost:8080/playwithutf`
- **Method:** `POST`
- **Content-Type:** `application/json`

### Request Body Schema

The request body should be a JSON object with two fields:

- `operation` (string): Must be either `"encode"` or `"decode"`.
- `input` (string): The string to encode or a comma-separated hexadecimal string to decode.

Example:

```json
{
  "operation": "encode",
  "input": "Hello, 世界!"
}
```

### Response Body Schema

The response will be a JSON object containing the result:

- `success` (boolean): `true` if the operation was successful.
- `code` (int): The HTTP status code.
- `data` (object): The result of the operation, including the input and output.
- `error` (string, optional): An error message if the operation failed.

Example:

```json
{
  "success": true,
  "code": 200,
  "data": {
    "operation": "encode",
    "input": "Hello, 世界!",
    "output": "SGVsbG8sIOS4lueVjCE="
  }
}
```

---

## Technical Details

- The application leverages Go's standard library to achieve its functionality.
- The `flag` package is used to handle command-line arguments and switch between modes.
- The `http` package powers the web server.
- The `encoding/json` package is used for handling API requests and responses.
- The `embed` package is used to compile the `index.html` file directly into the final Go binary, creating a single, self-contained executable.
