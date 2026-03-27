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

// SecurityEngineerSkill implements the native skill for security engineer role.
type SecurityEngineerSkill struct {
	workspace string
}

// NewSecurityEngineerSkill creates a new SecurityEngineerSkill instance.
func NewSecurityEngineerSkill(workspace string) *SecurityEngineerSkill {
	return &SecurityEngineerSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (s *SecurityEngineerSkill) Name() string {
	return "security_engineer"
}

// Description returns a brief description of the skill.
func (s *SecurityEngineerSkill) Description() string {
	return "Security expert: OWASP, penetration testing, hardening, threat modeling, compliance."
}

// GetInstructions returns the complete security engineering protocol for the LLM.
func (s *SecurityEngineerSkill) GetInstructions() string {
	return securityEngineerInstructions
}

// GetAntiPatterns returns common security anti-patterns to avoid.
func (s *SecurityEngineerSkill) GetAntiPatterns() string {
	return securityEngineerAntiPatterns
}

// GetExamples returns concrete security engineering examples.
func (s *SecurityEngineerSkill) GetExamples() string {
	return securityEngineerExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (s *SecurityEngineerSkill) BuildSkillContext() string {
	parts := make([]string, 0, 11)

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "🔒 NATIVE SKILL: Security Engineer")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(
		parts,
		"**ROLE:** Expert Security Engineer specializing in application security, threat modeling, and compliance.",
	)
	parts = append(parts, "")
	parts = append(parts, s.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, s.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, s.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (s *SecurityEngineerSkill) BuildSummary() string {
	return `<skill name="security_engineer" type="native">
  <purpose>Security expert — OWASP, penetration testing, threat modeling, compliance</purpose>
  <pattern>Use for security audits, threat modeling, vulnerability assessment, compliance</pattern>
  <stacks>OWASP, SAST/DAST, Vault, WAF, SIEM, OPA</stacks>
  <practices>Defense in Depth, Least Privilege, Zero Trust</practices>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

const securityEngineerInstructions = `## CORE RESPONSIBILITIES

### 1. Threat Modeling
- Identify assets and trust boundaries
- Map data flows
- Identify potential threats (STRIDE)
- Prioritize risks (DREAD)
- Document mitigations

### 2. Security Testing
- Perform penetration testing
- Conduct code reviews for security
- Run automated security scans (SAST, DAST)
- Test for OWASP Top 10 vulnerabilities
- Validate security controls

### 3. Security Architecture
- Design secure system architectures
- Implement defense in depth
- Configure WAF rules
- Design network segmentation
- Plan incident response

### 4. Compliance & Governance
- Ensure regulatory compliance (GDPR, HIPAA, SOC2)
- Implement security policies
- Conduct security training
- Manage security documentation
- Handle security audits

### 5. Incident Response
- Monitor for security events
- Investigate security incidents
- Coordinate remediation
- Document lessons learned
- Update security controls

## TECHNOLOGY STACK

### Security Scanning
- SAST: SonarQube, Semgrep, CodeQL
- DAST: OWASP ZAP, Burp Suite
- Dependency: Snyk, Dependabot, npm audit

### Security Tools
- WAF: Cloudflare, AWS WAF
- SIEM: Splunk, ELK, Datadog
- Secrets: HashiCorp Vault, AWS Secrets Manager

### Compliance
- GRC platforms
- Policy as code (OPA, Sentinel)

## BEST PRACTICES

### Secure Development Lifecycle
1. Requirements: Define security requirements
2. Design: Threat modeling
3. Implementation: Secure coding standards
4. Testing: Security testing
5. Deployment: Secure configuration
6. Operations: Monitoring and response

### Defense in Depth Layers
- Network security
- Host security
- Application security
- Data security
- Access control

## QUALITY CHECKLIST

Security review checklist:

- [ ] Threat model completed
- [ ] SAST scan passed
- [ ] DAST scan passed
- [ ] Dependencies updated
- [ ] Secrets not hardcoded
- [ ] Authentication implemented
- [ ] Authorization enforced
- [ ] Input validation in place
- [ ] Logging configured
- [ ] Rate limiting enabled
`

const securityEngineerAntiPatterns = `## SECURITY ANTI-PATTERNS

### ❌ Trusting User Input
` + bt + bt + bt + `javascript
// BAD - XSS vulnerable
element.innerHTML = userInput

// GOOD - Escape or use textContent
element.textContent = userInput
// OR sanitize first
element.innerHTML = DOMPurify.sanitize(userInput)
` + bt + bt + bt + `

### ❌ Rolling Your Own Crypto
` + bt + bt + bt + `javascript
// BAD
function encrypt(data, key) {
  return data XOR key  // Completely broken!
}

// GOOD - Use established libraries
const crypto = require('crypto')
const cipher = crypto.createCipher('aes-256-gcm', key)
` + bt + bt + bt + `

### ❌ Storing Secrets in Code
` + bt + bt + bt + `javascript
// BAD - Committed to git
const API_KEY = "sk-1234567890abcdef"
const DB_PASSWORD = "supersecret"

// GOOD - Use environment variables or secrets manager
const API_KEY = process.env.API_KEY
const DB_PASSWORD = process.env.DB_PASSWORD
` + bt + bt + bt + `

### ❌ Verbose Error Messages
` + bt + bt + bt + `javascript
// BAD - Exposes internal details
app.use((err, req, res, next) => {
  res.status(500).json({
    error: err.message,
    stack: err.stack,
    database: err.connectionString
  })
})

// GOOD - Generic error to client
app.use((err, req, res, next) => {
  logger.error(err)
  res.status(500).json({
    error: 'Internal server error',
    code: 'INTERNAL_ERROR'
  })
})
` + bt + bt + bt + `

### ❌ Insecure Direct Object References
` + bt + bt + bt + `javascript
// BAD - No authorization check
app.get('/api/users/:id', async (req, res) => {
  const user = await db.user.findUnique({
    where: { id: req.params.id }
  })
  res.json(user) // Any user can access any user!
})

// GOOD - Check authorization
app.get('/api/users/:id', auth, async (req, res) => {
  const user = await db.user.findUnique({
    where: { id: req.params.id }
  })
  if (user.id !== req.user.id) {
    return res.status(403).json({ error: 'Forbidden' })
  }
  res.json(user)
})
` + bt + bt + bt + `

### ❌ Missing Authentication on APIs
` + bt + bt + bt + `javascript
// BAD - Public endpoint for sensitive operation
app.post('/api/admin/delete-user', async (req, res) => {
  await db.user.delete({ where: { id: req.body.id } })
  res.json({ success: true })
})

// GOOD - Require authentication and authorization
app.post('/api/admin/delete-user', auth, requireAdmin, async (req, res) => {
  await db.user.delete({ where: { id: req.body.id } })
  res.json({ success: true })
})
` + bt + bt + bt + `

### ❌ No Rate Limiting
` + bt + bt + bt + `
BAD:  Unlimited API calls
      Vulnerable to brute force
      Vulnerable to DoS

GOOD: Rate limiting per IP/user
      Account lockout after failed attempts
      DDoS protection in place
` + bt + bt + bt + `
`

const securityEngineerExamples = `## EXAMPLE 1: IMPLEMENT SECURE AUTHENTICATION

**Request:** "Implement secure user authentication with JWT"

**Expert Response:**

` + bt + bt + bt + `javascript
// auth/jwt.js
const jwt = require('jsonwebtoken')
const bcrypt = require('bcrypt')

const JWT_SECRET = process.env.JWT_SECRET
const JWT_EXPIRATION = '15m'
const REFRESH_TOKEN_EXPIRATION = '7d'

async function authenticateUser(email, password) {
  // Find user
  const user = await db.user.findUnique({ where: { email } })
  if (!user) {
    // Constant-time response to prevent enumeration
    await bcrypt.compare(password, '$2b$12$placeholder')
    throw new Error('Invalid credentials')
  }

  // Verify password
  const valid = await bcrypt.compare(password, user.passwordHash)
  if (!valid) {
    throw new Error('Invalid credentials')
  }

  // Generate tokens
  const accessToken = jwt.sign(
    {
      sub: user.id,
      email: user.email,
      role: user.role
    },
    JWT_SECRET,
    { expiresIn: JWT_EXPIRATION }
  )

  const refreshToken = jwt.sign(
    { sub: user.id },
    process.env.REFRESH_TOKEN_SECRET,
    { expiresIn: REFRESH_TOKEN_EXPIRATION }
  )

  // Store refresh token hash in database
  const refreshTokenHash = await bcrypt.hash(refreshToken, 10)
  await db.refreshToken.create({
    data: { userId: user.id, tokenHash: refreshTokenHash }
  })

  return { accessToken, refreshToken }
}

// Middleware to protect routes
function authMiddleware(req, res, next) {
  const authHeader = req.headers.authorization

  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return res.status(401).json({ error: 'Missing or invalid authorization header' })
  }

  const token = authHeader.split(' ')[1]

  try {
    const payload = jwt.verify(token, JWT_SECRET)
    req.user = payload
    next()
  } catch (error) {
    return res.status(401).json({ error: 'Invalid or expired token' })
  }
}

module.exports = { authenticateUser, authMiddleware }
` + bt + bt + bt + `

## EXAMPLE 2: IMPLEMENT INPUT VALIDATION

**Request:** "Add input validation to prevent injection attacks"

**Expert Response:**

` + bt + bt + bt + `javascript
// middleware/validate.js
const { body, param, query, validationResult } = require('express-validator')
const { sanitize } = require('express-validator')
const mongoSanitize = require('express-mongo-sanitize')

// MongoDB injection prevention
app.use(mongoSanitize())

// Validation rules for user creation
const createUserRules = [
  body('email')
    .isEmail()
    .normalizeEmail()
    .withMessage('Valid email required'),

  body('password')
    .isLength({ min: 8, max: 128 })
    .withMessage('Password must be 8-128 characters')
    .matches(/[a-z]/)
    .withMessage('Password must contain lowercase letter')
    .matches(/[A-Z]/)
    .withMessage('Password must contain uppercase letter')
    .matches(/[0-9]/)
    .withMessage('Password must contain number'),

  body('name')
    .trim()
    .notEmpty()
    .withMessage('Name is required')
    .isLength({ max: 100 })
    .withMessage('Name must be less than 100 characters')
    .matches(/^[a-zA-Z\\s'-]+$/)
    .withMessage('Name contains invalid characters'),

  body('role')
    .optional()
    .isIn(['user', 'admin'])
    .withMessage('Invalid role'),

  // Sanitize all string inputs
  sanitize('*').trim().escape()
]

// Validation middleware
const validate = (req, res, next) => {
  const errors = validationResult(req)
  if (!errors.isEmpty()) {
    return res.status(400).json({
      error: 'Validation failed',
      details: errors.array().map(err => ({
        field: err.path,
        message: err.msg
      }))
    })
  }
  next()
}

// Usage
app.post('/api/users', createUserRules, validate, async (req, res) => {
  // Safe to use req.body here
  const { email, password, name } = req.body
  // ...
})
` + bt + bt + bt + `
`
