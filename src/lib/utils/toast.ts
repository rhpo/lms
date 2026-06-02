/**
 * Centralized toast helpers built on top of svelte-sonner.
 *
 * Usage:
 *   import { showToast } from '$lib/utils/toast';
 *   showToast.success('Opération réussie !');
 *   showToast.error('Une erreur est survenue.');
 *   showToast.info('Information importante.');
 *   showToast.loading('Chargement…');
 *   showToast.promise(myPromise, { loading: '...', success: '...', error: '...' });
 */

import { toast } from 'svelte-sonner';


const BASE_DURATION = 4000;

export const showToast = {
	/** Green success toast */
	success: (message: string, description?: string) =>
		toast.success(message, {
			description,
			duration: BASE_DURATION,
		}),

	/** Red error toast */
	error: (message: string, description?: string) =>
		toast.error(message, {
			description,
			duration: BASE_DURATION + 2000,
		}),

	/** Neutral info toast */
	info: (message: string, description?: string) =>
		toast.info(message, {
			description,
			duration: BASE_DURATION,
		}),

	/** Warning toast */
	warning: (message: string, description?: string) =>
		toast.warning(message, {
			description,
			duration: BASE_DURATION,
		}),

	/** Persistent loading toast - returns the toast id so you can dismiss it */
	loading: (message: string) =>
		toast.loading(message, {
			duration: Infinity,
		}),

	/** Promise-aware toast. Automatically switches from loading -> success/error. */
	promise: <T>(
		promise: Promise<T>,
		messages: { loading: string; success: string; error: string }
	) =>
		toast.promise(promise, {
			loading: messages.loading,
			success: messages.success,
			error: messages.error,
		}),

	/** Dismiss a specific toast by id */
	dismiss: (id?: string | number) => toast.dismiss(id),
};
