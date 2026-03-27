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
#   - Runs pre-release checks (make check, make test, make lint)
#   - Dispatches GitHub Actions release workflow (workflow crea el tag)
#   - Supports pre-releases (beta, alpha, rc)
#   - Dry-run mode for testing
#
# IMPORTANT: El binario compilado se llama picoclaw-agents (no picoclaw).
# No modificar BINARY_NAME ni binary: en .goreleaser.yaml sin revisar
# Dockerfile.goreleaser primero.

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

# Get last version tag (local first, then remote fallback)
# Uses git describe to get semantically latest tag, not alphabetically last.
get_last_version() {
    local ver
    ver=$(git describe --tags --abbrev=0 2>/dev/null || true)
    if [ -z "$ver" ]; then
        # Fallback: fetch from remote, strip ^{} dereference lines
        ver=$(git ls-remote --tags origin 2>/dev/null \
            | grep -v '\^{}' \
            | awk '{print $2}' \
            | sed 's|refs/tags/||' \
            | sort -V \
            | tail -1)
    fi
    echo "${ver:-v0.0.0}"
}

# Analyze commits since LAST_VERSION and set SUGGESTED_VERSION + CHANGE_TYPE
analyze_commits() {
    print_info "Analizando commits desde $LAST_VERSION..."

    COMMITS=$(git log "$LAST_VERSION"..HEAD --oneline 2>/dev/null || true)

    if [ -z "$COMMITS" ]; then
        print_warning "No hay commits nuevos desde $LAST_VERSION"
        SUGGESTED_VERSION="$LAST_VERSION"
        CHANGE_TYPE="NONE"
        return
    fi

    COMMIT_COUNT=$(echo "$COMMITS" | wc -l | tr -d ' ')
    print_info "Encontrados $COMMIT_COUNT commits"

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

    # Parse last version numbers
    local ver_clean
    ver_clean=$(echo "$LAST_VERSION" | sed 's/^v//' | sed 's/-.*//')
    LAST_MAJOR=$(echo "$ver_clean" | cut -d. -f1)
    LAST_MINOR=$(echo "$ver_clean" | cut -d. -f2)
    LAST_PATCH=$(echo "$ver_clean" | cut -d. -f3)

    # Determine bump type
    if [ "$BREAKING_COUNT" -gt 0 ]; then
        CHANGE_TYPE="MAJOR"
        print_info "➡️  BREAKING CHANGE detectado → MAJOR"
        SUGGESTED_VERSION="v$((LAST_MAJOR + 1)).0.0"
    elif [ "$NATIVE_SKILLS" -gt 0 ] || [ "$NEW_PROVIDER" -gt 0 ] || [ "$FEAT_COUNT" -gt 0 ]; then
        CHANGE_TYPE="MINOR"
        print_info "➡️  Nuevas features detectadas → MINOR"
        SUGGESTED_VERSION="v$LAST_MAJOR.$((LAST_MINOR + 1)).0"
    else
        CHANGE_TYPE="PATCH"
        print_info "➡️  Solo bug fixes → PATCH"
        SUGGESTED_VERSION="v$LAST_MAJOR.$LAST_MINOR.$((LAST_PATCH + 1))"
    fi

    print_info "➡️  Versión sugerida: $SUGGESTED_VERSION ($CHANGE_TYPE)"
}

# -------------------------------------------------------------------
# Validate arguments and determine mode
# -------------------------------------------------------------------
MODE="manual"
DRY_RUN=false

if [ -z "$VERSION" ]; then
    MODE="interactive"
elif [ "$VERSION" = "auto" ]; then
    MODE="auto"
elif [ "$VERSION" = "--dry-run" ] || [ "$VERSION" = "-d" ]; then
    MODE="dry-run"
    DRY_RUN=true
fi

echo ""
case "$MODE" in
    dry-run)
        print_info "╔═══════════════════════════════════════════════════════════╗"
        print_info "║     PicoClaw-Agents Release Script (DRY-RUN)              ║"
        print_info "║     MODO PRUEBA - SIN CAMBIOS REALES                      ║"
        print_info "╚═══════════════════════════════════════════════════════════╝"
        ;;
    interactive)
        print_info "╔═══════════════════════════════════════════════════════════╗"
        print_info "║     PicoClaw-Agents Release Script (Interactivo)          ║"
        print_info "╚═══════════════════════════════════════════════════════════╝"
        ;;
    auto)
        print_info "╔═══════════════════════════════════════════════════════════╗"
        print_info "║     PicoClaw-Agents Release Script (Automático)           ║"
        print_info "╚═══════════════════════════════════════════════════════════╝"
        ;;
    *)
        print_info "╔═══════════════════════════════════════════════════════════╗"
        print_info "║     PicoClaw-Agents Release Script (Manual)               ║"
        print_info "╚═══════════════════════════════════════════════════════════╝"
        # Validate SemVer format
        if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
            print_error "Formato de versión inválido: $VERSION"
            echo "El formato debe ser: vMAJOR.MINOR.PATCH o vMAJOR.MINOR.PATCH-LABEL"
            echo "Ejemplos: v1.0.1, v1.1.0, v2.0.0-beta, v1.0.0-rc1"
            exit 1
        fi
        ;;
esac
echo ""

# -------------------------------------------------------------------
# Step 1: Check prerequisites
# -------------------------------------------------------------------
print_info "Paso 1/5: Verificando prerequisitos..."

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

# -------------------------------------------------------------------
# Step 2: Check current branch and working tree
# -------------------------------------------------------------------
print_info "Paso 2/5: Verificando branch..."

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

UNCOMMITTED=$(git status --porcelain)
if [ -n "$UNCOMMITTED" ]; then
    print_warning "Hay cambios sin commitear:"
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

# -------------------------------------------------------------------
# Step 2.5: Get last version + analyze commits
# (required for interactive, auto, and dry-run modes)
# -------------------------------------------------------------------
if [ "$MODE" = "interactive" ] || [ "$MODE" = "auto" ] || [ "$MODE" = "dry-run" ]; then
    print_info "Paso 2.5/5: Obteniendo última versión y analizando commits..."

    # FIX: LAST_VERSION se obtiene AQUÍ, antes de usarlo en analyze_commits
    LAST_VERSION=$(get_last_version)
    print_info "Última versión detectada: $LAST_VERSION"

    analyze_commits
    echo ""

    if [ "$MODE" = "auto" ]; then
        VERSION="$SUGGESTED_VERSION"
        print_success "Versión automática: $VERSION"

    elif [ "$MODE" = "dry-run" ]; then
        print_info "➡️  Versión sugerida: $SUGGESTED_VERSION ($CHANGE_TYPE)"
        print_info ""
        print_info "═══════════════════════════════════════════════════════"
        print_info "MODO DRY-RUN — Lo que HARÍA:"
        print_info "═══════════════════════════════════════════════════════"
        print_info "1. ✅ Ejecutar: make check"
        print_info "2. ✅ Ejecutar: make test"
        print_info "3. ✅ Ejecutar: make lint"
        print_info "4. 🚀 Disparar: gh workflow run release.yml \\"
        print_info "       --field tag=$SUGGESTED_VERSION \\"
        print_info "       --field prerelease=false \\"
        print_info "       --field draft=false"
        print_info "   (El workflow crea el tag y ejecuta GoReleaser)"
        print_info "═══════════════════════════════════════════════════════"
        print_info ""
        print_info "Para ejecutar REALMENTE, usa:"
        print_info "  ./scripts/create-release.sh $SUGGESTED_VERSION"
        print_info ""
        print_info "O para interactivo:"
        print_info "  ./scripts/create-release.sh"
        exit 0

    else
        # Interactive: show suggestion and ask
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

# -------------------------------------------------------------------
# Step 3 (solo modo manual): verificar versión y confirmar
# -------------------------------------------------------------------
if [ "$MODE" = "manual" ]; then
    print_info "Paso 3/5: Verificando versión..."

    LAST_VERSION=$(get_last_version)
    print_info "Última versión: $LAST_VERSION"
    print_info "Nueva versión:  $VERSION"

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

# -------------------------------------------------------------------
# Step 4: Run verifications
# -------------------------------------------------------------------
print_info "Paso 4/5: Ejecutando verificaciones de código..."

print_info "Ejecutando make check (vet + fmt + test)..."
if ! make check > /dev/null 2>&1; then
    print_error "make check falló"
    print_info "Ejecuta 'make check' manualmente para ver los errores"
    exit 1
fi
print_success "make check passed"

print_info "Ejecutando make lint..."
if ! make lint > /dev/null 2>&1; then
    print_warning "make lint encontró problemas"
    print_info "Revisa manualmente con 'make lint'"
    # No exit — lint warnings no bloquean el release
fi
print_success "Lint check completed"
echo ""

# -------------------------------------------------------------------
# Step 5: Dispatch GitHub Actions release workflow
# El workflow (release.yml) crea el tag y ejecuta GoReleaser.
# El script NO crea el tag localmente para evitar doble creación.
# -------------------------------------------------------------------
print_info "Paso 5/5: Disparando release workflow..."
print_info "  Tag:        $VERSION"
print_info "  Pre-release: $PRERELEASE"
print_info "  Binary:     picoclaw-agents"
echo ""

if gh workflow run release.yml \
    --field tag="$VERSION" \
    --field prerelease="$PRERELEASE" \
    --field draft=false; then
    print_success "Release workflow disparado exitosamente"
    print_info "El workflow creará el tag y ejecutará GoReleaser automáticamente."
else
    print_error "No se pudo disparar el workflow"
    print_info "Puedes dispararlo manualmente en:"
    print_info "https://github.com/comgunner/picoclaw-agents/actions/workflows/release.yml"
    exit 1
fi
echo ""

# -------------------------------------------------------------------
# Summary
# -------------------------------------------------------------------
echo ""
print_info "╔═══════════════════════════════════════════════════════════╗"
print_info "║                    Resumen                               ║"
print_info "╚═══════════════════════════════════════════════════════════╝"
echo ""
print_success "Release $VERSION iniciado exitosamente"
echo ""
print_info "Monitorea el progreso en:"
print_info "  - Actions: https://github.com/comgunner/picoclaw-agents/actions"
print_info "  - Release: https://github.com/comgunner/picoclaw-agents/releases/tag/$VERSION"
print_info "  - Compare: https://github.com/comgunner/picoclaw-agents/compare/$LAST_VERSION...$VERSION"
echo ""

if [ "$PRERELEASE" = "true" ]; then
    print_warning "Este es un pre-release"
else
    print_success "Release estable: $VERSION"
fi

echo ""
print_info "Próximos pasos:"
print_info "  1. Monitorea el workflow en GitHub Actions"
print_info "  2. Verifica que los binarios picoclaw-agents se generaron correctamente"
print_info "  3. Revisa que el CHANGELOG.md refleja todos los cambios"
echo ""
