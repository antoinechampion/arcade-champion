import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { handleGlobalKeydown, registerGlobalShortcuts } from '../keyboard'
import * as client from '../api/client'

describe('keyboard shortcuts', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('triggers exitApp and prevents default on Ctrl+Shift+Delete', () => {
    const exitSpy = vi.spyOn(client, 'exitApp').mockResolvedValue(undefined)
    const event = new KeyboardEvent('keydown', {
      key: 'Delete',
      ctrlKey: true,
      shiftKey: true,
      cancelable: true,
    })
    const preventDefaultSpy = vi.spyOn(event, 'preventDefault')

    handleGlobalKeydown(event)

    expect(preventDefaultSpy).toHaveBeenCalled()
    expect(exitSpy).toHaveBeenCalled()
  })

  it('triggers exitApp and prevents default on Ctrl+Shift+Del', () => {
    const exitSpy = vi.spyOn(client, 'exitApp').mockResolvedValue(undefined)
    const event = new KeyboardEvent('keydown', {
      key: 'Del',
      ctrlKey: true,
      shiftKey: true,
      cancelable: true,
    })
    const preventDefaultSpy = vi.spyOn(event, 'preventDefault')

    handleGlobalKeydown(event)

    expect(preventDefaultSpy).toHaveBeenCalled()
    expect(exitSpy).toHaveBeenCalled()
  })

  it('does not trigger exitApp when modifier keys are missing', () => {
    const exitSpy = vi.spyOn(client, 'exitApp').mockResolvedValue(undefined)
    const event = new KeyboardEvent('keydown', {
      key: 'Delete',
      ctrlKey: false,
      shiftKey: true,
    })

    handleGlobalKeydown(event)
    expect(exitSpy).not.toHaveBeenCalled()
  })

  it('prevents default on navigation arrow keys outside input fields', () => {
    const event = new KeyboardEvent('keydown', {
      key: 'ArrowDown',
      cancelable: true,
    })
    const preventDefaultSpy = vi.spyOn(event, 'preventDefault')

    handleGlobalKeydown(event)
    expect(preventDefaultSpy).toHaveBeenCalled()
  })

  it('registers and unregisters keydown listener', () => {
    const exitSpy = vi.spyOn(client, 'exitApp').mockResolvedValue(undefined)
    const unregister = registerGlobalShortcuts()

    document.dispatchEvent(new KeyboardEvent('keydown', {
      key: 'Delete',
      ctrlKey: true,
      shiftKey: true,
    }))
    expect(exitSpy).toHaveBeenCalledTimes(1)

    unregister()

    document.dispatchEvent(new KeyboardEvent('keydown', {
      key: 'Delete',
      ctrlKey: true,
      shiftKey: true,
    }))
    expect(exitSpy).toHaveBeenCalledTimes(1)
  })
})
