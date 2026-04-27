<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <main class="auth-shell">
    <section class="auth-shell__intro" aria-labelledby="auth-title">
      <RouterLink class="auth-shell__brand" to="/">Beehive Blog</RouterLink>
      <h1 id="auth-title">Identity gateway access</h1>
      <p>
        Sign in through the v3 gateway. The UI keeps access tokens in memory and stores refresh
        tokens only through the storage adapter.
      </p>
    </section>
    <section class="auth-shell__panel">
      <RouterView />
    </section>
  </main>
</template>

<style scoped>
.auth-shell {
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(320px, 440px);
  gap: 32px;
  align-items: center;
  width: min(1120px, calc(100% - 32px));
  margin: 0 auto;
  padding: 40px 0;
}

.auth-shell__brand {
  display: inline-flex;
  margin-bottom: 32px;
  color: var(--bb-color-primary);
  font-weight: 800;
  text-decoration: none;
}

.auth-shell__intro h1 {
  max-width: 620px;
  margin: 0;
  font-size: clamp(2rem, 7vw, 4.2rem);
  line-height: 1.02;
}

.auth-shell__intro p {
  max-width: 620px;
  color: var(--bb-color-muted);
}

.auth-shell__panel {
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 24px;
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-panel);
}

@media (max-width: 820px) {
  .auth-shell {
    grid-template-columns: 1fr;
  }
}
</style>
