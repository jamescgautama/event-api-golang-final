# Event Booking API

## Description & Project Summary
This project is a high-performance **Event Booking System API** built with Go. It allows users to authenticate, view available events, and book tickets. The system features built-in rate limiting to prevent abuse and an asynchronous worker pool for processing booking confirmations, ensuring a responsive user experience.

### Key Features:
- **Authentication:** JWT-based secure login.
- **Event Management:** List events with real-time quota tracking.
- **Ticket Booking:** Secure booking process with concurrency safety.
- **Rate Limiting:** Protects `/api/login` and `/api/bookings` from excessive requests.
- **Asynchronous Processing:** A worker pool handles post-booking tasks (e.g., sending notifications).
- **In-Memory Storage:** Efficient data handling using thread-safe in-memory repositories.

## Tech Stack
- **Language:** [Go](https://go.dev/) (1.24+)
- **Routing:** Standard Library `net/http`
- **Authentication:** [JWT-Go](https://github.com/golang-jwt/jwt)
- **Rate Limiting:** `golang.org/x/time`
- **Containerization:** Docker & Docker Compose

## Start Instructions

### Prerequisites
- [Go](https://go.dev/doc/install) (if running locally)
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/) (if running with Docker)

### Running Locally
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd golangfinal
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run the application:
   ```bash
   go run main.go
   ```
   The server will start at `http://localhost:8080`.

### Running with Docker
1. Build and start the containers:
   ```bash
   docker-compose up --build
   ```
2. The server will be accessible at `http://localhost:8080`.

## API Endpoints
- `POST /api/login`: Authenticate and receive a JWT.
- `GET /api/events`: List all available events.
- `POST /api/bookings`: Book a ticket for an event (Auth required).
- `GET /api/bookings/history`: View your booking history (Auth required).
