import { requestJson } from '@/shared/api/httpClient';
import { GatewayHttpError } from '@/shared/api/httpClient';
import type { ApiMode } from '@/shared/config/env';
import { appConfig } from '@/shared/config/env';
import type {
  ContentListResponse,
  ContentPublicBySlugResponse,
  ContentStatus,
  ContentSummaryView,
  ContentType,
  PublicContentQuery,
  StudioContentListQuery,
} from '@/shared/api/types';

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

interface ContentPreviewApi {
  listPublicContent(query?: PublicContentQuery): Promise<ContentListResponse>;
  getPublicContentBySlug(slug: string): Promise<ContentPublicBySlugResponse>;
  listStudioContents(query?: StudioContentListQuery, accessToken?: string): Promise<ContentListResponse>;
}

function buildQueryString(params: Record<string, string | number | undefined>): string {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === '' || value === null) {
      continue;
    }
    search.set(key, String(value));
  }
  return search.toString() === '' ? '' : `?${search.toString()}`;
}

function normalizeContentTypeFilter(value?: string): ContentType | undefined {
  if (typeof value !== 'string' || value.length === 0) {
    return undefined;
  }
  return value;
}

function mockQueryFilter<T extends { type?: string; status?: string; visibility?: string; keyword?: string }>(
  items: ContentSummaryView[],
  query: T,
): ContentSummaryView[] {
  const keyword = query.keyword?.trim().toLowerCase();
  const type = normalizeContentTypeFilter(query.type);
  const status = normalizeContentTypeFilter(query.status) as ContentStatus | undefined;
  const visibility = normalizeContentTypeFilter(query.visibility) as string | undefined;

  return items.filter((item) => {
    const matchesType = type === undefined || item.type === type;
    const matchesStatus = status === undefined || item.status === status;
    const matchesVisibility = visibility === undefined || item.visibility === visibility;
    const matchesKeyword =
      keyword === undefined ||
      item.title.toLowerCase().includes(keyword) ||
      item.summary.toLowerCase().includes(keyword) ||
      item.slug.toLowerCase().includes(keyword);
    return matchesType && matchesStatus && matchesVisibility && matchesKeyword;
  });
}

function paginateItems<T>(items: T[], page: number, pageSize: number): T[] {
  const start = Math.max(1, page) - 1;
  return items.slice(start * pageSize, start * pageSize + pageSize);
}

function buildMockContentList(query: PublicContentQuery = {}): { items: ContentSummaryView[]; total: number; page: number; page_size: number } {
  const page = query.page ?? 1;
  const pageSize = query.page_size ?? 20;
  const filtered = mockQueryFilter(mockItems, query);
  return {
    items: paginateItems(filtered, page, pageSize),
    total: filtered.length,
    page,
    page_size: pageSize,
  };
}

function buildMockContentBySlug(slug: string): ContentPublicBySlugResponse {
  const found = mockItems.find((item) => item.slug === slug);
  if (!found) {
    throw new GatewayHttpError(404, 'content not found', {
      code: 120501,
      message: 'content not found',
      reference: 'mock',
      request_id: 'mock-request',
    });
  }

  return {
    content: {
      ...found,
      body_markdown: `# ${found.title}\n\n${found.summary}`,
      body_json: JSON.stringify({ type: 'doc', children: [{ type: 'paragraph', content: found.summary }] }),
      owner_user_id: 'user_mock_001',
      author_user_id: 'user_mock_001',
      source_type: 'manual',
      current_revision_id: 'rev_mock_001',
      comment_enabled: true,
      is_featured: false,
      sort_order: 0,
    },
  };
}

function createMockContentPreviewApi(): ContentPreviewApi {
  return {
    async listPublicContent(query) {
      return buildMockContentList(query);
    },
    async getPublicContentBySlug(slug) {
      return buildMockContentBySlug(slug);
    },
    async listStudioContents(query) {
      return buildMockContentList(query);
    },
  };
}

function createLiveContentPreviewApi(): ContentPreviewApi {
  return {
    listPublicContent(query = {}) {
      return requestJson<ContentListResponse>(
        `/api/v3/public/content/items${buildQueryString({
          page: query.page ?? 1,
          page_size: query.page_size ?? 20,
          type: query.type,
          keyword: query.keyword,
        })}`,
        {
          method: 'GET',
        },
      );
    },
    getPublicContentBySlug(slug) {
      return requestJson<ContentPublicBySlugResponse>(`/api/v3/public/content/items/${encodeURIComponent(slug)}`, {
        method: 'GET',
      });
    },
    listStudioContents(query = {}, accessToken) {
      const requestOptions: { method: 'GET'; accessToken?: string } = {
        method: 'GET',
      };
      if (accessToken) {
        requestOptions.accessToken = accessToken;
      }
      return requestJson<ContentListResponse>(
        `/api/v3/studio/content/items${buildQueryString({
          page: query.page ?? 1,
          page_size: query.page_size ?? 20,
          type: query.type,
          status: query.status,
          visibility: query.visibility,
          keyword: query.keyword,
        })}`,
        requestOptions,
      );
    },
  };
}

export function createContentPreviewApi(mode: ApiMode = appConfig.apiMode): ContentPreviewApi {
  return mode === 'live' ? createLiveContentPreviewApi() : createMockContentPreviewApi();
}

export const contentPreviewApi = createContentPreviewApi();
