# 🍔 Go-Food Delivery Microservices

A full-fledged, event-driven food delivery system built with **Go (Golang)** and **Microservices Architecture**. This project follows **Clean Architecture**, **DDD**, and **CQRS** patterns.

## 🏗 System Architecture
The system consists of 6 independent services communicating via HTTP and RabbitMQ:

1.  **Identity Service**: Handles User Registration, Login (Bcrypt), and JWT Issuance.
2.  **Catalog Service**: Manages the product menu with Role-Based Access Control (Admin only).
3.  **Inventory Service**: Manages stock levels, reacting to new products and orders.
4.  **Ordering Service**: Orchestrates orders, validating price (Catalog) and stock (Inventory).
5.  **Search Service**: A read-optimized projection for fast, paginated product searches.
6.  **Notification Service**: Background worker for simulated email confirmations.
7.  **API Gateway**: The single entry point (Port 8000) using Echo Proxy.

## 🚀 Technologies Used
- **Language**: Go 1.24
- **Web Framework**: Echo (v4)
- **Database**: PostgreSQL with GORM
- **Messaging**: RabbitMQ (AMQP)
- **Security**: JWT (JSON Web Tokens) with RBAC (Role-Based Access Control)
- **Reliability**: Outbox Pattern for guaranteed message delivery
- **Deployment**: Docker & Docker Compose
- **Testing**: Unit Tests with Testify and SQLMock
- **CI/CD**: GitHub Actions

## 🚦 How to Run
1. Clone the repository.
2. Ensure Docker Desktop is running.
3. Run the entire system with one command:
   ```bash
   docker-compose up --build