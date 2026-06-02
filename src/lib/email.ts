



//


//



//


import type { Profile } from './types';
import { showToast } from './utils/toast';

const EMAIL_ENDPOINT = '/email';


export type MailRecipient = Pick<Profile, 'email'> & Partial<Pick<Profile, 'full_name'>>;

export interface MailPayload {
  to: string | string[];
  bcc?: string | string[];
  subject: string;
  html: string;
  from?: string;
  replyTo?: string;
}

interface ResendResponse {
  id?: string;
  message?: string;
  name?: string;
}

/** "ramy.hadid42@esst-sup.com" -> "Ramy Hadid" (non-letters stripped). */
const prettifyEmail = (email: string): string =>
  email
    .split('@')[0]
    .split(/[._-]+/)
    .map((w) => w.replace(/[^a-zA-Z]/g, ''))
    .filter(Boolean)
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1).toLowerCase())
    .join(' ');

const applyTemplate = (html: string, name: string): string =>
  html.replace(/\$name/g, name).replace(/\{\{\s*name\s*\}\}/gi, name);

/** Raw POST to the server proxy. No toast, no templating. */
async function postMail(payload: MailPayload): Promise<ResendResponse> {
  const res = await fetch(EMAIL_ENDPOINT, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });

  const data = (await res.json().catch(() => ({}))) as ResendResponse;
  if (!res.ok) {
    throw new Error(data?.message || `Email send failed (${res.status})`);
  }
  return data;
}

async function withToast<T>(fn: () => Promise<T>, successMsg: string): Promise<T> {
  try {
    const res = await fn();
    showToast.success(successMsg);
    return res;
  } catch (err) {
    showToast.error("Échec de l'envoi de l'email", String(err));
    throw err;
  }
}

/**
 * Lowest-level send.
 * - Single string `to`: `$name` in the html is replaced by that email.
 * - Array `to`: sent as-is (all visible to each other). Prefer `sendMailToMany`
 *   which uses BCC instead.
 */
export function sendMail(payload: MailPayload): Promise<ResendResponse> {
  const html =
    typeof payload.to === 'string'
      ? applyTemplate(payload.html, prettifyEmail(payload.to))
      : payload.html;

  return withToast(
    () => postMail({ ...payload, html }),
    'Email envoyé avec succès !',
  );
}

/**
 * Bulk send via BCC - one Resend call, recipients are hidden from each other.
 * `to` defaults to the `from` address (Resend requires `to` to be non-empty).
 */
export function sendMailToMany(
  emails: string[],
  subject: string,
  html: string,
  opts: { from?: string; replyTo?: string; to?: string } = {},
): Promise<ResendResponse> {
  return withToast(
    () =>
      postMail({
        to: opts.to ?? opts.from ?? 'noreply@codiha.com',
        bcc: emails,
        subject,
        html,
        from: opts.from,
        replyTo: opts.replyTo,
      }),
    `Email envoyé à ${emails.length} destinataire${emails.length > 1 ? 's' : ''}.`,
  );
}

/** Send to a single User/Profile. `$name` -> full_name, else email. */
export function sendMailTo(
  user: MailRecipient,
  subject: string,
  html: string,
  opts: { from?: string; replyTo?: string } = {},
): Promise<ResendResponse> {
  const personalisedHtml = applyTemplate(
    html,
    user.full_name ?? prettifyEmail(user.email),
  );
  return withToast(
    () => postMail({ to: user.email, subject, html: personalisedHtml, ...opts }),
    'Email envoyé avec succès !',
  );
}

/**
 * Bulk send to many User/Profile objects via BCC - single API call.
 * Per-recipient `$name` substitution is not possible with BCC; any `$name`
 * tokens are left as a literal placeholder.
 */
export function sendMailToUsers(
  users: MailRecipient[],
  subject: string,
  html: string,
  opts: { from?: string; replyTo?: string } = {},
): Promise<ResendResponse> {
  return sendMailToMany(
    users.map((u) => u.email),
    subject,
    html,
    opts,
  );
}
