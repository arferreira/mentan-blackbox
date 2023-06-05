# Mentan.ai - Blackbox

## Description

This is a RESTful API built with Golang and Gin that provides endpoints to fetch data about blackboxes. It can also create, update, and delete existing blackboxes.

## Installation/Setup

1. Clone this repository to your local machine.
2. Install Golang onto your local machine following the [official installation guide](https://golang.org/doc/install).
3. Install Gin by running the command `go get -u github.com/gin-gonic/gin`.
4. Set up environment variables `BlackboxUsername` and `BlackboxPassword` on your system.
5. Run command `go run main.go` in the terminal to start the server.

## Usage

### Endpoints

| HTTP Method | Endpoint               | Request Body                                                        | Success Status Codes | Description                                              |
| ----------- | ---------------------- | ------------------------------------------------------------------- | -------------------- | -------------------------------------------------------- |
| GET         | `/api/v1/blackbox`     |                                                                     | 200 OK               | Returns a list of all blackboxes.                        |
| POST        | `/api/v1/blackbox`     | `{"name": "string", "description": "string", "is_enabled": "bool"}` | 201 Created          | Creates a new blackbox.                                  |
| GET         | `/api/v1/blackbox/:id` |                                                                     | 200 OK               | Returns details of a specific blackbox identified by ID. |
| PUT         | `/api/v1/blackbox/:id` | `{"name": "string", "description": "string", "is_enabled": "bool"}` | 200 OK               | Updates details of a specific blackbox identified by ID. |
| DELETE      | `/api/v1/blackbox/:id` |                                                                     | 204 No Content       | Deletes a specific blackbox identified by ID             |

## Contributors

- Antonio Souza <arfs.antonio@gmail.com>
