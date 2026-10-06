const IDLE_MS = 10000
const MIN_MOVE_PX = 10

export function hideIdleCursor(): () => void {
  const root = document.documentElement
  let timer: ReturnType<typeof setTimeout>
  // Position where the hidden cursor was first seen; jitter around it is ignored.
  let anchor: { x: number; y: number } | null = null

  const hide = () => {
    root.classList.add('cursor-hidden')
    anchor = null
  }

  const onMouseMove = (e: MouseEvent) => {
    if (root.classList.contains('cursor-hidden')) {
      if (!anchor) {
        anchor = { x: e.clientX, y: e.clientY }
        return
      }
      if (Math.hypot(e.clientX - anchor.x, e.clientY - anchor.y) < MIN_MOVE_PX) return
      root.classList.remove('cursor-hidden')
    }
    clearTimeout(timer)
    timer = setTimeout(hide, IDLE_MS)
  }

  hide()
  document.addEventListener('mousemove', onMouseMove)
  return () => {
    clearTimeout(timer)
    document.removeEventListener('mousemove', onMouseMove)
    root.classList.remove('cursor-hidden')
  }
}
