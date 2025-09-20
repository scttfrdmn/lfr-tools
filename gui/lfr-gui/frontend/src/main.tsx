import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'

console.log('main.tsx executing')
console.log('React:', React)
console.log('ReactDOM:', ReactDOM)

const rootElement = document.getElementById('root')
console.log('Root element:', rootElement)

if (!rootElement) {
  console.error('Root element not found!')
  document.body.innerHTML = '<div style="padding: 20px; color: red;">ERROR: Root element not found!</div>'
} else {
  try {
    console.log('Creating React root...')
    const root = ReactDOM.createRoot(rootElement as HTMLElement)
    console.log('Rendering App...')
    root.render(
      <React.StrictMode>
        <App />
      </React.StrictMode>,
    )
    console.log('App rendered successfully')
  } catch (error) {
    console.error('Error rendering React app:', error)
    document.body.innerHTML = `<div style="padding: 20px; color: red;">React Error: ${error}</div>`
  }
}
