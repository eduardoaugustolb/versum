import { defineConfig } from '@kubb/core'
import { pluginOas } from '@kubb/plugin-oas'
import { pluginTs } from '@kubb/plugin-ts'
import { pluginClient } from '@kubb/plugin-client'
import { pluginReactQuery } from '@kubb/plugin-react-query'
import { pluginZod } from '@kubb/plugin-zod'
import { pluginFaker } from '@kubb/plugin-faker'

export default defineConfig({
  root: '.',
  input: {
    path: 'https://versum-api.squareweb.app/openapi.json',
  },
  output: {
    path: './src/lib/kubb/gen',
    clean: true,
  },
  plugins: [
    pluginOas(),
    pluginTs({
      output: { path: 'models' },
    }),
    pluginClient({
      output: { path: 'clients' },
      client: "fetch"
    }),
    pluginReactQuery({
      output: { path: 'hooks' },
    }),
    pluginZod({
      output: { path: 'zod' },
    }),
    pluginFaker({
      output: { path: 'mocks' },
    }),
  ],
})
