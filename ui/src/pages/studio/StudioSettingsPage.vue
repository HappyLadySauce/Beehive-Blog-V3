<script setup lang="ts">
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import FormField from '@/shared/components/FormField.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
</script>

<template>
  <section class="settings-page">
    <PageHeader
      eyebrow="Studio"
      title="Settings"
      description="Workspace configuration placeholder with accessible form controls."
    />
    <StatusAlert tone="info" title="Mock mode">
      Settings are local placeholders until the gateway exposes the live management contract.
    </StatusAlert>
    <form class="settings-page__form" aria-label="Studio settings">
      <FormField label="Site title" for-id="site-title">
        <BaseInput id="site-title" model-value="Beehive Blog" />
      </FormField>
      <FormField label="Default author" for-id="default-author">
        <BaseInput id="default-author" model-value="Admin Editor" />
      </FormField>
      <BaseButton type="button">Save settings</BaseButton>
    </form>
  </section>
</template>

<style scoped>
.settings-page {
  display: grid;
  gap: 24px;
}

.settings-page__form {
  display: grid;
  max-width: 560px;
  gap: 16px;
}
</style>
