


import { json, error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';
import type { RequestHandler } from './$types';

const RESEND_ENDPOINT = 'https://api.resend.com/emails';

export const POST: RequestHandler = async ({ request }) => {
  const key = env.RESEND_KEY;
  if (!key) throw error(500, 'RESEND_KEY is not set on the server');

  const body = (await request.json()) as {
    to: string | string[];
    bcc?: string | string[];
    subject: string;
    html: string;
    from?: string;
    replyTo?: string;
  };

  const from = body.from ?? env.RESEND_FROM ?? 'PFE <noreply@esst-sup.com>';

  const res = await fetch(RESEND_ENDPOINT, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${key}`,
    },
    body: JSON.stringify({
      from,
      to: body.to,
      subject: body.subject,
      html: body.html,
      ...(body.bcc ? { bcc: body.bcc } : {}),
      ...(body.replyTo ? { reply_to: body.replyTo } : {}),
    }),
  });

  const data = (await res.json().catch(() => ({}))) as Record<string, unknown>;
  if (!res.ok) {
    throw error(res.status, (data?.message as string) || `Resend error ${res.status}`);
  }
  return json(data);
};
