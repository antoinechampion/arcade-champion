import { describe, it, expect } from 'vitest'
import { readGamepad } from '..'

function fakeGamepad(overrides: {
  axes?: number[]
  buttons?: Partial<Record<number, { pressed: boolean }>>
} = {}): Gamepad {
  const buttons = Array.from({ length: 17 }, (_, i) => ({
    pressed: overrides.buttons?.[i]?.pressed ?? false,
    touched: false,
    value: 0,
  }))

  return {
    axes: overrides.axes ?? [0, 0, 0, 0],
    buttons,
    connected: true,
    id: 'test',
    index: 0,
    mapping: 'standard',
    timestamp: 0,
    hapticActuators: [],
    vibrationActuator: null,
  } as unknown as Gamepad
}

describe('readGamepad', () => {
  it('returns all false for neutral state', () => {
    const result = readGamepad(fakeGamepad())
    expect(result).toEqual({
      up: false, down: false, left: false, right: false, confirm: false,
    })
  })

  it('detects d-pad buttons', () => {
    expect(readGamepad(fakeGamepad({ buttons: { 12: { pressed: true } } })).up).toBe(true)
    expect(readGamepad(fakeGamepad({ buttons: { 13: { pressed: true } } })).down).toBe(true)
    expect(readGamepad(fakeGamepad({ buttons: { 14: { pressed: true } } })).left).toBe(true)
    expect(readGamepad(fakeGamepad({ buttons: { 15: { pressed: true } } })).right).toBe(true)
  })

  it('detects stick directions past deadzone', () => {
    expect(readGamepad(fakeGamepad({ axes: [0, -0.6] })).up).toBe(true)
    expect(readGamepad(fakeGamepad({ axes: [0, 0.6] })).down).toBe(true)
    expect(readGamepad(fakeGamepad({ axes: [-0.6, 0] })).left).toBe(true)
    expect(readGamepad(fakeGamepad({ axes: [0.6, 0] })).right).toBe(true)
  })

  it('ignores stick within deadzone', () => {
    const result = readGamepad(fakeGamepad({ axes: [0.4, -0.4] }))
    expect(result.up).toBe(false)
    expect(result.down).toBe(false)
    expect(result.left).toBe(false)
    expect(result.right).toBe(false)
  })

  it('detects face button presses as confirm', () => {
    expect(readGamepad(fakeGamepad({ buttons: { 0: { pressed: true } } })).confirm).toBe(true)
    expect(readGamepad(fakeGamepad({ buttons: { 3: { pressed: true } } })).confirm).toBe(true)
  })

  it('does not treat d-pad as confirm', () => {
    expect(readGamepad(fakeGamepad({ buttons: { 12: { pressed: true } } })).confirm).toBe(false)
    expect(readGamepad(fakeGamepad({ buttons: { 15: { pressed: true } } })).confirm).toBe(false)
  })
})
