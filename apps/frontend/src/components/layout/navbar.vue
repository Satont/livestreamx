<script setup lang="ts">
import { useProfile } from '@/api/profile.ts'
import NavbarProfile from '@/components/navbar-profile.vue'
import ThemeSwitcher from '@/components/theme-switcher.vue'

const { data: profile } = useProfile().useData()
</script>

<template>
  <div
    class="h-14 bg-background w-full flex items-center px-4 py-4 justify-between border-b-2 border-border"
  >
    <div class="flex items-center gap-4">
      <h1 class="font-bold text-primary">LiveStreamX</h1>
      <nav
        class="flex items-center max-lg:space-x-4 space-x-6 text-sm font-medium"
      >
        <router-link
          class="transition-colors hover:text-foreground/80 text-foreground/60"
          to="/"
          >Streams</router-link
        >
        <router-link
          v-if="profile?.userProfile.user.isAdmin"
          class="transition-colors hover:text-foreground/80 text-foreground/60"
          to="/dashboard"
        >
          Dashboard
        </router-link>
      </nav>
    </div>

    <div class="flex items-center">
      <ThemeSwitcher />
      <NavbarProfile />
    </div>
  </div>
</template>
