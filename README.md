# PocoClinic EMR

PocoClinic is an open-source Electronic Medical Records (EMR) system designed specifically for non-profit healthcare organizations. The system provides a secure, easy-to-use platform for managing patient records while maintaining high standards of data privacy and security.

## 🌟 Features

- Secure authentication with 64-bit key + 4-digit PIN system
- Patient demographics management
- Modern React-based user interface with Mantine components
- Robust Go backend with modular monolith architecture
- CockroachDB for reliable and scalable data persistence
- Comprehensive audit logging
- HIPAA-compliant data handling

## 🏗️ Architecture

PocoClinic follows a modular monolith architecture with vertical slices using the mediator pattern. The system is designed with the following key components:

- **Frontend**: React-based SPA with TypeScript and Mantine UI
- **Backend**: Go-based API server
- **Database**: CockroachDB
- **Authentication**: Custom secure authentication system
- **Testing**: Vitest for frontend, Go testing for backend

For detailed architecture decisions, please refer to the [Architecture Decision Records](./docs/adr/README.md).

## 🚀 Getting Started

### Prerequisites

- Node.js (v18 or higher)
- Go (v1.21 or higher)
- CockroachDB
- Docker (optional, for containerized development)

### Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/PocoClinic.git
   cd PocoClinic
   ```

2. Install and run frontend (Windows):
   ```bash
   # Using the provided batch file (recommended)
   run-frontend.bat
   ```
   This will:
   - Install dependencies
   - Run tests
   - Start the development server (only if tests pass)

   Alternative manual steps:
   ```bash
   cd frontend
   npm install
   npm test run    # Run tests
   npm run dev     # Start development server
   ```

3. Install backend dependencies:
   ```bash
   cd backend
   go mod tidy
   ```

4. Set up the database:
   ```bash
   # Instructions for CockroachDB setup will be provided
   ```

5. Start the backend server:
   ```bash
   cd backend
   go run cmd/main.go
   ```

## 📁 Project Structure

```
PocoClinic/
├── docs/
│   ├── adr/          # Architecture Decision Records
│   └── api/          # API Documentation
├── frontend/
│   ├── src/
│   │   ├── components/   # Reusable UI components
│   │   ├── features/     # Feature-specific components
│   │   ├── mocks/        # MSW API mocks
│   │   └── shared/       # Shared utilities
│   └── package.json
├── backend/
│   ├── cmd/
│   ├── internal/
│   │   ├── domain/      # Domain models
│   │   ├── features/    # Feature implementations
│   │   └── shared/      # Shared utilities
│   └── go.mod
└── run-frontend.bat    # Windows frontend setup script
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](./CONTRIBUTING.md) for details on how to submit pull requests, report issues, and contribute to the project.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## 🔒 Security

If you discover any security-related issues, please email security@pococodo.com instead of using the issue tracker.

## 📫 Contact

- Project Maintainers: [List of maintainers]
- Community: [Links to community channels]
