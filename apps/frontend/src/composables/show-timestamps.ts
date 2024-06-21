import { useLocalStorage } from "@vueuse/core";

export const showTimestamps = useLocalStorage('stream-show-timestamps', false)