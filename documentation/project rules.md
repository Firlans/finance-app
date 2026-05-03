# Aturan Penamaan View

- Semua file view harus menggunakan `PascalCase`.
- Nama file view wajib diakhiri dengan `Page`, misalnya:
  - `HomePage.vue`
  - `LoginPage.vue`
  - `ForgetPasswordPage.vue`
- Rute otomatis dihasilkan dari nama file view dengan:
  - menghapus sufiks `Page`
  - mengonversi dari PascalCase ke `kebab-case`
  - `LandingPage.vue` menjadi `/` sebagai landing page default
  - `ForgetPasswordPage.vue` menjadi `/forget-password`

Contoh:

- `DashboardPage.vue` → `dashboard`
- `ResetPasswordPage.vue` → `reset-password`
