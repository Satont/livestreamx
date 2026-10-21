import type { MediaPlayerElement } from 'vidstack/elements'
import { onBeforeUnmount, watch, type Ref } from 'vue'

// Kept out of `.vue` files on purpose: the vidstack vite plugin scans them for
// `<media-*>` tags and `vds-*` theme classes to generate `vidstack/bundle`, and
// a `vds-*` name that isn't one of its theme classes stalls that scan, so tags
// after it (e.g. `<media-video-layout>`) are never registered.
const SETTINGS_MENU_SELECTOR = 'media-video-layout .vds-settings-menu'

/**
 * Returns a host element kept right before the default layout settings (gear)
 * menu, so custom control buttons can be teleported next to it.
 *
 * The default layout has no slots for extra buttons, and vidstack re-renders
 * the controls when the layout changes (e.g. large/small), which drops foreign
 * nodes: a mutation observer puts the host back.
 */
export function useSettingsMenuSlot(player: Ref<MediaPlayerElement | null>) {
  const host = document.createElement('div')
  host.style.display = 'contents'

  function place() {
    const settingsMenu = player.value?.querySelector(SETTINGS_MENU_SELECTOR)
    if (!settingsMenu) {
      host.remove()
      return
    }
    if (settingsMenu.previousElementSibling !== host) {
      settingsMenu.before(host)
    }
  }

  const observer = new MutationObserver(place)

  watch(player, (element) => {
    observer.disconnect()
    if (element) {
      observer.observe(element, { childList: true, subtree: true })
    }
    place()
  })

  onBeforeUnmount(() => {
    observer.disconnect()
    host.remove()
  })

  return host
}
