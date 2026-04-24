import {
  BookOpen,
  FolderKanban,
  Home,
  LayoutDashboard,
  Settings,
  Tags,
  UserRound,
} from 'lucide-vue-next';

export const publicNavItems = [
  { label: '首页', to: '/', icon: Home },
  { label: '文章', to: '/articles', icon: BookOpen },
  { label: '项目', to: '/projects', icon: FolderKanban },
  { label: '时间线', to: '/timeline', icon: Tags },
  { label: '关于', to: '/about', icon: UserRound },
] as const;

export const studioNavItems = [
  { label: '仪表盘', to: '/studio', icon: LayoutDashboard },
  { label: '内容中心', to: '/studio/content', icon: BookOpen },
  { label: '设置', to: '/studio/settings', icon: Settings },
] as const;
