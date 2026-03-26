#!/bin/bash
# PicoClaw-Agents Release Script
# Usage: ./scripts/create-release.sh <VERSION> [PRERELEASE]
# Example: ./scripts/create-release.sh v1.0.1 false
#          ./scripts/create-release.sh v1.1.0-beta true

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Arguments
VERSION=${1:-""}
PRERELEASE=${2:-false}

# Helper functions
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Validate arguments
if [ -z "$VERSION" ]; then
    print_error "Error: Debes especificar la versión"
    echo ""
    echo "Uso: $0 <VERSION> [PRERELEASE]"
    echo "Ejemplos:"
    echo "  $0 v1.0.1 false          # Release normal"
    echo "  $0 v1.1.0-beta true      # Pre-release"
    echo "  $0 v1.1.0-rc1 true       # Release candidate"
    exit 1
fi

# Validate version format (SemVer)
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
    print_error "Error: Formato de versión inválido"
    echo "El formato debe ser: vMAJOR.MINOR.PATCH o vMAJOR.MINOR.PATCH-LABEL"
    echo "Ejemplos: v1.0.1, v1.1.0, v2.0.0-beta, v1.0.0-rc1"
    exit 1
fi

echo ""
print_info "╔═══════════════════════════════════════════════════════════╗"
print_info "║     PicoClaw-Agents Release Script                        ║"
print_info "╚═══════════════════════════════════════════════════════════╝"
echo ""

# Step 1: Check prerequisites
print_info "Paso 1/6: Verificando prerequisitos..."

if ! command -v git &> /dev/null; then
    print_error "git no está instalado"
    exit 1
fi

if ! command -v gh &> /dev/null; then
    print_error "GitHub CLI (gh) no está instalado"
    print_info "Instala con: brew install gh (macOS) o visita https://cli.github.com/"
    exit 1
fi

if ! command -v make &> /dev/null; then
    print_error "make no está instalado"
    exit 1
fi

print_success "Prerequisitos verificados (git, gh, make)"
echo ""

# Step 2: Check current branch
print_info "Paso 2/6: Verificando branch..."

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ]; then
    print_warning "Estás en la branch '$CURRENT_BRANCH', no en 'main'"
    read -p "¿Continuar de todas formas? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Operación cancelada"
        exit 0
    fi
fi

# Check for uncommitted changes
UNCOMMITTED=$(git status --porcelain)
if [ -n "$UNCOMMITTED" ]; then
    print_warning "Hay cambios sin commitear"
    echo "$UNCOMMITTED"
    read -p "¿Continuar de todas formas? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Operación cancelada"
        exit 0
    fi
fi

print_success "Branch verificada: $CURRENT_BRANCH"
echo ""

# Step 3: Verify last version
print_info "Paso 3/6: Verificando última versión..."

LAST_VERSION=$(git ls-remote --tags origin | tail -1 | cut -f2 | sed 's|refs/tags/||')
if [ -z "$LAST_VERSION" ]; then
    print_warning "No se encontraron tags remotos"
    LAST_VERSION="v1.0.0"
fi

print_info "Última versión: $LAST_VERSION"
print_info "Nueva versión: $VERSION"

# Determine version type
if [[ $VERSION == *"-beta"* ]] || [[ $VERSION == *"-alpha"* ]]; then
    print_info "Tipo: Pre-release"
elif [[ $VERSION == *"-rc"* ]]; then
    print_info "Tipo: Release Candidate"
else
    print_info "Tipo: Release estable"
fi

read -p "¿Continuar con la creación de esta versión? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_info "Operación cancelada"
    exit 0
fi
echo ""

# Step 4: Run verifications
print_info "Paso 4/6: Ejecutando verificaciones de código..."

print_info "Ejecutando make check..."
if ! make check > /dev/null 2>&1; then
    print_error "make check falló"
    print_info "Ejecuta 'make check' manualmente para ver los errores"
    exit 1
fi
print_success "make check passed"

print_info "Ejecutando make test..."
if ! make test > /dev/null 2>&1; then
    print_error "make test falló"
    print_info "Ejecuta 'make test' manualmente para ver los errores"
    exit 1
fi
print_success "make test passed"

print_info "Ejecutando make security-check..."
if ! make security-check > /dev/null 2>&1; then
    print_warning "make security-check encontró problemas"
    print_info "Revisa manualmente con 'make security-check'"
fi
print_success "Security check completed"
echo ""

# Step 5: Create and push tag
print_info "Paso 5/6: Creando y subiendo tag..."

# Check if tag already exists
TAG_EXISTS=$(git tag -l "$VERSION")
if [ -n "$TAG_EXISTS" ]; then
    print_error "El tag $VERSION ya existe localmente"
    exit 1
fi

REMOTE_TAG_EXISTS=$(git ls-remote --tags origin "$VERSION")
if [ -n "$REMOTE_TAG_EXISTS" ]; then
    print_error "El tag $VERSION ya existe en el remoto"
    exit 1
fi

print_info "Creando tag $VERSION..."
git tag -a "$VERSION" -m "Release $VERSION"

print_info "Subiendo tag $VERSION..."
git push origin "$VERSION"

print_success "Tag $VERSION creado y subido exitosamente"
echo ""

# Step 6: Trigger release workflow
print_info "Paso 6/6: Disparando release workflow..."

if gh workflow run release.yml \
    --field tag="$VERSION" \
    --field prerelease="$PRERELEASE" \
    --field draft=false; then
    print_success "Release workflow disparado exitosamente"
else
    print_warning "No se pudo disparar el workflow"
    print_info "Puedes dispararlo manualmente en:"
    print_info "https://github.com/comgunner/picoclaw-agents/actions/workflows/release.yml"
fi
echo ""

# Summary
echo ""
print_info "╔═══════════════════════════════════════════════════════════╗"
print_info "║                    Resumen                               ║"
print_info "╚═══════════════════════════════════════════════════════════╝"
echo ""
print_success "Versión $VERSION creada exitosamente"
echo ""
print_info "Enlaces útiles:"
print_info "  - Releases: https://github.com/comgunner/picoclaw-agents/releases/tag/$VERSION"
print_info "  - Actions: https://github.com/comgunner/picoclaw-agents/actions"
print_info "  - Compare: https://github.com/comgunner/picoclaw-agents/compare/$LAST_VERSION...$VERSION"
echo ""

if [ "$PRERELEASE" = "true" ]; then
    print_warning "⚠️  Esto es un pre-release"
else
    print_success "✅ Release estable publicado"
fi

echo ""
print_info "Próximos pasos:"
print_info "  1. Monitorea el workflow en GitHub Actions"
print_info "  2. Verifica que los binarios se generaron correctamente"
print_info "  3. Actualiza CHANGELOG.md si es necesario"
print_info "  4. Anuncia el release en issues/PRs relacionados"
echo ""
