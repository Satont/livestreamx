<script setup lang="ts">
import { useProfile } from "@/api/profile.ts";
import { Button } from "@/components/ui/button";

const { data: profile, error } = useProfile();
const loginUri = `/api/auth/twitch?redirectUri=${window.location.origin}`
</script>

<template>
  <div class="flex items-center">
    <Button class="cursor-pointer" v-if="!profile || error" as="a" size="xs" :href="loginUri">Login</Button>
    <img v-else :src="profile.userProfile.avatarUrl" alt="avatar" class="size-8 rounded-full" />
  </div>
</template>
