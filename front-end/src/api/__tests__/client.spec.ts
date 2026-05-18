import { describe, it, expect, vi, beforeEach } from 'vitest'
import { fetchAllGames, fetchGame, createGame, deleteGame, searchPlatformGames, imageUrl } from '../client'
import type { GameInput } from '../types'

beforeEach(() => {
  vi.restoreAllMocks()
})

describe('fetchAllGames', () => {
  it('calls GET /api/games and returns games', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      json: async () => [
        { id: '1', title: 'SF3', platform: 'Fightcade', releaseYear: 1999, developer: 'Capcom', coverFilename: 'c.jpg', bannerFilename: 'b.jpg', appId: 'sfiii3n' },
      ],
    } as Response)

    const games = await fetchAllGames()
    expect(fetch).toHaveBeenCalledWith('/api/games')
    expect(games[0].appId).toBe('sfiii3n')
    expect(games[0].coverFilename).toBe('c.jpg')
  })

  it('passes query param when provided', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      json: async () => [],
    } as Response)

    await fetchAllGames('tekken')
    expect(fetch).toHaveBeenCalledWith('/api/games?query=tekken')
  })
})

describe('fetchGame', () => {
  it('returns undefined on 404', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: false,
      status: 404,
    } as Response)

    const game = await fetchGame('999')
    expect(game).toBeUndefined()
  })
})

describe('createGame', () => {
  it('sends multipart form data to POST /api/games', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      json: async () => ({ id: '11', title: 'New', platform: 'Steam', releaseYear: 2024, developer: 'Dev', coverFilename: 'x.jpg', bannerFilename: 'y.jpg', appId: '123' }),
    } as Response)

    const input: GameInput = {
      title: 'New',
      platform: 'Steam',
      releaseYear: 2024,
      developer: 'Dev',
      appId: '123',
      cover: new Blob(['cover'], { type: 'image/jpeg' }),
      banner: new Blob(['banner'], { type: 'image/jpeg' }),
    }

    const game = await createGame(input)
    expect(fetch).toHaveBeenCalledWith('/api/games', expect.objectContaining({ method: 'POST' }))
    expect(game.id).toBe('11')
  })
})

describe('deleteGame', () => {
  it('sends DELETE to /api/games/:id', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({ ok: true } as Response)
    await deleteGame('5')
    expect(fetch).toHaveBeenCalledWith('/api/games/5', { method: 'DELETE' })
  })
})

describe('searchPlatformGames', () => {
  it('calls the backend and maps the response', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      json: async () => [
        { game: 'Street Fighter 6', appId: '1364780' },
      ],
    } as Response)

    const results = await searchPlatformGames('Steam', 'fighter')
    expect(fetch).toHaveBeenCalledWith('/api/search?platform=steam&query=fighter')
    expect(results).toEqual([{ name: 'Street Fighter 6', platformId: '1364780' }])
  })

  it('throws on non-ok response', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({ ok: false, status: 404 } as Response)
    await expect(searchPlatformGames('Steam', 'xyz')).rejects.toThrow('Search failed: 404')
  })
})

describe('imageUrl', () => {
  it('returns /images/ prefixed path', () => {
    expect(imageUrl('1_cover_abc.jpg')).toBe('/images/1_cover_abc.jpg')
  })
})
