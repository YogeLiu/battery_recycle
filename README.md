# 废旧电池回收进销存（ERP）系统

A comprehensive ERP system for used battery recycling management, handling receipt, sales, and inventory operations.

## Features

- **Purchase Management**: Handle battery receipt with weighing, pricing, and document generation
- **Sales Management**: Process battery sales with inventory validation and order tracking  
- **Inventory Management**: Real-time inventory tracking with automatic record creation
- **Reporting**: Generate daily/monthly/yearly business analytics and reports
- **User Management**: Role-based access control (Super Admin, Normal User)

## Tech Stack

- **Backend**: Go 1.23 with Gin framework
- **Frontend**: React with Vite
- **Database**: MySQL 8.0+
- **Authentication**: JWT tokens
- **Deployment**: Docker & Docker Compose

## Project Structure

```
/battery-erp-system/
├── backend/          # Go backend service
│   ├── cmd/server/   # Application entry point
│   ├── internal/     # Internal packages
│   └── Dockerfile    # Backend container config
├── frontend/         # React frontend application
│   ├── src/          # Source code
│   ├── public/       # Static assets
│   └── Dockerfile    # Frontend container config
└── docker-compose.yml # Service orchestration
```

## Quick Start

1. **Prerequisites**
   - Docker & Docker Compose
   - Go 1.23+ (for development)
   - Node.js 20+ (for development)

2. **Run with Docker**
   ```bash
   docker-compose up --build
   ```

3. **Development Setup**
   ```bash
   # Backend
   cd backend
   go mod download
   go run cmd/server/main.go

   # Frontend  
   cd frontend
   npm install
   npm run dev
   ```

## API Documentation

All APIs use the `/jxc/v1` prefix and require JWT authentication (except login).

- **Authentication**: `POST /jxc/v1/auth/login`
- **Users**: `GET|POST /jxc/v1/users`
- **Categories**: `GET|POST /jxc/v1/categories`
- **Inbound**: `GET|POST /jxc/v1/inbound/orders`
- **Outbound**: `GET|POST /jxc/v1/outbound/orders`
- **Inventory**: `GET /jxc/v1/inventory`
- **Reports**: `GET /jxc/v1/reports/summary`

## Development

This project follows a modular architecture with clear separation between frontend and backend services. All business operations use atomic transactions to ensure data consistency.

## License

Private Project