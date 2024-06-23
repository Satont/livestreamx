import { useMutation, useQuery, useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { ref, watch } from 'vue'

import { FragmentType, graphql } from '@/gql'
import { useFragment } from '@/gql/fragment-masking.ts'
import { ChatEmote_FragmentFragment } from '@/gql/graphql.js'

export const ChatMessage_Fragment = graphql(`
  fragment ChatMessage_Fragment on ChatMessage {
    id
    segments {
      type
      content
      ... on MessageSegmentMention {
        user {
          id
          color
          displayName
        }
      }
      ... on MessageSegmentEmote {
        emote {
          ...ChatEmote_Fragment
        }
      }
    }
    sender {
      id
      avatarUrl
      color
      createdAt
      name
      displayName
    }
    createdAt
    reactions {
      ...ChatReaction_Fragment
    }
    replyTo
  }
`)

export const ChatEmote_Fragment = graphql(`
  fragment ChatEmote_Fragment on Emote {
    id
    name
    url
    width
    height
  }
`)

export const ChatReaction_Fragment = graphql(`
  fragment ChatReaction_Fragment on ChatMessageReaction {
    type
    user {
      id
      displayName
      color
      avatarUrl
    }
    reaction
    messageId
    ... on ChatMessageReactionEmote {
      emote {
        name
        url
      }
    }
  }
`)

export const useChat = createGlobalState(() => {
  const messages = ref<FragmentType<typeof ChatMessage_Fragment>[]>([])

  const sub = useSubscription({
    query: graphql(`
      subscription NewChatMessages {
        chatMessages {
          ...ChatMessage_Fragment
        }
      }
    `),
    variables: {}
  })

  watch(sub.data, (data) => {
    if (!data) return

    messages.value = [...messages.value, data.chatMessages]
  })

  const initialMessages = useQuery({
    query: graphql(`
      query ChatMessages {
        chatMessagesLatest {
          ...ChatMessage_Fragment
        }
      }
    `),
    variables: {}
  })

  watch(initialMessages.data, (data) => {
    if (!data) return

    messages.value = data.chatMessagesLatest
  })

  const useSendMessage = () =>
    useMutation(
      graphql(`
        mutation SendMessage($opts: SendMessageInput!) {
          sendMessage(input: $opts)
        }
      `)
    )

  const emotes = ref<ChatEmote_FragmentFragment[]>([])

  const useQueryEmotes = useQuery({
    query: graphql(`
      query ChatEmotes {
        getEmotes {
          ...ChatEmote_Fragment
        }
      }
    `),
    variables: {}
  })

  watch(useQueryEmotes.data, (data) => {
    if (!data) return

    const fragments = useFragment(ChatEmote_Fragment, data.getEmotes)
    emotes.value = fragments
  })

  const useReactionAddMutation = () =>
    useMutation(
      graphql(`
        mutation AddReaction($messageId: String!, $content: String!) {
          addReaction(messageId: $messageId, content: $content)
        }
      `)
    )

  const newReactionSub = useSubscription({
    query: graphql(`
      subscription NewMessageReaction {
        reactionAdd {
          ...ChatReaction_Fragment
        }
      }
    `),
    variables: {}
  })

  watch(newReactionSub.data, (data) => {
    if (!data) return

    const fragment = useFragment(ChatReaction_Fragment, data.reactionAdd)
    const message = messages.value.find(
      (m) => useFragment(ChatMessage_Fragment, m).id === fragment.messageId
    )
    if (!message) return

    useFragment(ChatMessage_Fragment, message).reactions.push(fragment)
  })

  return {
    messages,
    useSendMessage,
    emotes,
    useReactionAddMutation
  }
})
