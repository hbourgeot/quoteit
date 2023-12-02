import { Telegraf } from "telegraf";
import dotenv from "dotenv";

dotenv.config();
const BOT_API = process.env.TLG_BOT_API_KEY;

const bot = new Telegraf(BOT_API??'');

bot.start((ctx) => ctx.reply('Jelou'));
bot.help((ctx) => ctx.reply('no'));

bot.launch();

process.once("SIGINT", () => bot.stop("SIGINT"));
process.once("SIGTERM", () => bot.stop("SIGTERM"));