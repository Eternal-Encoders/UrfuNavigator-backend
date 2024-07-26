import express from 'express'
import payload from 'payload'
import https from "https";
import fs from "fs";
import path from "path";

require('dotenv').config()

function loadEnvOrFile(name: string): string {
  let data = process.env[name]
  if (!data) {
      const path = process.env[`${name}_FILE`]
      if (!path) {
          return ''
      }
      
      try {
          data = fs.readFileSync(path, 'utf8');
      } catch (err) {
          return ''
      }
  }

  return data
}
process.env['PAYLOAD_SECRET'] = loadEnvOrFile('PAYLOAD_SECRET')
process.env['DATABASE_URI'] = loadEnvOrFile('DATABASE_URI')

const PORT = process.env.PORT_ENV || 5000;

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
