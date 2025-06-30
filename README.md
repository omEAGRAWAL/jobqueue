
# Go Job Queue System

## 🔧 Setup

1. Install Go, PostgreSQL, and set up the following environment variable:

```
export DATABASE_URL="postgresql://neondb_owner:npg_VzoF2mdZ0hXp@ep-green-morning-a82bpjj4-pooler.eastus2.azure.neon.tech/neondb?sslmode=require&channel_binding=require"
```

2. Create the `jobs` table:

```sql
CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    payload TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    result TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

3. Run the app:

```bash
go run cmd/main.go
```

## 📦 API Endpoints
- `GET /`      - Simple HTML UI with live logs and updates
- `POST /jobs` – Submit a job
- `GET /jobs/{id}` – Get status/result
- `GET /jobs` – List jobs
