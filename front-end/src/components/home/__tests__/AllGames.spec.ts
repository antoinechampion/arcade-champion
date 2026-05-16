import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import AllGames from '../AllGames.vue'

const mockGames = [
  { id: '1', title: 'Alpha Game', platform: 'Steam', releaseYear: 2020, developer: 'Dev A', imageUrl: '/a.jpg' },
  { id: '2', title: 'Alpha Two', platform: 'MAME', releaseYear: 2021, developer: 'Dev B', imageUrl: '/b.jpg' },
  { id: '3', title: 'Beta Game', platform: 'Fightcade', releaseYear: 2022, developer: 'Dev C', imageUrl: '/c.jpg' },
  { id: '4', title: 'Gamma Game', platform: 'Steam', releaseYear: 2023, developer: 'Dev D', imageUrl: '/d.jpg' },
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
  it('renders all games', async () => {
    const wrapper = await mountAllGames()
    expect(wrapper.findAll('.game-card')).toHaveLength(4)
  })

  it('groups games by first letter', async () => {
    const wrapper = await mountAllGames()
    const secs = wrapper.findAll('[data-section]')
    expect(secs.map((s) => s.attributes('data-section'))).toEqual(['A', 'B', 'G'])
  })

  it('no card is selected without navigation context', async () => {
    const wrapper = await mountAllGames()
    const selected = wrapper.findAll('.game-card').findIndex((c) => c.classes().includes('selected'))
    expect(selected).toBe(-1)
  })

  it('filters games when searching', async () => {
    const wrapper = await mountAllGames()
    const input = wrapper.find('.arcade-input')
    await input.setValue('Alpha')
    await input.trigger('keydown', { key: 'Enter' })
    await flushPromises()
    expect(wrapper.findAll('.game-card')).toHaveLength(2)
    expect(wrapper.findAll('[data-section]')).toHaveLength(1)
  })

  it('shows empty state when no results', async () => {
    const wrapper = await mountAllGames()
    const input = wrapper.find('.arcade-input')
    await input.setValue('zzzzz')
    await input.trigger('keydown', { key: 'Enter' })
    await flushPromises()
    expect(wrapper.findAll('.game-card')).toHaveLength(0)
    expect(wrapper.text()).toContain('No games found.')
  })
})
