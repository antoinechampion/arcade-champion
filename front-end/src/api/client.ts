import type { Game, GameInput, PlatformSearchResult, Platform, Settings } from './types'

const BASE = import.meta.env.VITE_BACKEND_URL

function api(path: string, init?: RequestInit): Promise<Response> {
  const url = `${BASE}${path}`
  return init ? fetch(url, init) : fetch(url)
}

function buildFormData(input: GameInput): FormData {
  const fd = new FormData()
  fd.append('title', input.title)
  fd.append('platform', input.platform)
  fd.append('releaseYear', String(input.releaseYear))
  fd.append('developer', input.developer)
  fd.append('appId', input.appId)
  if (input.cover) fd.append('cover', input.cover, 'cover.jpg')
  if (input.banner) fd.append('banner', input.banner, 'banner.jpg')
  return fd
}

export function imageUrl(filename: string): string {
  return `${BASE}/images/${filename}`
}

export async function fetchRecentlyPlayed(): Promise<Game[]> {
  const res = await api('/api/games/recent')
  if (!res.ok) throw new Error(`Failed to fetch recently played: ${res.status}`)
  return res.json()
}

export async function fetchAllGames(query?: string): Promise<Game[]> {
  const params = query ? `?query=${encodeURIComponent(query)}` : ''
  const res = await api(`/api/games${params}`)
  if (!res.ok) throw new Error(`Failed to fetch games: ${res.status}`)
  return res.json()
}

export async function fetchGame(id: string): Promise<Game | undefined> {
  const res = await api(`/api/games/${id}`)
  if (res.status === 404) return undefined
  if (!res.ok) throw new Error(`Failed to fetch game: ${res.status}`)
  return res.json()
}

export async function createGame(input: GameInput): Promise<Game> {
  const res = await api('/api/games', {
    method: 'POST',
    body: buildFormData(input),
  })
  if (!res.ok) throw new Error(`Failed to create game: ${res.status}`)
  return res.json()
}

export async function updateGame(id: string, input: GameInput): Promise<Game> {
  const res = await api(`/api/games/${id}`, {
    method: 'PUT',
    body: buildFormData(input),
  })
  if (!res.ok) throw new Error(`Failed to update game: ${res.status}`)
  return res.json()
}

export async function deleteGame(id: string): Promise<void> {
  const res = await api(`/api/games/${id}`, { method: 'DELETE' })
  if (!res.ok) throw new Error(`Failed to delete game: ${res.status}`)
}

interface BackendSearchResult {
  game: string
  appId: string
}

export async function searchPlatformGames(platform: Platform, query: string): Promise<PlatformSearchResult[]> {
  const params = new URLSearchParams({ platform: platform, query })
  const res = await api(`/api/search?${params}`)
  if (!res.ok) throw new Error(`Search failed: ${res.status}`)
  const data: BackendSearchResult[] = await res.json()
  return data.map((r) => ({ name: r.game, platformId: r.appId }))
}

export async function launchGame(platform: Platform, appId: string, launchOptions?: Record<string, string>, signal?: AbortSignal): Promise<void> {
  const res = await api('/api/launch', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ platform, appId, launchOptions }),
    signal,
  })
  if (!res.ok) throw new Error(`Failed to launch game: ${res.status}`)
}

export async function fetchSettings(): Promise<Settings> {
  const res = await api('/api/settings')
  if (!res.ok) throw new Error(`Failed to fetch settings: ${res.status}`)
  return res.json()
}

export async function updateSettings(settings: Settings): Promise<Settings> {
  const res = await api('/api/settings', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(settings),
  })
  if (!res.ok) throw new Error(`Failed to update settings: ${res.status}`)
  return res.json()
}
