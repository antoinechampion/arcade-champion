import { describe, it, expect, vi, beforeEach } from 'vitest'
import { fetchRecentlyPlayed, searchPlatformGames } from '../client'

describe('fetchRecentlyPlayed', () => {
  it('returns a list of games', async () => {
    const games = await fetchRecentlyPlayed()
    expect(games.length).toBeGreaterThan(0)
  })

  it('returns games with required fields', async () => {
    const games = await fetchRecentlyPlayed()
    for (const game of games) {
      expect(game.id).toBeTruthy()
      expect(game.title).toBeTruthy()
      expect(game.platform).toBeTruthy()
      expect(game.releaseYear).toBeGreaterThan(0)
      expect(game.developer).toBeTruthy()
      expect(game.imageUrl).toBeTruthy()
    }
  })

  it('includes a banner URL on the first game', async () => {
    const games = await fetchRecentlyPlayed()
    expect(games[0].bannerUrl).toBeTruthy()
  })
})

describe('searchPlatformGames', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('calls the backend and maps the response', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      json: async () => [
        { game: 'Street Fighter 6', appId: '1364780' },
        { game: 'Tekken 8', appId: '1778820' },
      ],
    } as Response)

    const results = await searchPlatformGames('Steam', 'fighter')

    expect(fetch).toHaveBeenCalledWith('/api/search?platform=steam&query=fighter')
    expect(results).toEqual([
      { name: 'Street Fighter 6', platformId: '1364780' },
      { name: 'Tekken 8', platformId: '1778820' },
    ])
  })

  it('throws on non-ok response', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: false,
      status: 404,
    } as Response)

    await expect(searchPlatformGames('Steam', 'xyz')).rejects.toThrow('Search failed: 404')
  })
})
