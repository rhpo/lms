<script lang="ts">
  import { enhance } from '$app/forms';
  import { ShieldAlert, User, CheckCircle2 } from 'lucide-svelte';

  let { profiles = [] } = $props();

  let selectedId = $state('');

  const roles: Record<string, { color: string; label: string }> = {
    admin: { color: '#ef4444', label: 'Admin' },
    teacher: { color: '#10b981', label: 'Teacher' },
    student: { color: '#3b82f6', label: 'Student' },
    company: { color: '#f59e0b', label: 'Company' }
  };
</script>

<div class="login-test-card">
  <div class="header">
    <ShieldAlert size="1.5rem" />
    <div class="title">
      <h3>Login Test (Dev Only)</h3>
      <p>Quickly switch between accounts for testing.</p>
    </div>
  </div>

  <div class="profiles-list">
    {#each profiles as profile}
      <form method="POST" action="?/testLogin" use:enhance class="profile-item">
        <input type="hidden" name="userId" value={profile.id} />
        
        <button type="submit" class="profile-btn">
          <div class="avatar">
            {#if profile.avatar_url}
              <img src={profile.avatar_url} alt={profile.full_name} />
            {:else}
              <div class="placeholder"><User size="1.2rem" /></div>
            {/if}
          </div>
          
          <div class="info">
            <span class="name">{profile.full_name}</span>
            <div class="role-badge" style:--role-color={roles[profile.role]?.color || '#94a3b8'}>
              {roles[profile.role]?.label || profile.role}
            </div>
          </div>

          <div class="check">
            <CheckCircle2 size="1.2rem" />
          </div>
        </button>
      </form>
    {:else}
      <div class="empty">
        <p>No profiles found. Run the seed script first!</p>
      </div>
    {/each}
  </div>
</div>

<style>
  .login-test-card {
    background: var(--color-background-100);
    border: 2px dashed var(--color-accent);
    border-radius: 1.5rem;
    padding: 1.5rem;
    width: 100%;
    max-width: 400px;
    margin-top: 2rem;
  }

  .header {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1.5rem;
    color: var(--color-accent);

    h3 {
      margin: 0;
      font-size: 1.1rem;
      font-weight: 700;
    }

    p {
      margin: 0;
      font-size: 0.85rem;
      color: var(--color-text-muted);
    }
  }

  .profiles-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    max-height: 300px;
    overflow-y: auto;
    padding-right: 0.5rem;

    &::-webkit-scrollbar {
      width: 4px;
    }
    &::-webkit-scrollbar-thumb {
      background: var(--color-background-200);
      border-radius: 99px;
    }
  }

  .profile-btn {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 1rem;
    background: var(--color-background-200);
    border: 1px solid transparent;
    padding: 0.75rem;
    border-radius: 1rem;
    cursor: pointer;
    transition: all var(--transition-fast);
    text-align: left;

    &:hover {
      background: var(--color-background-300);
      border-color: var(--color-accent);
      transform: translateX(4px);
    }

    &:hover .check {
      opacity: 1;
    }
  }

  .avatar {
    width: 40px;
    height: 40px;
    border-radius: 12px;
    overflow: hidden;
    background: var(--color-background-300);

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .placeholder {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      color: var(--color-text-muted);
    }
  }

  .info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.2rem;

    .name {
      font-size: 0.95rem;
      font-weight: 600;
      color: var(--color-text);
    }
  }

  .role-badge {
    font-size: 0.7rem;
    text-transform: uppercase;
    font-weight: 700;
    color: var(--role-color);
    letter-spacing: 0.5px;
  }

  .check {
    opacity: 0;
    color: var(--color-accent);
    transition: opacity 0.2s;
  }

  .empty {
    text-align: center;
    padding: 2rem;
    color: var(--color-text-muted);
    font-style: italic;
  }
</style>
