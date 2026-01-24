# REST API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require authentication using JWT tokens.

### Headers

```
Authorization: Bearer <access_token>
```

## Endpoints

### Authentication

#### Register

```http
POST /auth/register
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "display_name": "User Name"
}
```

**Response:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "username": "username",
  "display_name": "User Name",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Login

```http
POST /auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900,
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "username",
    "display_name": "User Name"
  }
}
```

#### Refresh Token

```http
POST /auth/refresh
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900
}
```

#### Logout

```http
POST /auth/logout
```

**Headers:** Authorization required

**Response:**
```json
{
  "message": "Logged out successfully"
}
```

### Users

#### Get Current User

```http
GET /users/me
```

**Headers:** Authorization required

**Response:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "username": "username",
  "display_name": "User Name",
  "avatar_url": "https://...",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Update Profile

```http
PUT /users/me
```

**Headers:** Authorization required

**Request Body:**
```json
{
  "display_name": "New Name",
  "avatar_url": "https://..."
}
```

### Workspaces

#### List Workspaces

```http
GET /workspaces
```

**Headers:** Authorization required

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "My Workspace",
      "description": "Description",
      "owner_id": "uuid",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}
```

#### Create Workspace

```http
POST /workspaces
```

**Headers:** Authorization required

**Request Body:**
```json
{
  "name": "My Workspace",
  "description": "Description",
  "settings": {
    "background_color": "#ffffff",
    "grid_enabled": true,
    "grid_size": 20,
    "snap_to_grid": false
  }
}
```

#### Get Workspace

```http
GET /workspaces/:id
```

**Headers:** Authorization required

**Response:**
```json
{
  "id": "uuid",
  "name": "My Workspace",
  "description": "Description",
  "owner_id": "uuid",
  "settings": {},
  "collaborators": [],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Update Workspace

```http
PUT /workspaces/:id
```

**Headers:** Authorization required

**Request Body:**
```json
{
  "name": "Updated Name",
  "description": "Updated Description",
  "settings": {}
}
```

#### Delete Workspace

```http
DELETE /workspaces/:id
```

**Headers:** Authorization required

**Response:**
```json
{
  "message": "Workspace deleted successfully"
}
```

### Canvas Elements

#### Get Elements

```http
GET /workspaces/:workspace_id/elements
```

**Headers:** Authorization required

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "workspace_id": "uuid",
      "element_type": "text",
      "position": { "x": 100, "y": 200 },
      "size": { "width": 300, "height": 100 },
      "rotation": 0,
      "z_index": 1,
      "content": {},
      "style": {},
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### Create Element

```http
POST /workspaces/:workspace_id/elements
```

**Headers:** Authorization required

**Request Body:**
```json
{
  "element_type": "text",
  "position": { "x": 100, "y": 200 },
  "size": { "width": 300, "height": 100 },
  "content": {
    "text": "Hello World"
  },
  "style": {
    "font_size": 16,
    "color": "#000000"
  }
}
```

#### Update Element

```http
PUT /workspaces/:workspace_id/elements/:element_id
```

**Headers:** Authorization required

#### Delete Element

```http
DELETE /workspaces/:workspace_id/elements/:element_id
```

**Headers:** Authorization required

### Assets

#### Upload Asset

```http
POST /assets
```

**Headers:** 
- Authorization required
- Content-Type: multipart/form-data

**Form Data:**
- `file`: File to upload
- `workspace_id`: Workspace UUID

**Response:**
```json
{
  "id": "uuid",
  "filename": "image.png",
  "content_type": "image/png",
  "size_bytes": 12345,
  "url": "https://...",
  "created_at": "2024-01-01T00:00:00Z"
}
```

## Error Responses

### Format

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {}
  }
}
```

### Common Status Codes

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Access denied
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Validation error
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

### Error Codes

- `INVALID_REQUEST` - Invalid request format
- `UNAUTHORIZED` - Authentication failed
- `FORBIDDEN` - Access denied
- `NOT_FOUND` - Resource not found
- `VALIDATION_ERROR` - Validation failed
- `RATE_LIMIT_EXCEEDED` - Too many requests
- `INTERNAL_ERROR` - Server error
