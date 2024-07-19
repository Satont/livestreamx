<script setup lang="ts">
  import Navbar from "~/layouts/default/navbar.vue";
</script>

<template>
  <div class="flex h-dvh w-full flex-col">
    <Navbar />
    <slot />
  </div>
</template>
