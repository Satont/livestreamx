<script setup lang="ts">
import { ConfigProvider } from 'radix-vue'

const useIdFunction = () => useId()
</script>

<template>
  <ConfigProvider :use-id="useIdFunction">
    <NuxtLoadingIndicator />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <UiVueSonner />
    <UiToastToaster />
  </ConfigProvider>
</template>
