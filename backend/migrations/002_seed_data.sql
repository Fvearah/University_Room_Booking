USE university_room_booking;

-- Password untuk semua akun demo: password
INSERT INTO users (id, name, email, password_hash, role, created_at) VALUES
(1, 'Admin Kampus', 'admin@univ.ac.id', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'ADMIN', NOW()),
(2, 'Farhan', 'farhan@student.univ.ac.id', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'USER', NOW()),
(3, 'Naufal', 'naufal@student.univ.ac.id', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'USER', NOW());

INSERT INTO rooms (id, name, code, description, capacity, location, status, created_at, updated_at) VALUES
(1, 'Aula Utama', 'AUL-01', 'Aula besar dilengkapi proyektor ganda, sound system, dan AC sentral.', 250, 'Gedung Rektorat Lantai 1', 'ACTIVE', NOW(), NOW()),
(2, 'Ruang Rapat Senat', 'RPT-201', 'Meja konferensi kapasitas menengah dengan layar smart TV 75 inci.', 30, 'Gedung Administrasi Lantai 2', 'ACTIVE', NOW(), NOW()),
(3, 'Laboratorium Rekayasa Perangkat Lunak', 'LAB-302', '40 unit PC spesifikasi tinggi, jaringan gigabit, dan papan tulis kaca.', 40, 'Gedung Fakultas Teknik Lantai 3', 'ACTIVE', NOW(), NOW()),
(4, 'Ruang Seminar Lama', 'SEM-104', 'Sedang dalam renovasi instalasi listrik dan pendingin ruangan.', 80, 'Gedung Pascasarjana Lantai 1', 'INACTIVE', NOW(), NOW());

-- Sampel transaksi: lampau (otomatis tampil COMPLETED), mendatang (CONFIRMED), dan batal (CANCELLED)
INSERT INTO bookings (id, user_id, room_id, title, purpose, booking_date, start_time, end_time, status, created_at, updated_at) VALUES
(1, 2, 1, 'Kuliah Umum Industri', 'Sesi berbagi pengalaman praktisi teknologi informasi', '2026-11-20', '09:00:00', '11:00:00', 'CONFIRMED', NOW(), NOW()),
(2, 2, 1, 'Seminar Teknologi', 'Kegiatan seminar mahasiswa tingkat fakultas', '2026-11-25', '09:00:00', '12:00:00', 'CONFIRMED', NOW(), NOW()),
(3, 3, 2, 'Rapat Koordinasi Panitia', 'Pembahasan teknis acara tahunan himpunan mahasiswa', '2026-11-25', '13:00:00', '15:00:00', 'CANCELLED', NOW(), NOW());