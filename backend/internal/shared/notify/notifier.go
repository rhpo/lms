// Package notify provides a single entry-point for sending in-app
// notifications.  Every call fans out to a list of Channels; today the
// channels are the SQLite insert and Resend email delivery.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"pfe-backend/internal/entity"
	"pfe-backend/internal/repository"
)

// ──────────────────────────────────────────────────────────────────────
// Types & notification kinds (keep in sync with the frontend constants)
// ──────────────────────────────────────────────────────────────────────

const (
	TypeValidationRequise = "validation_requise"
	TypeAffectation       = "affectation"
	TypeJury              = "jury"
	TypeDisponibilite     = "disponibilite"
	TypeSujet             = "sujet"
)

// Message is the payload written into each notification row.
type Message struct {
	Message string `json:"message"`
}

// Channel is any transport that can deliver a notification.
// Implement this interface to add websockets, email, etc.
type Channel interface {
	Send(n *entity.Notification) error
}

// ──────────────────────────────────────────────────────────────────────
// Notifier — the only thing callers interact with
// ──────────────────────────────────────────────────────────────────────

// Notifier holds all channels and exposes high-level helpers.
type Notifier struct {
	channels    []Channel
	profileRepo *repository.ProfileRepository
}

// New creates a Notifier that persists to SQLite, sends email via Resend,
// and can resolve admin IDs for "notify all admins" helpers.
func New(notifRepo *repository.NotificationRepository, profileRepo *repository.ProfileRepository, resendAPIKey string) *Notifier {
	channels := []Channel{&dbChannel{repo: notifRepo}}

	// Add email channel if Resend API key is configured and not a test key
	if resendAPIKey != "" && resendAPIKey != "test-resend-key" {
		channels = append(channels, &emailChannel{
			apiKey:      resendAPIKey,
			profileRepo: profileRepo,
		})
	}

	return &Notifier{
		channels:    channels,
		profileRepo: profileRepo,
	}
}

// AddChannel registers an extra delivery channel (websocket, etc.).
func (n *Notifier) AddChannel(ch Channel) {
	n.channels = append(n.channels, ch)
}

// ──────────────────────────────────────────────────────────────────────
// Low-level send
// ──────────────────────────────────────────────────────────────────────

// Send dispatches a notification to every registered channel.
// Errors are logged but never propagated — notifications must not
// break the main request flow.
func (n *Notifier) Send(recipientID int64, notifType, message string) {
	payload, _ := json.Marshal(Message{Message: message})
	notif := &entity.Notification{
		RecipientID: recipientID,
		Type:        notifType,
		Payload:     string(payload),
	}
	for _, ch := range n.channels {
		if err := ch.Send(notif); err != nil {
			log.Printf("[notify] channel error: %v", err)
		}
	}
}

// SendMany sends the same notification to multiple recipients.
func (n *Notifier) SendMany(recipientIDs []int64, notifType, message string) {
	for _, id := range recipientIDs {
		n.Send(id, notifType, message)
	}
}

// ──────────────────────────────────────────────────────────────────────
// High-level helpers
// ──────────────────────────────────────────────────────────────────────

// AdminIDs returns the profile IDs of all admin users.
func (n *Notifier) AdminIDs() []int64 {
	profiles, err := n.profileRepo.FindAll()
	if err != nil {
		log.Printf("[notify] failed to fetch admin IDs: %v", err)
		return nil
	}
	var ids []int64
	for _, p := range profiles {
		if p.Role == "admin" {
			ids = append(ids, p.ID)
		}
	}
	return ids
}

// NotifyAdmins sends a notification to every admin.
func (n *Notifier) NotifyAdmins(notifType, message string) {
	n.SendMany(n.AdminIDs(), notifType, message)
}

// ──────────────────────────────────────────────────────────────────────
// dbChannel — SQLite persistence (the default channel)
// ──────────────────────────────────────────────────────────────────────

type dbChannel struct {
	repo *repository.NotificationRepository
}

func (d *dbChannel) Send(n *entity.Notification) error {
	return d.repo.Insert(n)
}

// ──────────────────────────────────────────────────────────────────────
// emailChannel — Resend email delivery
// ──────────────────────────────────────────────────────────────────────

type emailChannel struct {
	apiKey      string
	profileRepo *repository.ProfileRepository
}

// notifTypeSubjects maps notification types to human-readable email subjects.
var notifTypeSubjects = map[string]string{
	TypeValidationRequise: "Action requise — Plateforme PFE",
	TypeAffectation:       "Mise à jour de votre PFE",
	TypeJury:              "Soutenance — Plateforme PFE",
	TypeDisponibilite:     "Suivi PFE — Nouvelle information",
	TypeSujet:             "Sujet PFE — Mise à jour",
}

func (e *emailChannel) Send(n *entity.Notification) error {
	// Resolve recipient email
	profile, err := e.profileRepo.FindByID(n.RecipientID)
	if err != nil || profile == nil || profile.Email == "" {
		return nil // silently skip — no email to send to
	}

	// Extract message from payload
	var msg Message
	_ = json.Unmarshal([]byte(n.Payload), &msg)
	if msg.Message == "" {
		return nil
	}

	subject := notifTypeSubjects[n.Type]
	if subject == "" {
		subject = "Notification — Plateforme PFE"
	}

	// Build Resend API request
	body := map[string]interface{}{
		"from":    "Plateforme PFE <noreply@codiha.com>",
		"to":      []string{profile.Email},
		"subject": subject,
		"html": fmt.Sprintf(`
			<div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 600px; margin: 0 auto; padding: 32px 24px;">
				<div style="background: #1a56db; padding: 20px 24px; border-radius: 12px 12px 0 0;">
					<h1 style="color: #ffffff; margin: 0; font-size: 18px; font-weight: 600;">Plateforme PFE</h1>
				</div>
				<div style="border: 1px solid #e5e7eb; border-top: none; border-radius: 0 0 12px 12px; padding: 24px;">
					<p style="color: #374151; font-size: 15px; line-height: 1.6; margin: 0 0 16px;">Bonjour <strong>%s</strong>,</p>
					<p style="color: #374151; font-size: 15px; line-height: 1.6; margin: 0 0 24px;">%s</p>
					<hr style="border: none; border-top: 1px solid #e5e7eb; margin: 24px 0;">
					<p style="color: #9ca3af; font-size: 12px; margin: 0;">Ce message a été envoyé automatiquement par la plateforme PFE. Ne répondez pas à cet email.</p>
				</div>
			</div>
		`, profile.FullName, msg.Message),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("email marshal error: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("email request error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("email send error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("[notify/email] Resend API returned %d for recipient %s", resp.StatusCode, profile.Email)
	}

	return nil
}
