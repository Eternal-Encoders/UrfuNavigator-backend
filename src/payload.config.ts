import path from 'path'

import { mongooseAdapter } from '@payloadcms/db-mongodb'
import { webpackBundler } from '@payloadcms/bundler-webpack'
import { slateEditor } from '@payloadcms/richtext-slate'
import { buildConfig } from 'payload/config'

import Users from './collections/Users'
import Insitutes from './collections/Institutes'
import Floors from './collections/Floors/Floors'
import { Media } from './collections/Media'
import Stairs from './collections/Stairs'
import GraphPoints from './collections/GraphPoint'
import DevOrTestText from './ui/befor_nav_link/dev_or_test'

export default buildConfig({
  admin: {
    user: Users.slug,
    bundler: webpackBundler(),
    components: {
      beforeNavLinks: [DevOrTestText]
    }
  },
  editor: slateEditor({}),
  collections: [
    Users,
    Insitutes,
    Floors,
    Media,
    Stairs,
    GraphPoints
  ],
  globals: [],
  cors: ('*'),
  typescript: {
    outputFile: path.resolve(__dirname, 'payload-types.ts'),
  },
  graphQL: {
    schemaOutputFile: path.resolve(__dirname, 'generated-schema.graphql'),
  },
  plugins: [
  ],
  db: mongooseAdapter({
    url: process.env.DATABASE_URI,
  }),
})
