#!/bin/bash
# PicoClaw-Agents Release Script
#
# Usage:
#   Manual:       ./scripts/create-release.sh <VERSION> [PRERELEASE]
#   Interactive:  ./scripts/create-release.sh
#   Auto:         ./scripts/create-release.sh auto
#   Dry-run:      ./scripts/create-release.sh --dry-run
#
# Examples:
#   ./scripts/create-release.sh v1.0.1 false        # Manual release
#   ./scripts/create-release.sh v1.1.0-beta true    # Manual pre-release
#   ./scripts/create-release.sh                     # Interactive (analyze & suggest)
#   ./scripts/create-release.sh auto                # Auto (analyze & create)
#   ./scripts/create-release.sh --dry-run           # Test (no changes)
#
# Features:
#   - Validates version format (SemVer)
#   - Analyzes commits to suggest version (interactive/auto mode)
#   - Runs pre-release checks (make check, test, security-check)
#   - Creates and pushes git tags
#   - Triggers GitHub Actions release workflow
#   - Supports pre-releases (beta, alpha, rc)
#   - Dry-run mode for testing

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

# Analyze commits and suggest version
analyze_commits() {
    print_info "Analizando commits desde $LAST_VERSION..."

    # Get commits
    COMMITS=$(git log "$LAST_VERSION"..HEAD --oneline 2>/dev/null || echo "")

    if [ -z "$COMMITS" ]; then
        print_warning "No hay commits desde $LAST_VERSION"
        SUGGESTED_VERSION="$LAST_VERSION"
        CHANGE_TYPE="NONE"
        return
    fi

    COMMIT_COUNT=$(echo "$COMMITS" | wc -l | tr -d ' ')
    print_info "Encontrados $COMMIT_COUNT commits"

    # Analyze commit messages
    FEAT_COUNT=$(git log "$LAST_VERSION"..HEAD --oneline --grep="^feat" | wc -l | tr -d ' ')
    FIX_COUNT=$(git log "$LAST_VERSION"..HEAD --oneline --grep="^fix" | wc -l | tr -d ' ')
    BREAKING_COUNT=$(git log "$LAST_VERSION"..HEAD --oneline --grep="BREAKING" | wc -l | tr -d ' ')
    DOCS_COUNT=$(git log "$LAST_VERSION"..HEAD --oneline --grep="^docs" | wc -l | tr -d ' ')
    NATIVE_SKILLS=$(git log "$LAST_VERSION"..HEAD --oneline --grep="[Nn]ative [Ss]kill" | wc -l | tr -d ' ')
    NEW_PROVIDER=$(git log "$LAST_VERSION"..HEAD --oneline --grep="[Pp]rovider" | wc -l | tr -d ' ')

    echo ""
    print_info "Cambios detectados:"
    print_info "  - Features (feat): $FEAT_COUNT"
    print_info "  - Bug fixes (fix): $FIX_COUNT"
    print_info "  - Breaking changes: $BREAKING_COUNT"
    print_info "  - Docs: $DOCS_COUNT"
    print_info "  - Native Skills: $NATIVE_SKILLS"
    print_info "  - Nuevos Proveedores: $NEW_PROVIDER"
    echo ""

    # Determine version type
    if [ "$BREAKING_COUNT" -gt 0 ]; then
        CHANGE_TYPE="MAJOR"
        NEW_MINOR=0
        NEW_PATCH=0
        print_info "➡️  BREAKING CHANGE detectado"
    elif [ "$NATIVE_SKILLS" -gt 0 ] || [ "$NEW_PROVIDER" -gt 0 ] || [ "$FEAT_COUNT" -gt 0 ]; then
        CHANGE_TYPE="MINOR"
        NEW_MINOR=$((LAST_MINOR + 1))
        NEW_PATCH=0
        print_info "➡️  Nuevas features detectadas"
    else
        CHANGE_TYPE="PATCH"
        NEW_MINOR=$LAST_MINOR
        NEW_PATCH=$((LAST_PATCH + 1))
        print_info "➡️  Solo bug fixes"
    fi

    # Build suggested version
    if [ "$CHANGE_TYPE" = "MAJOR" ]; then
        NEW_MAJOR=$((LAST_MAJOR + 1))
        SUGGESTED_VERSION="v$NEW_MAJOR.0.0"
    elif [ "$CHANGE_TYPE" = "MINOR" ]; then
        SUGGESTED_VERSION="v$LAST_MAJOR.$NEW_MINOR.0"
    else
        SUGGESTED_VERSION="v$LAST_MAJOR.$LAST_MINOR.$NEW_PATCH"
    fi

    print_info "➡️  Versión sugerida: $SUGGESTED_VERSION ($CHANGE_TYPE)"
}

# Validate arguments and determine mode
MODE="manual"
DRY_RUN=false

if [ -z "$VERSION" ]; then
    # Interactive mode - analyze and suggest
    MODE="interactive"
elif [ "$VERSION" = "auto" ]; then
    # Auto mode - analyze and create without confirmation
    MODE="auto"
elif [ "$VERSION" = "--dry-run" ] || [ "$VERSION" = "-d" ]; then
    # Dry-run mode - test without changes
    MODE="dry-run"
    DRY_RUN=true
fi

echo ""
if [ "$MODE" = "dry-run" ]; then
    print_info "╔═══════════════════════════════════════════════════════════╗"
    print_info "║     PicoClaw-Agents Release Script (DRY-RUN)              ║"
    print_info "║     MODO PRUEBA - SIN CAMBIOS REALES                      ║"
    print_info "╚═══════════════════════════════════════════════════════════╝"
elif [ "$MODE" = "interactive" ]; then
    print_info "╔═══════════════════════════════════════════════════════════╗"
    print_info "║     PicoClaw-Agents Release Script (Interactivo)          ║"
    print_info "╚═══════════════════════════════════════════════════════════╝"
elif [ "$MODE" = "auto" ]; then
    print_info "╔═══════════════════════════════════════════════════════════╗"
    print_info "║     PicoClaw-Agents Release Script (Automático)           ║"
    print_info "╚═══════════════════════════════════════════════════════════╝"
else
    print_info "╔═══════════════════════════════════════════════════════════╗"
    print_info "║     PicoClaw-Agents Release Script (Manual)               ║"
    print_info "╚═══════════════════════════════════════════════════════════╝"

    # Validate version format (SemVer)
    if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
        print_error "Error: Formato de versión inválido"
        echo "El formato debe ser: vMAJOR.MINOR.PATCH o vMAJOR.MINOR.PATCH-LABEL"
        echo "Ejemplos: v1.0.1, v1.1.0, v2.0.0-beta, v1.0.0-rc1"
        exit 1
    fi
fi
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

# Step 2.5: Analyze commits (for interactive/auto/dry-run mode)
if [ "$MODE" = "interactive" ] || [ "$MODE" = "auto" ] || [ "$MODE" = "dry-run" ]; then
    print_info "Paso 2.5/6: Analizando commits..."

    # Parse last version to numbers
    LAST_VERSION_CLEAN=$(echo "$LAST_VERSION" | sed 's/^v//' | sed 's/-.*//')
    LAST_MAJOR=$(echo "$LAST_VERSION_CLEAN" | cut -d. -f1)
    LAST_MINOR=$(echo "$LAST_VERSION_CLEAN" | cut -d. -f2)
    LAST_PATCH=$(echo "$LAST_VERSION_CLEAN" | cut -d. -f3)

    # Analyze commits
    analyze_commits

    echo ""

    if [ "$MODE" = "auto" ]; then
        VERSION="$SUGGESTED_VERSION"
        print_success "Versión automática: $VERSION"
    elif [ "$MODE" = "dry-run" ]; then
        print_info "➡️  Versión sugerida: $SUGGESTED_VERSION ($CHANGE_TYPE)"
        print_info ""
        print_info "═══════════════════════════════════════════════════════"
        print_info "MODO DRY-RUN - Lo que HARÍA:"
        print_info "═══════════════════════════════════════════════════════"
        print_info "1. ✅ Ejecutar: make check"
        print_info "2. ✅ Ejecutar: make test"
        print_info "3. ✅ Ejecutar: make security-check"
        print_info "4. 🏷️  Crear tag: git tag $SUGGESTED_VERSION"
        print_info "5. 📤 Subir tag: git push origin $SUGGESTED_VERSION"
        print_info "6. 🚀 Disparar: gh workflow run release.yml"
        print_info "═══════════════════════════════════════════════════════"
        print_info ""
        print_info "Para ejecutar REALMENTE, usa:"
        print_info "  ./scripts/create-release.sh $SUGGESTED_VERSION"
        print_info ""
        print_info "O para interactivo:"
        print_info "  ./scripts/create-release.sh"
        exit 0
    else
        # Interactive mode - show suggestion and ask
        print_info "Versión sugerida: $SUGGESTED_VERSION ($CHANGE_TYPE)"
        echo ""
        echo "Opciones:"
        echo "  1. Usar versión sugerida ($SUGGESTED_VERSION)"
        echo "  2. Especificar versión manual"
        echo "  3. Cancelar"
        echo ""
        read -p "Elige una opción (1/2/3): " -n 1 -r
        echo

        if [[ $REPLY =~ ^[3]$ ]]; then
            print_info "Operación cancelada"
            exit 0
        elif [[ $REPLY =~ ^[2]$ ]]; then
            read -p "Ingresa la versión (ej: v1.0.1): " VERSION
            if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
                print_error "Formato de versión inválido"
                exit 1
            fi
        else
            VERSION="$SUGGESTED_VERSION"
        fi

        # Ask for pre-release
        echo ""
        read -p "¿Es un pre-release? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            PRERELEASE="true"
            read -p "Tipo de pre-release (beta/alpha/rc): " LABEL
            VERSION="$VERSION-$LABEL"
        fi
    fi

    print_info "Versión final: $VERSION"
    echo ""
fi

# Step 3: Verify last version (skip if already analyzed)
if [ "$MODE" != "interactive" ] && [ "$MODE" != "auto" ]; then
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
fi

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
