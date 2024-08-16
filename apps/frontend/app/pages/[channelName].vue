<script setup lang="ts">
import Chat from '~/components/Stream/Chat.vue'
import Player from '~/components/Stream/Player.vue'
import { reverseStreamChatDirection } from '~/composables/reverse-stream-chat-direction'
</script>

<template>
  <div
    class="h-full flex flex-col overflow-y-hidden"
    :class="{
      'xl:flex-row': !reverseStreamChatDirection,
      'xl:flex-row-reverse': reverseStreamChatDirection
    }"
  >
    <main class="h-[40dvh] xl:h-full xl:flex-1 overflow-y-hidden w-full">
      <ClientOnly>
        <Player class="w-full h-full" />
      </ClientOnly>
    </main>

    <aside class="h-[60dvh] w-full xl:h-full xl:w-96">
      <ClientOnly>
        <Chat />
      </ClientOnly>
    </aside>
  </div>
</template>
