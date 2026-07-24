/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_PI_TOUCH?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
