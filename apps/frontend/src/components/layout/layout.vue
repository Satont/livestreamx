<script setup lang="ts">
import Navbar from '@/components/layout/navbar.vue'
</script>

<template>
  <div class="w-full h-dvh flex flex-col">
    <Navbar />
    <RouterView />
  </div>
</template>
