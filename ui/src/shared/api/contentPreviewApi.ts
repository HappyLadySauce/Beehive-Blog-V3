import type { ContentListResponse, ContentSummaryView, PublicContentQuery } from '@/shared/api/types';

const now = Math.floor(Date.now() / 1000);

const mockItems: ContentSummaryView[] = [
  {
    content_id: 'mock_article_001',
    type: 'article',
    title: '把个人知识系统整理成可演进的平台',
    slug: 'personal-knowledge-platform',
    summary: '从公开表达、Studio 工作台和 AI 协作三个角度，梳理 Beehive Blog v3 的产品结构。',
    cover_image_url: '',
    status: 'published',
    visibility: 'public',
    ai_access: 'allowed',
    published_at: now - 3600 * 20,
    archived_at: 0,
    created_at: now - 3600 * 36,
    updated_at: now - 3600 * 3,
    tags: [
      {
        tag_id: 'tag_product',
        name: 'Product',
        slug: 'product',
        description: '产品设计',
        color: '#2a9d99',
        created_at: now,
        updated_at: now,
      },
    ],
  },
  {
    content_id: 'mock_project_001',
    type: 'project',
    title: 'Gateway-first 的前后端联调路径',
    slug: 'gateway-first-integration',
    summary: '前端只面向 gateway HTTP 契约，后端领域服务继续保持 RPC 边界。',
    cover_image_url: '',
    status: 'published',
    visibility: 'public',
    ai_access: 'allowed',
    published_at: now - 3600 * 64,
    archived_at: 0,
    created_at: now - 3600 * 96,
    updated_at: now - 3600 * 12,
    tags: [
      {
        tag_id: 'tag_arch',
        name: 'Architecture',
        slug: 'architecture',
        description: '架构设计',
        color: '#5e6ad2',
        created_at: now,
        updated_at: now,
      },
    ],
  },
  {
    content_id: 'mock_note_001',
    type: 'note',
    title: '从笔记到文章的演进链路',
    slug: 'note-to-article-flow',
    summary: '笔记、草稿、审阅、发布之间保持清晰状态，避免 AI 输出直接越权发布。',
    cover_image_url: '',
    status: 'published',
    visibility: 'public',
    ai_access: 'denied',
    published_at: now - 3600 * 120,
    archived_at: 0,
    created_at: now - 3600 * 160,
    updated_at: now - 3600 * 40,
    tags: [
      {
        tag_id: 'tag_ai',
        name: 'AI',
        slug: 'ai',
        description: 'AI 协作',
        color: '#18a058',
        created_at: now,
        updated_at: now,
      },
    ],
  },
];

export const contentPreviewApi = {
  async listPublicContent(query: PublicContentQuery = {}): Promise<ContentListResponse> {
    const page = query.page ?? 1;
    const pageSize = query.page_size ?? 20;
    const keyword = query.keyword?.trim().toLowerCase();
    const type = query.type;
    const filtered = mockItems.filter((item) => {
      const matchesType = type === undefined || item.type === type;
      const matchesKeyword =
        keyword === undefined ||
        item.title.toLowerCase().includes(keyword) ||
        item.summary.toLowerCase().includes(keyword);
      return matchesType && matchesKeyword;
    });

    return {
      items: filtered.slice((page - 1) * pageSize, page * pageSize),
      total: filtered.length,
      page,
      page_size: pageSize,
    };
  },
};
