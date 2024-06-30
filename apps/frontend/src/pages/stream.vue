<script setup lang="ts">
import Chat from '@/components/chat.vue'
import Player from '@/components/player/player.vue'
import { reverseStreamChatDirection } from '@/composables/reverse-stream-chat-direction.ts'
</script>

<template>
  <div
    class="h-full flex flex-col overflow-y-hidden"
    :class="{
      'xl:flex-row': !reverseStreamChatDirection,
      'xl:flex-row-reverse': reverseStreamChatDirection
    }"
  >
    <main class="h-[40dvh] xl:h-full xl:flex-1 overflow-y-auto w-full">
      <Player class="w-full h-full" />
    </main>

    <aside class="h-[60dvh] w-full xl:h-full xl:w-96">
      <Chat />
    </aside>
  </div>
</template>

<style scoped></style>
