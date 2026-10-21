<script setup lang="ts">
import Chat from '~/components/Stream/Chat.vue'
import Player from '~/components/Stream/Player.vue'
import { reverseStreamChatDirection } from '~/composables/reverse-stream-chat-direction'
</script>

<template>
  <ClientOnly>
    <div
      class="h-full flex flex-col overflow-y-hidden"
      :class="{
        'xl:flex-row': !reverseStreamChatDirection,
        'xl:flex-row-reverse': reverseStreamChatDirection
      }"
    >
      <main class="shrink-0 overflow-y-auto w-full xl:h-full xl:min-w-0 xl:flex-1 xl:shrink">
        <Player class="w-full xl:max-h-full" />
      </main>

      <aside class="min-h-0 flex-1 w-full xl:h-full xl:w-96 xl:flex-none">
        <Chat />
      </aside>
    </div>
  </ClientOnly>
</template>
