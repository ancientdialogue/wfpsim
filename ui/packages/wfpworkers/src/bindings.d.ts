export {};

declare global {
  const WASM_R2: R2Bucket; //bucket
  const WFPSIM_KV: KVNamespace; //kv
  const ASSETS_ENDPOINT: string;
  const PREVIEW_ENDPOINT: string;
  const AUTH_KEY: string;
}
