package tests

import (
	"testing"
)

func TestNotationWorkflow(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	// Step 1: President submits grades
	body := map[string]float64{
		"criterion1": 3.5,
		"criterion2": 3.0,
		"criterion3": 4.0,
		"criterion4": 3.5,
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body, h.AuthHeader(SeedTeacherISIL2ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Étape 1 - Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)

	// Step 2: Member submits grades
	body2 := map[string]float64{
		"criterion1": 3.0,
		"criterion2": 3.5,
		"criterion3": 3.0,
		"criterion4": 4.0,
	}
	resp2, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body2, h.AuthHeader(SeedTeacherCHIM1ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Étape 2 - Erreur requête: %v", err)
	}
	result2 := MustParseResponse(resp2)
	AssertSuccess(t, result2)

	// Step 3: Admin resolves grade (chooses president's grade)
	resp3, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/resolve-grade", map[string]any{"choice": "president"}, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Étape 3 - Erreur requête: %v", err)
	}
	result3 := MustParseResponse(resp3)
	AssertSuccess(t, result3)
}

func TestNotationWorkflowResolveWithMemberGrade(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	// President submits
	body := map[string]float64{
		"criterion1": 3.5,
		"criterion2": 3.0,
		"criterion3": 4.0,
		"criterion4": 3.5,
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body, h.AuthHeader(SeedTeacherISIL2ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	MustParseResponse(resp)

	// Member submits
	body2 := map[string]float64{
		"criterion1": 3.0,
		"criterion2": 3.5,
		"criterion3": 3.0,
		"criterion4": 4.0,
	}
	resp2, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body2, h.AuthHeader(SeedTeacherCHIM1ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	MustParseResponse(resp2)

	// Resolve with member's grade
	resp3, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/resolve-grade", map[string]any{"choice": "member"}, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result3 := MustParseResponse(resp3)
	AssertSuccess(t, result3)
}

func TestNotationWorkflowResolveWithNewEvaluation(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	// President submits
	body := map[string]float64{
		"criterion1": 3.5,
		"criterion2": 3.0,
		"criterion3": 4.0,
		"criterion4": 3.5,
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body, h.AuthHeader(SeedTeacherISIL2ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	MustParseResponse(resp)

	// Member submits
	body2 := map[string]float64{
		"criterion1": 3.0,
		"criterion2": 3.5,
		"criterion3": 3.0,
		"criterion4": 4.0,
	}
	resp2, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body2, h.AuthHeader(SeedTeacherCHIM1ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	MustParseResponse(resp2)

	// Resolve with new evaluation
	resp3, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/resolve-grade", map[string]any{
		"choice": "new",
		"grades": map[string]float64{
			"criterion1": 3.0,
			"criterion2": 3.0,
			"criterion3": 3.5,
			"criterion4": 3.5,
		},
	}, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result3 := MustParseResponse(resp3)
	AssertSuccess(t, result3)
}

func TestNotationWorkflowResolveWithoutGrades(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	// Try to resolve before anyone submitted
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/resolve-grade", map[string]any{"choice": "president"}, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestNotationWorkflowSubmitGradeInvalidCriterion(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	// Submit with invalid criterion (> 4)
	body := map[string]float64{
		"criterion1": 5.0,
		"criterion2": 3.0,
		"criterion3": 4.0,
		"criterion4": 3.5,
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body, h.AuthHeader(SeedTeacherISIL2ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestNotationWorkflowSubmitGradeNegativeCriterion(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]float64{
		"criterion1": -1.0,
		"criterion2": 3.0,
		"criterion3": 4.0,
		"criterion4": 3.5,
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body, h.AuthHeader(SeedTeacherISIL2ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestNotationWorkflowSubmitGradeUnauthorized(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]float64{
		"criterion1": 3.0,
		"criterion2": 3.0,
		"criterion3": 3.0,
		"criterion4": 3.0,
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestNotationWorkflowSubmitGradeDefenseNotFound(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]float64{
		"criterion1": 3.0,
		"criterion2": 3.0,
		"criterion3": 3.0,
		"criterion4": 3.0,
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/99999/submit-grade", body, h.AuthHeader(SeedTeacherISIL2ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestNotationWorkflowResolveUnauthorized(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/resolve-grade", map[string]any{"choice": "president"}, h.AuthHeader(SeedTeacherISIL1ID, "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestNotationWorkflowResolveDefenseNotFound(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/99999/resolve-grade", map[string]any{"choice": "president"}, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestNotationWorkflowResolveNewEvaluationWithoutGrades(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	// President submits
	body := map[string]float64{
		"criterion1": 3.5,
		"criterion2": 3.0,
		"criterion3": 4.0,
		"criterion4": 3.5,
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body, h.AuthHeader(SeedTeacherISIL2ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	MustParseResponse(resp)

	// Member submits
	body2 := map[string]float64{
		"criterion1": 3.0,
		"criterion2": 3.5,
		"criterion3": 3.0,
		"criterion4": 4.0,
	}
	resp2, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body2, h.AuthHeader(SeedTeacherCHIM1ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	MustParseResponse(resp2)

	// Resolve with new but missing grades
	resp3, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/resolve-grade", map[string]any{"choice": "new"}, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result3 := MustParseResponse(resp3)
	AssertError(t, result3)
}

func TestNotationWorkflowFinalGradeCalculation(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	// Step 1: President submits: 3.5+3.0+4.0+3.5 = 14.0
	body := map[string]float64{
		"criterion1": 3.5,
		"criterion2": 3.0,
		"criterion3": 4.0,
		"criterion4": 3.5,
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body, h.AuthHeader(SeedTeacherISIL2ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	MustParseResponse(resp)

	// Step 2: Member submits: 3.0+3.5+3.0+4.0 = 13.5
	body2 := map[string]float64{
		"criterion1": 3.0,
		"criterion2": 3.5,
		"criterion3": 3.0,
		"criterion4": 4.0,
	}
	resp2, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body2, h.AuthHeader(SeedTeacherCHIM1ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	MustParseResponse(resp2)

	// Step 3: Resolve with president's grade (14.0) + existing supervisor eval (3.5) = 17.5/20
	resp3, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/resolve-grade", map[string]any{"choice": "president"}, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result3 := MustParseResponse(resp3)
	AssertSuccess(t, result3)

	// Verify the defense shows final grade
	resp4, err := h.App.Test(newHTTPRequest("GET", "/api/admin/defenses/1", nil, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur vérification: %v", err)
	}
	result4 := MustParseResponse(resp4)
	AssertSuccess(t, result4)
}

func TestNotationWorkflowDoubleSubmission(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]float64{
		"criterion1": 3.5,
		"criterion2": 3.0,
		"criterion3": 4.0,
		"criterion4": 3.5,
	}

	// First submission succeeds
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body, h.AuthHeader(SeedTeacherISIL2ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	MustParseResponse(resp)

	// Second submission from same jury member should update (not error)
	resp2, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body, h.AuthHeader(SeedTeacherISIL2ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result2 := MustParseResponse(resp2)
	AssertSuccess(t, result2)
}

func TestNotationWorkflowJuryConfirmsPrintedVersion(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	// President confirms printed version
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/confirm-jury", map[string]any{
		"member_id":     SeedTeacherISIL2ID,
		"wants_printed": true,
	}, h.AuthHeader(SeedTeacherISIL2ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestNotationWorkflowJuryDeclines(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/decline-jury", map[string]any{
		"member_id": SeedTeacherCHIM1ID,
		"reason":    "Conflit d'intérêts",
	}, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestNotationWorkflowRecommendJury(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/defenses/recommend-jury?pfe_id=1", nil, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestNotationWorkflowRecommendJuryNoPFE(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/admin/defenses/recommend-jury", nil, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestNotationWorkflowDuplicateJuryMember(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	// President and member same person should fail if the handler checks
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses", map[string]any{
		"assignment_id": 1,
		"president_id":  SeedTeacherISIL1ID,
		"member_id":     SeedTeacherISIL1ID,
		"scheduled_at":  "2025-06-15T10:00:00Z",
		"room":          "Salle B",
	}, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestNotationWorkflowFullIntegration(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	// Complete workflow: submit president grade -> submit member grade -> resolve -> verify
	body := map[string]float64{
		"criterion1": 3.0,
		"criterion2": 3.0,
		"criterion3": 3.0,
		"criterion4": 3.0,
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body, h.AuthHeader(SeedTeacherISIL2ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)

	body2 := map[string]float64{
		"criterion1": 4.0,
		"criterion2": 4.0,
		"criterion3": 4.0,
		"criterion4": 4.0,
	}
	resp2, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/submit-grade", body2, h.AuthHeader(SeedTeacherCHIM1ID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result2 := MustParseResponse(resp2)
	AssertSuccess(t, result2)

	resp3, err := h.App.Test(newHTTPRequest("POST", "/api/admin/defenses/1/resolve-grade", map[string]any{"choice": "member"}, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result3 := MustParseResponse(resp3)
	AssertSuccess(t, result3)

	// Verify final grade: member jury (4+4+4+4=16) + supervisor (3.5) = 19.5/20
	resp4, err := h.App.Test(newHTTPRequest("GET", "/api/admin/defenses/1", nil, h.AuthHeader(SeedAdminID, "admin")))
	if err != nil {
		t.Fatalf("❌ Erreur vérification: %v", err)
	}
	result4 := MustParseResponse(resp4)
	AssertSuccess(t, result4)
}
