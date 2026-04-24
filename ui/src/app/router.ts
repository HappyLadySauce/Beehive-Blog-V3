import { createRouter, createWebHistory } from 'vue-router';

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/shared/layouts/PublicLayout.vue'),
      children: [
        { path: '', name: 'public-home', component: () => import('@/pages/public/HomePage.vue') },
        { path: 'articles', name: 'public-articles', component: () => import('@/pages/public/ArticleListPage.vue') },
        { path: 'articles/:slug', name: 'public-article-detail', component: () => import('@/pages/public/ArticleDetailPage.vue') },
        { path: 'projects', name: 'public-projects', component: () => import('@/pages/public/ProjectsPage.vue') },
        { path: 'timeline', name: 'public-timeline', component: () => import('@/pages/public/TimelinePage.vue') },
        { path: 'about', name: 'public-about', component: () => import('@/pages/public/AboutPage.vue') },
      ],
    },
    {
      path: '/',
      component: () => import('@/shared/layouts/AppShell.vue'),
      children: [
        { path: 'login', name: 'auth-login', component: () => import('@/pages/auth/LoginPage.vue') },
        { path: 'register', name: 'auth-register', component: () => import('@/pages/auth/RegisterPage.vue') },
      ],
    },
    {
      path: '/studio',
      component: () => import('@/shared/layouts/StudioLayout.vue'),
      children: [
        { path: '', name: 'studio-dashboard', component: () => import('@/pages/studio/StudioDashboardPage.vue') },
        { path: 'content', name: 'studio-content', component: () => import('@/pages/studio/StudioContentPage.vue') },
        { path: 'settings', name: 'studio-settings', component: () => import('@/pages/studio/StudioSettingsPage.vue') },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
  scrollBehavior() {
    return { top: 0 };
  },
});
