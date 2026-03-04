#  Go-Food Microservices Ecosystem (Full-Stack)

A complete, high-performance microservice architecture built in **Go (Golang)**. This system demonstrates advanced patterns for **reliability, security, and AI-driven data enrichment**, modeled after modern enterprise distributed systems.

---

##  Key Features & Architectural Patterns

###  AI-Powered Semantic Search 
Unlike traditional keyword matching, this system uses an **Asynchronous AI Pipeline**. When a product is created, an LLM (via **OpenRouter**) analyzes the description to generate semantic tags (e.g., "energy boost," "vegan-friendly"). These are indexed in a specialized **Search Service**, allowing users to find food based on intent and meaning.

###  Reliability & Security
*   **Transactional Outbox Pattern:** Solves the "dual-write" problem in the Catalog and Ordering services. Guarantees 100% event delivery to RabbitMQ even during network failures.
*   **CQRS:** Segregates the **Write Model** (Catalog) from the **Read Model** (Search) to optimize performance and scalability.
*   **JWT & RBAC:** Implements stateless authentication and Role-Based Access Control to secure administrative routes.

---

## System Architecture
The system consists of 6 independent services communicating via HTTP and RabbitMQ:

1.  **Identity Service**: Handles User Registration, Login (Bcrypt), and JWT Issuance.
2.  **Catalog Service**: Manages the product menu with Role-Based Access Control (Admin only).
3.  **Inventory Service**: Manages stock levels, reacting to new products and orders.
4.  **Ordering Service**: Orchestrates orders, validating price (Catalog) and stock (Inventory).
5.  **Search Service**: A read-optimized projection for fast, paginated product searches.
6.  **Notification Service**: Background worker for simulated email confirmations.
7.  **API Gateway**: The single entry point (Port 8000) using Echo Proxy.

## Technologies Used
- **Language**: Go 1.24
- **Web Framework**: Echo (v4)
- **Database**: PostgreSQL with GORM
- **Messaging**: RabbitMQ (AMQP)
- **Security**: JWT (JSON Web Tokens) with RBAC (Role-Based Access Control)
- **AI**: OpenRouter (LLM Integration for Semantic Tagging)
- **Reliability**: Outbox Pattern for guaranteed message delivery
- **Deployment**: Docker & Docker Compose
- **Testing**: Unit Tests with Testify and SQLMock
- **CI/CD**: GitHub Actions

## How to Run
1. Clone the repository.
2. Ensure Docker Desktop is running.
3. Run the entire system with one command:
   ```bash
   docker-compose up --build