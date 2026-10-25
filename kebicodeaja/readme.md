# keBCodeAja

Media Pembelajaran Interaktif Informatika Berbasis LMS untuk Siswa Kelas X (Fase E).

- backend  : Golang (Gin, GORM, JWT) + MySQL
- frontend : React (Vite) + Tailwind CSS
- docs     : ERD dan dokumentasi

## Kebutuhan

- Go 1.21 atau lebih baru
- Node.js 20.19 atau lebih baru
- MySQL (XAMPP atau Laragon), nyalakan dulu sebelum menjalankan backend

## Menjalankan

Database `kebcodeaja` dibuat otomatis oleh backend. Kalau MySQL punya password, isi `DB_PASS` di `backend/.env`.

Backend (terminal 1):

    cd backend
    go mod tidy
    go run ./cmd/server

Frontend (terminal 2):

    cd frontend
    npm install
    npm run dev

Buka http://localhost:5173

Backend harus dinyalakan ulang setiap kali file Go diubah.

## Fitur saat ini

- Registrasi dan login (JWT, bcrypt)
- Dashboard guru dan siswa dengan data dari database (kosong tampil 0)
- Guru membuat kelas, kode kelas dibuat otomatis
- Siswa bergabung ke kelas dengan kode kelas
