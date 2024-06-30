/// <reference types="vite/client" />
/// <reference types="vite-svg-loader" />

interface ImportMetaEnv {
  readonly VITE_STREAM_URL?: string
  // more env variables...
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
