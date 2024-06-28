import { useMutation, useQuery, useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed, ref, watch } from 'vue'

import { FragmentType, graphql } from '@/gql'
import { useFragment } from '@/gql/fragment-masking.ts'
import { ChatEmote_FragmentFragment } from '@/gql/graphql.js'
import { useRoute } from "vue-router";

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
  const routerParams = useRoute()
  const channelName = computed(() => {
    if (typeof routerParams.params.channelName != 'string') return null

    return routerParams.params.channelName
  })

  const channelRequest = useQuery({
    query: graphql(`
        query ChannelPageChannel($channelName: String!) {
          fetchUserByName(name: $channelName) {
            id
            name
            avatarUrl
          }
        }
    `),
    get variables() {
      return {
        channelName: channelName.value
      }
    }
  })

  const channelData = computed(() => {
    if (!channelRequest.data) return null

    return channelRequest.data.value
  })

  const messages = ref<FragmentType<typeof ChatMessage_Fragment>[]>([])

  const sub = useSubscription({
    query: graphql(`
      subscription NewChatMessages($channelID: UUID!) {
        chatMessages(channelID: $channelID) {
          ...ChatMessage_Fragment
        }
      }
    `),
    get variables() {
      return {
        channelID: channelData.value?.fetchUserByName.id
      }
    },
    pause() {
      return !channelData.value?.fetchUserByName.id
    },
  })

  watch(sub.data, (data) => {
    if (!data) return

    messages.value = [...messages.value, data.chatMessages]
  })

  const initialMessages = useQuery({
    query: graphql(`
      query ChatMessages($channelID: UUID!) {
        chatMessagesLatest(channelID: $channelID) {
          ...ChatMessage_Fragment
        }
      }
    `),
    get variables() {
      return {
        channelID: channelData.value?.fetchUserByName.id
      }
    },
    pause() {
      return !channelData.value?.fetchUserByName.id
    },
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
      query ChatEmotes($channelID: UUID!) {
        getEmotes(channelID: $channelID) {
          ...ChatEmote_Fragment
        }
      }
    `),
    get variables() {
      return {
        channelID: channelData.value?.fetchUserByName.id
      }
    },
    pause() {
      return !channelData.value?.fetchUserByName.id
    }
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
      subscription NewMessageReaction($channelID: UUID!) {
        reactionAdd(channelID: $channelID) {
          ...ChatReaction_Fragment
        }
      }
    `),
    get variables() {
      return {
        channelID: channelData.value?.fetchUserByName.id
      }
    },
    pause() {
      return !channelData.value?.fetchUserByName.id
    },
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
    useReactionAddMutation,
    channelData,
  }
})
