<script setup>
const sections = [
  { id: 1, name: 'Dashboard', link: '/dashboard' },
  { id: 2, name: 'Transactions', link: '/transactions' },
  { id: 3, name: 'Budgets', link: '/budgets' },
  { id: 4, name: 'Reports', link: '/reports' },
  { id: 5, name: 'Settings', link: '/settings' }
]

</script>

<template>
  <nav>
    <ul class="menu-sections">
      <li v-for="section in sections" :key="section.id" class="menu-section">
        <a :href="section.link">{{ section.name }}</a>
      </li>
    </ul>
  </nav>
</template>
