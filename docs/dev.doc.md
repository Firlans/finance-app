module
- entity.go = mendefinisikan struktur data domain dan type utama. Biasanya berisi struct seperti `User`, `RegisterRequest`, `LoginResponse`, dan tipe data yang dipakai di module.
- repository.go = menampung data access layer. Berisi interface repository dan implementasi method yang berinteraksi dengan database, seperti `Save`, `FindByEmail`, `Update`, dan operasi CRUD lain.
- usecase.go = menampung business logic dan workflow aplikasi. Berisi interface usecase serta implementasi method yang mengorkestrasi repo, validasi, mailer, dan aturan bisnis.
- handler.go = (opsional tetapi direkomendasikan) menampung HTTP controller atau route handler. Berisi parsing request, pemanggilan usecase, dan respon HTTP.
- errors.go = (opsional) mendefinisikan error domain dan konstanta error yang dapat digunakan ulang di module.

Tujuan pembagian ini:
- `entity.go` = data / model
- `repository.go` = persistence / database
- `usecase.go` = logic aplikasi / business rules
- `handler.go` = interface HTTP / route
- `errors.go` = definisi error khusus module

Dengan aturan ini, developer baru akan lebih mudah memahami struktur module dan tanggung jawab setiap file. 