// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)

package skills

import (
	"strings"
)

// FullStackDeveloperSkill implements native skill for full-stack development assistance.
// It provides comprehensive development patterns, best practices, and code examples.
type FullStackDeveloperSkill struct {
	workspace string
}

// NewFullStackDeveloperSkill creates a new FullStackDeveloperSkill instance.
func NewFullStackDeveloperSkill(workspace string) *FullStackDeveloperSkill {
	return &FullStackDeveloperSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (f *FullStackDeveloperSkill) Name() string {
	return "fullstack_developer"
}

// Description returns a brief description of the skill.
func (f *FullStackDeveloperSkill) Description() string {
	return "Expert full-stack development assistant with patterns for frontend, backend, database, testing, DevOps, and security."
}

// GetInstructions returns the complete development guidelines.
func (f *FullStackDeveloperSkill) GetInstructions() string {
	return fullstackDeveloperInstructions
}

// GetAntiPatterns returns common development anti-patterns.
func (f *FullStackDeveloperSkill) GetAntiPatterns() string {
	return fullstackDeveloperAntiPatterns
}

// GetExamples returns concrete development examples.
func (f *FullStackDeveloperSkill) GetExamples() string {
	return fullstackDeveloperExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (f *FullStackDeveloperSkill) BuildSkillContext() string {
	var parts []string

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "🚀 NATIVE SKILL: Full-Stack Developer")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(parts, "**PURPOSE:** Expert full-stack development assistance with industry best practices.")
	parts = append(parts, "")
	parts = append(parts, f.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, f.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, f.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (f *FullStackDeveloperSkill) BuildSummary() string {
	return `<skill name="fullstack_developer" type="native">
  <purpose>Expert full-stack development assistance</purpose>
  <pattern>Use for coding, debugging, architecture, best practices</pattern>
  <stacks>Frontend (React/Vue), Backend (Node/Python/Go), Database (SQL/NoSQL)</stacks>
  <practices>TDD, CI/CD, Security, Code Review</practices>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

const fullstackDeveloperInstructions = `## DEVELOPMENT WORKFLOW

### 1. Requirements Analysis
- Understand user story/requirement
- Identify acceptance criteria
- Consider edge cases

### 2. Architecture Design
- Choose appropriate stack
- Design API contracts
- Plan database schema

### 3. Implementation
- Write tests first (TDD)
- Implement feature
- Refactor for clarity

### 4. Code Review
- Check for bugs
- Verify best practices
- Ensure documentation

### 5. Deployment
- CI/CD pipeline
- Environment configuration
- Monitoring setup

## FRONTEND PATTERNS

### Component Structure (React)
` + bt + bt + bt + `jsx
function Component({ prop1, prop2 }) {
  // 1. Hooks
  const [state, setState] = useState(initialValue)

  // 2. Effects
  useEffect(() => {
    // Side effects
  }, [dependencies])

  // 3. Handlers
  const handleClick = () => {
    // Logic
  }

  // 4. Render
  return <JSX />
}
` + bt + bt + bt + `

### State Management
- Local state: ` + bt + `useState` + bt + `
- Global state: Redux/Zustand/Context
- Server state: React Query/SWR

### Event Handling Pattern
` + bt + bt + bt + `jsx
// Parent component
function Parent() {
  const handleSave = async (data) => {
    try {
      await api.save(data)
      // Handle success
    } catch (error) {
      // Handle error
    }
  }

  return <Child onSave={handleSave} />
}

// Child component
function Child({ onSave }) {
  const handleSubmit = () => {
    onSave(formData)
  }

  return <button onClick={handleSubmit}>Save</button>
}
` + bt + bt + bt + `

## BACKEND PATTERNS

### REST API Structure
` + bt + bt + bt + `
GET    /resource          - List
POST   /resource          - Create
GET    /resource/:id      - Get one
PUT    /resource/:id      - Update
DELETE /resource/:id      - Delete
` + bt + bt + bt + `

### Error Handling (Node.js)
` + bt + bt + bt + `javascript
try {
  // Operation
} catch (error) {
  // Handle error
  logger.error(error)
  return res.status(500).json({ error: 'Message' })
}
` + bt + bt + bt + `

### Middleware Pattern (Express)
` + bt + bt + bt + `javascript
// Authentication middleware
const auth = (req, res, next) => {
  const token = req.headers.authorization?.split(' ')[1]
  if (!token) return res.status(401).json({ error: 'Unauthorized' })

  try {
    req.user = jwt.verify(token, process.env.JWT_SECRET)
    next()
  } catch (error) {
    res.status(401).json({ error: 'Invalid token' })
  }
}

// Usage
router.get('/protected', auth, handler)
` + bt + bt + bt + `

## DATABASE PATTERNS

### Repository Pattern
` + bt + bt + bt + `typescript
interface IRepository<T> {
  findById(id: string): Promise<T | null>
  findAll(): Promise<T[]>
  create(data: Partial<T>): Promise<T>
  update(id: string, data: Partial<T>): Promise<T>
  delete(id: string): Promise<void>
}

class UserRepository implements IRepository<User> {
  async findById(id: string): Promise<User | null> {
    return db.user.findUnique({ where: { id } })
  }
  // ... other methods
}
` + bt + bt + bt + `

### SQL Best Practices
- Use transactions for multi-step operations
- Index frequently queried columns
- Use prepared statements (prevent SQL injection)
- Implement migrations

### NoSQL Best Practices
- Design for query patterns
- Embed vs. reference decisions
- Handle eventual consistency

## TESTING PATTERNS

### Unit Test Structure (AAA)
` + bt + bt + bt + `javascript
describe('Component', () => {
  it('should do something', () => {
    // Arrange
    const input = 'value'

    // Act
    const result = function(input)

    // Assert
    expect(result).toBe(expected)
  })
})
` + bt + bt + bt + `

### Test Pyramid
- **Unit Tests** (70%) - Fast, isolated
- **Integration Tests** (20%) - Component interaction
- **E2E Tests** (10%) - Full user flows

### Mocking Pattern
` + bt + bt + bt + `javascript
// Mock API call
jest.mock('../api')
api.fetchUser.mockResolvedValue({ id: 1, name: 'John' })

// Mock timer
jest.useFakeTimers()
jest.runAllTimers()
` + bt + bt + bt + `

## SECURITY CHECKLIST

### Input/Output
- [ ] Input validation (all user inputs)
- [ ] Output encoding (prevent XSS)
- [ ] Sanitize HTML content
- [ ] Validate file uploads

### Authentication/Authorization
- [ ] Authentication (JWT/OAuth2)
- [ ] Authorization (role-based access)
- [ ] Session management
- [ ] Password hashing (bcrypt, argon2)

### Infrastructure
- [ ] Rate limiting (prevent abuse)
- [ ] HTTPS (encrypt in transit)
- [ ] Secrets management (env vars, vaults)
- [ ] Dependency scanning (npm audit, etc.)
- [ ] CORS configuration

### Database
- [ ] SQL injection prevention (prepared statements)
- [ ] NoSQL injection prevention
- [ ] Query parameterization
- [ ] Access control

## GIT WORKFLOW

### Branch Naming
- ` + bt + `feature/description` + bt + `
- ` + bt + `fix/description` + bt + `
- ` + bt + `hotfix/description` + bt + `

### Commit Messages (Conventional Commits)
` + bt + `type(scope): description` + bt + `

Types: ` + bt + `feat` + bt + `, ` + bt + `fix` + bt + `, ` + bt + `docs` + bt + `, ` + bt + `style` + bt + `, ` + bt + `refactor` + bt + `, ` + bt + `test` + bt + `, ` + bt + `chore` + bt + `

Example:
` + bt + bt + bt + `
feat(auth): add JWT token refresh endpoint

- Implement refresh token generation
- Add token expiration handling
- Update authentication middleware

Closes #123
` + bt + bt + bt + `

### Pull Request Guidelines
- Small, focused changes
- Descriptive title and description
- Link to related issues
- Include tests
- Request review from teammates

## CODE QUALITY

### Linting
- ESLint (JavaScript/TypeScript)
- Pylint (Python)
- golangci-lint (Go)

### Formatting
- Prettier (JavaScript/TypeScript/CSS)
- Black (Python)
- gofmt (Go)

### Type Safety
- TypeScript for JavaScript projects
- mypy for Python projects
- Static typing in Go

### Documentation
- JSDoc for JavaScript/TypeScript
- docstrings for Python
- godoc for Go
- README with setup instructions
`

const fullstackDeveloperAntiPatterns = `## CODE SMELLS

### ❌ Long Functions (>50 lines)
**Problem:** Hard to understand, test, reuse
**Solution:** Extract into smaller functions

` + bt + bt + bt + `javascript
// BAD
function processUser(userData) {
  // 100 lines of code doing everything
  // Validate, save, send email, log, etc.
}

// GOOD
function processUser(userData) {
  validateUser(userData)
  const user = saveUser(userData)
  sendWelcomeEmail(user)
  logUserCreation(user)
}
` + bt + bt + bt + `

### ❌ Deep Nesting (>3 levels)
**Problem:** Cognitive load, hard to follow
**Solution:** Early returns, extract methods

` + bt + bt + bt + `javascript
// BAD
if (user) {
  if (user.isActive) {
    if (user.hasPermission) {
      // Do something
    }
  }
}

// GOOD
if (!user || !user.isActive || !user.hasPermission) return
// Do something
` + bt + bt + bt + `

### ❌ Magic Numbers
**Problem:** Unclear meaning, hard to change
**Solution:** Named constants

` + bt + bt + bt + `javascript
// BAD
if (status === 1) { /* ... */ }
setTimeout(callback, 86400000)

// GOOD
const STATUS_ACTIVE = 1
const MILLISECONDS_PER_DAY = 86400000

if (status === STATUS_ACTIVE) { /* ... */ }
setTimeout(callback, MILLISECONDS_PER_DAY)
` + bt + bt + bt + `

### ❌ Duplicated Code
**Problem:** Maintenance nightmare
**Solution:** Extract common logic

## SECURITY ANTI-PATTERNS

### ❌ SQL Injection
` + bt + bt + bt + `javascript
// BAD
db.query("SELECT * FROM users WHERE id = " + userId)

// GOOD
db.query('SELECT * FROM users WHERE id = ?', [userId])
` + bt + bt + bt + `

### ❌ XSS Vulnerability
` + bt + bt + bt + `javascript
// BAD
element.innerHTML = userInput

// GOOD
element.textContent = userInput
// OR sanitize first
element.innerHTML = DOMPurify.sanitize(userInput)
` + bt + bt + bt + `

### ❌ Hardcoded Secrets
` + bt + bt + bt + `javascript
// BAD
const API_KEY = 'sk-1234567890'
const DB_PASSWORD = 'secret123'

// GOOD
const API_KEY = process.env.API_KEY
const DB_PASSWORD = process.env.DB_PASSWORD
` + bt + bt + bt + `

### ❌ Weak Password Hashing
` + bt + bt + bt + `javascript
// BAD
const hash = md5(password)
const hash = sha1(password)

// GOOD
const hash = await bcrypt.hash(password, 12)
const hash = await argon2.hash(password)
` + bt + bt + bt + `

## TESTING ANTI-PATTERNS

### ❌ Testing Implementation Details
` + bt + bt + bt + `javascript
// BAD
expect(component.state.value).toBe('x')

// GOOD
expect(screen.getByText('Expected')).toBeInTheDocument()
` + bt + bt + bt + `

### ❌ Skipping Tests
` + bt + bt + bt + `javascript
// BAD
it.skip('important test', () => {})
xdescribe('critical suite', () => {})

// GOOD
// Either implement test or remove placeholder
` + bt + bt + bt + `

### ❌ No Assertions
` + bt + bt + bt + `javascript
// BAD
it('does something', () => {
  someFunction()
})

// GOOD
it('does something', () => {
  const result = someFunction()
  expect(result).toBe(expected)
})
` + bt + bt + bt + `

## GIT ANTI-PATTERNS

### ❌ Committing to Main
**Problem:** Breaks CI/CD, no code review
**Solution:** Always use feature branches

### ❌ Vague Commit Messages
` + bt + bt + bt + `bash
# BAD
git commit -m "fix stuff"
git commit -m "WIP"

# GOOD
git commit -m "fix(auth): resolve JWT expiration bug"
git commit -m "feat(api): add user profile endpoint"
` + bt + bt + bt + `

### ❌ Large Commits
**Problem:** Hard to review, hard to revert
**Solution:** Small, atomic commits

### ❌ No Code Review
**Problem:** Bugs slip through, knowledge silos
**Solution:** Require PR reviews before merge

## API ANTI-PATTERNS

### ❌ Inconsistent Error Responses
` + bt + bt + bt + `javascript
// BAD - Inconsistent
{ error: 'Message' }
{ errors: [...] }
{ message: 'Error' }

// GOOD - Consistent
{
  error: {
    code: 'VALIDATION_ERROR',
    message: 'Invalid input',
    details: [...]
  }
}
` + bt + bt + bt + `

### ❌ No Rate Limiting
**Problem:** API abuse, DoS vulnerability
**Solution:** Implement rate limiting per user/IP

### ❌ Exposing Internal Errors
` + bt + bt + bt + `javascript
// BAD
return res.status(500).json({
  error: err.stack,
  database: err.connectionString
})

// GOOD
return res.status(500).json({
  error: 'Internal server error',
  code: 'INTERNAL_ERROR'
})
` + bt + bt + bt + `
`

const fullstackDeveloperExamples = `## EXAMPLE 1: CREATE REACT COMPONENT

**Request:** "Create a user profile card component"

**Expert Response:**

1. **Component File:** ` + bt + `src/components/UserProfile.jsx` + bt + `
` + bt + bt + bt + `jsx
import React, { useState, useEffect } from 'react'
import PropTypes from 'prop-types'
import './UserProfile.css'

function UserProfile({ userId, onEdit }) {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    async function fetchUser() {
      try {
        const response = await fetch("/api/users/" + userId)
        if (!response.ok) throw new Error('User not found')
        const data = await response.json()
        setUser(data)
      } catch (err) {
        setError(err.message)
      } finally {
        setLoading(false)
      }
    }
    fetchUser()
  }, [userId])

  if (loading) return <div className="loading">Loading...</div>
  if (error) return <div className="error">Error: {error}</div>
  if (!user) return null

  return (
    <div className="user-profile">
      <img src={user.avatar} alt={user.name} className="avatar" />
      <h2 className="name">{user.name}</h2>
      <p className="email">{user.email}</p>
      <button onClick={onEdit} className="edit-btn">Edit Profile</button>
    </div>
  )
}

UserProfile.propTypes = {
  userId: PropTypes.string.isRequired,
  onEdit: PropTypes.func
}

export default UserProfile
` + bt + bt + bt + `

2. **Test File:** ` + bt + `src/components/UserProfile.test.jsx` + bt + `
` + bt + bt + bt + `jsx
import { render, screen, waitFor } from '@testing-library/react'
import UserProfile from './UserProfile'

global.fetch = jest.fn()

describe('UserProfile', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('shows loading state initially', () => {
    render(<UserProfile userId="123" />)
    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })

  it('displays user data when loaded', async () => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        name: 'John Doe',
        email: 'john@example.com',
        avatar: 'url'
      })
    })

    render(<UserProfile userId="123" />)

    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument()
      expect(screen.getByText('john@example.com')).toBeInTheDocument()
    })
  })

  it('shows error state on failure', async () => {
    fetch.mockRejectedValueOnce(new Error('Not found'))

    render(<UserProfile userId="123" />)

    await waitFor(() => {
      expect(screen.getByText(/Error:/)).toBeInTheDocument()
    })
  })
})
` + bt + bt + bt + `

## EXAMPLE 2: CREATE REST API ENDPOINT

**Request:** "Create Express endpoint to create users"

**Expert Response:**

1. **Route Handler:** ` + bt + `routes/users.js` + bt + `
` + bt + bt + bt + `javascript
const express = require('express')
const router = express.Router()
const { body, validationResult } = require('express-validator')
const User = require('../models/User')
const { hashPassword } = require('../utils/crypto')

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
      const passwordHash = await hashPassword(req.body.password)

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

2. **Test:** ` + bt + `tests/routes/users.test.js` + bt + `
` + bt + bt + bt + `javascript
const request = require('supertest')
const app = require('../../app')
const User = require('../../models/User')

describe('POST /api/users', () => {
  beforeEach(async () => {
    await User.deleteMany({})
  })

  it('creates a new user', async () => {
    const response = await request(app)
      .post('/api/users')
      .send({
        email: 'test@example.com',
        password: 'password123',
        name: 'Test User'
      })

    expect(response.status).toBe(201)
    expect(response.body).toHaveProperty('id')
    expect(response.body.email).toBe('test@example.com')
    expect(response.body.name).toBe('Test User')
    expect(response.body).not.toHaveProperty('passwordHash')
  })

  it('validates email format', async () => {
    const response = await request(app)
      .post('/api/users')
      .send({
        email: 'invalid-email',
        password: 'password123',
        name: 'Test User'
      })

    expect(response.status).toBe(400)
    expect(response.body.error).toBe('Validation failed')
  })

  it('validates password length', async () => {
    const response = await request(app)
      .post('/api/users')
      .send({
        email: 'test@example.com',
        password: 'short',
        name: 'Test User'
      })

    expect(response.status).toBe(400)
  })

  it('prevents duplicate emails', async () => {
    await User.create({
      email: 'existing@example.com',
      passwordHash: 'hash',
      name: 'Existing'
    })

    const response = await request(app)
      .post('/api/users')
      .send({
        email: 'existing@example.com',
        password: 'password123',
        name: 'Duplicate'
      })

    expect(response.status).toBe(409)
  })
})
` + bt + bt + bt + `

## EXAMPLE 3: DATABASE MIGRATION

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
CREATE INDEX idx_users_created_at ON users(created_at);

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

-- Down migration (for rollback)
-- DROP TRIGGER IF EXISTS update_users_updated_at ON users;
-- DROP FUNCTION IF EXISTS update_updated_at_column();
-- DROP TABLE IF EXISTS users;
` + bt + bt + bt + `

## EXAMPLE 4: DOCKER CONFIGURATION

**Request:** "Create Docker setup for Node.js app"

**Expert Response:**

1. **Dockerfile:** ` + bt + `Dockerfile` + bt + `
` + bt + bt + bt + `dockerfile
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies (production only)
RUN npm ci --only=production

# Copy source code
COPY . .

# Build if needed
# RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \\
    adduser -S nodejs -u 1001

# Copy from builder
COPY --from=builder --chown=nodejs:nodejs /app .

# Switch to non-root user
USER nodejs

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \\
  CMD node healthcheck.js

# Start application
CMD ["node", "src/index.js"]
` + bt + bt + bt + `

2. **Docker Compose:** ` + bt + `docker-compose.yml` + bt + `
` + bt + bt + bt + `yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
      - DATABASE_URL=postgresql://user:pass@db:5432/mydb
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
      - POSTGRES_DB=mydb
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
` + bt + bt + bt + `

## EXAMPLE 5: GITHUB ACTIONS CI/CD

**Request:** "Create CI/CD pipeline for Node.js project"

**Expert Response:**

` + bt + bt + bt + `yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: test_db
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Lint
        run: npm run lint

      - name: Test
        run: npm test
        env:
          DATABASE_URL: postgresql://test:test@localhost:5432/test_db

      - name: Build
        run: npm run build

  deploy:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'

    steps:
      - uses: actions/checkout@v4

      - name: Deploy to production
        run: |
          echo "Deploying to production..."
          # Add deployment commands here
` + bt + bt + bt + `
`
