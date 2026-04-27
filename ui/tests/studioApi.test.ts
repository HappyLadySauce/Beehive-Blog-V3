import { describe, expect, it } from 'vitest'

import { createMockStudioApi } from '@/features/studio/api/mockStudioApi'

describe('studio api facade', () => {
  it('soft deletes users and hides them from default user lists', async () => {
    const api = createMockStudioApi()
    const before = await api.listUsers()
    const target = before.items.find((user) => user.user_id !== 'user_mock_admin')

    expect(target).toBeDefined()

    await api.deleteUser(target!.user_id)
    const after = await api.listUsers()
    const withDeleted = await api.listUsers({ include_deleted: true })

    expect(after.items.some((user) => user.user_id === target!.user_id)).toBe(false)
    expect(withDeleted.items.find((user) => user.user_id === target!.user_id)?.status).toBe('deleted')
  })

  it('creates and lists content drafts through the studio api', async () => {
    const api = createMockStudioApi()
    const created = await api.createContent({
      type: 'article',
      title: 'Integration draft',
      slug: 'integration-draft',
      body_markdown: 'Integration draft',
      body_json: '{"type":"doc","content":[{"type":"paragraph"}]}',
      visibility: 'private',
      ai_access: 'denied',
      tag_ids: [],
    })

    const list = await api.listContents({ keyword: 'integration' })

    expect(created.content.status).toBe('draft')
    expect(list.items.some((content) => content.content_id === created.content.content_id)).toBe(true)
  })
})
