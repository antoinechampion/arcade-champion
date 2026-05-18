import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import RecentlyPlayed from '../RecentlyPlayed.vue'

vi.mock('@/api/client', () => ({
  imageUrl: (filename: string) => `/images/${filename}`,
}))

const games = [
  { id: '1', title: 'Game A', platform: 'Steam', releaseYear: 2020, developer: 'Dev A', coverFilename: 'a.jpg', bannerFilename: 'a_b.jpg' },
  { id: '2', title: 'Game B', platform: 'MAME', releaseYear: 2021, developer: 'Dev B', coverFilename: 'b.jpg', bannerFilename: 'b_b.jpg' },
  { id: '3', title: 'Game C', platform: 'Fightcade', releaseYear: 2022, developer: 'Dev C', coverFilename: 'c.jpg', bannerFilename: 'c_b.jpg' },
]

describe('RecentlyPlayed', () => {
  it('renders a card for each game', () => {
    const wrapper = mount(RecentlyPlayed, { props: { games } })
    const cards = wrapper.findAll('.game-card')
    expect(cards).toHaveLength(3)
  })

  it('no card is selected without navigation context', () => {
    const wrapper = mount(RecentlyPlayed, { props: { games } })
    const cards = wrapper.findAll('.game-card')
    expect(cards[0].classes()).not.toContain('selected')
    expect(cards[1].classes()).not.toContain('selected')
  })
})
