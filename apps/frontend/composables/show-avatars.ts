import { useLocalStorage } from '@vueuse/core'

export const showAvatars = useLocalStorage('stream-show-avatars', false)
