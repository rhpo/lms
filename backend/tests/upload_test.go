package tests

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func createMultipartRequest(app *fiber.App, method, url, fieldName, fileName, content string, authToken string) (*http.Response, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur création form file: %v", err)
	}
	_, err = io.WriteString(part, content)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur écriture contenu: %v", err)
	}
	writer.Close()

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur création requête: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	return app.Test(req)
}

func TestUploadAvatar(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()


	resp, err := createMultipartRequest(h.App, "POST", "/api/upload/avatar", "file", "avatar.png", "fake-png-content", h.AuthToken(SeedAdminID, "admin"))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestUploadAvatarUnauthorized(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := createMultipartRequest(h.App, "POST", "/api/upload/avatar", "file", "avatar.png", "fake-png-content", "")
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestUploadAvatarInvalidType(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := createMultipartRequest(h.App, "POST", "/api/upload/avatar", "file", "avatar.txt", "fake-text-content", h.AuthToken(SeedAdminID, "admin"))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestUploadCompanyLogo(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := createMultipartRequest(h.App, "POST", "/api/upload/company-logo", "file", "logo.png", "fake-png-content", h.AuthToken(SeedCompany1ID, "company"))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestUploadCompanyLogoUnauthorized(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := createMultipartRequest(h.App, "POST", "/api/upload/company-logo", "file", "logo.png", "fake-png-content", "")
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestUploadCompanyLogoWrongRole(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := createMultipartRequest(h.App, "POST", "/api/upload/company-logo", "file", "logo.png", "fake-png-content", h.AuthToken(SeedStudentISIL1ID, "student"))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestUploadCompanyLogoInvalidType(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := createMultipartRequest(h.App, "POST", "/api/upload/company-logo", "file", "logo.gif", "fake-gif-content", h.AuthToken(SeedCompany1ID, "company"))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestUploadMemoire(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := createMultipartRequest(h.App, "POST", "/api/upload/memoire", "file", "memoire.pdf", "fake-pdf-content", h.AuthToken(SeedStudentISIL1ID, "student"))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestUploadMemoireUnauthorized(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := createMultipartRequest(h.App, "POST", "/api/upload/memoire", "file", "memoire.pdf", "fake-pdf-content", "")
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestUploadMemoireWrongRole(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := createMultipartRequest(h.App, "POST", "/api/upload/memoire", "file", "memoire.pdf", "fake-pdf-content", h.AuthToken(SeedTeacherISIL1ID, "teacher"))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestUploadMemoireInvalidType(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := createMultipartRequest(h.App, "POST", "/api/upload/memoire", "file", "memoire.txt", "fake-text-content", h.AuthToken(SeedStudentISIL1ID, "student"))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestUploadProfilePhoto(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := createMultipartRequest(h.App, "POST", "/api/profile/avatar", "file", "profile.png", "fake-png-content", h.AuthToken(SeedAdminID, "admin"))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertSuccess(t, result)
}

func TestUploadProfilePhotoUnauthorized(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := createMultipartRequest(h.App, "POST", "/api/profile/avatar", "file", "profile.png", "fake-png-content", "")
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}

func TestUploadProfilePhotoInvalidType(t *testing.T) {
	h := NewTestHelper()
	defer h.Close()

	resp, err := createMultipartRequest(h.App, "POST", "/api/profile/avatar", "file", "profile.bmp", "fake-bmp-content", h.AuthToken(SeedAdminID, "admin"))
	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	result := MustParseResponse(resp)
	AssertError(t, result)
}
