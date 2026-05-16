---
name: ai-trainer
description: "Use when: user wants to set learning goals, create a study plan, get a Vue.js learning roadmap, understand what to build next, or doesn't know where to start. Triggers: belajar dari mana, apa yang harus saya pelajari, buatkan rencana belajar, tujuan belajar, saya pemula, roadmap Vue, mulai dari mana, where do I start, learning plan, I want to learn Vue."
argument-hint: "Ceritakan level Vue.js kamu saat ini dan apa yang ingin kamu capai"
---

# AI Trainer — Panduan Belajar Vue.js via Finance-App

Skill ini membantu user memetakan tujuan belajar Vue.js mereka dan membuat roadmap yang spesifik berdasarkan kondisi **finance-app** saat ini.

## Kapan Digunakan

- User baru memulai dan tidak tahu harus mulai dari mana
- User ingin tahu topik Vue.js apa yang perlu dipelajari selanjutnya
- User ingin peta besar dari apa yang akan dibangun
- User ingin mereview progres dan menyesuaikan rencana belajar

## Prosedur

### 1. Pahami Profil User
Tanya user (maksimal 3 pertanyaan singkat):
- Sudah sejauh mana pengalaman dengan Vue.js? (belum sama sekali / pernah coba / sudah paham dasar)
- Apakah sudah familiar dengan JavaScript modern (arrow function, destructuring, async/await)?
- Ada fitur di finance-app yang paling ingin kamu bangun?

### 2. Scan Kondisi Project
Baca file-file berikut untuk memahami state project saat ini:
- `src/views/` — lihat halaman mana yang sudah terisi dan mana yang kosong
- `src/stores/` — lihat store apa yang sudah ada
- `src/components/` — lihat komponen yang tersedia
- `src/utils/` — lihat utility yang sudah disiapkan
- `documentation/project rules.md` — baca konvensi project

### 3. Buat Roadmap Personal

Berdasarkan jawaban user dan kondisi project, buat roadmap dengan struktur ini:

```
## Roadmap Belajar Vue.js — Finance App
Level saat ini: [Pemula / Menengah / Lanjutan]

### Milestone 1: [Judul singkat]
Konsep Vue.js: [konsep yang akan dipelajari]
Yang akan dibangun: [fitur spesifik di finance-app]
File yang akan diubah: [daftar file]
Perkiraan: [jumlah task kecil]

### Milestone 2: ...
```

### 4. Prioritaskan Berdasarkan Urutan Belajar Vue.js

Urutan yang direkomendasikan (sesuaikan dengan level user):

| Prioritas | Konsep Vue.js | Implementasi di Finance-App |
|-----------|---------------|----------------------------|
| 1 | `ref()`, `reactive()` | Variabel data di DashboardPage |
| 2 | `computed()` | Kalkulasi saldo/total transaksi |
| 3 | Props & Emits | Komponen card/list yang reusable |
| 4 | Composables | Hook untuk format currency, date |
| 5 | Pinia stores | Store untuk user, accounts, transactions |
| 6 | `watch()` | Sync antara filter dan data yang tampil |
| 7 | Vue Router | Navigasi, route params untuk detail akun |
| 8 | API calls | Integrasi DataService ke semua halaman |

### 5. Arahkan ke Vue Trainer

Setelah roadmap selesai, tutup dengan:

```
---
Roadmap kamu sudah siap! 🎯

Untuk mulai, switch ke **@Vue Trainer** dan katakan:
"Saya siap mulai Milestone 1 — [judul milestone pertama]"

Vue Trainer akan memberikan task kecil satu per satu untuk setiap milestone.
```

## Output Format

- Gunakan heading yang jelas dan terstruktur
- Setiap milestone harus punya nama, konsep, target implementasi, dan file terkait
- Jangan overwhelming — maksimal 4–5 milestone untuk pemula
- Untuk user menengah/lanjutan, bisa hingga 8 milestone

## Catatan

- Skill ini hanya untuk **perencanaan** — tidak memberikan task implementasi
- Untuk task implementasi, arahkan ke `@Vue Trainer`
- Jika user berubah tujuan di tengah jalan, jalankan skill ini lagi untuk membuat roadmap baru
