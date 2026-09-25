const POPUP_WIDTH = 420
const POPUP_HEIGHT = 800

export function useChatPopout() {
  const router = useRouter()

  function chatPopoutPath(channelName: string) {
    return router.resolve({
      name: 'chat-channelName',
      params: { channelName }
    }).href
  }

  function openChatPopout(channelName: string) {
    const left =
      window.screenX + Math.max(0, (window.outerWidth - POPUP_WIDTH) / 2)
    const top =
      window.screenY + Math.max(0, (window.outerHeight - POPUP_HEIGHT) / 2)

    const features = [
      'popup=yes',
      `width=${POPUP_WIDTH}`,
      `height=${POPUP_HEIGHT}`,
      `left=${Math.round(left)}`,
      `top=${Math.round(top)}`,
      'resizable=yes',
      'scrollbars=yes'
    ].join(',')

    const popup = window.open(
      chatPopoutPath(channelName),
      `livestreamx-chat-${channelName}`,
      features
    )

    if (!popup) {
      useToast().toast({
        title: 'Popup blocked',
        description: 'Allow popups for this site to open the chat in a popup',
        variant: 'destructive'
      })
      return
    }

    popup.focus()
  }

  return {
    openChatPopout
  }
}
