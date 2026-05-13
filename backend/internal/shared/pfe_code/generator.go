package pfe_code

import (
	"fmt"
	"strings"
)

// Generate génère un code PFE au format PFE-[CODE_SPECIALITE]-[ANNEE]-[NNN].
// Exemple : PFE-ISIL-2025-001
func Generate(specialityCode string, academicYearLabel string, sequence int) string {
	// Extrait l'année de début du label (ex: "2024-2025" -> "2025").
	year := extractEndYear(academicYearLabel)
	return fmt.Sprintf("PFE-%s-%s-%03d", strings.ToUpper(specialityCode), year, sequence)
}

// extractEndYear extrait l'année de fin d'un label comme "2024-2025".
// Retourne l'année complète si le format n'est pas reconnu.
func extractEndYear(label string) string {
	parts := strings.Split(label, "-")
	if len(parts) >= 2 {
		return strings.TrimSpace(parts[1])
	}
	return label
}
