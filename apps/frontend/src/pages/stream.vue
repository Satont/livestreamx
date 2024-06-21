<script setup lang="ts">
import Chat from "@/components/chat.vue";
import Player from "@/components/player.vue";
</script>

<template>
  <div class="w-full min-h-screen flex flex-col h-screen">
    <div class="flex-1 flex flex-row overflow-y-hidden">
      <main class="flex-1 overflow-y-auto">
        <Player />
      </main>

      <aside class="w-96">
        <Chat />
      </aside>
    </div>
  </div>
</template>

<style scoped>

</style>