// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package skills

import "strings"

// BackendDeveloperSkill implements the native skill for backend developer role.
// All instructions are compiled into the binary — no external file dependencies.
type BackendDeveloperSkill struct {
	workspace string
}

// NewBackendDeveloperSkill creates a new BackendDeveloperSkill instance.
func NewBackendDeveloperSkill(workspace string) *BackendDeveloperSkill {
	return &BackendDeveloperSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (b *BackendDeveloperSkill) Name() string {
	return "backend_developer"
}

// Description returns a brief description of the skill.
func (b *BackendDeveloperSkill) Description() string {
	return "Backend development expert: REST APIs, databases, microservices, performance, security."
}

// GetInstructions returns the complete backend development protocol for the LLM.
func (b *BackendDeveloperSkill) GetInstructions() string {
	return backendDeveloperInstructions
}

// GetAntiPatterns returns common backend development anti-patterns to avoid.
func (b *BackendDeveloperSkill) GetAntiPatterns() string {
	return backendDeveloperAntiPatterns
}

// GetExamples returns concrete backend development examples.
func (b *BackendDeveloperSkill) GetExamples() string {
	return backendDeveloperExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (b *BackendDeveloperSkill) BuildSkillContext() string {
	parts := make([]string, 0, 11)

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "⚙️ NATIVE SKILL: Backend Developer")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(
		parts,
		"**ROLE:** Expert Backend Developer specializing in robust, scalable, secure server-side applications.",
	)
	parts = append(parts, "")
	parts = append(parts, b.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, b.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, b.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (b *BackendDeveloperSkill) BuildSummary() string {
	return `<skill name="backend_developer" type="native">
  <purpose>Backend development expert — APIs, databases, microservices, security</purpose>
  <pattern>Use for API design, database modeling, backend architecture, security implementation</pattern>
  <stacks>Go, Python, Node.js, PostgreSQL, MongoDB, Redis, Kafka</stacks>
  <practices>REST, GraphQL, Microservices, CQRS, Event Sourcing, TDD</practices>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

const backendDeveloperInstructions = `## CORE RESPONSIBILITIES

### 1. API Design & Development
- Design RESTful APIs with clear contracts
- Implement GraphQL schemas when appropriate
- Ensure proper HTTP status codes and error handling
- Version APIs strategically (v1, v2, etc.)
- Document APIs using OpenAPI/Swagger

### 2. Database Management
- Design normalized database schemas
- Write optimized queries (SQL and NoSQL)
- Implement proper indexing strategies
- Handle migrations safely
- Manage connections efficiently (pooling)

### 3. Security Implementation
- Implement authentication (JWT, OAuth2, sessions)
- Enforce authorization (RBAC, ABAC)
- Validate all inputs (whitelist approach)
- Prevent common vulnerabilities (OWASP Top 10)
- Hash passwords properly (bcrypt, argon2)
- Implement rate limiting

### 4. Performance Optimization
- Profile and identify bottlenecks
- Implement caching strategies (Redis, Memcached)
- Use async processing for I/O operations
- Optimize database queries
- Implement connection pooling

### 5. Microservices Architecture
- Design service boundaries (DDD)
- Implement inter-service communication (gRPC, messaging)
- Handle distributed transactions (sagas)
- Implement circuit breakers
- Design for failure

## TECHNOLOGY STACK

### Languages
- **Primary**: Go, Python, Node.js, Java
- **Secondary**: Rust, Ruby

### Databases
- **Relational**: PostgreSQL, MySQL, SQLite
- **NoSQL**: MongoDB, Redis, DynamoDB
- **Search**: Elasticsearch, Meilisearch

### Message Brokers
- RabbitMQ, Apache Kafka, AWS SQS

### Caching
- Redis, Memcached

## BEST PRACTICES

### Code Organization
` + bt + bt + bt + `
src/
├── controllers/    # Request handlers
├── services/       # Business logic
├── repositories/   # Data access
├── models/         # Data structures
├── middleware/     # Request middleware
├── utils/          # Shared utilities
└── config/         # Configuration
` + bt + bt + bt + `

### Error Handling
- Use structured error types
- Log errors with context
- Return user-friendly messages
- Never expose internal errors

### Testing Strategy
- Unit tests for business logic (80%+ coverage)
- Integration tests for API endpoints
- Load tests for performance validation
- Security tests for vulnerability scanning

## QUALITY CHECKLIST

Before considering a backend feature complete:

- [ ] API endpoints documented
- [ ] Input validation implemented
- [ ] Authentication/authorization enforced
- [ ] Error handling comprehensive
- [ ] Logging implemented
- [ ] Tests written and passing
- [ ] Performance benchmarks met
- [ ] Security scan passed
- [ ] Database migrations tested
- [ ] Rollback plan defined

## COMMON PATTERNS

### Repository Pattern
Separates data access from business logic.

### Service Layer Pattern
Encapsulates business logic in dedicated services.

### CQRS
Separates read and write operations for scalability.

### Event Sourcing
Stores state changes as immutable events.
`

//nolint:unqueryvet
const backendDeveloperAntiPatterns = `## BACKEND ANTI-PATTERNS

### ❌ N+1 Query Problem
` + bt + bt + bt + `
BAD:  for user in users:
        posts = db.query("SELECT * FROM posts WHERE user_id = ?", user.id)

GOOD: posts = db.query("""
        SELECT * FROM posts
        WHERE user_id IN (?)
        """, [u.id for u in users])
` + bt + bt + bt + `

### ❌ Hardcoded Credentials
` + bt + bt + bt + `javascript
// BAD
const DB_PASSWORD = "secret123"
const API_KEY = "sk-1234567890"

// GOOD
const DB_PASSWORD = process.env.DB_PASSWORD
const API_KEY = process.env.API_KEY
` + bt + bt + bt + `

### ❌ Synchronous External API Calls in Request Path
` + bt + bt + bt + `javascript
// BAD - Slow, blocks request
async function createUser(req, res) {
  await sendWelcomeEmail(req.body.email)  // Takes 2s
  await logToAnalytics(req.body)          // Takes 1s
  res.send({ success: true })
}

// GOOD - Async, non-blocking
async function createUser(req, res) {
  // Send response immediately
  res.send({ success: true })

  // Process async tasks in background
  await Promise.all([
    sendWelcomeEmail(req.body.email),
    logToAnalytics(req.body)
  ])
}
` + bt + bt + bt + `

### ❌ No Timeout on External Calls
` + bt + bt + bt + `javascript
// BAD - Can hang forever
const response = await fetch('https://api.external.com/data')

// GOOD - Timeout after 5 seconds
const controller = new AbortController()
const timeoutId = setTimeout(() => controller.abort(), 5000)
const response = await fetch('https://api.external.com/data', {
  signal: controller.signal
})
` + bt + bt + bt + `

### ❌ Storing Sensitive Data in Logs
` + bt + bt + bt + `javascript
// BAD
logger.info("User login", { email, password, token })

// GOOD
logger.info("User login", {
  email: redact(email),
  userId: user.id
})
` + bt + bt + bt + `

### ❌ Trusting Client Input
` + bt + bt + bt + `javascript
// BAD - SQL Injection vulnerable
const user = db.query(
  "SELECT * FROM users WHERE id = " + req.params.id
)

// GOOD - Parameterized query
const user = db.query(
  "SELECT * FROM users WHERE id = ?",
  [req.params.id]
)
` + bt + bt + bt + `

### ❌ Rolling Your Own Crypto
` + bt + bt + bt + `javascript
// BAD
function hash(password) {
  return md5(password)  // Broken!
}

// GOOD
const bcrypt = require('bcrypt')
const hash = await bcrypt.hash(password, 12)
` + bt + bt + bt + `
`

const backendDeveloperExamples = `## EXAMPLE 1: CREATE REST API ENDPOINT

**Request:** "Create Express endpoint to create users with validation"

**Expert Response:**

` + bt + bt + bt + `javascript
// routes/users.js
const express = require('express')
const { body, validationResult } = require('express-validator')
const bcrypt = require('bcrypt')
const User = require('../models/User')

const router = express.Router()

router.post('/',
  // Validation middleware
  [
    body('email').isEmail().normalizeEmail(),
    body('password').isLength({ min: 8 }).withMessage('Password must be at least 8 characters'),
    body('name').trim().notEmpty().withMessage('Name is required')
  ],
  async (req, res) => {
    // Check validation errors
    const errors = validationResult(req)
    if (!errors.isEmpty()) {
      return res.status(400).json({
        error: 'Validation failed',
        details: errors.array()
      })
    }

    try {
      // Check if user exists
      const existingUser = await User.findOne({ email: req.body.email })
      if (existingUser) {
        return res.status(409).json({ error: 'Email already registered' })
      }

      // Hash password
      const passwordHash = await bcrypt.hash(req.body.password, 12)

      // Create user
      const user = new User({
        email: req.body.email,
        passwordHash,
        name: req.body.name
      })

      await user.save()

      // Return success (without sensitive data)
      res.status(201).json({
        id: user._id,
        email: user.email,
        name: user.name,
        createdAt: user.createdAt
      })
    } catch (error) {
      console.error('Create user error:', error)
      res.status(500).json({ error: 'Internal server error' })
    }
  }
)

module.exports = router
` + bt + bt + bt + `

## EXAMPLE 2: DATABASE MIGRATION

**Request:** "Create PostgreSQL migration for users table"

**Expert Response:**

` + bt + bt + bt + `sql
-- migrations/001_create_users_table.sql

-- Up migration
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  role VARCHAR(50) DEFAULT 'user',
  email_verified BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for common queries
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();
` + bt + bt + bt + `

## EXAMPLE 3: IMPLEMENT RATE LIMITING

**Request:** "Add rate limiting to prevent API abuse"

**Expert Response:**

` + bt + bt + bt + `javascript
// middleware/rateLimiter.js
const rateLimit = require('express-rate-limit')
const RedisStore = require('rate-limit-redis')
const Redis = require('ioredis')

const redis = new Redis(process.env.REDIS_URL)

// General API rate limiter
const apiLimiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100, // Limit each IP to 100 requests per windowMs
  store: new RedisStore({
    sendCommand: (...args) => redis.call(...args),
  }),
  standardHeaders: true,
  legacyHeaders: false,
  message: { error: 'Too many requests, please try again later' }
})

// Stricter limiter for auth endpoints
const authLimiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 5, // Limit each IP to 5 requests per windowMs
  store: new RedisStore({
    sendCommand: (...args) => redis.call(...args),
  }),
  message: { error: 'Too many authentication attempts' }
})

module.exports = { apiLimiter, authLimiter }

// Usage in app.js
app.use('/api/', apiLimiter)
app.use('/api/auth/', authLimiter)
` + bt + bt + bt + `
`
