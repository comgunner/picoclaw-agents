#!/bin/bash
# Script para actualizar CHANGELOG.md automáticamente
# Uso: ./scripts/update-changelog.sh [major|minor|patch] "[Descripción]"

set -e

VERSION=${1:-"patch"}  # major, minor, patch
DESCRIPTION=${2:-"Release"}

# Obtener última versión
LAST_VERSION=$(grep -E '^## \[[0-9]+\.[0-9]+\.[0-9]+\]' CHANGELOG.md | head -1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' || echo "3.4.5")

# Parsear versión
MAJOR=$(echo $LAST_VERSION | cut -d. -f1)
MINOR=$(echo $LAST_VERSION | cut -d. -f2)
PATCH=$(echo $LAST_VERSION | cut -d. -f3)

# Incrementar versión
case $VERSION in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
esac

NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
TODAY=$(date +%Y-%m-%d)

# Insertar nueva entrada en CHANGELOG
TEMP_FILE=$(mktemp)
echo "## [$NEW_VERSION] - $TODAY" > $TEMP_FILE
echo "" >> $TEMP_FILE
echo "### 📝 Descripción" >> $TEMP_FILE
echo "- $DESCRIPTION" >> $TEMP_FILE
echo "" >> $TEMP_FILE
cat CHANGELOG.md >> $TEMP_FILE
mv $TEMP_FILE CHANGELOG.md

echo "✅ CHANGELOG.md actualizado a v$NEW_VERSION"
echo ""
echo "Próximos pasos:"
echo "  1. Editar CHANGELOG.md para agregar detalles específicos"
echo "  2. Commit: git add CHANGELOG.md && git commit -m 'chore: update CHANGELOG for v$NEW_VERSION'"
