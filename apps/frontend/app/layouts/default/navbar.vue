<script setup lang="ts">
import NavbarProfile from '~/layouts/default/navbar/navbar-profile.vue'
import NavbarThemeSwitcher from '~/layouts/default/navbar/navbar-theme-switcher.vue'
</script>

<template>
  <div
    class="flex h-14 w-full items-center justify-between border-b-2 border-border bg-background px-4 py-4"
  >
    <div class="flex items-center gap-4">
      <NuxtLink
        to="/"
        class="font-bold text-primary"
      >
        LiveStreamX
      </NuxtLink>
      <nav
        class="flex items-center space-x-6 text-sm font-medium max-lg:space-x-4"
      >
        <NuxtLink
          class="text-foreground/60 transition-colors hover:text-foreground/80"
          to="/"
        >
          Streams
        </NuxtLink>
      </nav>
    </div>

    <!-- there is not gap because there is weird bug when DialogOrSheet provided additional div to this flex, and it resulting content shift -->
    <div class="flex items-center">
      <NavbarThemeSwitcher class="mr-2" />
      <NavbarProfile />
    </div>
  </div>
</template>
