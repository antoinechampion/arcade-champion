<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { fetchAllGames, deleteGame } from '@/api/client'
import type { Game } from '@/api/types'

const router = useRouter()
const games = ref<Game[]>([])

async function loadGames() {
  games.value = await fetchAllGames()
}

async function remove(game: Game) {
  if (!confirm(`Delete "${game.title}"?`)) return
  await deleteGame(game.id)
  await loadGames()
}

onMounted(loadGames)
</script>

<template>
  <div class="backoffice px-12 py-8">
    <RouterLink to="/" class="back-link">&larr; Back to Home</RouterLink>

    <div class="flex items-center justify-between mt-6 mb-6">
      <h1 class="text-2xl font-bold">Back Office</h1>
      <RouterLink to="/backoffice/add" class="add-btn">+ Add Game</RouterLink>
    </div>

    <table class="game-table">
      <thead>
        <tr>
          <th>Title</th>
          <th>Platform</th>
          <th>Year</th>
          <th>Developer</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="game in games" :key="game.id">
          <td>{{ game.title }}</td>
          <td>{{ game.platform }}</td>
          <td>{{ game.releaseYear }}</td>
          <td>{{ game.developer }}</td>
          <td class="actions">
            <RouterLink :to="`/backoffice/edit/${game.id}`" class="action-link">Edit</RouterLink>
            <button class="action-link danger" @click="remove(game)">Delete</button>
          </td>
        </tr>
      </tbody>
    </table>

    <p v-if="!games.length" class="opacity-50 text-sm mt-4">No games in library.</p>
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

.add-btn {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  background: var(--color-primary);
  color: var(--color-text);
  font-size: 0.875rem;
  font-weight: 500;
  text-decoration: none;
  transition: background 0.2s ease;
}

.add-btn:hover {
  background: var(--color-primary-light);
}

.game-table {
  width: 100%;
  border-collapse: collapse;
}

.game-table th {
  text-align: left;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  opacity: 0.5;
  padding: 0.5rem 0.75rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.game-table td {
  padding: 0.625rem 0.75rem;
  font-size: 0.875rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.actions {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
}

.action-link {
  background: none;
  border: none;
  color: var(--color-primary-light);
  font-size: 0.8125rem;
  cursor: pointer;
  text-decoration: none;
  opacity: 0.7;
  transition: opacity 0.2s ease;
}

.action-link:hover {
  opacity: 1;
}

.action-link.danger {
  color: #f87171;
}
</style>
