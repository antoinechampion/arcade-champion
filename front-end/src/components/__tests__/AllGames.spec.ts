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

function getContainer(wrapper: ReturnType<typeof mount>) {
  return wrapper.find('[tabindex]')
}

function selectedCard(wrapper: ReturnType<typeof mount>) {
  return wrapper.findAll('.game-card').findIndex((c) => c.classes().includes('selected'))
}

async function press(container: ReturnType<typeof getContainer>, key: string, times = 1) {
  for (let i = 0; i < times; i++) {
    await container.trigger('keydown', { key })
  }
}

// Sections: A[Alpha Game, Alpha Two], B[Beta Game], G[Gamma Game]
// Flat DOM order: 0=Alpha Game, 1=Alpha Two, 2=Beta Game, 3=Gamma Game

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

  it('no card is selected when not focused', async () => {
    const wrapper = await mountAllGames()
    expect(selectedCard(wrapper)).toBe(-1)
  })

  it('selects the first card on focus', async () => {
    const wrapper = await mountAllGames()
    await getContainer(wrapper).trigger('focus')
    expect(selectedCard(wrapper)).toBe(0)
  })

  it('ArrowRight within a section', async () => {
    const wrapper = await mountAllGames()
    const c = getContainer(wrapper)
    await c.trigger('focus')
    await press(c, 'ArrowRight')
    expect(selectedCard(wrapper)).toBe(1)
  })

  it('ArrowRight crosses into next section', async () => {
    const wrapper = await mountAllGames()
    const c = getContainer(wrapper)
    await c.trigger('focus')
    await press(c, 'ArrowRight', 2)
    expect(selectedCard(wrapper)).toBe(2)
  })

  it('ArrowLeft crosses into previous section', async () => {
    const wrapper = await mountAllGames()
    const c = getContainer(wrapper)
    await c.trigger('focus')
    await press(c, 'ArrowRight', 2)
    await press(c, 'ArrowLeft')
    expect(selectedCard(wrapper)).toBe(1)
  })

  it('ArrowLeft does not go before first card', async () => {
    const wrapper = await mountAllGames()
    const c = getContainer(wrapper)
    await c.trigger('focus')
    await press(c, 'ArrowLeft')
    expect(selectedCard(wrapper)).toBe(0)
  })

  it('ArrowRight does not go past last card', async () => {
    const wrapper = await mountAllGames()
    const c = getContainer(wrapper)
    await c.trigger('focus')
    await press(c, 'ArrowRight', 10)
    expect(selectedCard(wrapper)).toBe(3)
  })

  it('ArrowDown moves to next section same column', async () => {
    const wrapper = await mountAllGames()
    const c = getContainer(wrapper)
    await c.trigger('focus')
    // A[0] → B[0]
    await press(c, 'ArrowDown')
    expect(selectedCard(wrapper)).toBe(2)
  })

  it('ArrowDown clamps column if next section is shorter', async () => {
    const wrapper = await mountAllGames()
    const c = getContainer(wrapper)
    await c.trigger('focus')
    await press(c, 'ArrowRight')
    // A[1] → B only has 1 game, clamp to B[0]
    await press(c, 'ArrowDown')
    expect(selectedCard(wrapper)).toBe(2)
  })

  it('ArrowUp moves to previous section last row', async () => {
    const wrapper = await mountAllGames()
    const c = getContainer(wrapper)
    await c.trigger('focus')
    await press(c, 'ArrowRight', 2)
    // B[0] → A last row col 0 = A[0]
    await press(c, 'ArrowUp')
    expect(selectedCard(wrapper)).toBe(0)
  })

  it('ArrowUp from first section stays put', async () => {
    const wrapper = await mountAllGames()
    const c = getContainer(wrapper)
    await c.trigger('focus')
    await press(c, 'ArrowUp')
    expect(selectedCard(wrapper)).toBe(0)
  })

  it('ArrowDown from last section stays put', async () => {
    const wrapper = await mountAllGames()
    const c = getContainer(wrapper)
    await c.trigger('focus')
    await press(c, 'ArrowRight', 3)
    await press(c, 'ArrowDown')
    expect(selectedCard(wrapper)).toBe(3)
  })

  it('filters games when searching', async () => {
    const wrapper = await mountAllGames()
    const input = wrapper.find('.search-input')
    await input.setValue('Alpha')
    await input.trigger('keydown', { key: 'Enter' })
    await flushPromises()
    expect(wrapper.findAll('.game-card')).toHaveLength(2)
    expect(wrapper.findAll('[data-section]')).toHaveLength(1)
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
    const c = getContainer(wrapper)
    await c.trigger('focus')
    expect(selectedCard(wrapper)).toBe(0)
    await c.trigger('blur')
    expect(selectedCard(wrapper)).toBe(-1)
  })
})
