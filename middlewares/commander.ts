import type { Context } from "telegraf";

export function commander() {
  return async (ctx: Context, next: () => Promise<void>) => {
    if (ctx.updateType === 'message' && (ctx.chat?.type === 'group' || ctx.chat?.type === 'supergroup')) {
      
    }
  }
}