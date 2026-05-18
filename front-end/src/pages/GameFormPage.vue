<script setup lang="ts">
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { fetchGame, createGame, updateGame, searchPlatformGames, imageUrl } from '@/api/client'
import type { Platform, GameInput } from '@/api/types'
import PlatformGameSearch from '@/components/backoffice/PlatformGameSearch.vue'
import ImageCropper from '@/components/backoffice/ImageCropper.vue'

const router = useRouter()
const route = useRoute()

const editId = computed(() => route.params.id as string | undefined)
const isEdit = computed(() => !!editId.value)

const platform = ref<Platform>('Steam')
const platformId = ref('')
const title = ref('')
const releaseYear = ref<number | undefined>()
const developer = ref('')
const coverSourceUrl = ref('')
const bannerSourceUrl = ref('')
const coverData = ref('')
const bannerData = ref('')

const platforms: Platform[] = ['Steam', 'Fightcade', 'MAME']

const platformLabels: Record<Platform, { label: string; placeholder: string }> = {
  Steam: { label: 'App ID', placeholder: 'Search Steam games…' },
  Fightcade: { label: 'Game ID', placeholder: 'Search Fightcade games…' },
  MAME: { label: 'Driver Name', placeholder: 'Search MAME games…' },
}

const currentLabel = computed(() => platformLabels[platform.value])
const searchFn = computed(() => (q: string) => searchPlatformGames(platform.value, q))

watch(platform, () => { platformId.value = '' })

onMounted(async () => {
  if (!editId.value) return
  const game = await fetchGame(editId.value)
  if (!game) { router.replace('/backoffice'); return }

  platform.value = game.platform
  title.value = game.title
  releaseYear.value = game.releaseYear
  developer.value = game.developer
  coverSourceUrl.value = imageUrl(game.coverFilename)
  bannerSourceUrl.value = imageUrl(game.bannerFilename)
  coverData.value = coverSourceUrl.value
  bannerData.value = bannerSourceUrl.value

  await nextTick()
  platformId.value = game.appId
})

function dataUrlToBlob(dataUrl: string): Blob {
  const [header, base64] = dataUrl.split(',')
  const mime = header.match(/:(.*?);/)![1]
  const bytes = atob(base64)
  const buffer = new Uint8Array(bytes.length)
  for (let i = 0; i < bytes.length; i++) buffer[i] = bytes.charCodeAt(i)
  return new Blob([buffer], { type: mime })
}

async function save() {
  const input: GameInput = {
    title: title.value,
    platform: platform.value,
    releaseYear: releaseYear.value ?? 0,
    developer: developer.value,
    appId: platformId.value,
    cover: dataUrlToBlob(coverData.value),
    banner: dataUrlToBlob(bannerData.value),
  }

  if (isEdit.value) {
    await updateGame(editId.value!, input)
  } else {
    await createGame(input)
  }

  router.push('/backoffice')
}
</script>

<template>
  <div class="form-page px-12 py-8">
    <RouterLink to="/backoffice" class="back-link">&larr; Back to Library</RouterLink>
    <h1 class="text-2xl font-bold mt-6 mb-6">{{ isEdit ? 'Edit Game' : 'Add Game' }}</h1>

    <form class="game-form" @submit.prevent="save">
      <fieldset>
        <legend>Launch Configuration</legend>

        <label>
          Platform
          <select v-model="platform">
            <option v-for="p in platforms" :key="p" :value="p">{{ p }}</option>
          </select>
        </label>

        <label>
          {{ currentLabel.label }}
          <PlatformGameSearch
            :key="platform"
            v-model="platformId"
            :search="searchFn"
            :placeholder="currentLabel.placeholder"
          />
        </label>
      </fieldset>

      <fieldset>
        <legend>Game Metadata</legend>

        <label>
          Title
          <input v-model="title" type="text" placeholder="Display name" required>
        </label>

        <label>
          Release Year
          <input v-model.number="releaseYear" type="number" min="1970" max="2099" placeholder="e.g. 1999" required>
        </label>

        <label>
          Developer
          <input v-model="developer" type="text" placeholder="e.g. Capcom" required>
        </label>

        <label>
          Cover Image URL
          <input v-model="coverSourceUrl" type="url" placeholder="https://..." required>
        </label>
        <ImageCropper
          v-if="coverSourceUrl"
          :url="coverSourceUrl"
          :frame-width="200"
          :frame-height="267"
          @cropped="coverData = $event"
        />

        <label>
          Banner Image URL
          <input v-model="bannerSourceUrl" type="url" placeholder="https://..." required>
        </label>
        <ImageCropper
          v-if="bannerSourceUrl"
          :url="bannerSourceUrl"
          :frame-width="640"
          :frame-height="360"
          @cropped="bannerData = $event"
        />
      </fieldset>

      <div class="form-actions">
        <button type="submit" class="save-btn">{{ isEdit ? 'Save Changes' : 'Add Game' }}</button>
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

.game-form {
  max-width: 680px;
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

input, select {
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

input:focus, select:focus {
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
