# Real-time Chat App (Go + WebSockets)

## 🔹 Features
- ✅ JWT-based authentication
- ✅ Real-time chat using WebSockets
- ✅ Message storage in PostgreSQL
- ✅ Chat history retrieval (REST API)
- ✅ Rate limiting (Max 5 messages per 10 seconds)
- ✅ Structured logging & error handling

## 🔹 API Endpoints
| Method | Endpoint | Description |
|--------|----------|------------|
| `POST` | `/register` | Register a new user |
| `POST` | `/login` | Get JWT token |
| `GET` | `/ws` | WebSocket connection |
| `GET` | `/chat/history?user1={}&user2={}&page={}&limit={}` | Fetch chat history |

## 🔹 Running the Project
1. Install Go and PostgreSQL
2. Clone the repo:  
