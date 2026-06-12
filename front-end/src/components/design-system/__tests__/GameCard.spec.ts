import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import GameCard from '../GameCard.vue'

const baseProps = {
  title: 'Street Fighter III',
  platform: 'Fightcade',
  releaseYear: 1999,
  developer: 'Capcom',
  imageUrl: '/sf3.jpg',
}

describe('GameCard', () => {
  it('renders game information', () => {
    const wrapper = mount(GameCard, { props: baseProps })
    expect(wrapper.text()).toContain('Street Fighter III')
    expect(wrapper.text()).toContain('Fightcade')
    expect(wrapper.text()).toContain('1999')
    expect(wrapper.text()).toContain('Capcom')
  })

  it('renders the cover image', () => {
    const wrapper = mount(GameCard, { props: baseProps })
    const img = wrapper.find('img')
    expect(img.attributes('src')).toBe('/sf3.jpg')
    expect(img.attributes('alt')).toBe('Street Fighter III')
  })

  it('applies selected class when selected', () => {
    const wrapper = mount(GameCard, { props: { ...baseProps, selected: true } })
    expect(wrapper.find('.game-card').classes()).toContain('selected')
  })

  it('does not apply selected class by default', () => {
    const wrapper = mount(GameCard, { props: baseProps })
    expect(wrapper.find('.game-card').classes()).not.toContain('selected')
  })

  it('shows a titled fallback when the cover image fails to load', async () => {
    const wrapper = mount(GameCard, { props: baseProps })
    await wrapper.find('img').trigger('error')
    expect(wrapper.find('img').exists()).toBe(false)
    expect(wrapper.find('.card-fallback').text()).toContain('Street Fighter III')
  })
})
