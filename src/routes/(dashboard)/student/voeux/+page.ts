import { student } from '$lib/api';

export async function load() {
  const [wishes, settings, pfe] = await Promise.all([
    student.listWishes().catch(() => []),
    student.settings().catch(() => null),
    student.getMyPFE().catch(() => null),
  ]);
  return {
    wishes: wishes ?? [],
    maxWishes: settings?.max_wishes ?? 5,
    hasPfe: pfe != null,
  };
}
