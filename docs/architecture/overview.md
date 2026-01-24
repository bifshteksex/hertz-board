# Architecture Overview

## High-Level Architecture

HertzBoard follows a microservices architecture with a clear separation between frontend and backend services.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CLIENTS                                         │
├─────────────────┬─────────────────┬─────────────────┬───────────────────────┤
│   Web (SvelteKit)│  Desktop (Tauri) │  Mobile (Capacitor)│  Mobile Native     │
└────────┬────────┴────────┬────────┴────────┬────────┴───────────┬───────────┘
         │                 │                 │                     │
         └─────────────────┴─────────────────┴─────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           LOAD BALANCER (Nginx)                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                    ┌───────────────┴───────────────┐
                    ▼                               ▼
┌─────────────────────────────────┐   ┌─────────────────────────────────────┐
│         HERTZ API GATEWAY       │   │        HERTZ WEBSOCKET SERVER       │
│  - REST API                     │   │  - Real-time collaboration          │
│  - Authentication               │   │  - Presence (cursors, selections)   │
│  - Rate limiting                │   │  - CRDT synchronization             │
└───────────────┬─────────────────┘   └───────────────┬─────────────────────┘
                │                                     │
                └─────────────────┬───────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         INTERNAL SERVICES (Kitex gRPC)                      │
├─────────────────┬─────────────────┬─────────────────┬───────────────────────┤
│  Workspace Svc  │    User Svc     │   Asset Svc     │   Collaboration Svc   │
└────────┬────────┴────────┬────────┴────────┬────────┴───────────┬───────────┘
         │                 │                 │                    │
         ▼                 ▼                 ▼                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              DATA LAYER                                     │
├─────────────────┬─────────────────┬─────────────────┬───────────────────────┤
│   PostgreSQL    │      Redis      │      MinIO      │    ClickHouse         │
└─────────────────┴─────────────────┴─────────────────┴───────────────────────┘
```

## Core Components

### 1. API Gateway (Hertz)
- Entry point for all HTTP requests
- Authentication & Authorization
- Rate limiting
- Request routing to internal services

### 2. WebSocket Server (Hertz)
- Real-time collaboration
- CRDT-based synchronization
- Presence system (cursors, selections)
- Room management

### 3. Microservices (Kitex)

#### User Service
- User registration and authentication
- Profile management
- OAuth integration (Google, GitHub)
- Password reset

#### Workspace Service
- CRUD operations for workspaces
- Permissions management
- Sharing and collaboration

#### Asset Service
- File upload and storage
- Image processing
- CDN synchronization

#### Collaboration Service
- CRDT operations
- Conflict resolution
- History and undo/redo

### 4. Data Layer

#### PostgreSQL
- User data
- Workspace metadata
- Permissions
- Canvas elements metadata

#### Redis
- Session management
- Caching
- Pub/Sub for real-time events
- Rate limiting

#### MinIO
- Image storage
- Export files
- Backups

#### ClickHouse
- Analytics
- Usage statistics
- Audit logs

## Design Principles

### 1. Clean Architecture
- Clear separation of concerns
- Dependency inversion
- Domain-driven design

### 2. CRDT-based Synchronization
- Conflict-free replicated data types
- Eventual consistency
- Offline-first support

### 3. Real-time Collaboration
- WebSocket for bidirectional communication
- Presence awareness
- Operational transformation

### 4. Scalability
- Horizontal scaling of services
- Stateless design
- Caching strategies

### 5. Security
- JWT-based authentication
- Role-based access control (RBAC)
- Input validation
- Rate limiting

## Data Flow

### 1. User Authentication Flow
```
User → API Gateway → User Service → PostgreSQL
                  ↓
              JWT Token
                  ↓
             Redis (Session)
```

### 2. Canvas Update Flow
```
User → WebSocket Server → CRDT Engine → Redis Pub/Sub
                        ↓
                  PostgreSQL (Persist)
                        ↓
              Broadcast to other users
```

### 3. Asset Upload Flow
```
User → API Gateway → Asset Service → MinIO
                  ↓
            PostgreSQL (Metadata)
```

## Technology Stack

### Backend
- **Hertz**: High-performance HTTP framework
- **Kitex**: gRPC microservices framework
- **PostgreSQL**: Primary database
- **Redis**: Cache and pub/sub
- **MinIO**: Object storage
- **ClickHouse**: Analytics database
- **NATS**: Message queue

### Frontend
- **SvelteKit**: Full-stack framework
- **TypeScript**: Type safety
- **Tailwind CSS**: Styling
- **Yjs**: CRDT library
- **TipTap**: Rich text editor

### DevOps
- **Docker**: Containerization
- **Kubernetes**: Orchestration
- **GitHub Actions**: CI/CD
- **Prometheus**: Monitoring
- **Grafana**: Visualization
- **Jaeger**: Distributed tracing

## Security Considerations

1. **Authentication**: JWT tokens with short expiry
2. **Authorization**: RBAC for workspace access
3. **Data Encryption**: TLS for all communications
4. **Input Validation**: Strict validation on all inputs
5. **Rate Limiting**: Prevent abuse
6. **CORS**: Proper CORS configuration
7. **XSS Protection**: Content sanitization
8. **SQL Injection**: Prepared statements

## Performance Optimizations

1. **Caching**: Redis for frequently accessed data
2. **CDN**: Asset delivery via CDN
3. **Lazy Loading**: Components and assets
4. **Database Indexing**: Optimized queries
5. **Connection Pooling**: Reuse database connections
6. **Compression**: gzip/brotli compression

## Monitoring and Observability

1. **Metrics**: Prometheus for system metrics
2. **Logging**: Structured logging with levels
3. **Tracing**: Jaeger for distributed tracing
4. **Alerting**: Alert rules for critical issues
5. **Dashboards**: Grafana dashboards
