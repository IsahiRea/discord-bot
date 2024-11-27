# Discord Bot (WIP)

Welcome to the repository for my Discord bot!

## ⚠️ Work In Progress

This bot is currently a **Work In Progress**.The purpose of this project is to learn and experiment with Discord bot development.

### Current Status

This is a Discord bot with a custom backend built using Go and a frontend utilizing discord.js for bot commands.

### Features

- **Authentication**: Secure login and refresh token system.
- **User Management**: Create, retrieve, and manage Discord user data.
- **Fun Commands**: Includes trivia, GIF picker, and more.
- **Bot Utility Commands**: Ban, kick, reload, etc.
- **API**: RESTful API for bot functionalities, including readiness checks and fun commands.

---

## Project Structure

### Bot (Frontend)
- **auth**: Handles authentication logic.
- **commands**: Contains the bot's command logic for various functionalities.
- **events**: Event handling for Discord interactions.
- **shared**: Common utilities shared across the bot.
- **deploy-commands.js**: Script to deploy bot commands to Discord.

### Backend
- **internal/auth**: Contains authentication-related logic in Go.
- **internal/database**: SQL files and database models to manage user and token data.
- **internal/sql**: Database schema and queries.
- **handlers**: Go files that handle HTTP requests for users, tokens, and trivia.
- **main.go**: Initializes the server and sets up routes.

---

## API Endpoints

### Authentication
- `POST /api/login`: Login to the bot.
- `GET /api/refresh`: Refresh the JWT token.
- `GET /api/revoke`: Revoke the current session.

### User Management
- `GET /api/users/{discordID}`: Get user by Discord ID.
- `POST /api/users`: Create a new user.

### Fun Commands
- `GET /api/trivia`: Get a random trivia question.
- `POST /api/stories`: Add on to created story
- `POST /api/images`: Generate requested image

---

### 🚧 Planned Features

Here are some features I plan to implement in the near future:

- User interaction features
- Admin tools
- Fun commands

