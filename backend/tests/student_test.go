package tests

import (
	"testing"
)

func TestStudentDashboard(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/student/dashboard", nil, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestStudentDashboardUnauthorized(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/student/dashboard", nil, nil))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestStudentDashboardWrongRole(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/student/dashboard", nil, h.AuthHeader(SeedTeacherISIL1ID, "teacher")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestStudentCatalogue(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/student/catalogue", nil, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
	data, ok := result["data"].([]any)
	if !ok {
		t.Fatal("❌ Student catalogue: data n'est pas un tableau")
	}
	if len(data) == 0 {
		t.Fatal("❌ Student catalogue: tableau vide inattendu")
	}
}

func TestStudentGetCatalogueSubject(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/student/catalogue/3", nil, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestStudentGetCatalogueSubjectNotFound(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/student/catalogue/99999", nil, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestStudentListWishes(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/student/wishes", nil, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestStudentCreateWish(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]any{"subject_id": 5}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/student/wishes", body, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestStudentCreateWishDuplicate(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]any{"subject_id": 3}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/student/wishes", body, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "déjà")
}

func TestStudentCreateWishInvalidSubject(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]any{"subject_id": 99999}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/student/wishes", body, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestStudentDeleteWish(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("DELETE", "/api/student/wishes/1", nil, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestStudentDeleteWishNotFound(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("DELETE", "/api/student/wishes/99999", nil, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "introuvable")
}

func TestStudentDeleteWishOtherStudent(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("DELETE", "/api/student/wishes/3", nil, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "non autorisé")
}

func TestStudentMyPfe(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/student/my-pfe", nil, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
	data, ok := result["data"].(map[string]any)
	if !ok {
		t.Fatal("❌ Student my-pfe: data manquant")
	}
	if data["pfe_code"] != "PFE-ISIL-2025-001" {
		t.Fatalf("❌ Student my-pfe: code PFE incorrect %v", data["pfe_code"])
	}
}

func TestStudentMyPfeNoAssignment(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/student/my-pfe", nil, h.AuthHeader(SeedStudentCHIM1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "aucun PFE")
}

func TestStudentListMeetings(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/student/my-pfe/meetings", nil, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestStudentCreateMeeting(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]any{
		"meeting_date": "2025-05-25T10:00:00Z",
		"duration":     30,
		"meeting_type": "presentiel",
		"topics":       "Discussion sur l'avancement",
	}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/student/my-pfe/meetings", body, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestStudentCreateMeetingValidation(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/student/my-pfe/meetings", map[string]string{}, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestStudentSubmitMemoire(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"memoire_url": "https://storage.supabase.co/memoires/test-memoire.pdf"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/student/my-pfe/memoire", body, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestStudentSubmitMemoireValidation(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("POST", "/api/student/my-pfe/memoire", map[string]string{}, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestStudentSubmitMemoireNoAssignment(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	body := map[string]string{"memoire_url": "https://storage.supabase.co/memoires/test.pdf"}
	resp, err := h.App.Test(newHTTPRequest("POST", "/api/student/my-pfe/memoire", body, h.AuthHeader(99999, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertErrorContains(t, result, "aucun PFE")
}

func TestStudentSoutenance(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/student/soutenance", nil, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestStudentNotifications(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := h.App.Test(newHTTPRequest("GET", "/api/student/notifications", nil, h.AuthHeader(SeedStudentISIL1ID, "student")))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}
