import { exitApp } from '@/api/client'

export function handleGlobalKeydown(e: KeyboardEvent): void {
  if (e.ctrlKey && e.shiftKey && (e.key === 'Delete' || e.key === 'Del')) {
    e.preventDefault()
    exitApp()
    return
  }
  if (['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight', ' '].includes(e.key)) {
    const tag = (e.target as HTMLElement | null)?.tagName
    if (tag !== 'INPUT' && tag !== 'TEXTAREA') e.preventDefault()
  }
}

export function registerGlobalShortcuts(): () => void {
  document.addEventListener('keydown', handleGlobalKeydown)
  return () => {
    document.removeEventListener('keydown', handleGlobalKeydown)
  }
}
