<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { fetchSettings, updateSettings } from '@/api/client'

const router = useRouter()

const fightcadeUsername = ref('')
const fightcadePassword = ref('')
const fightcadeCookie = ref('')

onMounted(async () => {
  const s = await fetchSettings()
  fightcadeUsername.value = s.fightcadeUsername
  fightcadePassword.value = s.fightcadePassword
  fightcadeCookie.value = s.fightcadeCookie
})

async function save() {
  await updateSettings({
    fightcadeUsername: fightcadeUsername.value,
    fightcadePassword: fightcadePassword.value,
    fightcadeCookie: fightcadeCookie.value,
  })
  router.push('/backoffice')
}
</script>

<template>
  <div class="settings-page px-12 py-8">
    <RouterLink to="/backoffice" class="back-link">&larr; Back to Library</RouterLink>
    <h1 class="text-2xl font-bold mt-6 mb-6">Settings</h1>

    <form class="settings-form" @submit.prevent="save">
      <fieldset>
        <legend>Fightcade</legend>

        <label>
          Username
          <input v-model="fightcadeUsername" type="text" placeholder="Fightcade username">
        </label>

        <label>
          Password
          <input v-model="fightcadePassword" type="password" placeholder="Fightcade password">
        </label>

        <label>
          Cookie
          <input v-model="fightcadeCookie" type="text" placeholder="Session cookie">
        </label>
      </fieldset>

      <div class="form-actions">
        <button type="submit" class="save-btn">Save Settings</button>
        <RouterLink to="/backoffice" class="cancel-link">Cancel</RouterLink>
      </div>
    </form>
  </div>
</template>

<style scoped>
.back-link {
  color: var(--color-text);
  opacity: 0.6;
  font-size: 0.875rem;
  text-decoration: none;
  transition: opacity 0.2s ease;
}

.back-link:hover,
.back-link:focus-visible {
  opacity: 1;
}

.settings-form {
  max-width: 480px;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

fieldset {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.875rem;
}

legend {
  font-size: 0.8125rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  opacity: 0.6;
  padding: 0 0.5rem;
}

label {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  font-size: 0.8125rem;
  opacity: 0.8;
}

input {
  padding: 0.5rem 0.625rem;
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  background: rgba(255, 255, 255, 0.06);
  color: var(--color-text);
  font-size: 0.875rem;
  font-family: inherit;
  outline: none;
  transition: border-color 0.2s ease;
}

input:focus {
  border-color: var(--color-primary-light);
}

input::placeholder {
  opacity: 0.3;
}

.form-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.save-btn {
  padding: 0.5rem 1.25rem;
  border-radius: 6px;
  border: none;
  background: var(--color-primary);
  color: var(--color-text);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s ease;
}

.save-btn:hover {
  background: var(--color-primary-light);
}

.cancel-link {
  color: var(--color-text);
  opacity: 0.5;
  font-size: 0.875rem;
  text-decoration: none;
  transition: opacity 0.2s ease;
}

.cancel-link:hover {
  opacity: 0.8;
}
</style>
