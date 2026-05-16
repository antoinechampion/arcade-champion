import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import FeaturedGame from '../FeaturedGame.vue'

const baseProps = {
  title: 'Guilty Gear Strive',
  platform: 'Steam',
  releaseYear: 2021,
  developer: 'Arc System Works',
  bannerUrl: '/ggs-banner.jpg',
}

describe('FeaturedGame', () => {
  it('renders game details', () => {
    const wrapper = mount(FeaturedGame, { props: baseProps })
    expect(wrapper.text()).toContain('Guilty Gear Strive')
    expect(wrapper.text()).toContain('Steam')
    expect(wrapper.text()).toContain('2021')
    expect(wrapper.text()).toContain('Arc System Works')
  })

  it('renders the banner image', () => {
    const wrapper = mount(FeaturedGame, { props: baseProps })
    const img = wrapper.find('img')
    expect(img.attributes('src')).toBe('/ggs-banner.jpg')
    expect(img.attributes('alt')).toBe('Guilty Gear Strive')
  })

  it('renders the play button', () => {
    const wrapper = mount(FeaturedGame, { props: baseProps })
    expect(wrapper.text()).toContain('Play')
  })
})
