import { company } from '$lib/api';
import { error } from '@sveltejs/kit';

export async function load({ params }) {
  const id = parseInt(params.id);
  if (isNaN(id)) throw error(400, 'ID invalide');

  try {
    const subject = await company.getSubject(id);
    return { subject };
  } catch {
    throw error(404, 'Sujet introuvable');
  }
}
