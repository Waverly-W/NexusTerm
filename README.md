# NexusTerm

**Mobile-First Web SSH Client**

[中文文档](README_ZH.md)

NexusTerm is a modern, responsive web-based SSH client designed for seamless remote server management on both desktop and mobile devices.

![Screenshot](docs/screenshot.png)

## Features

-   **📱 Mobile-First Design**: Optimized touch controls, virtual keyboard support, and responsive layout.
-   **📑 Multi-Tab Interface**: Manage multiple SSH sessions simultaneously in a tabbed view.
-   **🔒 Secure Host Management**: Save hosts with encrypted credentials (AES-256).
-   **🚀 Simple Deployment**: One-click deployment via Docker.
-   **🛠️ Admin Controls**: Bootstrapping admin account and toggling registration.

## Getting Started

The easiest way to run NexusTerm is using Docker Compose.

### Prerequisites

-   Docker & Docker Compose

### Deployment

1.  Clone the repository:
    ```bash
    git clone https://github.com/yourusername/nexusterm.git
    cd nexusterm
    ```

2.  Run with one command:
    ```bash
    docker compose up -d
    ```

3.  Access the application:
    Open your browser and navigate to `http://localhost:8080`.

## Configuration

You can configure the application using environment variables in `docker-compose.yml` or a `.env` file.

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PORT` | Service port | `8080` |
| `DB_PATH` | SQLite database path | `/data/nexus.db` |
| `ADMIN_USER` | Admin username (created on startup) | `admin` |
| `ADMIN_PASSWORD` | Admin password | (Empty) |
| `DISABLE_REGISTRATION` | Disable new user registration | `false` |

### Example `.env`

```ini
ADMIN_USER=admin
ADMIN_PASSWORD=my_secure_password
DISABLE_REGISTRATION=true
```

## Development

### Backend (Go)

```bash
cd backend
go run main.go
```

### Frontend (Vue 3 + Vite)

```bash
cd frontend
npm install
npm run dev
```

## License

MIT License
