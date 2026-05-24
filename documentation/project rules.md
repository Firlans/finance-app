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

---

# Aturan Penamaan Section

- Semua file section disimpan di folder `src/views/sections/`.
- Nama file section menggunakan `PascalCase` dan wajib diakhiri dengan `Section`, misalnya:
  - `MenuSection.vue`
  - `AccountsSection.vue`
  - `CategoriesSection.vue`
- Section adalah komponen UI yang dipakai di dalam sebuah `Page`, bukan sebagai route mandiri.
- Section bertanggung jawab atas logika dan tampilan bagian tertentu dari halaman (misalnya: form CRUD, tabel data, navigasi).
- Section tidak memiliki route — tidak perlu diakhiri `Page` dan tidak akan ter-register sebagai rute otomatis.
- Page yang menggunakan section cukup melakukan `import` dan merender komponen section tersebut.

Contoh:

- `SettingsPage.vue` menggunakan `AccountsSection.vue` dan `CategoriesSection.vue`
- `DashboardPage.vue` dapat menggunakan `SummarySection.vue`, `RecentTransactionsSection.vue`, dst.
