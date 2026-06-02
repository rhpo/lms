<script lang="ts">
  import { admin } from "$lib/api";
  import { showToast } from "$lib/utils/toast";
  import { sendMail, sendMailToUsers } from "$lib/email";
  import { Mail, Send, Users, GraduationCap, BookOpen } from "lucide-svelte";

  import Page from "$lib/components/ui/Page.svelte";
  import Input from "$lib/components/ui/Input/Input.svelte";
  import Button from "$lib/components/ui/Button.svelte";

  let recipientType = $state<
    "single" | "all_students" | "all_teachers" | "all"
  >("single");

  let sending = $state(false);
  let emailTo = $state("");
  let emailBody = $state("");
  let emailSubject = $state("");

  async function sendEmail() {
    sending = true;

    try {
      if (recipientType === "single")
        await sendMail({
          to: emailTo,
          html: emailBody,
          subject: emailSubject,
        });
      else {
        const users = await admin.listUsers();
        const targets = users.filter((u) => {
          if (recipientType === "all_students") return u.role === "student";
          if (recipientType === "all_teachers") return u.role === "teacher";
          return u.role === "student" || u.role === "teacher";
        });

        if (targets.length === 0) {
          showToast.error("Aucun destinataire trouvé.");
          return;
        }

        await sendMailToUsers(targets, emailSubject, emailBody);
      }


      emailTo = "";
      emailSubject = "";
      emailBody = "";
    } catch {

    } finally {
      sending = false;
    }
  }
</script>

<Page
  title="Envoi d'emails"
  subtitle="Envoyer un email à un utilisateur ou à un groupe d'utilisateurs."
>
  <div class="email-page">
    <div class="card">
      <div class="card-header">
        <Mail size={20} />
        <h2>Nouvel email</h2>
      </div>

      <form
        class="email-form"
        onsubmit={(e) => {
          e.preventDefault();
          sendEmail();
        }}
      >
        <div class="form-group">
          <label for="recipient-type">Destinataire</label>
          <select
            id="recipient-type"
            bind:value={recipientType}
            class="select-input"
          >
            <option value="single">Un utilisateur</option>
            <option value="all_students">Tous les étudiants</option>
            <option value="all_teachers">Tous les enseignants</option>
            <option value="all">Tous (étudiants + enseignants)</option>
          </select>
        </div>

        {#if recipientType === "single"}
          <Input
            label="Adresse email"
            bind:value={emailTo}
            placeholder="exemple@esst-sup.com"
            type="email"
            required
          />
        {:else}
          <div class="info-banner">
            {#if recipientType === "all_students"}
              <GraduationCap size={16} />
              <span>L'email sera envoyé à tous les étudiants.</span>
            {:else if recipientType === "all_teachers"}
              <BookOpen size={16} />
              <span>L'email sera envoyé à tous les enseignants.</span>
            {:else}
              <Users size={16} />
              <span
                >L'email sera envoyé à tous les utilisateurs (étudiants +
                enseignants).</span
              >
            {/if}
          </div>
        {/if}

        <Input
          label="Objet"
          bind:value={emailSubject}
          placeholder="Objet de l'email"
          required
        />

        <div class="form-group">
          <label for="email-body">Corps du message (HTML SUPPORT)</label>
          <textarea
            id="email-body"
            class="textarea-input"
            bind:value={emailBody}
            placeholder="Bonjour <b>$name</b>, je vous prie de passer de choisir un sujet PFE."
            rows="10"
            required
          ></textarea>
        </div>

        <div class="form-actions">
          <Button
            variant="primary"
            label={sending ? "Envoi en cours..." : "Envoyer"}
            Icon={Send}
            type="submit"
            disabled={sending ||
              !emailSubject ||
              !emailBody ||
              (recipientType === "single" && !emailTo)}
          />
        </div>
      </form>
    </div>
  </div>
</Page>

<style>
  .card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: 1.5rem;
  }

  .card-header {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    margin-bottom: 1.25rem;
    color: var(--color-text);

    h2 {
      margin: 0;
      font-size: var(--text-lg);
      font-weight: 600;
    }
  }

  .email-form {
    display: flex;
    flex-direction: column;
    gap: 0.85rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;

    label {
      font-size: 13px;
      font-weight: 500;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: var(--color-text);
    }
  }

  .select-input {
    padding: 0.8rem 1rem;
    border: 1px solid var(--color-background-200);
    border-radius: 1rem;
    background: var(--color-background-100);
    color: var(--color-text);
    font-family: var(--font-sans);
    font-size: 1rem;
    cursor: pointer;

    &:focus {
      outline: none;
      border-color: var(--color-accent);
    }
  }

  .textarea-input {
    width: 100%;
    padding: 0.8rem 1rem;
    border: 1px solid var(--color-background-200);
    border-radius: 1rem;
    background: var(--color-background-100);
    color: var(--color-text);
    font-family: var(--font-sans);
    font-size: 0.95rem;
    resize: vertical;
    min-height: 160px;

    &:focus {
      outline: none;
      border-color: var(--color-accent);
    }

    &::placeholder {
      color: var(--color-text-muted);
    }
  }

  .info-banner {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.65rem 0.85rem;
    background: color-mix(
      in srgb,
      var(--color-accent) 8%,
      var(--color-surface)
    );
    border: 1px solid color-mix(in srgb, var(--color-accent) 15%, transparent);
    border-radius: 8px;
    font-size: var(--text-sm);
    color: var(--color-text);
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
    margin-top: 0.5rem;
  }
</style>
