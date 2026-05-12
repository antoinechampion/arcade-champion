import { describe, it, expect } from 'vitest'
import { fetchRecentlyPlayed } from '../client'

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
