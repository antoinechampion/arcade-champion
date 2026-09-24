import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import BackOfficePage from '../BackOfficePage.vue'
import * as client from '@/api/client'

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn() }),
  RouterLink: {
    template: '<a><slot /></a>',
  },
}))

describe('BackOfficePage', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
    vi.spyOn(client, 'fetchAllGames').mockResolvedValue([])
  })

  it('renders the Quit App button in the header', () => {
    const wrapper = mount(BackOfficePage)
    const quitBtn = wrapper.find('button.quit-btn')
    expect(quitBtn.exists()).toBe(true)
    expect(quitBtn.text()).toBe('Quit App')
  })

  it('calls exitApp when Quit App button is clicked', async () => {
    const exitSpy = vi.spyOn(client, 'exitApp').mockResolvedValue(undefined)
    const wrapper = mount(BackOfficePage)
    await flushPromises()

    const quitBtn = wrapper.find('button.quit-btn')
    await quitBtn.trigger('click')

    expect(exitSpy).toHaveBeenCalledTimes(1)
  })
})
