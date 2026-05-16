<script setup lang="ts">
import { ref } from 'vue'
import type { PlatformSearchResult } from '@/api/types'

const props = defineProps<{
  placeholder?: string
  search: (query: string) => Promise<PlatformSearchResult[]>
}>()

const model = defineModel<string>({ required: true })
const query = ref('')
const results = ref<PlatformSearchResult[]>([])
const selectedName = ref('')

async function doSearch() {
  results.value = await props.search(query.value)
}

function select(result: PlatformSearchResult) {
  model.value = result.platformId
  selectedName.value = result.name
  query.value = ''
  results.value = []
}
</script>

<template>
  <div class="platform-search">
    <div class="search-row">
      <input
        v-model="query"
        type="text"
        :placeholder="placeholder"
        @keydown.enter.prevent="doSearch"
      >
      <button type="button" class="search-btn" @click="doSearch">Search</button>
    </div>

    <p v-if="selectedName" class="selected-name">
      {{ selectedName }}
    </p>

    <table v-if="results.length" class="results-table">
      <thead>
        <tr>
          <th>Name</th>
          <th>ID</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="r in results"
          :key="r.platformId"
          :class="{ active: r.platformId === model }"
          @click="select(r)"
        >
          <td>{{ r.name }}</td>
          <td><code>{{ r.platformId }}</code></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.platform-search {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.search-row {
  display: flex;
  gap: 0.5rem;
}

.search-row input {
  flex: 1;
}

.search-btn {
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  background: rgba(255, 255, 255, 0.06);
  color: var(--color-text);
  font-size: 0.8125rem;
  cursor: pointer;
  transition: background 0.2s ease;
}

.search-btn:hover {
  background: rgba(255, 255, 255, 0.12);
}

.selected-name {
  font-size: 0.8125rem;
  opacity: 0.7;
}

.results-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8125rem;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  overflow: hidden;
}

.results-table th {
  text-align: left;
  font-size: 0.6875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  opacity: 0.5;
  padding: 0.375rem 0.625rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.results-table td {
  padding: 0.375rem 0.625rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.results-table tbody tr {
  cursor: pointer;
  transition: background 0.15s ease;
}

.results-table tbody tr:hover {
  background: rgba(255, 255, 255, 0.06);
}

.results-table tbody tr.active {
  background: rgba(124, 92, 224, 0.15);
}

.results-table code {
  background: rgba(255, 255, 255, 0.08);
  padding: 0.0625rem 0.3rem;
  border-radius: 3px;
  font-size: 0.75rem;
}
</style>
