import { createRouter, createWebHistory } from 'vue-router'

const modules = import.meta.glob('../views/*.{vue,pattern.js}')

const routes = Object.entries(modules).map(([filePath, loader]) => {
  const fileName = filePath.replace('../views/', '')
  const name = fileName
    .replace(/\.vue$/, '')
    .replace(/\.pattern\.js$/, '')
    .replace(/Page$/, '')
    .replace(/([a-z0-9])([A-Z])/g, '$1-$2')
    .replace(/([A-Z])([A-Z][a-z])/g, '$1-$2')
    .toLowerCase()

  return {
    path: name === 'landing' ? '/' : `/${name}`,
    name,
    component: loader
  }
})
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes
})

export default router
