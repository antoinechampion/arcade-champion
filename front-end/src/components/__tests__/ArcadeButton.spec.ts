import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ArcadeButton from '../ArcadeButton.vue'

describe('ArcadeButton', () => {
  it('renders a label', () => {
    const wrapper = mount(ArcadeButton, { props: { label: 'Play' } })
    expect(wrapper.text()).toBe('Play')
  })

  it('renders an icon slot', () => {
    const wrapper = mount(ArcadeButton, {
      slots: { icon: '<svg data-testid="icon"></svg>' },
    })
    expect(wrapper.find('[data-testid="icon"]').exists()).toBe(true)
  })

  it('renders both icon and label together', () => {
    const wrapper = mount(ArcadeButton, {
      props: { label: 'Settings' },
      slots: { icon: '<svg data-testid="icon"></svg>' },
    })
    expect(wrapper.text()).toBe('Settings')
    expect(wrapper.find('[data-testid="icon"]').exists()).toBe(true)
  })

  it('hides label span when no label provided', () => {
    const wrapper = mount(ArcadeButton, {
      slots: { icon: '<svg></svg>' },
    })
    expect(wrapper.findAll('span').length).toBe(1)
  })
})
