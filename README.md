Football League Simulation API
This project is a REST API developed in Go, designed to simulate a football league and predict team championship probabilities using advanced algorithms and concurrent processing.
Features

League Table: Displays the current league standings, including goal differences and championship probability predictions.
Weekly Match Simulation: Simulates matches for the current week and updates team standings accordingly.
Full League Simulation: Automatically simulates all remaining weeks to complete the season.
League Reset: Resets team statistics and match fixtures to start a new season.
Championship Predictions: Utilizes Monte Carlo simulation algorithms, executed concurrently using Go's goroutines for multithreaded performance, to calculate championship probabilities based on remaining matches.
Advanced Logging: Implements detailed logging for info, warnings, and errors to streamline development and debugging.
Automated Database Setup and Migration: Automatically checks and creates the required database schema and tables on startup, initializing fixture data as needed.

Prerequisites
To set up and run this project locally, ensure you have the following installed:

Go: Version 1.20 or higher. Install Go
Docker Desktop: Required to run the MSSQL database. Download Docker Desktop
Git: Used to clone the project repository. Download Git

Installation and Setup
Follow these steps to set up and run the project locally:
1. Clone the Repository
Clone the project repository to your local machine:
git clone https://github.com/muzaffertuna/football-league-sim.git
cd football-league-sim

2. Configure Environment Variables
The API uses environment variables for configuration, such as database connection details and server address. Create a .env file in the project root directory (football-league-sim) and add the following content, replacing the password with a strong one of your choice:
SA_PASSWORD="YourStrongPassword123"
DB_CONNECTION_STRING="sqlserver://sa:${SA_PASSWORD}@localhost:1433?database=FootballLeagueSim&TrustServerCertificate=true"
SERVER_ADDRESS=":8080"

3. Set Up and Start the Database
The project includes a docker-compose.yml file in the root directory to configure and run an MSSQL database container, utilizing the SA_PASSWORD from the .env file.
Start the database with:
docker-compose up -d

4. Run the Go Application
With the database running and the .env file configured, start the API application:
cd cmd/api
go run main.go

Upon successful startup, you should see output similar to:
2025/05/27 16:07:56 logger.go:39: [INFO] Successfully connected to MSSQL
2025/05/27 16:07:56 logger.go:39: [INFO] Starting server on :8080

5. Accessing and Testing the API
The API will be available at http://localhost:8080. You can interact with the endpoints using the following methods:
Method 1: Swagger UI (Recommended)
Swagger UI provides a user-friendly interface to explore and test all API endpoints directly in your browser. Visit:
http://localhost:8080/swagger/index.html#/

Expand each endpoint to view details, click "Try it out" to send requests, and review responses instantly.
Method 2: Postman, cURL, or Similar Tools
Alternatively, use tools like Postman, Insomnia, or cURL to send requests to the API endpoints.
API Endpoints
For the application to function correctly, you must reset the league before running simulations for the first time.
POST /reset-league
Description: Resets all team statistics and match fixtures to start a new season. This endpoint must be called before other simulation operations on first use to initialize the database. Subsequent uses can proceed without resetting.
cURL Example:
curl -X POST http://localhost:8080/reset-league

POST /play-week
Description: Simulates matches for the current week and updates team standings. Each call advances to the next week.
cURL Example:
curl -X POST http://localhost:8080/play-week

GET /league-table
Description: Retrieves the current league standings, including goal differences and championship probabilities calculated using multithreaded Monte Carlo simulations.
cURL Example:
curl -X GET http://localhost:8080/league-table

POST /simulate-all-weeks
Description: Simulates all remaining weeks to complete the season.
cURL Example:
curl -X POST http://localhost:8080/simulate-all-weeks

Code Snippets
Below are the key code snippets for easy reference and copying.
Environment Variables (.env)
SA_PASSWORD="YourStrongPassword123"
DB_CONNECTION_STRING="sqlserver://sa:${SA_PASSWORD}@localhost:1433?database=FootballLeagueSim&TrustServerCertificate=true"
SERVER_ADDRESS=":8080"

Clone the Repository
git clone https://github.com/muzaffertuna/football-league-sim.git
cd football-league-sim

Start the Database
docker-compose up -d

Run the Application
cd cmd/api
go run main.go

Notes

Ensure the .env file is correctly configured before starting the database and application.
The /reset-league endpoint is mandatory for first-time setup to avoid database errors.
Use Swagger UI for a seamless testing experience, or leverage tools like Postman for manual testing.

