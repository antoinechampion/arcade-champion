import { describe, it, expect, vi, afterEach } from 'vitest'
import { hideIdleCursor } from '../cursor'

const move = (x: number, y: number) =>
  document.dispatchEvent(new MouseEvent('mousemove', { clientX: x, clientY: y }))

describe('idle cursor', () => {
  afterEach(() => {
    vi.useRealTimers()
  })

  it('ignores jitter, shows on a real move, hides again when idle', () => {
    vi.useFakeTimers()
    const root = document.documentElement
    const stop = hideIdleCursor()
    expect(root.classList.contains('cursor-hidden')).toBe(true)

    move(100, 100)
    move(103, 102)
    expect(root.classList.contains('cursor-hidden')).toBe(true)

    move(120, 100)
    expect(root.classList.contains('cursor-hidden')).toBe(false)

    vi.advanceTimersByTime(10000)
    expect(root.classList.contains('cursor-hidden')).toBe(true)

    stop()
  })
})
