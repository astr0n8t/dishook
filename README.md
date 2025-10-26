# dishook

**dishook** is a lightweight Discord bot that converts Discord slash commands into HTTP webhook requests. It allows you to easily integrate Discord with external services by mapping slash commands to HTTP endpoints with customizable parameters, headers, and request bodies.

## Features

- **Simple webhook integration** - Map Discord slash commands to HTTP endpoints
- **Flexible configuration** - YAML-based configuration with hot-reloading
- **Support for subcommands** - Create nested command structures with subcommand groups
- **Template support** - Use Go templates to dynamically populate request data
- **Authentication** - Built-in support for authentication headers
- **Docker support** - Run as a containerized application
- **Custom arguments** - Support for string, int, float, and boolean argument types
- **Discord context** - Access Discord user information in your webhooks

## Installation

### Using Docker (Recommended)

```bash
docker pull astr0n8t/dishook:latest
docker run -v /path/to/config.yaml:/config.yaml astr0n8t/dishook
```

### Building from Source

#### Prerequisites

- Go 1.25.3 or later
- Git

#### Build Steps

```bash
# Clone the repository
git clone https://github.com/astr0n8t/dishook.git
cd dishook

# Build the binary
make build

# Run the application
./bin/dishook
```

## Configuration

dishook uses a YAML configuration file to define your Discord bot token, guild ID, and slash commands. Create a `config.yaml` file in the same directory as the executable:

### Basic Configuration Structure

```yaml
token: YOUR_DISCORD_BOT_TOKEN
guild_id: YOUR_GUILD_ID  # Optional: leave empty for global commands

commands:
  command_name:
    description: Description of your command
    response: Response message shown in Discord
    url: https://your-webhook-endpoint.com
    method: POST  # Optional: defaults to POST
    # ... additional options
```

### Configuration Options

#### Command Properties

- **`description`** (required): Description shown in Discord
- **`response`** (required): Message shown to the user after command execution
- **`url`** (required): HTTP endpoint to send the webhook request
- **`method`**: HTTP method (GET, POST, PUT, DELETE, etc.) - defaults to POST
- **`response_code`**: Expected HTTP response code - defaults to 200
- **`auth_header_name`**: Name of the authentication header
- **`auth_header_value`**: Value of the authentication header
- **`headers`**: Array of custom headers
- **`arguments`**: Array of command arguments
- **`data`**: Request body data (supports Go templates)
- **`subcommand`**: Map of subcommands (one level deep)
- **`subcommand_group`**: Map of subcommand groups (two levels deep)

#### Argument Properties

- **`name`**: Argument identifier
- **`description`**: Description shown in Discord
- **`type`**: Argument type (`string`, `int`, `float`, `bool`)
- **`required`**: Whether the argument is required
- **`default`**: Default value if not provided
- **`discord`**: Set to `true` to access Discord context (e.g., user information)

### Configuration Examples

#### Simple Command

```yaml
commands:
  hello:
    description: Say hello to the world
    response: Hello sent!
    url: https://example.com/hello
    method: GET
```

#### Command with Arguments

```yaml
commands:
  notify:
    description: Send a notification
    response: Notification sent!
    url: https://api.example.com/notify
    auth_header_name: Authorization
    auth_header_value: Bearer YOUR_API_TOKEN
    arguments:
      - name: message
        type: string
        description: The notification message
        required: true
      - name: priority
        type: int
        description: Priority level
        default: 0
        required: false
    data:
      message: "{{ .message }}"
      priority: "{{ .priority }}"
```

#### Command with Discord Context

```yaml
commands:
  register:
    description: Register a user
    response: User registered!
    url: https://api.example.com/register
    arguments:
      - name: discord_username
        type: string
        discord: true
      - name: email
        type: string
        description: Your email address
        required: true
    data:
      username: "{{ .discord_username }}"
      email: "{{ .email }}"
```

#### Nested Subcommands

```yaml
commands:
  admin:
    description: Admin commands
    subcommand:
      ban:
        description: Ban a user
        response: User banned!
        url: https://api.example.com/ban
        arguments:
          - name: user_id
            type: string
            description: User to ban
            required: true
        data:
          user_id: "{{ .user_id }}"
      
  settings:
    description: Settings commands
    subcommand_group:
      notifications:
        description: Notification settings
        subcommand:
          enable:
            description: Enable notifications
            response: Notifications enabled!
            url: https://api.example.com/settings/notifications/enable
```

#### Custom Headers

```yaml
commands:
  webhook:
    description: Send a webhook request
    response: Request sent!
    url: https://api.example.com/webhook
    headers:
      - name: X-Custom-Header
        value: custom-value
      - name: X-API-Version
        value: v1
    arguments:
      - name: payload
        type: string
        description: JSON payload
        required: true
    data:
      payload: "{{ .payload }}"
```

## Usage

### Setting Up Your Discord Bot

1. Create a Discord application at https://discord.com/developers/applications
2. Create a bot user and copy the bot token
3. Enable the "applications.commands" scope in OAuth2
4. Invite the bot to your server with the appropriate permissions
5. Copy your server (guild) ID

### Running dishook

#### With Binary

```bash
# Place your config.yaml in the same directory
./dishook
```

#### With Docker

```bash
docker run -v $(pwd)/config.yaml:/config.yaml astr0n8t/dishook:latest
```

#### Using Docker Compose

Create a `docker-compose.yml`:

```yaml
version: '3.8'

services:
  dishook:
    image: astr0n8t/dishook:latest
    volumes:
      - ./config.yaml:/config.yaml
    restart: unless-stopped
```

Run with:

```bash
docker-compose up -d
```

### Hot Reloading

dishook watches the configuration file for changes and automatically reloads when modifications are detected. Simply edit your `config.yaml` and the bot will update its commands without needing a restart.

## Development

### Project Structure

```
dishook/
├── cmd/            # CLI command definitions
├── config/         # Configuration management
├── internal/       # Core application logic
│   ├── http.go     # HTTP request handling
│   ├── internal.go # Main application logic
│   ├── slash.go    # Discord slash command handling
│   └── types.go    # Type definitions
├── version/        # Version information
├── config.yaml     # Configuration file
├── main.go         # Application entry point
└── Makefile        # Build automation
```

### Building

```bash
# Build for current platform
make build

# Build for Alpine Linux
make build-alpine

# Build Docker image
make package

# Tag Docker image
make tag

# Run tests
make test

# Clean build artifacts
make clean
```

### Running Tests

```bash
go test ./...
```

Or using make:

```bash
make test
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Credits

- **Development Lead**: Nathan Higley ([astr0n8t](https://github.com/astr0n8t))
- See [AUTHORS.md](AUTHORS.md) for a full list of contributors

## Support

If you encounter any issues or have questions:

- Open an issue: https://github.com/astr0n8t/dishook/issues
- Check existing issues for solutions

## Acknowledgments

- Built with [discordgo](https://github.com/bwmarrin/discordgo)
- Uses [Cobra](https://github.com/spf13/cobra) for CLI
- Configuration powered by [Viper](https://github.com/spf13/viper)
