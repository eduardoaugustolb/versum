type IsDevHandler = () => boolean;
type IsDebugHandler = () => boolean;

interface EnvHandlers {
  isDev?: IsDevHandler;
  isDebug?: IsDebugHandler;
}

const defaultIsDev: IsDevHandler = () =>
  process.env.BUN_ENV !== "production" && process.env.BUN_ENV !== "staging";

const defaultIsDebug: IsDebugHandler = () => !!process.env.DEBUG;

let currentIsDev: IsDevHandler = defaultIsDev;
let currentIsDebug: IsDebugHandler = defaultIsDebug;

export function getIsDev(): boolean {
  return currentIsDev();
}

export function getIsDebugEnabled(): boolean {
  return currentIsDebug();
}

export function setEnvHandlers(handlers: EnvHandlers): void {
  if (handlers.isDev) currentIsDev = handlers.isDev;
  if (handlers.isDebug) currentIsDebug = handlers.isDebug;
}

export function resetEnvHandlers(): void {
  currentIsDev = defaultIsDev;
  currentIsDebug = defaultIsDebug;
}

export type { IsDevHandler, IsDebugHandler, EnvHandlers };
