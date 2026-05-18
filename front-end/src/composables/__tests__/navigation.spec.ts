import { describe, it, expect } from 'vitest'
import { defineComponent, h, nextTick } from 'vue'
import { mount } from '@vue/test-utils'
import { usePageNavigation, useComponentNavigation, type NavCommand } from '../navigation'

function createTestApp(options: {
  zoneOrder: string[]
  zones: Record<string, {
    onCommand: (cmd: NavCommand) => boolean
    onEnter?: (from: 'up' | 'down') => void
  }>
}) {
  const calls: string[] = []

  const childComponents = Object.entries(options.zones).map(([id, handler]) =>
    defineComponent({
      name: `Zone_${id}`,
      setup() {
        const { active } = useComponentNavigation(id, {
          onCommand(cmd) {
            calls.push(`${id}:${cmd}`)
            return handler.onCommand(cmd)
          },
          onEnter(from) {
            calls.push(`${id}:enter:${from}`)
            handler.onEnter?.(from)
          },
        })
        return { active }
      },
      render() {
        return h('div', { 'data-zone': id, 'data-active': this.active })
      },
    })
  )

  const App = defineComponent({
    setup() {
      usePageNavigation(options.zoneOrder)
    },
    render() {
      return h('div', childComponents.map((C) => h(C)))
    },
  })

  const wrapper = mount(App, { attachTo: document.body })

  function activeZone() {
    const el = wrapper.findAll('[data-active="true"]')
    return el.length === 1 ? el[0].attributes('data-zone') : null
  }

  function press(key: string) {
    window.dispatchEvent(new KeyboardEvent('keydown', { key }))
    return nextTick()
  }

  return { wrapper, calls, activeZone, press }
}

describe('usePageNavigation + useComponentNavigation', () => {
  it('first zone is active by default', () => {
    const { activeZone } = createTestApp({
      zoneOrder: ['a', 'b'],
      zones: {
        a: { onCommand: () => false },
        b: { onCommand: () => false },
      },
    })
    expect(activeZone()).toBe('a')
  })

  it('moves to next zone on unconsumed down', async () => {
    const { activeZone, press } = createTestApp({
      zoneOrder: ['a', 'b'],
      zones: {
        a: { onCommand: () => false },
        b: { onCommand: () => false },
      },
    })
    await press('ArrowDown')
    expect(activeZone()).toBe('b')
  })

  it('moves to previous zone on unconsumed up', async () => {
    const { activeZone, press } = createTestApp({
      zoneOrder: ['a', 'b'],
      zones: {
        a: { onCommand: () => false },
        b: { onCommand: () => false },
      },
    })
    await press('ArrowDown')
    await press('ArrowUp')
    expect(activeZone()).toBe('a')
  })

  it('stays on current zone when command is consumed', async () => {
    const { activeZone, press } = createTestApp({
      zoneOrder: ['a', 'b'],
      zones: {
        a: { onCommand: () => true },
        b: { onCommand: () => false },
      },
    })
    await press('ArrowDown')
    expect(activeZone()).toBe('a')
  })

  it('does not move past first or last zone', async () => {
    const { activeZone, press } = createTestApp({
      zoneOrder: ['a', 'b'],
      zones: {
        a: { onCommand: () => false },
        b: { onCommand: () => false },
      },
    })
    await press('ArrowUp')
    expect(activeZone()).toBe('a')
    await press('ArrowDown')
    await press('ArrowDown')
    expect(activeZone()).toBe('b')
  })

  it('does not change zone on unconsumed left/right', async () => {
    const { activeZone, press } = createTestApp({
      zoneOrder: ['a', 'b'],
      zones: {
        a: { onCommand: () => false },
        b: { onCommand: () => false },
      },
    })
    await press('ArrowLeft')
    expect(activeZone()).toBe('a')
    await press('ArrowRight')
    expect(activeZone()).toBe('a')
  })

  it('calls onEnter with direction when switching zones', async () => {
    const { calls, press } = createTestApp({
      zoneOrder: ['a', 'b'],
      zones: {
        a: { onCommand: () => false },
        b: { onCommand: () => false },
      },
    })
    await press('ArrowDown')
    expect(calls).toContain('b:enter:up')
    await press('ArrowUp')
    expect(calls).toContain('a:enter:down')
  })

  it('handles confirm command', async () => {
    const { calls, press } = createTestApp({
      zoneOrder: ['a'],
      zones: {
        a: { onCommand: (cmd) => cmd === 'confirm' },
      },
    })
    await press(' ')
    expect(calls).toContain('a:confirm')
  })

  it('works with three zones', async () => {
    const { activeZone, press } = createTestApp({
      zoneOrder: ['a', 'b', 'c'],
      zones: {
        a: { onCommand: () => false },
        b: { onCommand: () => false },
        c: { onCommand: () => false },
      },
    })
    await press('ArrowDown')
    expect(activeZone()).toBe('b')
    await press('ArrowDown')
    expect(activeZone()).toBe('c')
    await press('ArrowUp')
    expect(activeZone()).toBe('b')
  })

  it('reclaims focus when an earlier zone registers late', async () => {
    const ZoneA = defineComponent({
      setup() {
        const { active } = useComponentNavigation('a', { onCommand: () => false })
        return { active }
      },
      render() {
        return h('div', { 'data-zone': 'a', 'data-active': this.active })
      },
    })

    const ZoneB = defineComponent({
      setup() {
        const { active } = useComponentNavigation('b', { onCommand: () => false })
        return { active }
      },
      render() {
        return h('div', { 'data-zone': 'b', 'data-active': this.active })
      },
    })

    const App = defineComponent({
      props: { showA: { type: Boolean, default: false } },
      setup() {
        usePageNavigation(['a', 'b'])
      },
      render() {
        return h('div', [
          this.showA ? h(ZoneA) : null,
          h(ZoneB),
        ])
      },
    })

    const wrapper = mount(App, { attachTo: document.body, props: { showA: false } })
    await nextTick()

    function activeZone() {
      const el = wrapper.findAll('[data-active="true"]')
      return el.length === 1 ? el[0].attributes('data-zone') : null
    }

    expect(activeZone()).toBe('b')

    await wrapper.setProps({ showA: true })
    await nextTick()
    expect(activeZone()).toBe('a')

    wrapper.unmount()
  })
})
