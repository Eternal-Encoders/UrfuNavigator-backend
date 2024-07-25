import express from 'express'
import payload from 'payload'
import https from "https";
import fs from "fs";
import path from "path";
import { initRedis } from '@aengz/payload-redis-cache';

require('dotenv').config()

const PORT = process.env.PORT_ENV || 5000;

initRedis({
  redisUrl: `redis://default:${process.env.REDIS_PASS}@${process.env.HOST}/cache`
})

const app = express()

// Redirect root to Admin panel
app.get('/', (_, res) => {
  res.redirect('/admin')
})

const start = async () => {
  // Initialize Payload
  await payload.init({
    secret: process.env.PAYLOAD_SECRET,
    express: app,
    onInit: async () => {
      payload.logger.info(`Payload Admin URL: ${payload.getAdminURL()}`)
    },
  })

  // Add your own express routes here
  app.listen(PORT);
}

start()
