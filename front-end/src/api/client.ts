import type { Game, GameInput, PlatformSearchResult, Platform, Settings } from './types'

const games: Game[] = [
  {
    id: '1',
    title: 'Street Fighter III: 3rd Strike',
    platform: 'Fightcade',
    releaseYear: 1999,
    developer: 'Capcom',
    coverFilename: '1_cover_mock.jpg',
    bannerFilename: '1_banner_mock.jpg',
    launchConfig: { gameId: 'sfiii3n' },
  },
  {
    id: '2',
    title: 'Marvel VS Capcom 2',
    platform: 'MAME',
    releaseYear: 2000,
    developer: 'Capcom',
    coverFilename: '2_cover_mock.jpg',
    bannerFilename: '2_banner_mock.jpg',
    launchConfig: { driverName: 'mvsc2' },
  },
  {
    id: '3',
    title: 'King of Fighters 98',
    platform: 'Fightcade',
    releaseYear: 1998,
    developer: 'SNK',
    coverFilename: '3_cover_mock.jpg',
    bannerFilename: '3_banner_mock.jpg',
    launchConfig: { gameId: 'kof98' },
  },
  {
    id: '4',
    title: 'Street Fighter 6',
    platform: 'Steam',
    releaseYear: 2023,
    developer: 'Capcom',
    coverFilename: '4_cover_mock.jpg',
    bannerFilename: '4_banner_mock.jpg',
    launchConfig: { appId: '1364780' },
  },
  {
    id: '5',
    title: 'Street Fighter 2: Champion Edition',
    platform: 'Fightcade',
    releaseYear: 1992,
    developer: 'Capcom',
    coverFilename: '5_cover_mock.jpg',
    bannerFilename: '5_banner_mock.jpg',
    launchConfig: { gameId: 'sf2ce' },
  },
  {
    id: '6',
    title: 'Guilty Gear Strive',
    platform: 'Steam',
    releaseYear: 2021,
    developer: 'Arc System Works',
    coverFilename: '6_cover_mock.jpg',
    bannerFilename: '6_banner_mock.jpg',
    launchConfig: { appId: '1384160' },
  },
  {
    id: '7',
    title: 'Tekken 8',
    platform: 'Steam',
    releaseYear: 2024,
    developer: 'Bandai Namco',
    coverFilename: '7_cover_mock.jpg',
    bannerFilename: '7_banner_mock.jpg',
    launchConfig: { appId: '1778820' },
  },
  {
    id: '8',
    title: 'Garou: Mark of the Wolves',
    platform: 'Fightcade',
    releaseYear: 1999,
    developer: 'SNK',
    coverFilename: '8_cover_mock.jpg',
    bannerFilename: '8_banner_mock.jpg',
    launchConfig: { gameId: 'garou' },
  },
  {
    id: '9',
    title: 'Samurai Shodown V Special',
    platform: 'Fightcade',
    releaseYear: 2004,
    developer: 'SNK',
    coverFilename: '9_cover_mock.jpg',
    bannerFilename: '9_banner_mock.jpg',
    launchConfig: { gameId: 'samsh5sp' },
  },
  {
    id: '10',
    title: 'Metal Slug 3',
    platform: 'MAME',
    releaseYear: 2000,
    developer: 'SNK',
    coverFilename: '10_cover_mock.jpg',
    bannerFilename: '10_banner_mock.jpg',
    launchConfig: { driverName: 'mslug3' },
  },
]

let nextId = 11

export function imageUrl(filename: string): string {
  return `/images/${filename}`
}

export async function fetchRecentlyPlayed(): Promise<Game[]> {
  return games
}

export async function fetchAllGames(query?: string): Promise<Game[]> {
  let result = [...games].sort((a, b) => a.title.localeCompare(b.title))
  if (query) {
    const q = query.toLowerCase()
    result = result.filter((g) => g.title.toLowerCase().includes(q))
  }
  return result
}

export async function fetchGame(id: string): Promise<Game | undefined> {
  return games.find((g) => g.id === id)
}

export async function createGame(input: GameInput): Promise<Game> {
  const newGame: Game = {
    id: String(nextId++),
    title: input.title,
    platform: input.platform,
    releaseYear: input.releaseYear,
    developer: input.developer,
    coverFilename: '',
    bannerFilename: '',
    launchConfig: input.launchConfig,
  }
  games.push(newGame)
  return newGame
}

export async function updateGame(id: string, input: GameInput): Promise<Game> {
  const idx = games.findIndex((g) => g.id === id)
  const updated: Game = {
    id,
    title: input.title,
    platform: input.platform,
    releaseYear: input.releaseYear,
    developer: input.developer,
    coverFilename: games[idx].coverFilename,
    bannerFilename: games[idx].bannerFilename,
    launchConfig: input.launchConfig,
  }
  games[idx] = updated
  return updated
}

export async function deleteGame(id: string): Promise<void> {
  const idx = games.findIndex((g) => g.id === id)
  if (idx !== -1) games.splice(idx, 1)
}

export async function openKeyboard(): Promise<void> {
  // TODO: replace with actual backend call that spawns steam-osk
}

interface BackendSearchResult {
  game: string
  appId: string
}

export async function searchPlatformGames(platform: Platform, query: string): Promise<PlatformSearchResult[]> {
  const params = new URLSearchParams({ platform: platform.toLowerCase(), query })
  const res = await fetch(`/api/search?${params}`)
  if (!res.ok) throw new Error(`Search failed: ${res.status}`)
  const data: BackendSearchResult[] = await res.json()
  return data.map((r) => ({ name: r.game, platformId: r.appId }))
}

export async function fetchSettings(): Promise<Settings> {
  const res = await fetch('/api/settings')
  if (!res.ok) throw new Error(`Failed to fetch settings: ${res.status}`)
  return res.json()
}

export async function updateSettings(settings: Settings): Promise<Settings> {
  const res = await fetch('/api/settings', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(settings),
  })
  if (!res.ok) throw new Error(`Failed to update settings: ${res.status}`)
  return res.json()
}
