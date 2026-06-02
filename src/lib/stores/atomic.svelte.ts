/**
 * Atomic reference-data store.
 *
 * Loaded once on first authenticated page mount.  Every page inside (dashboard)
 * can import { atomic } and read .domains / .departments / .specialities / .niveaux
 * synchronously - the data is guaranteed to be present (the layout gate blocks
 * rendering until all four are resolved, and throws if any fetch fails).
 */

import { ref } from '$lib/api';
import type { Domain, Department, Speciality, YearType } from '$lib/types';



let _domains = $state<Domain[]>([]);
let _departments = $state<Department[]>([]);
let _specialities = $state<Speciality[]>([]);
let _ready = $state(false);
let _error = $state<string | null>(null);

/** Static - never changes, but co-located here so every page grabs from one place. */
const NIVEAUX: YearType[] = ['licence', 'master'];



export const atomic = {
  /* ── reactive getters ─────────────────────────────────────────────────── */
  get domains() { return _domains; },
  get departments() { return _departments; },
  get specialities() { return _specialities; },
  get niveaux() { return NIVEAUX; },
  get ready() { return _ready; },
  get error() { return _error; },

  /**
   * Fetch all reference data in parallel.
   * Throws (and sets .error) if ANY request fails.
   * Safe to call multiple times - subsequent calls are no-ops once loaded.
   */
  async load() {
    if (_ready) return;

    try {
      const [domains, departments, specialities] = await Promise.all([
        ref.domains(),
        ref.departments(),
        ref.specialities(),
      ]);

      if (domains === null || domains === undefined ||
        departments === null || departments === undefined ||
        specialities === null || specialities === undefined) {
        throw new Error('Données de référence manquantes - le serveur a renvoyé une réponse vide.');
      }

      _domains = domains;
      _departments = departments;
      _specialities = specialities;
      _ready = true;
      _error = null;
    } catch (err) {
      _error = err instanceof Error
        ? err.message
        : 'Impossible de charger les données de référence.';
      throw err;
    }
  },

  /** Force a full reload (e.g. after admin edits settings). */
  async reload() {
    _ready = false;
    _error = null;
    await this.load();
  },
};
