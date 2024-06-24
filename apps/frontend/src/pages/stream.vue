<script setup lang="ts">
import Chat from '@/components/chat.vue'
import Player from '@/components/player/player.vue'
</script>

<template>
  <div class="h-full flex flex-col xl:flex-row overflow-y-hidden">
    <main class="h-[30vh] xl:h-full xl:flex-1 overflow-y-auto w-full">
      <Player class="w-full h-full" />
    </main>

    <aside class="h-[70vh] w-full xl:h-full xl:w-96">
      <Chat />
    </aside>
  </div>
</template>

<style scoped></style>
