<script setup lang="ts">
const modes = [
  { icon: 'lucide:sun', title: 'Light', value: 'light' },
  { icon: 'lucide:moon', title: 'Dark', value: 'dark' },
  { icon: 'lucide:laptop', title: 'System', value: 'system' }
]

const colorMode = useColorMode()
const setTheme = (val: string) => {
  colorMode.preference = val
}

const currentIcon = computed(() => {
  return modes.find((m) => m.value === colorMode?.preference)?.icon
})

defineOptions({
  inheritAttrs: false
})
</script>

<template>
  <UiDropdownMenu>
    <UiDropdownMenuTrigger as-child>
      <UiButton
        class="h-9 w-9"
        v-bind="$attrs"
        variant="ghost"
        size="icon"
        title="Theme switcher"
      >
        <span class="sr-only">Theme switcher</span>
        <Icon
          :name="currentIcon || 'lucide:sun'"
          class="h-[18px] w-[18px]"
        />
      </UiButton>
    </UiDropdownMenuTrigger>
    <UiDropdownMenuContent
      align="end"
      :side-offset="5"
    >
      <UiDropdownMenuItem
        class="cursor-pointer"
        v-for="(m, i) in modes"
        :key="i"
        :icon="m.icon"
        :title="m.title"
        @click="setTheme(m.value)"
      />
    </UiDropdownMenuContent>
  </UiDropdownMenu>
</template>

<style scoped></style>
