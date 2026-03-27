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

// QAEngineerSkill implements the native skill for QA engineer role.
type QAEngineerSkill struct {
	workspace string
}

// NewQAEngineerSkill creates a new QAEngineerSkill instance.
func NewQAEngineerSkill(workspace string) *QAEngineerSkill {
	return &QAEngineerSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (q *QAEngineerSkill) Name() string {
	return "qa_engineer"
}

// Description returns a brief description of the skill.
func (q *QAEngineerSkill) Description() string {
	return "QA expert: testing strategies, test automation, coverage analysis, quality gates."
}

// GetInstructions returns the complete QA engineering protocol for the LLM.
func (q *QAEngineerSkill) GetInstructions() string {
	return qaEngineerInstructions
}

// GetAntiPatterns returns common QA anti-patterns to avoid.
func (q *QAEngineerSkill) GetAntiPatterns() string {
	return qaEngineerAntiPatterns
}

// GetExamples returns concrete QA engineering examples.
func (q *QAEngineerSkill) GetExamples() string {
	return qaEngineerExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (q *QAEngineerSkill) BuildSkillContext() string {
	parts := make([]string, 0, 11)

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "✅ NATIVE SKILL: QA Engineer")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(
		parts,
		"**ROLE:** Expert QA Engineer specializing in test automation, quality assurance, and continuous testing.",
	)
	parts = append(parts, "")
	parts = append(parts, q.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, q.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, q.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (q *QAEngineerSkill) BuildSummary() string {
	return `<skill name="qa_engineer" type="native">
  <purpose>QA expert — test automation, coverage analysis, quality gates</purpose>
  <pattern>Use for test strategy, test automation, quality assurance, E2E testing</pattern>
  <stacks>Jest, Playwright, Cypress, pytest, k6, Selenium</stacks>
  <practices>Test Pyramid, TDD, BDD, Continuous Testing</practices>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

const qaEngineerInstructions = `## CORE RESPONSIBILITIES

### 1. Test Strategy
- Define testing pyramid approach
- Identify test scenarios
- Prioritize test cases
- Estimate testing effort
- Define quality metrics

### 2. Test Automation
- Write unit tests
- Develop integration tests
- Create E2E test suites
- Implement visual regression tests
- Build performance tests

### 3. Test Infrastructure
- Set up test environments
- Configure CI/CD integration
- Manage test data
- Implement test reporting
- Maintain test frameworks

### 4. Quality Gates
- Define acceptance criteria
- Implement automated quality checks
- Monitor code coverage
- Track defect metrics
- Report quality status

### 5. Exploratory Testing
- Perform manual testing sessions
- Identify edge cases
- Document bugs clearly
- Verify fixes
- Share learnings

## TECHNOLOGY STACK

### Testing Frameworks
- **Unit**: Jest, Vitest, pytest, Go testing
- **Integration**: Supertest, pytest, httptest
- **E2E**: Playwright, Cypress, Selenium
- **Performance**: k6, JMeter, Locust
- **Visual**: Percy, Chromatic

### Coverage Tools
- Istanbul/nyc, coverage.py, gocov

### Test Management
- TestRail, Xray, Zephyr

## BEST PRACTICES

### Test Pyramid
` + bt + bt + bt + `
        /\\
       /  \\      E2E Tests (10%)
      /----\\
     /      \\    Integration Tests (20%)
    /--------\\
   /          \\  Unit Tests (70%)
  /------------\\
` + bt + bt + bt + `

### AAA Pattern
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

### Test Data Management
- Use factories, not fixtures
- Isolate test data
- Clean up after tests
- Use realistic data

## QUALITY METRICS

### Coverage Targets
- Unit tests: 80%+ line coverage
- Integration tests: Critical paths covered
- E2E tests: Happy paths + critical edge cases

### Defect Metrics
- Defect density
- Mean time to detection
- Defect escape rate
- Bug fix time

## QUALITY CHECKLIST

Before considering a feature complete:

- [ ] Unit tests written and passing
- [ ] Integration tests covering API
- [ ] E2E test for happy path
- [ ] Edge cases tested
- [ ] Error scenarios tested
- [ ] Performance acceptable
- [ ] Accessibility verified
- [ ] Cross-browser tested (if UI)
- [ ] Documentation updated
- [ ] Code coverage meets target
`

const qaEngineerAntiPatterns = `## QA ANTI-PATTERNS

### ❌ Testing Implementation Details
` + bt + bt + bt + `javascript
// BAD - Tests break on refactoring
test('sets state correctly', () => {
  expect(component.state.value).toBe('x')
})

// GOOD - Test behavior
test('displays the correct value', () => {
  expect(screen.getByText('Expected Value')).toBeInTheDocument()
})
` + bt + bt + bt + `

### ❌ Interdependent Tests
` + bt + bt + bt + `javascript
// BAD - Test order matters
test('creates user', () => {
  global.userId = createUser()
})

test('updates user', () => {
  updateUser(global.userId) // Depends on previous test!
})

// GOOD - Each test is independent
test('creates user', () => {
  const user = createUser()
  expect(user.id).toBeDefined()
})

test('updates user', () => {
  const user = createUser()
  const updated = updateUser(user.id)
  expect(updated.name).toBe('New Name')
})
` + bt + bt + bt + `

### ❌ Flaky Tests
` + bt + bt + bt + `javascript
// BAD - Timing dependent, sometimes fails
test('loads data', async () => {
  clickButton()
  expect(screen.getByText('Loaded')).toBeInTheDocument() // May fail!
})

// GOOD - Wait for condition
test('loads data', async () => {
  clickButton()
  await waitFor(() => {
    expect(screen.getByText('Loaded')).toBeInTheDocument()
  })
})
` + bt + bt + bt + `

### ❌ No Assertions
` + bt + bt + bt + `javascript
// BAD - Test always passes
test('calls API', () => {
  api.getData()
})

// GOOD - Verify behavior
test('calls API', () => {
  api.getData()
  expect(api.get).toHaveBeenCalledTimes(1)
  expect(api.get).toHaveBeenCalledWith('/endpoint')
})
` + bt + bt + bt + `

### ❌ Skipping Tests
` + bt + bt + bt + `javascript
// BAD - Temporarily skipped becomes permanent
it.skip('important test', () => {})
xdescribe('critical suite', () => {})

// GOOD - Either fix test or remove
it('important test', () => {
  // Properly implemented
})
` + bt + bt + bt + `

### ❌ Manual Regression Only
` + bt + bt + bt + `
BAD:  All regression testing manual
      Takes 2 weeks per release
      Human error inevitable

GOOD: Automated regression suite
      Runs in 30 minutes
      Consistent, repeatable
` + bt + bt + bt + `

### ❌ Testing in Production Without Safeguards
` + bt + bt + bt + `javascript
// BAD - Tests affect real data
test('creates user', () => {
  await api.createUser({ email: 'test@test.com' })
})

// GOOD - Use isolated test environment
test('creates user', () => {
  // Test runs against test database
  // Data cleaned up after test
  await api.createUser({ email: 'test@test.com' })
})
` + bt + bt + bt + `
`

const qaEngineerExamples = `## EXAMPLE 1: WRITE UNIT TEST

**Request:** "Write unit tests for a user service"

**Expert Response:**

` + bt + bt + bt + `javascript
// tests/services/user.service.test.js
const UserService = require('../../src/services/user.service')
const UserRepository = require('../../src/repositories/user.repository')
const bcrypt = require('bcrypt')

jest.mock('../../src/repositories/user.repository')
jest.mock('bcrypt')

describe('UserService', () => {
  let userService
  let mockUserRepo

  beforeEach(() => {
    mockUserRepo = new UserRepository()
    UserService.prototype.repo = mockUserRepo
    userService = new UserService()
    jest.clearAllMocks()
  })

  describe('createUser', () => {
    it('should create a user successfully', async () => {
      // Arrange
      const userData = {
        email: 'test@example.com',
        password: 'password123',
        name: 'Test User'
      }

      const hashedPassword = 'hashed123'
      const savedUser = { id: 1, ...userData, passwordHash: hashedPassword }

      mockUserRepo.findByEmail.mockResolvedValue(null)
      bcrypt.hash.mockResolvedValue(hashedPassword)
      mockUserRepo.save.mockResolvedValue(savedUser)

      // Act
      const result = await userService.createUser(userData)

      // Assert
      expect(mockUserRepo.findByEmail).toHaveBeenCalledWith(userData.email)
      expect(bcrypt.hash).toHaveBeenCalledWith(userData.password, 12)
      expect(mockUserRepo.save).toHaveBeenCalledWith({
        email: userData.email,
        passwordHash: hashedPassword,
        name: userData.name
      })
      expect(result).toEqual({
        id: savedUser.id,
        email: savedUser.email,
        name: savedUser.name
      })
    })

    it('should throw error if email already exists', async () => {
      // Arrange
      const userData = { email: 'existing@example.com', password: 'pass', name: 'Test' }
      mockUserRepo.findByEmail.mockResolvedValue({ id: 999, email: userData.email })

      // Act & Assert
      await expect(userService.createUser(userData))
        .rejects
        .toThrow('Email already registered')

      expect(mockUserRepo.save).not.toHaveBeenCalled()
    })
  })

  describe('getUserById', () => {
    it('should return user without sensitive data', async () => {
      // Arrange
      const user = {
        id: 1,
        email: 'test@example.com',
        passwordHash: 'hash123',
        name: 'Test'
      }
      mockUserRepo.findById.mockResolvedValue(user)

      // Act
      const result = await userService.getUserById(1)

      // Assert
      expect(result).toEqual({
        id: user.id,
        email: user.email,
        name: user.name
      })
      expect(result.passwordHash).toBeUndefined()
    })

    it('should throw error if user not found', async () => {
      // Arrange
      mockUserRepo.findById.mockResolvedValue(null)

      // Act & Assert
      await expect(userService.getUserById(999))
        .rejects
        .toThrow('User not found')
    })
  })
})
` + bt + bt + bt + `

## EXAMPLE 2: WRITE E2E TEST

**Request:** "Write E2E test for user login flow"

**Expert Response:**

` + bt + bt + bt + `javascript
// tests/e2e/login.spec.js
const { test, expect } = require('@playwright/test')

test.describe('Login Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('should login successfully with valid credentials', async ({ page }) => {
    // Navigate to login
    await page.click('[data-testid="login-button"]')
    await expect(page).toHaveURL('/login')

    // Fill login form
    await page.fill('[name="email"]', 'test@example.com')
    await page.fill('[name="password"]', 'password123')

    // Submit
    await page.click('[type="submit"]')

    // Verify successful login
    await expect(page).toHaveURL('/dashboard')
    await expect(page.locator('[data-testid="user-menu"]')).toBeVisible()
    await expect(page.locator('[data-testid="welcome-message"]'))
      .toContainText('Welcome')
  })

  test('should show error with invalid credentials', async ({ page }) => {
    // Navigate to login
    await page.click('[data-testid="login-button"]')

    // Fill invalid credentials
    await page.fill('[name="email"]', 'wrong@example.com')
    await page.fill('[name="password"]', 'wrongpassword')

    // Submit
    await page.click('[type="submit"]')

    // Verify error message
    await expect(page.locator('[data-testid="error-message"]'))
      .toContainText('Invalid credentials')

    // Should stay on login page
    await expect(page).toHaveURL('/login')
  })

  test('should validate required fields', async ({ page }) => {
    // Navigate to login
    await page.click('[data-testid="login-button"]')

    // Submit empty form
    await page.click('[type="submit"]')

    // Verify validation errors
    await expect(page.locator('[data-testid="email-error"]')).toBeVisible()
    await expect(page.locator('[data-testid="password-error"]')).toBeVisible()
  })

  test('should redirect to dashboard if already logged in', async ({ page }) => {
    // Setup authenticated state
    await page.context().addCookies([{
      name: 'session',
      value: 'valid-session-token',
      domain: 'localhost',
      path: '/'
    }])

    // Try to access login page
    await page.goto('/login')

    // Should redirect to dashboard
    await expect(page).toHaveURL('/dashboard')
  })
})
` + bt + bt + bt + `
`
