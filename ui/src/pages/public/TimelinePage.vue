<template>
  <main class="mx-auto max-w-900px px-4 py-10 sm:px-6 lg:px-8">
    <header class="mb-8">
      <p class="m-0 text-13px font-700 text-brand-honey">Timeline</p>
      <h1 class="m-0 mt-2 text-32px font-900">经历时间线</h1>
      <p class="m-0 mt-3 text-15px leading-7 text-brand-muted">首版以静态事件展示页面节奏，后续接入经历和项目关系。</p>
    </header>
    <ol class="grid gap-4 border-l border-brand-line pl-5">
      <li class="bb-panel relative p-5">
        <span class="absolute -left-7 top-6 h-3 w-3 rounded-full bg-brand-honey" />
        <p class="m-0 text-13px font-700 text-brand-muted">2026 Q2</p>
        <h2 class="m-0 mt-1 text-20px font-900">v3 前后端联调启动</h2>
        <p class="m-0 mt-2 text-14px leading-6 text-brand-muted">以 gateway 契约为唯一 HTTP 入口，前端先完成工程化与 auth 链路。</p>
      </li>
      <li class="bb-panel relative p-5">
        <span class="absolute -left-7 top-6 h-3 w-3 rounded-full bg-brand-leaf" />
        <p class="m-0 text-13px font-700 text-brand-muted">2026 Q1</p>
        <h2 class="m-0 mt-1 text-20px font-900">内容模型收口</h2>
        <p class="m-0 mt-2 text-14px leading-6 text-brand-muted">文章、项目、经历、笔记和 AI 协作边界进入 v3 文档体系。</p>
      </li>
    </ol>
  </main>
</template>
