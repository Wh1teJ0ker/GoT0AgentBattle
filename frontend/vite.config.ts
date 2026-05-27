import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

const directivePackages = [
  '/node_modules/antd/',
  '/node_modules/@ant-design/',
  '/node_modules/rc-',
  '/node_modules/framer-motion/',
  '/node_modules/motion-dom/',
  '/node_modules/motion-utils/',
]

function shouldStripUseClient(id: string) {
  return directivePackages.some((segment) => id.includes(segment))
}

function stripUseClientDirective() {
  return {
    name: 'strip-third-party-use-client',
    enforce: 'pre' as const,
    transform(code: string, id: string) {
      if (!shouldStripUseClient(id)) {
        return null
      }

      const next = code.replace(/^[\t ]*["']use client["'];?\s*/gm, '')
      if (next === code) {
        return null
      }

      return {
        code: next,
        map: null,
      }
    },
  }
}

function shouldIgnoreModuleDirectiveWarning(warning: {code?: string; id?: string}) {
  return (
    warning.code === 'MODULE_LEVEL_DIRECTIVE' &&
    typeof warning.id === 'string' &&
    shouldStripUseClient(warning.id)
  )
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [stripUseClientDirective(), react()],
  build: {
    chunkSizeWarningLimit: 700,
    rollupOptions: {
      onwarn(warning, warn) {
        if (shouldIgnoreModuleDirectiveWarning(warning)) {
          return
        }
        warn(warning)
      },
      output: {
        manualChunks(id) {
          if (!id.includes('node_modules')) {
            return
          }

          if (
            id.includes('/node_modules/react/') ||
            id.includes('/node_modules/react-dom/') ||
            id.includes('/node_modules/scheduler/')
          ) {
            return 'react-vendor'
          }

          if (
            id.includes('/node_modules/antd/es/form/') ||
            id.includes('/node_modules/antd/es/input/') ||
            id.includes('/node_modules/antd/es/input-number/') ||
            id.includes('/node_modules/antd/es/select/') ||
            id.includes('/node_modules/antd/es/switch/') ||
            id.includes('/node_modules/antd/es/button/') ||
            id.includes('/node_modules/rc-field-form/') ||
            id.includes('/node_modules/rc-input/') ||
            id.includes('/node_modules/rc-textarea/') ||
            id.includes('/node_modules/rc-select/') ||
            id.includes('/node_modules/rc-switch/') ||
            id.includes('/node_modules/rc-input-number/')
          ) {
            return 'antd-controls'
          }

          if (
            id.includes('/node_modules/antd/es/layout/') ||
            id.includes('/node_modules/antd/es/grid/') ||
            id.includes('/node_modules/antd/es/flex/') ||
            id.includes('/node_modules/antd/es/card/') ||
            id.includes('/node_modules/antd/es/avatar/') ||
            id.includes('/node_modules/antd/es/badge/') ||
            id.includes('/node_modules/antd/es/divider/') ||
            id.includes('/node_modules/antd/es/list/') ||
            id.includes('/node_modules/antd/es/progress/') ||
            id.includes('/node_modules/antd/es/space/') ||
            id.includes('/node_modules/antd/es/statistic/') ||
            id.includes('/node_modules/antd/es/tag/') ||
            id.includes('/node_modules/antd/es/typography/')
          ) {
            return 'antd-display'
          }

          if (
            id.includes('/node_modules/@ant-design/') ||
            id.includes('/node_modules/@rc-component/') ||
            id.includes('/node_modules/antd/es/')
          ) {
            return 'antd-shared'
          }

          if (
            id.includes('/node_modules/framer-motion/') ||
            id.includes('/node_modules/motion-dom/') ||
            id.includes('/node_modules/motion-utils/')
          ) {
            return 'motion-vendor'
          }

          if (
            id.includes('/node_modules/dayjs/') ||
            id.includes('/node_modules/@babel/runtime/')
          ) {
            return 'shared-vendor'
          }

          return 'vendor'
        },
      },
    },
  },
})
