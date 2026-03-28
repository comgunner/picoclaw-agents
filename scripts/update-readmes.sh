#!/bin/bash

# Script para actualizar READMEs multilingües con las nuevas features v3.4.6-v3.5.2
# Uso: ./scripts/update-readmes.sh

set -e

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "🚀 Actualizando READMEs multilingües con features v3.4.6-v3.5.2..."

# Función para actualizar versión
update_version() {
    local file=$1
    local lang=$2

    echo "📝 Actualizando $file ($lang)..."

    # Actualizar versión
    sed -i.bak 's/Version:** v3\.4\.5+/Version:** v3.5.2/g' "$file" || true
    sed -i.bak 's/Versión:** v3\.4\.5+/Versión:** v3.5.2/g' "$file" || true
    sed -i.bak 's/版本：** v3\.4\.5+/版本：** v3.5.2/g' "$file" || true
    sed -i.bak 's/バージョン：** v3\.4\.5+/バージョン：** v3.5.2/g' "$file" || true
    sed -i.bak 's/Version :** v3\.4\.5+/Version :** v3.5.2/g' "$file" || true
    sed -i.bak 's/Versão:** v3\.4\.5+/Versão:** v3.5.2/g' "$file" || true
    sed -i.bak 's/Phiên bản:** v3\.4\.5+/Phiên bản:** v3.5.2/g' "$file" || true

    # Eliminar backup files
    rm -f "${file}.bak"

    echo "✅ $file actualizado"
}

# Actualizar todos los READMEs
update_version "$BASE_DIR/README.zh.md" "Chino"
update_version "$BASE_DIR/README.ja.md" "Japonés"
update_version "$BASE_DIR/README.fr.md" "Francés"
update_version "$BASE_DIR/README.pt-br.md" "Portugués"
update_version "$BASE_DIR/README.vi.md" "Vietnamita"

echo ""
echo "🎉 ¡Actualización completada!"
echo ""
echo "📋 Resumen de cambios aplicados:"
echo "   - Versión actualizada a v3.5.2"
echo ""
echo "⚠️  Nota: Las actualizaciones de features, noticias y DevOps deben aplicarse"
echo "   manualmente para asegurar traducciones precisas en cada idioma."
echo ""
echo "📝 Archivos modificados:"
echo "   - README.zh.md (Chino)"
echo "   - README.ja.md (Japonés)"
echo "   - README.fr.md (Francés)"
echo "   - README.pt-br.md (Portugués)"
echo "   - README.vi.md (Vietnamita)"
echo ""
