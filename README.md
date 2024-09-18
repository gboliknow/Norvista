# Norvista Movie Reservation Service

Norvista is a backend system for managing movie reservations, allowing users to sign up, log in, browse movies, reserve seats, and manage reservations. The service includes user authentication, movie and showtime management, seat reservation functionality, and reporting on reservations.

## Table of Contents
- [Project Overview](#project-overview)
- [Features](#features)
- [Architecture](#architecture)
- [Setup](#setup)
- [API Endpoints](#api-endpoints)
- [Future Improvements](#future-improvements)
- [Contributing](#contributing)
- [License](#license)

## Project Overview

Norvista aims to provide a comprehensive backend solution for movie reservations. It supports:
- **User Authentication and Authorization**: Users can sign up, log in, and perform actions based on their roles.
- **Movie Management**: Admins can manage movies, including adding, updating, and deleting them.
- **Showtime Management**: Admins can manage showtimes for movies.
- **Seat Reservation**: Users can reserve seats for specific showtimes, view available seats, and cancel reservations.
- **Reporting**: Admins can view all reservations, capacity, and revenue reports.

## Features

- **User Authentication**: Sign up, log in, and role-based access control.
- **Movie Management**: Add, update, delete movies, and manage genres and posters.
- **Showtime Management**: Create and manage showtimes for movies.
- **Seat Reservation**: Reserve, view, and cancel seat reservations.
- **Reporting**: Admins can generate reports on reservations, capacity, and revenue.

## Architecture

The system is built using Go with GORM and PostgreSQL. It follows a layered architecture with the following components:

- **Models**: Defines the data structure and relationships.
- **Store Layer**: Handles database interactions and business logic.
- **API Layer**: Provides endpoints for interaction with the system.
- **Server**: Manages HTTP requests and responses.

### Data Models

- **User**: Represents users of the system with roles.
- **Movie**: Represents movies with details like title, description, and genre.
- **Showtime**: Represents showtimes for movies.
- **Seat**: Represents seats in a theater.
- **Reservation**: Represents seat reservations by users.

### Store Functions

- `CreateUser(user *models.User) (*models.User, error)`
- `FindUserByEmail(email string, user *models.User) error`
- `GetAllUsers() ([]models.User, error)`
- `CreateMovie(movie *models.Movie) (*models.Movie, error)`
- `UpdateMovie(movie *models.Movie) error`
- `CreateReservation(reservation *models.Reservation) error`
- `GetSeatsByShowtimeID(showtimeID string) ([]models.Seat, error)`
- `CancelReservation(reservationID string) error`

## Setup

### Prerequisites

- Go 1.XX or higher
- PostgreSQL
- `gorm` and `postgres` Go modules

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/gboliknow/norvista.git
   cd norvista
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up the environment variables:**
   Create a `.env` file in the root directory and add your PostgreSQL connection string:
   ```
   DB_URL=postgres://username:password@localhost:5432/norvista
   ```

4. **Run the application:**
   ```bash
   go run main.go
   ```

## API Endpoints

### User Endpoints

- `POST /users/signup`: Sign up a new user.
- `POST /users/login`: Log in a user.
- `GET /users/{id}`: Get user details.

### Movie Endpoints

- `POST /movies`: Add a new movie.
- `PUT /movies/{id}`: Update a movie.
- `DELETE /movies/{id}`: Delete a movie.
- `GET /movies`: Get all movies.

### Showtime Endpoints

- `POST /showtimes`: Create a new showtime.
- `PUT /showtimes/{id}`: Update a showtime.
- `DELETE /showtimes/{id}`: Delete a showtime.
- `GET /showtimes`: Get all showtimes.

### Reservation Endpoints

- `POST /reservations`: Reserve seats.
- `GET /reservations/user/{id}`: Get reservations for a user.
- `DELETE /reservations/{id}`: Cancel a reservation.
- `GET /seats/showtime/{id}`: Get seats for a specific showtime.

## Future Improvements

- **Frontend Integration**: Develop a frontend application to interact with the backend service.
- **Enhanced Reporting**: Add more detailed reporting and analytics.
- **Performance Optimization**: Improve the performance of queries and data handling.
- **Scalability**: Implement features for horizontal scaling and load balancing.
- **Testing**: Write comprehensive tests for all components and features.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or fix.
3. Commit your changes and push to your fork.
4. Create a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Feel free to adjust the details based on your specific needs and project setup!
https://roadmap.sh/projects/movie-reservation-system
