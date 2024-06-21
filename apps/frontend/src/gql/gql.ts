/* eslint-disable */
import * as types from './graphql';
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 */
const documents = {
    "\n\t\t\tsubscription NewChatMessages {\n\t\t\t\tchatMessages {\n\t\t\t\t\t...MessageFragment\n\t\t\t\t}\n\t\t\t}\n\t\t": types.NewChatMessagesDocument,
    "\n\t\t\tquery ChatMessages {\n\t\t\t\tchatMessagesLatest {\n\t\t\t\t\t...MessageFragment\n\t\t\t\t}\n\t\t\t}\n\t\t": types.ChatMessagesDocument,
    "\n\t\tmutation SendMessage($opts: SendMessageInput!) {\n\t\t\tsendMessage(input: $opts)\n\t\t}\n\t": types.SendMessageDocument,
    "\n\t\t\tquery UserProfile {\n\t\t\t\tuserProfile {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tdisplayName\n\t\t\t\t\tisBanned\n\t\t\t\t\tcreatedAt\n\t\t\t\t\tcolor\n\t\t\t\t\tavatarUrl\n\t\t\t\t}\n\t\t\t}\n\t\t": types.UserProfileDocument,
    "\n    fragment MessageFragment on ChatMessage {\n\t\t\tid\n\t\t\tsegments {\n\t\t\t\t\ttype\n\t\t\t\t\tcontent\n\t\t\t\t\t...on MessageSegmentMention {\n\t\t\t\t\t\tuser {\n\t\t\t\t\t\t\tcolor\n\t\t\t\t\t\t\tdisplayName\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t}\n\t\t\tsender {\n\t\t\t\t\tid\n\t\t\t\t\tavatarUrl\n\t\t\t\t\tcolor\n\t\t\t\t\tcreatedAt\n\t\t\t\t\tname\n\t\t\t\t\tdisplayName\n\t\t\t}\n\t\t\tcreatedAt\n    }\n": types.MessageFragmentFragmentDoc,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = graphql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function graphql(source: string): unknown;

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tsubscription NewChatMessages {\n\t\t\t\tchatMessages {\n\t\t\t\t\t...MessageFragment\n\t\t\t\t}\n\t\t\t}\n\t\t"): (typeof documents)["\n\t\t\tsubscription NewChatMessages {\n\t\t\t\tchatMessages {\n\t\t\t\t\t...MessageFragment\n\t\t\t\t}\n\t\t\t}\n\t\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tquery ChatMessages {\n\t\t\t\tchatMessagesLatest {\n\t\t\t\t\t...MessageFragment\n\t\t\t\t}\n\t\t\t}\n\t\t"): (typeof documents)["\n\t\t\tquery ChatMessages {\n\t\t\t\tchatMessagesLatest {\n\t\t\t\t\t...MessageFragment\n\t\t\t\t}\n\t\t\t}\n\t\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tmutation SendMessage($opts: SendMessageInput!) {\n\t\t\tsendMessage(input: $opts)\n\t\t}\n\t"): (typeof documents)["\n\t\tmutation SendMessage($opts: SendMessageInput!) {\n\t\t\tsendMessage(input: $opts)\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tquery UserProfile {\n\t\t\t\tuserProfile {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tdisplayName\n\t\t\t\t\tisBanned\n\t\t\t\t\tcreatedAt\n\t\t\t\t\tcolor\n\t\t\t\t\tavatarUrl\n\t\t\t\t}\n\t\t\t}\n\t\t"): (typeof documents)["\n\t\t\tquery UserProfile {\n\t\t\t\tuserProfile {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\tdisplayName\n\t\t\t\t\tisBanned\n\t\t\t\t\tcreatedAt\n\t\t\t\t\tcolor\n\t\t\t\t\tavatarUrl\n\t\t\t\t}\n\t\t\t}\n\t\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n    fragment MessageFragment on ChatMessage {\n\t\t\tid\n\t\t\tsegments {\n\t\t\t\t\ttype\n\t\t\t\t\tcontent\n\t\t\t\t\t...on MessageSegmentMention {\n\t\t\t\t\t\tuser {\n\t\t\t\t\t\t\tcolor\n\t\t\t\t\t\t\tdisplayName\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t}\n\t\t\tsender {\n\t\t\t\t\tid\n\t\t\t\t\tavatarUrl\n\t\t\t\t\tcolor\n\t\t\t\t\tcreatedAt\n\t\t\t\t\tname\n\t\t\t\t\tdisplayName\n\t\t\t}\n\t\t\tcreatedAt\n    }\n"): (typeof documents)["\n    fragment MessageFragment on ChatMessage {\n\t\t\tid\n\t\t\tsegments {\n\t\t\t\t\ttype\n\t\t\t\t\tcontent\n\t\t\t\t\t...on MessageSegmentMention {\n\t\t\t\t\t\tuser {\n\t\t\t\t\t\t\tcolor\n\t\t\t\t\t\t\tdisplayName\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t}\n\t\t\tsender {\n\t\t\t\t\tid\n\t\t\t\t\tavatarUrl\n\t\t\t\t\tcolor\n\t\t\t\t\tcreatedAt\n\t\t\t\t\tname\n\t\t\t\t\tdisplayName\n\t\t\t}\n\t\t\tcreatedAt\n    }\n"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;