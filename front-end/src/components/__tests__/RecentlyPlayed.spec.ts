import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import RecentlyPlayed from '../RecentlyPlayed.vue'

const games = [
  { id: '1', title: 'Game A', platform: 'Steam', releaseYear: 2020, developer: 'Dev A', imageUrl: '/a.jpg' },
  { id: '2', title: 'Game B', platform: 'MAME', releaseYear: 2021, developer: 'Dev B', imageUrl: '/b.jpg' },
  { id: '3', title: 'Game C', platform: 'Fightcade', releaseYear: 2022, developer: 'Dev C', imageUrl: '/c.jpg' },
]

describe('RecentlyPlayed', () => {
  it('renders a card for each game', () => {
    const wrapper = mount(RecentlyPlayed, { props: { games } })
    const cards = wrapper.findAll('.game-card')
    expect(cards).toHaveLength(3)
  })

  it('selects the first card by default', () => {
    const wrapper = mount(RecentlyPlayed, { props: { games } })
    const cards = wrapper.findAll('.game-card')
    expect(cards[0].classes()).toContain('selected')
    expect(cards[1].classes()).not.toContain('selected')
  })

  it('moves selection right on ArrowRight', async () => {
    const wrapper = mount(RecentlyPlayed, { props: { games } })
    await window.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight' }))
    await wrapper.vm.$nextTick()
    const cards = wrapper.findAll('.game-card')
    expect(cards[0].classes()).not.toContain('selected')
    expect(cards[1].classes()).toContain('selected')
  })

  it('moves selection left on ArrowLeft', async () => {
    const wrapper = mount(RecentlyPlayed, { props: { games } })
    await window.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight' }))
    await window.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowLeft' }))
    await wrapper.vm.$nextTick()
    const cards = wrapper.findAll('.game-card')
    expect(cards[0].classes()).toContain('selected')
  })

  it('does not go below zero', async () => {
    const wrapper = mount(RecentlyPlayed, { props: { games } })
    await window.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowLeft' }))
    await wrapper.vm.$nextTick()
    const cards = wrapper.findAll('.game-card')
    expect(cards[0].classes()).toContain('selected')
  })

  it('does not go past the last card', async () => {
    const wrapper = mount(RecentlyPlayed, { props: { games } })
    for (let i = 0; i < 5; i++) {
      await window.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight' }))
    }
    await wrapper.vm.$nextTick()
    const cards = wrapper.findAll('.game-card')
    expect(cards[2].classes()).toContain('selected')
  })
})
