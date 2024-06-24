import { useMutation, useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { graphql } from '@/gql'

export const useProfile = createGlobalState(() => {
  return useQuery({
    query: graphql(`
      query UserProfile {
        userProfile {
          id
          name
          displayName
          isBanned
          createdAt
          color
          avatarUrl
          isAdmin
          __typename
        }
      }
    `),
    variables: {}
  })
})

export const useProfileUpdate = createGlobalState(() => {
  return useMutation(
    graphql(`
      mutation UpdateProfile($input: UpdateUserProfileInput!) {
        updateUserProfile(input: $input) {
          __typename
          id
        }
      }
    `)
  )
})
