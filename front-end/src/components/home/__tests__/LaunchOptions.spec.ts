import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import LaunchOptions from '../LaunchOptions.vue'

describe('LaunchOptions', () => {
  it('renders all mode options with first option focused', () => {
    const wrapper = mount(LaunchOptions)
    expect(wrapper.text()).toContain('Online')
    expect(wrapper.text()).toContain('Training')
    expect(wrapper.text()).toContain('Arcade')

    const buttons = wrapper.findAllComponents({ name: 'ArcadeButton' })
    expect(buttons).toHaveLength(3)
    expect(buttons[0].props('focused')).toBe(true)
    expect(buttons[1].props('focused')).toBe(false)
    expect(buttons[2].props('focused')).toBe(false)
    wrapper.unmount()
  })

  it('navigates options with ArrowRight and ArrowLeft', async () => {
    const wrapper = mount(LaunchOptions)
    const buttons = wrapper.findAllComponents({ name: 'ArcadeButton' })

    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight' }))
    await wrapper.vm.$nextTick()
    expect(buttons[0].props('focused')).toBe(false)
    expect(buttons[1].props('focused')).toBe(true)

    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight' }))
    await wrapper.vm.$nextTick()
    expect(buttons[2].props('focused')).toBe(true)

    // Clamps at end
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight' }))
    await wrapper.vm.$nextTick()
    expect(buttons[2].props('focused')).toBe(true)

    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowLeft' }))
    await wrapper.vm.$nextTick()
    expect(buttons[1].props('focused')).toBe(true)

    wrapper.unmount()
  })

  it('emits select with selected mode on space confirm', async () => {
    const wrapper = mount(LaunchOptions)
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight' }))
    await wrapper.vm.$nextTick()

    window.dispatchEvent(new KeyboardEvent('keydown', { key: ' ' }))
    expect(wrapper.emitted('select')).toEqual([['training']])
    wrapper.unmount()
  })

  it('emits cancel on Escape', () => {
    const wrapper = mount(LaunchOptions)
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
    expect(wrapper.emitted('cancel')).toHaveLength(1)
    wrapper.unmount()
  })

  it('swallows vertical navigation keys (ArrowUp, ArrowDown) without bubbling or cancelling', () => {
    const wrapper = mount(LaunchOptions)
    const upEvent = new KeyboardEvent('keydown', { key: 'ArrowUp', cancelable: true })
    const downEvent = new KeyboardEvent('keydown', { key: 'ArrowDown', cancelable: true })

    window.dispatchEvent(upEvent)
    window.dispatchEvent(downEvent)

    expect(wrapper.emitted('select')).toBeUndefined()
    expect(wrapper.emitted('cancel')).toBeUndefined()
    expect(upEvent.defaultPrevented).toBe(true)
    expect(downEvent.defaultPrevented).toBe(true)
    wrapper.unmount()
  })

  it('clicking an option emits select', async () => {
    const wrapper = mount(LaunchOptions)
    const buttons = wrapper.findAllComponents({ name: 'ArcadeButton' })
    await buttons[2].trigger('click')
    expect(wrapper.emitted('select')).toEqual([['arcade']])
    wrapper.unmount()
  })
})
