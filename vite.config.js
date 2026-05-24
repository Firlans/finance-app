import { fileURLToPath, URL } from 'node:url'
import { compileTemplate } from '@vue/compiler-sfc'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

function vuePatternPlugin() {
  const templateRE = /export\s+const\s+template\s*=\s*(`([\s\S]*?)`|'([\s\S]*?)'|"([\s\S]*?)")/m

  return {
    name: 'vite-plugin-vue-pattern',
    enforce: 'pre',
    transform(code, id) {
      if (!id.endsWith('.pattern.js')) return null

      const match = code.match(templateRE)
      if (!match) {
        throw new Error(`[vite-plugin-vue-pattern] Missing export const template in ${id}`)
      }

      const template = match[2] ?? match[3] ?? match[4] ?? ''
      const source = code.replace(templateRE, '')
      const compiled = compileTemplate({ source: template, filename: id, id })

      if (compiled.errors && compiled.errors.length) {
        throw compiled.errors[0]
      }

      return `${source}
import { defineComponent, onMounted } from 'vue'
${compiled.code}
const __pattern_setup = typeof setup === 'function' ? setup : () => ({})
export default defineComponent({
  setup(props, ctx) {
    const bindings = __pattern_setup(props, ctx)
    onMounted(async () => {
      if (typeof preProcess === 'function') await preProcess()
      if (typeof process === 'function') await process()
    })
    return bindings
  },
  render
})`
    }
  }
}

// https://vite.dev/config/
export default defineConfig({
  base: '/finance-app/',
  plugins: [
    vuePatternPlugin(),
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '@packages/components': fileURLToPath(new URL('../../packages/components', import.meta.url)),
      '@packages/utils': fileURLToPath(new URL('../../packages/utils', import.meta.url))
    },
  },
})
