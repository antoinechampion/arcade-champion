import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import AllGames from '../AllGames.vue'

const mockGames = [
  { id: '1', title: 'Alpha Game', platform: 'Steam', releaseYear: 2020, developer: 'Dev A', imageUrl: '/a.jpg' },
  { id: '2', title: 'Beta Game', platform: 'MAME', releaseYear: 2021, developer: 'Dev B', imageUrl: '/b.jpg' },
  { id: '3', title: 'Gamma Game', platform: 'Fightcade', releaseYear: 2022, developer: 'Dev C', imageUrl: '/c.jpg' },
  { id: '4', title: 'Delta Game', platform: 'Steam', releaseYear: 2023, developer: 'Dev D', imageUrl: '/d.jpg' },
]

vi.mock('@/api/client', () => ({
  fetchAllGames: vi.fn((query?: string) => {
    let games = [...mockGames]
    if (query) {
      const q = query.toLowerCase()
      games = games.filter((g) => g.title.toLowerCase().includes(q))
    }
    return Promise.resolve(games)
  }),
}))

async function mountAllGames() {
  const wrapper = mount(AllGames)
  await flushPromises()
  return wrapper
}

describe('AllGames', () => {
  it('renders all games sorted from the API', async () => {
    const wrapper = await mountAllGames()
    const cards = wrapper.findAll('.game-card')
    expect(cards).toHaveLength(4)
  })

  it('no card is selected when not focused', async () => {
    const wrapper = await mountAllGames()
    const cards = wrapper.findAll('.game-card')
    expect(cards[0].classes()).not.toContain('selected')
    expect(cards[1].classes()).not.toContain('selected')
  })

  it('selects the first card on focus', async () => {
    const wrapper = await mountAllGames()
    const grid = wrapper.find('.games-grid')
    await grid.trigger('focus')
    const cards = wrapper.findAll('.game-card')
    expect(cards[0].classes()).toContain('selected')
    expect(cards[1].classes()).not.toContain('selected')
  })

  it('moves selection right on ArrowRight', async () => {
    const wrapper = await mountAllGames()
    const grid = wrapper.find('.games-grid')
    await grid.trigger('focus')
    await grid.trigger('keydown', { key: 'ArrowRight' })
    const cards = wrapper.findAll('.game-card')
    expect(cards[0].classes()).not.toContain('selected')
    expect(cards[1].classes()).toContain('selected')
  })

  it('moves selection left on ArrowLeft', async () => {
    const wrapper = await mountAllGames()
    const grid = wrapper.find('.games-grid')
    await grid.trigger('focus')
    await grid.trigger('keydown', { key: 'ArrowRight' })
    await grid.trigger('keydown', { key: 'ArrowLeft' })
    const cards = wrapper.findAll('.game-card')
    expect(cards[0].classes()).toContain('selected')
  })

  it('does not go below zero', async () => {
    const wrapper = await mountAllGames()
    const grid = wrapper.find('.games-grid')
    await grid.trigger('focus')
    await grid.trigger('keydown', { key: 'ArrowLeft' })
    const cards = wrapper.findAll('.game-card')
    expect(cards[0].classes()).toContain('selected')
  })

  it('does not go past the last card', async () => {
    const wrapper = await mountAllGames()
    const grid = wrapper.find('.games-grid')
    await grid.trigger('focus')
    for (let i = 0; i < 10; i++) {
      await grid.trigger('keydown', { key: 'ArrowRight' })
    }
    const cards = wrapper.findAll('.game-card')
    expect(cards[3].classes()).toContain('selected')
  })

  it('filters games when searching', async () => {
    const wrapper = await mountAllGames()
    const input = wrapper.find('.search-input')
    await input.setValue('Alpha')
    await input.trigger('keydown', { key: 'Enter' })
    await flushPromises()
    const cards = wrapper.findAll('.game-card')
    expect(cards).toHaveLength(1)
    expect(wrapper.text()).toContain('Alpha Game')
  })

  it('shows empty state when no results', async () => {
    const wrapper = await mountAllGames()
    const input = wrapper.find('.search-input')
    await input.setValue('zzzzz')
    await input.trigger('keydown', { key: 'Enter' })
    await flushPromises()
    expect(wrapper.findAll('.game-card')).toHaveLength(0)
    expect(wrapper.text()).toContain('No games found.')
  })

  it('clears selection on blur', async () => {
    const wrapper = await mountAllGames()
    const grid = wrapper.find('.games-grid')
    await grid.trigger('focus')
    expect(wrapper.findAll('.game-card')[0].classes()).toContain('selected')
    await grid.trigger('blur')
    expect(wrapper.findAll('.game-card')[0].classes()).not.toContain('selected')
  })
})
