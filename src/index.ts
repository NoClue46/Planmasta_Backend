import { serve } from "@hono/node-server";
import { Hono } from "hono";
import OpenAI from "openai";
import dotenv from "dotenv";
dotenv.config();

const app = new Hono();

const openai = new OpenAI({
  apiKey: process.env.OPENAI_API_KEY,
});

app.post("/chat", async (c) => {
  const { message } = await c.req.json();
  const response = await openai.chat.completions.create({
    model: "gpt-4.1",
    messages: [{ role: "user", content: message }],
  });
  return c.json({ response: response });
});

serve(
  {
    fetch: app.fetch,
    port: 3000,
  },
  (info) => {
    console.log(`Server is running on port :${info.port}`);
  }
);
