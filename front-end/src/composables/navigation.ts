import { ref, provide, inject, onMounted, onUnmounted, watchEffect, type InjectionKey, type Ref } from 'vue'

const navigationLocked = ref(false)
export function lockNavigation() { navigationLocked.value = true }
export function unlockNavigation() { navigationLocked.value = false }

export type NavDirection = 'up' | 'down' | 'left' | 'right'
export type NavCommand = NavDirection | 'confirm'

export interface NavigationHandler {
  onCommand(command: NavCommand): boolean
  onEnter?(from: 'up' | 'down'): void
}

interface NavigationManager {
  register(id: string, handler: NavigationHandler): void
  unregister(id: string): void
  activeZoneId: Ref<string | null>
  emitCommand(command: NavCommand): void
}

const NAV_KEY: InjectionKey<NavigationManager> = Symbol('navigation')

const KEY_MAP: Record<string, NavCommand> = {
  ArrowUp: 'up',
  ArrowDown: 'down',
  ArrowLeft: 'left',
  ArrowRight: 'right',
  ' ': 'confirm',
}

export function usePageNavigation(zoneOrder: string[]) {
  const handlers = new Map<string, NavigationHandler>()
  const activeZoneId = ref<string | null>(zoneOrder[0] ?? null)

  function findNextZone(from: number, direction: 1 | -1): string | null {
    for (let i = from + direction; i >= 0 && i < zoneOrder.length; i += direction) {
      if (handlers.has(zoneOrder[i])) return zoneOrder[i]
    }
    return null
  }

  function emitCommand(command: NavCommand) {
    const activeId = activeZoneId.value
    if (!activeId) return

    const handler = handlers.get(activeId)
    if (!handler) {
      if (command !== 'up' && command !== 'down') return
      const next = findNextZone(zoneOrder.indexOf(activeId), command === 'down' ? 1 : -1)
      if (next) {
        activeZoneId.value = next
        handlers.get(next)!.onEnter?.(command === 'down' ? 'up' : 'down')
      }
      return
    }

    const consumed = handler.onCommand(command)
    if (consumed) return

    if (command !== 'up' && command !== 'down') return

    const next = findNextZone(zoneOrder.indexOf(activeId), command === 'down' ? 1 : -1)
    if (next) {
      activeZoneId.value = next
      handlers.get(next)!.onEnter?.(command === 'down' ? 'up' : 'down')
    }
  }

  function onKeydown(e: KeyboardEvent) {
    if (navigationLocked.value) return
    const command = KEY_MAP[e.key]
    if (!command) return
    e.preventDefault()
    emitCommand(command)
  }

  onMounted(() => window.addEventListener('keydown', onKeydown))
  onUnmounted(() => window.removeEventListener('keydown', onKeydown))

  const manager: NavigationManager = {
    register(id, handler) {
      handlers.set(id, handler)
      if (!activeZoneId.value || !handlers.has(activeZoneId.value)) {
        activeZoneId.value = id
      } else {
        const activeIdx = zoneOrder.indexOf(activeZoneId.value)
        const newIdx = zoneOrder.indexOf(id)
        if (newIdx >= 0 && newIdx < activeIdx) {
          activeZoneId.value = id
        }
      }
    },
    unregister(id) { handlers.delete(id) },
    activeZoneId,
    emitCommand,
  }

  provide(NAV_KEY, manager)
}

export function useComponentNavigation(id: string, handler: NavigationHandler) {
  const manager = inject(NAV_KEY, null)
  const active = ref(false)

  if (manager) {
    onMounted(() => manager.register(id, handler))
    onUnmounted(() => manager.unregister(id))

    watchEffect(() => {
      active.value = manager.activeZoneId.value === id
    })
  }

  return { active }
}
