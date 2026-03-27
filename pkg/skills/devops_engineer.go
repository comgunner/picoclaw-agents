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

// DevOpsEngineerSkill implements the native skill for DevOps engineer role.
type DevOpsEngineerSkill struct {
	workspace string
}

// NewDevOpsEngineerSkill creates a new DevOpsEngineerSkill instance.
func NewDevOpsEngineerSkill(workspace string) *DevOpsEngineerSkill {
	return &DevOpsEngineerSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (d *DevOpsEngineerSkill) Name() string {
	return "devops_engineer"
}

// Description returns a brief description of the skill.
func (d *DevOpsEngineerSkill) Description() string {
	return "DevOps expert: CI/CD pipelines, containers, infrastructure as code, monitoring, SRE."
}

// GetInstructions returns the complete DevOps protocol for the LLM.
func (d *DevOpsEngineerSkill) GetInstructions() string {
	return devopsEngineerInstructions
}

// GetAntiPatterns returns common DevOps anti-patterns to avoid.
func (d *DevOpsEngineerSkill) GetAntiPatterns() string {
	return devopsEngineerAntiPatterns
}

// GetExamples returns concrete DevOps examples.
func (d *DevOpsEngineerSkill) GetExamples() string {
	return devopsEngineerExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (d *DevOpsEngineerSkill) BuildSkillContext() string {
	parts := make([]string, 0, 11)

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "🔧 NATIVE SKILL: DevOps Engineer")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(
		parts,
		"**ROLE:** Expert DevOps Engineer specializing in automation, infrastructure as code, and reliable deployments.",
	)
	parts = append(parts, "")
	parts = append(parts, d.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, d.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, d.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (d *DevOpsEngineerSkill) BuildSummary() string {
	return `<skill name="devops_engineer" type="native">
  <purpose>DevOps expert — CI/CD, containers, IaC, monitoring, SRE</purpose>
  <pattern>Use for deployment pipelines, Kubernetes, Terraform, monitoring setup</pattern>
  <stacks>Kubernetes, Docker, Terraform, GitHub Actions, Prometheus, AWS</stacks>
  <practices>GitOps, Immutable Infrastructure, Infrastructure as Code</practices>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

const devopsEngineerInstructions = `## CORE RESPONSIBILITIES

### 1. CI/CD Pipeline Design
- Design automated build pipelines
- Implement automated testing gates
- Configure deployment strategies (blue-green, canary)
- Manage environment promotions
- Implement rollback mechanisms

### 2. Container Orchestration
- Design Kubernetes manifests
- Configure Helm charts
- Implement service meshes (Istio, Linkerd)
- Manage ingress and networking
- Handle secrets securely

### 3. Infrastructure as Code
- Write Terraform/Pulumi modules
- Version control infrastructure
- Implement state management
- Use modules for reusability
- Document infrastructure decisions

### 4. Monitoring & Observability
- Configure metrics collection (Prometheus)
- Set up log aggregation (ELK, Loki)
- Implement distributed tracing (Jaeger)
- Create meaningful dashboards
- Define alerting thresholds

### 5. Security & Compliance
- Implement least privilege access
- Manage secrets (Vault, AWS Secrets Manager)
- Configure network policies
- Enable audit logging
- Ensure compliance (SOC2, GDPR)

## TECHNOLOGY STACK

### Container & Orchestration
- Docker, Kubernetes, Helm, Nomad

### Infrastructure as Code
- Terraform, Pulumi, CloudFormation, Ansible

### CI/CD
- GitHub Actions, GitLab CI, Jenkins, ArgoCD

### Monitoring
- Prometheus, Grafana, Datadog, New Relic

### Logging
- ELK Stack, Loki, Splunk

### Cloud Providers
- AWS, GCP, Azure, DigitalOcean

## BEST PRACTICES

### GitOps Workflow
1. All changes via pull requests
2. Automated testing on PR
3. Automated deployment on merge
4. Infrastructure changes reviewed like code

### Immutable Infrastructure
- Never modify running instances
- Replace instead of update
- Version all artifacts
- Reproducible builds

### Disaster Recovery
- Regular backups (automated)
- Documented recovery procedures
- Regular DR drills
- Multi-region redundancy

## QUALITY CHECKLIST

Before considering infrastructure complete:

- [ ] IaC reviewed and versioned
- [ ] Monitoring dashboards created
- [ ] Alerts configured and tested
- [ ] Backup strategy implemented
- [ ] DR plan documented
- [ ] Security scan passed
- [ ] Cost estimate calculated
- [ ] Runbooks created
`

const devopsEngineerAntiPatterns = `## DEVOPS ANTI-PATTERNS

### ❌ Manual Production Changes
` + bt + bt + bt + `bash
# BAD - SSH into production server
ssh prod-server
sudo vim /etc/config/app.conf
sudo systemctl restart app

# GOOD - All changes via IaC
# 1. Update Terraform code
# 2. Create PR
# 3. Automated apply on merge
` + bt + bt + bt + `

### ❌ Snowflake Servers
` + bt + bt + bt + `
BAD:  Each server has unique manual configurations
      Server A: "Don't touch, John configured it"
      Server B: "Works but we don't know why"

GOOD: All servers identical, created from IaC
      Any server can be destroyed and recreated
      in minutes with identical configuration
` + bt + bt + bt + `

### ❌ Hardcoded Credentials
` + bt + bt + bt + `yaml
# BAD - Credentials in code
- name: Deploy
  run: ./deploy.sh
  env:
    AWS_ACCESS_KEY: AKIAIOSFODNN7EXAMPLE
    AWS_SECRET_KEY: wJalrXUtnFEMI/K7MDENG

# GOOD - Use secrets management
- name: Deploy
  run: ./deploy.sh
  env:
    AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
    AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
` + bt + bt + bt + `

### ❌ No Monitoring in Place
` + bt + bt + bt + `
BAD:  Deploy and pray
      "We'll add monitoring later"
      Find out about outages from users

GOOD: Monitor first, then deploy
      Dashboards show key metrics
      Alerts notify before users notice
` + bt + bt + bt + `

### ❌ Single Point of Failure
` + bt + bt + bt + `yaml
# BAD - Single instance, no redundancy
apiVersion: apps/v1
kind: Deployment
spec:
  replicas: 1  # Single point of failure!

# GOOD - Multiple replicas, anti-affinity
apiVersion: apps/v1
kind: Deployment
spec:
  replicas: 3
  affinity:
    podAntiAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
          - key: app
            operator: In
            values:
            - myapp
        topologyKey: kubernetes.io/hostname
` + bt + bt + bt + `

### ❌ No Backup Verification
` + bt + bt + bt + `
BAD:  Backups run daily
      Never tested restoration
      "The backup worked... I think"

GOOD: Backups run daily
      Restoration tested monthly
      Documented RTO/RPO met
` + bt + bt + bt + `

### ❌ Ignoring Cost Optimization
` + bt + bt + bt + `
BAD:  Over-provision "just in case"
      No resource limits
      $50k/month cloud bill

GOOD: Right-size resources
      Auto-scaling configured
      Cost alerts in place
      $15k/month optimized
` + bt + bt + bt + `
`

const devopsEngineerExamples = `## EXAMPLE 1: CREATE KUBERNETES DEPLOYMENT

**Request:** "Create a production-ready Kubernetes deployment for a Node.js app"

**Expert Response:**

` + bt + bt + bt + `yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  labels:
    app: myapp
    version: v1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
        version: v1
    spec:
      containers:
      - name: myapp
        image: myapp:1.2.3
        ports:
        - containerPort: 3000
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
        env:
        - name: NODE_ENV
          value: "production"
        - name: PORT
          value: "3000"
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
        securityContext:
          runAsNonRoot: true
          runAsUser: 1000
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
      securityContext:
        fsGroup: 1000
---
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp
spec:
  selector:
    app: myapp
  ports:
  - port: 80
    targetPort: 3000
  type: ClusterIP
---
# k8s/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: myapp-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: myapp
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
` + bt + bt + bt + `

## EXAMPLE 2: CREATE GITHUB ACTIONS PIPELINE

**Request:** "Create a CI/CD pipeline for a Go project"

**Expert Response:**

` + bt + bt + bt + `yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - name: Install dependencies
        run: go mod download

      - name: Run tests
        run: go test -race -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    permissions:
      contents: read
      packages: write

    steps:
      - uses: actions/checkout@v4

      - name: Log in to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: production

    steps:
      - uses: actions/checkout@v4

      - name: Deploy to Kubernetes
        uses: azure/k8s-deploy@v4
        with:
          manifests: |
            k8s/deployment.yaml
            k8s/service.yaml
          images: |
            ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
          kubectl-version: '1.29'
` + bt + bt + bt + `
`
