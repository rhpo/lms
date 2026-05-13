package tests

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAdminDashboard(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/dashboard", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
	data, ok := result["data"].(map[string]any)
	if !ok {
		t.Fatal("❌ Admin dashboard: data manquant")
	}
	if _, exists := data["total_users"]; !exists {
		t.Fatal("❌ Admin dashboard: total_users manquant")
	}
}

func TestAdminDashboardUnauthorized(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/dashboard", nil, nil))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestAdminListUsers(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/accounts/users", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
	data, ok := result["data"].([]any)
	if !ok {
		t.Fatal("❌ Admin list users: data n'est pas un tableau")
	}
	if len(data) == 0 {
		t.Fatal("❌ Admin list users: tableau vide inattendu")
	}
}

func TestAdminCreateUser(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"id": "test-new-user", "role": "teacher", "full_name": "Nouvel Enseignant", "email": "new.teacher@test.dz"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/accounts/users", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
	data, ok := result["data"].(map[string]any)
	if !ok {
		t.Fatal("❌ Admin create user: data manquant")
	}
	if data["full_name"] != "Nouvel Enseignant" {
		t.Fatalf("❌ Admin create user: nom incorrect %v", data["full_name"])
	}
}

func TestAdminCreateUserValidation(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"id": "", "role": "", "full_name": "", "email": ""}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/accounts/users", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "Tous les champs sont requis")
}

func TestAdminGetUser(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/accounts/users/seed-admin-001", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
	data, ok := result["data"].(map[string]any)
	if !ok {
		t.Fatal("❌ Admin get user: data manquant")
	}
	if data["full_name"] != "Admin Test" {
		t.Fatalf("❌ Admin get user: nom incorrect %v", data["full_name"])
	}
}

func TestAdminGetUserNotFound(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/accounts/users/invalid-id", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestAdminUpdateUser(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"full_name": "Admin Modifié", "email": "admin.modified@test.dz"}
	resp, err := h.App.Test(newHTTPRequest("PATCH", "/api/admin/accounts/users/seed-admin-001", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminUserActionDeactivate(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"action": "deactivate"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/accounts/users/seed-teacher-isil-001/action", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminUserActionInvalid(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"action": "invalid_action"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/accounts/users/seed-admin-001/action", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestAdminListCompanies(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/accounts/companies", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminCompanyAction(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"action": "validate"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/accounts/companies/seed-company-001/action", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminListReports(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/reports", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminListSubjects(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/subjects", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
	data, ok := result["data"].([]any)
	if !ok {
		t.Fatal("❌ Admin list subjects: data n'est pas un tableau")
	}
	if len(data) < 6 {
		t.Fatalf("❌ Admin list subjects: attendu >=6 sujets, obtenu %d", len(data))
	}
}

func TestAdminGetSubject(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/subjects/seed-subject-001", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminGetSubjectNotFound(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/subjects/invalid-subject", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestAdminSubjectActionAssignValidators(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"action": "assign-validators", "validator_id": "seed-teacher-isil-002"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/subjects/seed-subject-001/action", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminSubjectActionInvalid(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/subjects/seed-subject-001/action", map[string]string{}, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestAdminListAssignments(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/pfe", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminGetAssignment(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/pfe/seed-assignment-001", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
	data, ok := result["data"].(map[string]any)
	if !ok {
		t.Fatal("❌ Admin get assignment: data manquant")
	}
	if data["pfe_code"] != "PFE-ISIL-2025-001" {
		t.Fatalf("❌ Admin get assignment: code PFE incorrect %v", data["pfe_code"])
	}
}

func TestAdminGetAssignmentNotFound(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/pfe/invalid", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestAdminListDefenses(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/defenses", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminCreateDefense(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{
		"assignment_id": "seed-assignment-001",
		"president_id":  "seed-teacher-isil-002",
		"member_id":     "seed-teacher-chim-001",
		"scheduled_at":  "2025-07-01T10:00:00Z",
		"room":          "Salle B",
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminCreateDefenseValidation(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses", map[string]string{}, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestAdminGetDefense(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/defenses/seed-defense-001", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminGetDefenseNotFound(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/defenses/invalid", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestAdminRecommendJury(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/defenses/recommend-jury?pfe_id=seed-assignment-001", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminRecommendJuryMissingParam(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/defenses/recommend-jury", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "pfe_id")
}

func TestAdminSubmitGrade(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]float64{"criterion1": 3.0, "criterion2": 3.5, "criterion3": 2.5, "criterion4": 3.0}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/seed-defense-001/submit-grade", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminResolveGrade(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]float64{"criterion1": 3.5, "criterion2": 3.0, "criterion3": 3.5, "criterion4": 3.0}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/seed-defense-001/resolve-grade", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminConfirmJury(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/seed-jury-001/confirm-jury", map[string]string{}, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminDeclineJury(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/seed-jury-001/decline-jury", map[string]string{}, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminListDeadlines(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/settings/deadlines", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminUpdateDeadlines(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]any{"submission_open_at": "2025-01-01", "submission_close_at": "2025-06-30", "max_wishes": 3}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/settings/deadlines", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminListSpecialities(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/settings/specialities", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminCreateSpeciality(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"id": "test-spec", "name": "Test Spec", "code": "TEST", "year_type": "master"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/settings/specialities", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminDeleteSpeciality(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("DELETE", "/api/admin/settings/specialities/test-spec", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminListDomains(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/settings/domains", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminCreateDomain(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"id": "test-domain", "name": "Test Domain"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/settings/domains", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminDeleteDomain(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("DELETE", "/api/admin/settings/domains/test-domain", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminListPromotions(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/settings/promotions", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminCreatePromotion(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"id": "test-promo", "label": "Test Promotion", "academic_year_id": "seed-ay-2425"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/settings/promotions", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminDeletePromotion(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("DELETE", "/api/admin/settings/promotions/test-promo", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminListAcademicYears(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/settings/academic-years", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminCreateAcademicYear(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"id": "test-ay", "label": "2025-2026", "status": "active"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/settings/academic-years", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminCloseAcademicYear(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/settings/academic-years/seed-ay-2324/close", map[string]string{}, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminStatistics(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/statistics", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminAuditLog(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/audit-log", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminExports(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	exports := []string{"/api/admin/exports/affectations", "/api/admin/exports/plannings", "/api/admin/exports/statistiques"}
	for _, endpoint := range exports {
		t.Run(endpoint, func(t *testing.T) {
			resp, err := h.App.Test(newHTTPRequest("GET", endpoint, nil, h.AuthHeader("seed-admin-001", "admin")))
			if err != nil {
				t.Fatalf("❌ Erreur requête: %v", err)
			}
			result := MustParseResponse(resp)
			AssertSuccess(t, result)
		})
	}
}

func TestAdminRoleProtection(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	token := h.AuthHeader("seed-student-isil-001", "student")
	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/dashboard", nil, token))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestAdminImportCSV(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	csvData := "role,full_name,email\nteacher,CSV Teacher 1,csv.teacher1@test.dz\nstudent,CSV Student 1,csv.student1@test.dz\n"
	body := map[string]string{"csv_data": csvData}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/accounts/users/import-csv", body, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

// Ensure unused import is used
var _ = json.Marshal
var _ = strings.NewReader
