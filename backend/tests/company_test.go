package tests

import (
	"testing"
)

func TestCompanyDashboard(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/company/dashboard", nil, h.AuthHeader(SeedCompany1ID, "company")))
	
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestCompanyDashboardUnauthorized(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/company/dashboard", nil, nil))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestCompanyDashboardWrongRole(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/company/dashboard", nil, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestCompanyListSubjects(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/company/subjects", nil, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestCompanyCreateSubject(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]any{
		"title":       "Sujet Entreprise Test",
		"description": "Description du sujet proposé par l'entreprise",
		"group_type":  "binome",
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/company/subjects", body, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestCompanyCreateSubjectValidation(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/company/subjects", map[string]string{}, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestCompanyGetSubject(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/company/subjects/5", nil, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestCompanyGetSubjectNotFound(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/company/subjects/99999", nil, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestCompanyUpdateSubject(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"title": "Sujet Entreprise Modifié"}
	resp, err := h.App.Test(newHTTPRequest("PATCH", "/api/company/subjects/5", body, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestCompanyUpdateSubjectNotOwned(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"title": "Tentative modification"}
	resp, err := h.App.Test(newHTTPRequest("PATCH", "/api/company/subjects/1", body, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "non autorisé")
}

func TestCompanyListCandidats(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/company/subjects/5/candidats", nil, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestCompanyAcceptCandidat(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]any{"student_id": SeedStudentISIL2ID, "action": "accept"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/company/subjects/5/candidats", body, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestCompanyListSupervisedPfes(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/company/supervised-pfes", nil, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestCompanyGetSupervisedPfe(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/company/supervised-pfes/2", nil, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestCompanyGetSupervisedPfeNotFound(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/company/supervised-pfes/99999", nil, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestCompanyCreateMeeting(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]any{
		"meeting_date": "2025-06-01T10:00:00Z",
		"duration":     45,
		"meeting_type": "visio",
		"topics":       "Suivi du projet en entreprise",
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/company/supervised-pfes/2/meetings", body, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestCompanyCreateMeetingValidation(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/company/supervised-pfes/2/meetings", map[string]string{}, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestCompanySubmitEvaluation(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]float64{"criterion5": 3.0}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/company/supervised-pfes/2/evaluation", body, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestCompanySubmitEvaluationInvalidCriterion(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]float64{"criterion5": -1.0}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/company/supervised-pfes/2/evaluation", body, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestCompanyListReports(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/company/reports", nil, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestCompanyCreateReport(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{
		"correction_type": "company_name",
		"description":     "Le nom de l'entreprise est incorrect",
		"requested_value": "TechCorp Algeria Updated",
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/company/reports", body, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestCompanyCreateReportValidation(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/company/reports", map[string]string{}, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestCompanyNotifications(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/company/notifications", nil, h.AuthHeader(SeedCompany1ID, "company")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}
