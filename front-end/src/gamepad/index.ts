export let inputLocked = false
export function setInputLocked(value: boolean) { inputLocked = value }

const STICK_DEADZONE = 0.5
const REPEAT_DELAY = 400
const REPEAT_INTERVAL = 150

const DIRECTIONS = ['up', 'down', 'left', 'right', 'confirm'] as const
type Direction = (typeof DIRECTIONS)[number]

const KEY_FOR_DIRECTION: Record<Direction, string> = {
  up: 'ArrowUp',
  down: 'ArrowDown',
  left: 'ArrowLeft',
  right: 'ArrowRight',
  confirm: ' ',
}

export function readGamepad(gp: Gamepad): Record<Direction, boolean> {
  const axisX = gp.axes[0] ?? 0
  const axisY = gp.axes[1] ?? 0

  return {
    up: (gp.buttons[12]?.pressed ?? false) || axisY < -STICK_DEADZONE,
    down: (gp.buttons[13]?.pressed ?? false) || axisY > STICK_DEADZONE,
    left: (gp.buttons[14]?.pressed ?? false) || axisX < -STICK_DEADZONE,
    right: (gp.buttons[15]?.pressed ?? false) || axisX > STICK_DEADZONE,
    confirm: gp.buttons.some((b, i) => i < 12 && b.pressed),
  }
}

export function startGamepadPolling(): () => void {
  let prev: Record<Direction, boolean> = { up: false, down: false, left: false, right: false, confirm: false }
  const timers = new Map<Direction, { timeout: number; interval: number | null }>()
  let running = true

  function emit(key: string) {
    window.dispatchEvent(new KeyboardEvent('keydown', { key }))
  }

  function poll() {
    if (!running) return

    const gp = navigator.getGamepads()[0]
    if (gp && !inputLocked) {
      const state = readGamepad(gp)

      for (const dir of DIRECTIONS) {
        if (state[dir] && !prev[dir]) {
          const key = KEY_FOR_DIRECTION[dir]
          emit(key)

          if (dir !== 'confirm') {
            const timeout = window.setTimeout(() => {
              const t = timers.get(dir)
              if (!t) return
              t.interval = window.setInterval(() => emit(key), REPEAT_INTERVAL)
            }, REPEAT_DELAY)
            timers.set(dir, { timeout, interval: null })
          }
        } else if (!state[dir] && prev[dir]) {
          const t = timers.get(dir)
          if (t) {
            clearTimeout(t.timeout)
            if (t.interval !== null) clearInterval(t.interval)
            timers.delete(dir)
          }
        }
      }

      prev = state
    }

    requestAnimationFrame(poll)
  }

  requestAnimationFrame(poll)

  return () => {
    running = false
    for (const [, t] of timers) {
      clearTimeout(t.timeout)
      if (t.interval !== null) clearInterval(t.interval)
    }
    timers.clear()
  }
}
