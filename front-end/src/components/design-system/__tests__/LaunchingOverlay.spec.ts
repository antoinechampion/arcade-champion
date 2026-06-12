import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import LaunchingOverlay from '../LaunchingOverlay.vue'

describe('LaunchingOverlay', () => {
  it('shows the default message', () => {
    const wrapper = mount(LaunchingOverlay)
    expect(wrapper.text()).toContain('Launching…')
  })

  it('shows a custom message', () => {
    const wrapper = mount(LaunchingOverlay, { props: { message: 'Matchmaking in progress…' } })
    expect(wrapper.text()).toContain('Matchmaking in progress…')
  })

  it('shows the cancel hint only when cancellable', () => {
    const plain = mount(LaunchingOverlay)
    expect(plain.text()).not.toContain('cancel')

    const cancellable = mount(LaunchingOverlay, { props: { cancellable: true } })
    expect(cancellable.text()).toContain('cancel')
  })

  it('emits cancel on Escape when cancellable', () => {
    const wrapper = mount(LaunchingOverlay, { props: { cancellable: true } })
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
    expect(wrapper.emitted('cancel')).toHaveLength(1)
    wrapper.unmount()
  })

  it('does not emit cancel when not cancellable', () => {
    const wrapper = mount(LaunchingOverlay)
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
    expect(wrapper.emitted('cancel')).toBeUndefined()
    wrapper.unmount()
  })
})
