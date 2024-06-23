<script setup lang="ts">
import Chat from '@/components/chat.vue'
import Player from '@/components/player.vue'
</script>

<template>
  <div class="w-full min-h-screen flex flex-col h-screen">
    <div class="h-full flex flex-col xl:flex-row overflow-y-hidden">
      <main class="h-[30vh] xl:h-full xl:flex-1 overflow-y-auto w-full">
        <Player />
      </main>

      <aside class="h-[70vh] w-full xl:h-full xl:w-96">
        <Chat />
      </aside>
    </div>
  </div>
</template>

<style scoped></style>
