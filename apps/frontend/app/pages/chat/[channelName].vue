<script setup lang="ts">
import { computed } from 'vue'

import { useChat } from '~/api/chat.js'
import Chat from '~/components/Stream/Chat.vue'

definePageMeta({
  layout: false
})

const { channelName } = useChat()

useHead({
  title: computed(() =>
    channelName.value ? `${channelName.value} chat` : 'Chat'
  )
})
</script>

<template>
  <ClientOnly>
    <div class="h-dvh w-full overflow-hidden">
      <Chat popout />
    </div>
  </ClientOnly>
</template>
