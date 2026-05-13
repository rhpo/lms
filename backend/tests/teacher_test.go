package tests

import (
	"testing"
)

func TestTeacherDashboard(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/dashboard", nil, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherDashboardUnauthorized(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/dashboard", nil, nil))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestTeacherDashboardWrongRole(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/dashboard", nil, h.AuthHeader("seed-student-isil-001", "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestTeacherListProposedSubjects(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/proposed-subjects", nil, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherCreateProposedSubject(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]any{
		"title":       "Nouveau Sujet Test",
		"description": "Description du nouveau sujet",
		"group_type":  "monome",
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/teacher/proposed-subjects", body, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherCreateProposedSubjectValidation(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/teacher/proposed-subjects", map[string]string{}, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestTeacherGetProposedSubject(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/proposed-subjects/seed-subject-001", nil, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherGetProposedSubjectNotFound(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/proposed-subjects/invalid", nil, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestTeacherUpdateProposedSubject(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"title": "Sujet Modifié"}
	resp, err := h.App.Test(newHTTPRequest("PATCH", "/api/teacher/proposed-subjects/seed-subject-001", body, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherListCandidats(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/proposed-subjects/seed-subject-001/candidats", nil, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherAcceptCandidat(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"student_id": "seed-student-isil-001", "action": "accept"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/teacher/proposed-subjects/seed-subject-005/candidats", body, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherListSubjectsToValidate(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/subjects-to-validate", nil, h.AuthHeader("seed-teacher-isil-002", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherGetSubjectToValidate(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/subjects-to-validate/seed-subject-003", nil, h.AuthHeader("seed-teacher-isil-002", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherValidateSubject(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"decision": "valide", "comment": "Sujet intéressant"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/teacher/subjects-to-validate/seed-subject-003", body, h.AuthHeader("seed-teacher-isil-002", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherValidateSubjectInvalidDecision(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"decision": "invalid"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/teacher/subjects-to-validate/seed-subject-003", body, h.AuthHeader("seed-teacher-isil-002", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestTeacherListSupervisedPfes(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/supervised-pfes", nil, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherGetSupervisedPfe(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/supervised-pfes/seed-assignment-001", nil, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
	data, ok := result["data"].(map[string]any)
	if !ok {
		t.Fatal("❌ Teacher get supervised pfe: data manquant")
	}
	if data["pfe_code"] != "PFE-ISIL-2025-001" {
		t.Fatalf("❌ Teacher get supervised pfe: code PFE incorrect %v", data["pfe_code"])
	}
}

func TestTeacherGetSupervisedPfeNotFound(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/supervised-pfes/invalid", nil, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestTeacherCreateMeeting(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]any{
		"meeting_date": "2025-05-20T14:00:00Z",
		"duration":     60,
		"meeting_type": "visio",
		"topics":       "Avancement du projet",
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/teacher/supervised-pfes/seed-assignment-001/meetings", body, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherCreateMeetingValidation(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/teacher/supervised-pfes/seed-assignment-001/meetings", map[string]string{}, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestTeacherSubmitEvaluation(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]float64{"criterion5": 3.5}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/teacher/supervised-pfes/seed-assignment-001/evaluation", body, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherSubmitEvaluationInvalidCriterion(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]float64{"criterion5": 5.0}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/teacher/supervised-pfes/seed-assignment-001/evaluation", body, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestTeacherListJuryDuties(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/jury-duties", nil, h.AuthHeader("seed-teacher-isil-002", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherGetJuryDuty(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/jury-duties/seed-defense-001", nil, h.AuthHeader("seed-teacher-isil-002", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherGetJuryDutyNotFound(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/jury-duties/invalid", nil, h.AuthHeader("seed-teacher-isil-002", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestTeacherUpdateAvailability(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"availability_status": "indisponible_jusqu_au", "unavailable_until": "2025-06-30"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/teacher/availability", body, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestTeacherUpdateAvailabilityInvalid(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"availability_status": "invalid"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/teacher/availability", body, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestTeacherNotifications(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/notifications", nil, h.AuthHeader("seed-teacher-isil-001", "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestAdminCanAccessTeacherEndpoints(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/teacher/dashboard", nil, h.AuthHeader("seed-admin-001", "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}
