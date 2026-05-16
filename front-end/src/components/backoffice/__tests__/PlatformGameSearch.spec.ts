import { describe, it, expect } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import PlatformGameSearch from '../PlatformGameSearch.vue'

const catalog = [
  { name: 'Street Fighter 6', platformId: '1364780' },
  { name: 'Street Fighter V', platformId: '310950' },
  { name: 'Tekken 8', platformId: '1778820' },
]

async function mockSearch(query: string) {
  const q = query.toLowerCase()
  return catalog.filter(
    (g) => g.name.toLowerCase().includes(q) || g.platformId.toLowerCase().includes(q),
  )
}

function mountSearch() {
  return mount(PlatformGameSearch, {
    props: { search: mockSearch, placeholder: 'Search…', modelValue: '' },
  })
}

async function searchFor(wrapper: ReturnType<typeof mountSearch>, query: string) {
  await wrapper.find('input').setValue(query)
  await wrapper.find('.search-btn').trigger('click')
  await flushPromises()
}

describe('PlatformGameSearch', () => {
  it('shows no results before searching', () => {
    const wrapper = mountSearch()
    expect(wrapper.find('.results-table').exists()).toBe(false)
  })

  it('shows results after clicking search', async () => {
    const wrapper = mountSearch()
    await searchFor(wrapper, 'street')

    const rows = wrapper.findAll('.results-table tbody tr')
    expect(rows).toHaveLength(2)
    expect(rows[0].text()).toContain('Street Fighter 6')
    expect(rows[0].text()).toContain('1364780')
  })

  it('emits selected platform ID on row click', async () => {
    const wrapper = mountSearch()
    await searchFor(wrapper, 'street')

    await wrapper.find('.results-table tbody tr').trigger('click')
    expect(wrapper.emitted('update:modelValue')![0]).toEqual(['1364780'])
  })

  it('shows game name and clears input after selection', async () => {
    const wrapper = mountSearch()
    await searchFor(wrapper, 'street')

    await wrapper.find('.results-table tbody tr').trigger('click')
    await flushPromises()

    expect(wrapper.find('.results-table').exists()).toBe(false)
    expect(wrapper.find('.selected-name').text()).toBe('Street Fighter 6')
    expect((wrapper.find('input').element as HTMLInputElement).value).toBe('')
  })

})
