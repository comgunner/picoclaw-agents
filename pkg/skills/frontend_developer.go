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

// FrontendDeveloperSkill implements the native skill for frontend developer role.
type FrontendDeveloperSkill struct {
	workspace string
}

// NewFrontendDeveloperSkill creates a new FrontendDeveloperSkill instance.
func NewFrontendDeveloperSkill(workspace string) *FrontendDeveloperSkill {
	return &FrontendDeveloperSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (f *FrontendDeveloperSkill) Name() string {
	return "frontend_developer"
}

// Description returns a brief description of the skill.
func (f *FrontendDeveloperSkill) Description() string {
	return "Frontend development expert: React, Vue, performance, accessibility, design systems."
}

// GetInstructions returns the complete frontend development protocol for the LLM.
func (f *FrontendDeveloperSkill) GetInstructions() string {
	return frontendDeveloperInstructions
}

// GetAntiPatterns returns common frontend development anti-patterns to avoid.
func (f *FrontendDeveloperSkill) GetAntiPatterns() string {
	return frontendDeveloperAntiPatterns
}

// GetExamples returns concrete frontend development examples.
func (f *FrontendDeveloperSkill) GetExamples() string {
	return frontendDeveloperExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (f *FrontendDeveloperSkill) BuildSkillContext() string {
	parts := make([]string, 0, 11)

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "🎨 NATIVE SKILL: Frontend Developer")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(
		parts,
		"**ROLE:** Expert Frontend Developer specializing in responsive, accessible, performant user interfaces.",
	)
	parts = append(parts, "")
	parts = append(parts, f.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, f.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, f.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (f *FrontendDeveloperSkill) BuildSummary() string {
	return `<skill name="frontend_developer" type="native">
  <purpose>Frontend development expert — React, Vue, performance, accessibility</purpose>
  <pattern>Use for UI components, state management, responsive design, performance optimization</pattern>
  <stacks>React, Vue, Svelte, TypeScript, Tailwind CSS, Next.js</stacks>
  <practices>Component-driven, TDD, Mobile-first, Accessibility-first</practices>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

const frontendDeveloperInstructions = `## CORE RESPONSIBILITIES

### 1. Component Development
- Build reusable, composable components
- Implement proper prop types/interfaces
- Manage component state effectively
- Follow single responsibility principle
- Document component APIs

### 2. State Management
- Choose appropriate state solution (local vs global)
- Implement state machines for complex flows
- Handle async state (loading, error, success)
- Optimize re-renders (memoization)
- Manage server state (React Query, SWR)

### 3. Performance Optimization
- Minimize bundle size (code splitting, tree shaking)
- Optimize images and assets
- Implement lazy loading
- Use virtual scrolling for large lists
- Monitor Core Web Vitals

### 4. Accessibility (a11y)
- Ensure keyboard navigation
- Implement proper ARIA labels
- Maintain focus management
- Support screen readers
- Meet WCAG 2.1 AA standards

### 5. Responsive Design
- Mobile-first approach
- Fluid typography and layouts
- Touch-friendly interactions
- Cross-browser compatibility
- Progressive enhancement

## TECHNOLOGY STACK

### Frameworks
- **Primary**: React, Vue.js, Svelte
- **Secondary**: Angular, Solid

### Styling
- **CSS-in-JS**: Styled Components, Emotion
- **Utility-first**: Tailwind CSS
- **Preprocessors**: Sass, Less
- **CSS Frameworks**: shadcn/ui, Radix UI

### Build Tools
- Vite, Webpack, esbuild, Rollup

### Testing
- Unit: Jest, Vitest
- Component: React Testing Library, Vue Test Utils
- E2E: Playwright, Cypress

## BEST PRACTICES

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

### File Organization
` + bt + bt + bt + `
src/
├── components/     # Reusable components
├── pages/          # Route components
├── hooks/          # Custom hooks
├── utils/          # Helper functions
├── styles/         # Global styles
├── assets/         # Images, fonts
└── types/          # TypeScript types
` + bt + bt + bt + `

### State Management Hierarchy
1. Local component state (useState)
2. Shared state (Context, Zustand)
3. Server state (React Query, SWR)
4. URL state (query params)

## QUALITY CHECKLIST

Before considering a frontend feature complete:

- [ ] Component is responsive
- [ ] Accessibility tested (keyboard, screen reader)
- [ ] Cross-browser tested
- [ ] Performance metrics met (LCP, FID, CLS)
- [ ] Unit tests written
- [ ] E2E tests passing
- [ ] Loading states handled
- [ ] Error states handled
- [ ] Documented in Storybook

## PERFORMANCE METRICS

### Core Web Vitals Targets
- **LCP** (Largest Contentful Paint): < 2.5s
- **FID** (First Input Delay): < 100ms
- **CLS** (Cumulative Layout Shift): < 0.1

### Bundle Size Budgets
- Initial JS: < 170KB (gzipped)
- Initial CSS: < 50KB (gzipped)
- Total page weight: < 500KB
`

const frontendDeveloperAntiPatterns = `## FRONTEND ANTI-PATTERNS

### ❌ Prop Drilling
` + bt + bt + bt + `jsx
// BAD - Passing props through multiple levels
function App() {
  return <Layout user={user}><Content user={user} /></Layout>
}

// GOOD - Use Context or composition
function App() {
  return (
    <UserProvider value={user}>
      <Layout><Content /></Layout>
    </UserProvider>
  )
}
` + bt + bt + bt + `

### ❌ Over-fetching Data
` + bt + bt + bt + `javascript
// BAD - Fetching entire user object
const { data } = useQuery('user', () => fetch('/api/user'))
const name = data.name // Only need name

// GOOD - Fetch only what's needed
const { data } = useQuery('userName', () => fetch('/api/user?fields=name'))
` + bt + bt + bt + `

### ❌ Unnecessary Re-renders
` + bt + bt + bt + `jsx
// BAD - New object reference on every render
function Parent() {
  return <Child config={{ theme: 'dark' }} />
}

// GOOD - Memoize or use useState
function Parent() {
  const [config] = useState({ theme: 'dark' })
  return <Child config={config} />
}
` + bt + bt + bt + `

### ❌ Inline Functions in JSX Props
` + bt + bt + bt + `jsx
// BAD - Creates new function on every render
function List({ items }) {
  return items.map(item => (
    <div onClick={() => handleClick(item)} />
  ))
}

// GOOD - Memoize with useCallback
function List({ items }) {
  const handleClickMemo = useCallback(
    (item) => handleClick(item),
    [handleClick]
  )
  return items.map(item => (
    <div onClick={() => handleClickMemo(item)} />
  ))
}
` + bt + bt + bt + `

### ❌ Ignoring Accessibility
` + bt + bt + bt + `jsx
// BAD - Not keyboard accessible
<div onClick={handleClick}>Click me</div>

// GOOD - Use semantic HTML
<button onClick={handleClick}>Click me</button>

// Or add keyboard support
<div
  onClick={handleClick}
  onKeyDown={(e) => e.key === 'Enter' && handleClick()}
  tabIndex={0}
  role="button"
>
  Click me
</div>
` + bt + bt + bt + `

### ❌ Blocking the Main Thread
` + bt + bt + bt + `javascript
// BAD - Heavy computation blocks UI
function processLargeData(data) {
  return data.map(expensiveOperation)
}

// GOOD - Use Web Workers or requestIdleCallback
function processLargeData(data) {
  const worker = new Worker('worker.js')
  worker.postMessage(data)
  worker.onmessage = (e) => console.log(e.data)
}
` + bt + bt + bt + `

### ❌ Storing Sensitive Data in localStorage
` + bt + bt + bt + `javascript
// BAD - Accessible via XSS
localStorage.setItem('token', authToken)

// GOOD - Use httpOnly cookies
// Set cookie server-side with httpOnly flag
` + bt + bt + bt + `
`

const frontendDeveloperExamples = `## EXAMPLE 1: CREATE REUSABLE COMPONENT

**Request:** "Create a user profile card component with loading and error states"

**Expert Response:**

` + bt + bt + bt + `jsx
// components/UserProfile.jsx
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

  if (loading) return <div className="loading" role="status">Loading...</div>
  if (error) return <div className="error" role="alert">Error: {error}</div>
  if (!user) return null

  return (
    <article className="user-profile" aria-labelledby="user-name">
      <img src={user.avatar} alt={user.name + "'s profile"} className="avatar" />
      <h2 id="user-name" className="name">{user.name}</h2>
      <p className="email">{user.email}</p>
      {onEdit && (
        <button onClick={onEdit} className="edit-btn" aria-label="Edit profile">
          Edit Profile
        </button>
      )}
    </article>
  )
}

UserProfile.propTypes = {
  userId: PropTypes.string.isRequired,
  onEdit: PropTypes.func
}

export default UserProfile
` + bt + bt + bt + `

## EXAMPLE 2: IMPLEMENT FORM WITH VALIDATION

**Request:** "Create a registration form with client-side validation"

**Expert Response:**

` + bt + bt + bt + `jsx
// components/RegistrationForm.jsx
import React, { useState } from 'react'

function RegistrationForm({ onSubmit }) {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    confirmPassword: ''
  })
  const [errors, setErrors] = useState({})
  const [submitting, setSubmitting] = useState(false)

  const validate = () => {
    const newErrors = {}

    if (!formData.email.includes('@')) {
      newErrors.email = 'Valid email required'
    }

    if (formData.password.length < 8) {
      newErrors.password = 'Password must be at least 8 characters'
    }

    if (formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = 'Passwords do not match'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e) => {
    e.preventDefault()

    if (!validate()) return

    setSubmitting(true)
    try {
      await onSubmit(formData)
    } catch (error) {
      setErrors({ submit: error.message })
    } finally {
      setSubmitting(false)
    }
  }

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    })
    // Clear error when user starts typing
    if (errors[e.target.name]) {
      setErrors({ ...errors, [e.target.name]: null })
    }
  }

  return (
    <form onSubmit={handleSubmit} noValidate>
      <div>
        <label htmlFor="email">Email</label>
        <input
          id="email"
          name="email"
          type="email"
          value={formData.email}
          onChange={handleChange}
          aria-invalid={!!errors.email}
          aria-describedby={errors.email ? 'email-error' : undefined}
        />
        {errors.email && (
          <span id="email-error" className="error" role="alert">
            {errors.email}
          </span>
        )}
      </div>

      <div>
        <label htmlFor="password">Password</label>
        <input
          id="password"
          name="password"
          type="password"
          value={formData.password}
          onChange={handleChange}
          aria-invalid={!!errors.password}
          aria-describedby={errors.password ? 'password-error' : undefined}
        />
        {errors.password && (
          <span id="password-error" className="error" role="alert">
            {errors.password}
          </span>
        )}
      </div>

      <button type="submit" disabled={submitting}>
        {submitting ? 'Registering...' : 'Register'}
      </button>

      {errors.submit && (
        <div className="error" role="alert">{errors.submit}</div>
      )}
    </form>
  )
}

export default RegistrationForm
` + bt + bt + bt + `
`
