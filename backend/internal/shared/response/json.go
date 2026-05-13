package response

import (
	"github.com/gofiber/fiber/v3"

	"pfe-backend/internal/shared/apperror"
)

// SuccessResponse est la structure standard pour toutes les réponses réussies.
type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

// ErrorResponse est la structure standard pour toutes les erreurs.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// PaginatedData est la structure pour les réponses paginées.
type PaginatedData struct {
	Items   any `json:"items"`
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// OK envoie une réponse JSON 200 avec les données fournies.
func OK(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// Created envoie une réponse JSON 201 avec les données fournies.
func Created(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// NoContent envoie une réponse 204 sans corps.
func NoContent(c fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// OKPaginated envoie une réponse JSON 200 avec des données paginées.
func OKPaginated(c fiber.Ctx, items any, total, page, perPage int) error {
	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Success: true,
		Data: PaginatedData{
			Items:   items,
			Total:   total,
			Page:    page,
			PerPage: perPage,
		},
	})
}

// Error envoie une réponse d'erreur structurée.
// Supporte à la fois les *apperror.Error et les erreurs génériques.
func Error(c fiber.Ctx, err error) error {
	if appErr, ok := err.(*apperror.Error); ok {
		return c.Status(appErr.StatusCode()).JSON(ErrorResponse{
			Success: false,
			Error:   appErr.Message,
		})
	}

	// Erreur inconnue -> 500
	return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
		Success: false,
		Error:   "Erreur interne du serveur",
	})
}

// ValidationError envoie une réponse 400 avec un message de validation.
func ValidationError(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
		Success: false,
		Error:   message,
	})
}

// Unauthorized envoie une réponse 401.
func Unauthorized(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
		Success: false,
		Error:   message,
	})
}

// Forbidden envoie une réponse 403.
func Forbidden(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
		Success: false,
		Error:   message,
	})
}

// NotFound envoie une réponse 404.
func NotFound(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		Success: false,
		Error:   message,
	})
}

// Conflict envoie une réponse 409.
func Conflict(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
		Success: false,
		Error:   message,
	})
}
