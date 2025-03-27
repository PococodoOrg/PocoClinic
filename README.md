# PocoClinic EMR

**Note: This project is in open development using generative AI. We are actively looking for testers and contributors to help improve and expand the system.**

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
   git clone https://github.com/PococodoOrg/PocoClinic.git
   cd PocoClinic
   ```

2. Install and run the entire system (Windows):
   ```bash
   # Using the provided batch file (recommended)
   run-all.bat
   ```
   This will:
   - Install frontend and backend dependencies
   - Run tests
   - Start both the frontend and backend servers

3. Install and run frontend (Windows):
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

4. Install backend dependencies:
   ```bash
   cd backend
   go mod tidy
   ```

5. Set up the database:
   ```bash
   # Instructions for CockroachDB setup will be provided
   ```

6. Start the backend server:
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

## Our Story

PocoClinic is an open-source Electronic Medical Record (EMR) system designed specifically for small, independent clinics, with a particular focus on non-profit healthcare organizations. We believe that every clinic deserves access to quality healthcare technology, regardless of their size or budget.

### Built for Local Healthcare

- **Small Clinic Focus**: Designed for independent clinics and small healthcare practices
- **Local Network Ready**: Runs on a Raspberry Pi, making it perfect for local intranet deployment
- **Complete Documentation**: Includes a comprehensive instruction binder for setup and maintenance
- **Human-Centered Design**: Features are designed with clinic staff and patient care in mind

### Privacy-First Approach

- **Staff Authentication**: Clinic staff access the system using secure credentials
- **Patient Privacy**: Patient records can be obfuscated using 64-bit QR code identifiers
- **End-to-End Security**: All data is encrypted both at rest and in transit
- **Local Control**: Your data stays within your clinic's network

### Why PocoClinic?

We believe that healthcare technology should be:
- Accessible to all clinics, regardless of size
- Simple and straightforward to use
- Secure and private
- Focused on patient care

PocoClinic is our contribution to the healthcare community - a complete EMR solution that prioritizes patient care while maintaining high standards of security and privacy. We offer it as an option for clinics looking for a simple, secure way to manage their patient records.
