````markdown
# Football League Simulation API

This project is a REST API developed in Go, designed to simulate a football league and predict team championship probabilities using advanced algorithms and concurrent processing.

## Features

* **League Table**: Displays the current league standings, including goal differences and championship probability predictions.
* **Weekly Match Simulation**: Simulates matches for the current week and updates team standings accordingly.
* **Full League Simulation**: Automatically simulates all remaining weeks to complete the season.
* **League Reset**: Resets team statistics and match fixtures to start a new season.
* **Championship Predictions**: Utilizes Monte Carlo simulation algorithms, executed concurrently using Go's goroutines for multithreaded performance, to calculate championship probabilities based on remaining matches.
* **Advanced Logging**: Implements detailed logging for info, warnings, and errors to streamline development and debugging. Log output is directed to both the console and an `app.log` file in the project root.
* **Automated Database Setup and Migration**: Automatically checks and creates the required database schema and tables on application startup. Additionally, initial fixture data is automatically populated into the database as needed.

## Prerequisites

To set up and run this project locally, ensure you have the following installed:

* **Go**: Version 1.20 or higher. [Install Go](https://go.dev/doc/install)
* **Docker Desktop**: Required to run the MSSQL database. [Download Docker Desktop](https://www.docker.com/products/docker-desktop/)
* **Git**: Used to clone the project repository. [Download Git](https://git-scm.com/downloads)

## Installation and Setup

Follow these steps to set up and run the project locally:

### 1. Clone the Repository

First, clone the project repository to your local machine:

```bash
git clone [https://github.com/muzaffertuna/football-league-sim.git](https://github.com/muzaffertuna/football-league-sim.git)
cd football-league-sim
````

### 2\. Configure Environment Variables

The API uses environment variables for configuration, such as database connection details and server address. Create a `.env` file in the project root directory (`football-league-sim`) and add the following content, replacing the password with a strong one of your choice:

```ini
SA_PASSWORD="YourStrongPassword123" # <<< REPLACE WITH YOUR OWN STRONG PASSWORD!
DB_CONNECTION_STRING="sqlserver://sa:${SA_PASSWORD}@localhost:1433?database=FootballLeagueSim&TrustServerCertificate=true"
SERVER_ADDRESS=":8080"
```

**Important**: Ensure the `SA_PASSWORD` value in your `.env` file **exactly matches** the strong password you will use for the MSSQL Server being brought up by Docker. This password will be used by both the Docker container and your Go application to connect to the database.

### 3\. Set Up and Start the Database

The project includes a `docker-compose.yml` file in the root directory. This file is configured to run an MSSQL database container and automatically uses the `SA_PASSWORD` from your `.env` file for database initialization.

Start the database container using Docker Compose:

```bash
docker-compose up -d
```

This command will download the Docker image (if not already present locally), start the MSSQL Server using your specified `SA_PASSWORD`, and map port `1433` from your host machine to the container's `1433` port. Please note that it may take a few minutes for the database to fully start and complete its initial setup.

### 4\. Run the Go Application

With the database running and your `.env` file correctly configured, you can now start the API application:

```bash
# Ensure you are in the project's root directory (football-league-sim)
cd cmd/api # Navigate to the directory containing your main.go application file

go run main.go
```

Upon successful startup, you should see output in your console similar to:

```
2025/05/27 16:07:56 logger.go:39: [INFO] Successfully connected to MSSQL
2025/05/27 16:07:56 logger.go:39: [INFO] Database initialization and migration complete.
2025/05/27 16:07:56 logger.go:39: [INFO] Starting server on :8080
```

The "Database initialization and migration complete." line confirms that the application has successfully set up the database schema and initial data.

### 5\. Accessing and Testing the API

The API will be available locally at `http://localhost:8080`. You can interact with its endpoints using one of the following methods:

#### Method 1: Swagger UI (Recommended)

Swagger UI provides a user-friendly, interactive interface to explore and test all API endpoints directly in your web browser. This is the easiest way to get started and understand the API's functionality.

  * Open your web browser and navigate to: `http://localhost:8080/swagger/index.html#/`

  * You will be greeted with an interface similar to the one shown below:

  * From this interface, you can expand each endpoint to view its details, click the "Try it out" button to send requests, and instantly review the responses.

#### Method 2: Postman, cURL, or Similar Tools

Alternatively, you can use popular API testing tools like Postman, Insomnia, or command-line utilities such as `cURL` to send requests to the API endpoints.

## API Endpoints and Usage Flow

For the application to function correctly and avoid database errors, you **must initialize the league** using the `POST /reset-league` endpoint before performing any other simulation operations for the first time.

### `POST /reset-league`

  * **Description**: Resets all team statistics and match fixtures to initiate a new season. **This endpoint must be called before other simulation operations on first use to ensure the database is properly initialized with fixture data.** Subsequent uses can proceed without resetting if you wish to continue the current simulation.
  * **cURL Example**:
    ```bash
    curl -X POST http://localhost:8080/reset-league
    ```

### `POST /play-week`

  * **Description**: Simulates matches for the current week and updates team standings accordingly. Each call advances the league to the next week.
  * **cURL Example**:
    ```bash
    curl -X POST http://localhost:8080/play-week
    ```

### `GET /league-table`

  * **Description**: Retrieves the current league standings, including goal differences and championship probabilities calculated using the multithreaded Monte Carlo simulations.
  * **cURL Example**:
    ```bash
    curl -X GET http://localhost:8080/league-table
    ```

### `POST /simulate-all-weeks`

  * **Description**: Automatically simulates all remaining weeks to complete the entire season.
  * **cURL Example**:
    ```bash
    curl -X POST http://localhost:8080/simulate-all-weeks
    ```

## Simulation Scenarios

After successfully running the API, you can try the following simulation scenarios:

1.  **Initialize the League**: Start by sending a request to the `POST /reset-league` endpoint to set up a new season.
2.  **Advance Week by Week**: Send successive requests to the `POST /play-week` endpoint to simulate matches week by week. Observe how the league standings evolve after each week.
3.  **Check League Standings**: After simulating a week, call the `GET /league-table` endpoint to see the updated standings and championship probabilities. Pay attention to how the predictions change as the league progresses.
4.  **Complete the Season**: If you wish to quickly finish the remaining part of the league, use the `POST /simulate-all-weeks` endpoint. This will automatically play out all remaining matches.

## Code Snippets

For quick reference and copying, here are the key code snippets mentioned in the setup:

### Environment Variables (`.env`) Example

```ini
SA_PASSWORD="YourStrongPassword123"
DB_CONNECTION_STRING="sqlserver://sa:${SA_PASSWORD}@localhost:1433?database=FootballLeagueSim&TrustServerCertificate=true"
SERVER_ADDRESS=":8080"
```

### Clone the Repository

```bash
git clone [https://github.com/muzaffertuna/football-league-sim.git](https://github.com/muzaffertuna/football-league-sim.git)
cd football-league-sim
```

### Start the Database

```bash
docker-compose up -d
```

### Run the Application

```bash
cd cmd/api
go run main.go
```

## Notes

  * Ensure the `.env` file is correctly configured with your chosen password before starting both the database and the application.
  * The `POST /reset-league` endpoint is crucial for the first-time setup to avoid database errors and populate initial league data.
  * Leverage Swagger UI for a seamless and interactive API testing experience, or utilize tools like Postman for manual testing if preferred.
  * If you make changes to API endpoints or their annotations, update the Swagger documentation by running `swag init` in the `cmd/api` directory.
  * Should you encounter any Go module-related issues (e.g., "missing modules" or import problems), execute `go mod tidy` in the project's root directory (`football-league-sim`).

<!-- end list -->

```
```