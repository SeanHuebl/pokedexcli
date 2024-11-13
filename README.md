# Pokedex CLI

The **Pokedex CLI** is a command-line application written in Go that enables users to explore Pokémon locations, view detailed Pokémon data, and attempt to catch Pokémon with randomized success rates. This project provided an opportunity to build a Go application from scratch, focusing on CLI functionality, API integration, and concurrency.

## Project Overview

This project was guided only at a high level, with broad goals and suggestions, but it lacked code samples or detailed pseudo code. I was directed to the API documentation, but the implementation, including how to make API calls and structure the data handling, was self-driven. This approach encouraged independent problem-solving and reinforced my understanding of Go, API integration, and memory-safe data access.

### Key Learning Outcomes

Throughout the development process, I gained practical experience in the following areas:

- **HTTP Requests and API Integration**: I learned how to make GET requests from the CLI, parse JSON responses into Go structs, and use those structs to manipulate and display data.
- **Caching for Optimized Access**: Implemented an in-memory cache to store API responses, reducing redundant calls to the API and speeding up data retrieval.
- **Go Routines and Concurrency**: Utilized Go routines to handle background tasks, such as managing cache expiration, allowing the main program to remain responsive.
- **Channel Communication**: Leveraged open channels to listen for continuous user input, enabling a real-time CLI experience.
- **Memory Safety with Mutexes**: Protected shared resources, specifically maps, using mutexes to prevent concurrent access issues, ensuring safe data manipulation.

## Features

- **Map Navigation**: Allows the user to view Pokémon locations by navigating to the next or previous set of locations.
- **Exploration**: Retrieves details about the Pokémon found in a specified location.
- **Catch Pokémon**: Enables users to attempt to catch a Pokémon, with the success rate based on the Pokémon's experience level.
- **Inspect Pokémon**: Displays detailed information about a caught Pokémon, including its stats and type.
- **View Pokedex**: Lists all Pokémon that have been successfully caught and saved in the Pokedex.

## Commands

- `help`: Displays a help message with descriptions for all available commands.
- `map`: Shows the next 20 Pokémon location areas.
- `mapb`: Shows the previous 20 Pokémon location areas.
- `explore <location-area>`: Displays the Pokémon available in the specified area.
- `catch <pokemon-name>`: Attempts to catch a specific Pokémon.
- `inspect <pokemon-name>`: Displays information about a caught Pokémon.
- `pokedex`: Lists all Pokémon saved in the Pokedex.
- `exit`: Exits the application.

## Technical Implementation

The application is built around Go’s capabilities in concurrent programming, API handling, and data caching. Below are some technical highlights:

- **API Requests and JSON Parsing**: The application sends GET requests to the PokeAPI, retrieves JSON data, parses it into Go structs, and uses this data for CLI operations.
- **Caching**: An in-memory cache stores API responses for faster access. This cache is periodically cleared using a Go routine to manage expiration.
- **Concurrency with Go Routines**: Background tasks, like cache management, run concurrently, allowing the CLI to remain responsive to user commands.
- **Channels for CLI Input**: The application listens for user input on an open channel, enabling a seamless CLI experience without blocking the main application loop.
- **Mutexes for Safe Data Access**: Mutexes ensure safe access to shared data structures, such as maps, protecting against race conditions in a concurrent environment.

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/pokedexcli.git
   cd pokedexcli
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the CLI**:
   ```bash
   go run main.go
   ```

## Usage

Run the application by executing `go run main.go`. Use the commands listed above to navigate the Pokedex data, explore areas, catch Pokémon, and view information.

## Acknowledgments

This project provided a valuable learning experience in API integration, concurrency, and memory-safe programming in Go. With minimal guidance, I implemented each component, gaining confidence in building a complex CLI application from scratch.

## License

This project is open-source and available under the MIT License.