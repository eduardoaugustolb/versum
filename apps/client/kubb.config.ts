import { defineConfig } from "@kubb/core";
import { pluginClient } from "@kubb/plugin-client";
import { pluginFaker } from "@kubb/plugin-faker";
import { pluginOas } from "@kubb/plugin-oas";
import { pluginReactQuery } from "@kubb/plugin-react-query";
import { pluginTs } from "@kubb/plugin-ts";
import { pluginZod } from "@kubb/plugin-zod";

export default defineConfig({
  root: ".",
  input: {
    path: "http://localhost:4002/openapi.json",
  },
  output: {
    path: "./src/lib/kubb/gen",
    clean: true,
  },
  plugins: [
    pluginOas(),
    pluginTs({
      output: { path: "models", override: true },
    }),
    pluginClient({
      output: { path: "clients", override: true },
      client: "fetch",
      baseURL: "https://versum-api.squareweb.app/",
    }),
    pluginReactQuery({
      client: {
        baseURL: "https://versum-api.squareweb.app/",
        client: "fetch",
      },
      output: { path: "hooks", override: true },
      infinite: {},
    }),
    pluginZod({
      output: { path: "zod", override: true },
    }),
    pluginFaker({
      output: { path: "mocks", override: true },
    }),
  ],
});
