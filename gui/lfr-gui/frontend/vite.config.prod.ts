import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Production optimized configuration
export default defineConfig({
  plugins: [
    react({
      // Enable React Fast Refresh in development
      fastRefresh: true,
    })
  ],

  // Build optimizations
  build: {
    target: 'esnext',
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true, // Remove console.log in production
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          // Separate vendor chunks for better caching
          vendor: ['react', 'react-dom'],
          cloudscape: ['@cloudscape-design/components', '@cloudscape-design/global-styles'],
          terminal: ['@xterm/xterm', '@xterm/addon-fit', '@xterm/addon-web-links'],
          router: ['react-router-dom'],
        },
      },
    },
    chunkSizeWarningLimit: 1000, // Increase limit for educational app
  },

  // Performance optimizations
  optimizeDeps: {
    include: [
      'react',
      'react-dom',
      '@cloudscape-design/components',
      '@xterm/xterm'
    ],
  },

  // Server configuration for development
  server: {
    port: 3000,
    strictPort: true,
    cors: {
      origin: ['http://localhost:3000', 'wails://localhost'],
      credentials: true,
    },
  },

  // Preview server configuration
  preview: {
    port: 4173,
    strictPort: true,
  },

  // Environment variables
  define: {
    'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'production'),
    'process.env.APP_VERSION': JSON.stringify(process.env.npm_package_version || '1.0.0'),
  },

  // Path resolution
  resolve: {
    alias: {
      '@': '/src',
      '@components': '/src/components',
      '@hooks': '/src/hooks',
      '@utils': '/src/utils',
    },
  },
})