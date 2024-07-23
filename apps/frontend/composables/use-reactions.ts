import { createGlobalState } from '@vueuse/core'
import type { FragmentType } from '~/gql'

import { ChatMessage_Fragment, useChat } from '~/api/chat'
import { useFragment } from '~/gql'
import { type ChatMessage_FragmentFragment } from '~/gql/graphql.js'

export const useReactions = createGlobalState(() => {
  const { channelData, useReactionAddMutation } = useChat()
  const reactionAddMutation = useReactionAddMutation()

  const currentMessage = ref<ChatMessage_FragmentFragment | null>(null)
  const dialogOpened = ref(false)

  function openDialog(
    message: FragmentType<typeof ChatMessage_Fragment> | null
  ) {
    currentMessage.value = useFragment(ChatMessage_Fragment, message)
    dialogOpened.value = true
  }

  async function addReaction(name: string, messageId?: string) {
    if (!channelData.value) return
    const msgId = messageId ?? currentMessage.value?.id
    if (!msgId) return

    try {
      const { error } = await reactionAddMutation.executeMutation({
        messageId: msgId,
        content: name,
        channelID: channelData.value!.fetchUserByName.id
      })

      if (error) {
        throw new Error(error.toString())
      }
    } catch (e) {
      console.log(e)
    }
  }

  return {
    currentMessage: readonly(currentMessage),
    dialogOpened,
    openDialog,
    addReaction
  }
})
