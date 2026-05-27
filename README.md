# Event Booking API

A high-performance Event Booking System API built with Go, featuring JWT-based authentication, rate limiting, and an asynchronous worker pool for post-booking tasks.

## Features

- **Secure Authentication**: Uses JWT (JSON Web Tokens) for securing bookings and session lookups.
- **Quota Tracking & Concurrency**: Manages ticket capacities safely across concurrent requests.
- **Built-in Rate Limiting**: Protects endpoint routes like `/api/login` and `/api/bookings` from brute-force attempts.
- **Asynchronous Task Workers**: Employs a worker pool to handle secondary booking operations (e.g. notifications) off the main request-response cycle.
- **In-Memory Store**: Thread-safe repositories to read and write events without database bottlenecks.

## Architecture Overview

```text
                                  +------------------------------+
                                  |         HTTP Client          |
                                  +------------------------------+
                                                 |
                                                 v
                                  +------------------------------+
                                  |    Standard net/http Router  |
                                  +------------------------------+
                                                 |
                                                 v
                                  +------------------------------+
                                  |    Rate Limiter / JWT Auth   |
                                  +------------------------------+
                                                 |
                                                 v
                                  +------------------------------+
                                  |    Event Booking Handler     |
                                  +------------------------------+
                                                 |
                       +-------------------------+-------------------------+
                       | (Synchronous Response)                            | (Asynchronous Queue)
                       v                                                   v
           +-----------------------+                         +----------------------------+
           | Ticket Confirmation / |                         |  Worker Pool (5 Goroutines)|
           | In-Memory Data Update |                         |  Handles notification logs |
           +-----------------------+                         +----------------------------+
```

## Tech Stack

- **Language**: Go 1.24+
- **HTTP Routing**: Go Standard Library `net/http`
- **Security**: `github.com/golang-jwt/jwt`
- **Rate Limiting**: `golang.org/x/time/rate`
- **Deployment**: Docker / Docker Compose

## Setup

### Prerequisites
- Go 1.24+ (if running bare-metal)
- Docker & Docker Compose (if containerizing)

### Local Development
1. Fetch project dependencies:
   ```bash
   go mod download
   ```
2. Start the application:
   ```bash
   go run main.go
   ```
   The API will listen at `http://localhost:8080`.

### Container Setup
1. Build and run the service via Compose:
   ```bash
   docker-compose up --build
   ```
   The container maps internal port `8080` to your host's `8080`.

## API Routes & expected responses

### Authenticate User
* **Route**: `POST /api/login`
* **Body (JSON)**:
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword"
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```

### List Events
* **Route**: `GET /api/events`
* **Response (200 OK)**:
  ```json
  [
    {
      "id": "event-1",
      "title": "Tech Conference 2026",
      "description": "Annual tech summit",
      "date": "2026-06-15T09:00:00Z",
      "location": "Convention Center",
      "category": "Technology",
      "capacity": 100,
      "available_seats": 98,
      "price": 250000
    }
  ]
  ```

### Book Ticket
* **Route**: `POST /api/bookings`
* **Headers**: `Authorization: Bearer <token>`
* **Body (JSON)**:
  ```json
  {
    "event_id": "event-1"
  }
  ```
* **Response (201 Created)**:
  ```json
  {
    "booking_id": "book-98765",
    "status": "pending",
    "message": "Booking is being processed"
  }
  ```

### View Booking History
* **Route**: `GET /api/bookings/history`
* **Headers**: `Authorization: Bearer <token>`
* **Response (200 OK)**:
  ```json
  [
    {
      "id": "book-98765",
      "status": "confirmed",
      "event": {
        "title": "Tech Conference 2026",
        "date": "2026-06-15T09:00:00Z",
        "location": "Convention Center"
      },
      "created_at": "2026-05-27T17:03:00Z"
    }
  ]
  ```
