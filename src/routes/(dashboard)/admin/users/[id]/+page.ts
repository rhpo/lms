import { admin } from '$lib/api';
import type { Profile, Department, Domain, Speciality, Promotion } from '$lib/types';
import { error } from '@sveltejs/kit';

export async function load({ params }) {
  const id = parseInt(params.id);
  if (isNaN(id)) {
    throw error(400, 'ID invalide');
  }

  try {
    const [profile, departments, domains, specialities, promotions] = await Promise.all([
      admin.getUser(id),
      admin.listDepartments(),
      admin.listDomains(),
      admin.listSpecialities(),
      admin.listPromotions(),
    ]);

    return {
      profile: profile as Profile,
      departments: departments as Department[],
      domains: domains as Domain[],
      specialities: specialities as Speciality[],
      promotions: promotions as Promotion[],
    };
  } catch (err) {
    throw error(404, 'Utilisateur introuvable ou erreur de chargement');
  }
}
