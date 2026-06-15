# Cinema Ticket Booking

ระบบจองตั๋วหนังออนไลน์ ทำด้วย Vue 3 + Go + MongoDB

---

## Tech Stack

**Frontend**
- Vue 3 + Vite
- Vue Router
- Axios
- Firebase Authentication

**Backend**
- Go + Fiber
- MongoDB
- Redis
- RabbitMQ

---

## วิธีรัน

### ต้องมีก่อน
- Docker + Docker Compose

### รันทั้งระบบ

```bash
docker compose up -d
```

เข้าใช้งานได้ที่ `http://localhost:5173`

---

## หน้าที่มีในระบบ (Frontend)

### สำหรับผู้ใช้ทั่วไป

| หน้า | URL | คำอธิบาย |
|------|-----|-----------|
| Login | `/login` | เข้าสู่ระบบด้วย Firebase |
| หน้าแรก | `/home` | ดูรายการหนังที่กำลังฉายและ Coming Soon กรองได้ตาม Genre |
| รายละเอียดหนัง | `/movies/:id` | ดูข้อมูลหนัง + รอบฉาย เลือกรอบแล้วกดจอง |
| เลือกที่นั่ง | `/showtimes/:id/seats` | เลือกที่นั่ง real-time ผ่าน WebSocket |
| ยืนยันการจอง | `/bookings/:id` | สรุปการจอง + เลือกวิธีชำระเงิน มี countdown 5 นาที |
| ประวัติการจอง | `/history` | ดูรายการจองทั้งหมดของตัวเอง |
| โปรไฟล์ | `/profile` | แก้ไขชื่อ เบอร์โทร วันเกิด |

### สำหรับ Admin

| หน้า | URL | คำอธิบาย |
|------|-----|-----------|
| จัดการหนัง | `/admin/home` | เพิ่ม แก้ไข ลบหนัง + จัดการรอบฉายแบบ accordion |
| ดูการจองทั้งหมด | `/admin/bookings` | ดูและจัดการการจองของทุก user |
| จัดการผู้ใช้ | `/admin/users` | ดู user ทั้งหมด + เปลี่ยน role |
| Event Log | `/admin/logs` | ดู audit log ของระบบ |

---

## โครงสร้างโฟลเดอร์ (Frontend)

```
src/
├── assets/styles/     # ไฟล์ CSS แยกตามหน้า
├── components/        # HeaderBar, AdminHeaderBar, Toast
├── composables/       # Logic แยกออกจาก View
│   ├── Login.js
│   ├── useMovies.js
│   ├── useMovieDetail.js
│   ├── useSeat.js
│   ├── useBookingConfirm.js
│   ├── useHistory.js
│   ├── useProfile.js
│   └── useAdminMovies.js
├── router/            # Vue Router
├── services/          # Axios instance + interceptor
├── views/             # หน้าต่างๆ
└── firebase.js        # Firebase config
```

---

## Environment Variables

สร้างไฟล์ `.env` ใน `/frontend`

```env
VITE_FIREBASE_API_KEY=
VITE_FIREBASE_AUTH_DOMAIN=
VITE_FIREBASE_PROJECT_ID=
VITE_FIREBASE_STORAGE_BUCKET=
VITE_FIREBASE_MESSAGING_SENDER_ID=
VITE_FIREBASE_APP_ID=
VITE_FIREBASE_MEASUREMENT_ID=
```

---

## Backend

ระบบหลังบ้านเขียนด้วย Go (Fiber) เชื่อมต่อกับ MongoDB, Redis, RabbitMQ และ Firebase สำหรับ Authentication

### ของที่ใช้และเอาไปใช้ทำอะไร

**MongoDB**
ใช้เป็นฐานข้อมูลหลักของทั้งระบบ เก็บข้อมูล users, movies, showtimes, seats, bookings, event_logs ทุก collection อยู่ใน database ชื่อ cinema

**Redis**
ใช้ทำ distributed lock ตอนจองที่นั่ง เวลา user กดจอง ระบบจะ lock ที่นั่งนั้นไว้ใน Redis ทันที (อายุ 5 นาที) ถ้ามีคนอื่นพยายามจองที่นั่งเดียวกันในช่วงนี้ ระบบจะตอบกลับว่าที่นั่งไม่ว่างทันที กันไม่ให้เกิดการจองซ้ำตอนมีคนกดพร้อมกันหลายคน

**RabbitMQ**
ใช้สำหรับส่ง event แบบ async ไปบันทึก audit log ตอนเกิดเหตุการณ์สำคัญ เช่น จองสำเร็จ, จองหมดเวลา, ที่นั่งถูกปลดล็อก ฯลฯ backend จะส่งข้อมูลเข้า queue แล้วมี worker (รันแยกเป็น goroutine อยู่เบื้องหลัง) คอยรับข้อมูลจาก queue มาบันทึกลง MongoDB collection event_logs อีกที ทำแบบนี้เพื่อให้ตอบ response กลับ user เร็ว ไม่ต้องรอเขียน log เสร็จก่อน

**Firebase Authentication**
ใช้สำหรับ login/register ฝั่ง frontend จะ login ผ่าน Firebase แล้วได้ token มา backend มีหน้าที่ verify token นี้อีกที (ผ่าน Firebase Admin SDK) ถ้า user คนนี้ยังไม่มีในระบบ backend จะสร้างให้อัตโนมัติใน MongoDB พร้อม role เริ่มต้นเป็น user

**WebSocket**
ใช้สำหรับอัพเดทสถานะที่นั่งแบบ real-time เวลามีคนจอง/จ่ายเงิน/ที่นั่งหลุดจอง ทุกคนที่เปิดหน้ารอบฉายเดียวกันอยู่จะเห็นสถานะที่นั่งเปลี่ยนทันทีโดยไม่ต้อง refresh และมี countdown แจ้งเวลาที่เหลือของ booking ที่ยังไม่จ่ายเงินด้วย

---

### โครงสร้างโฟลเดอร์ (Backend)

```
backend/
├── config/          # โหลดค่าจาก .env
├── database/        # การเชื่อมต่อ MongoDB, Redis, RabbitMQ
├── internal/
│   ├── handlers/    # ฟังก์ชันรับ request ของแต่ละ endpoint
│   ├── middleware/  # ตรวจสอบ Firebase token และเช็ค role admin
│   ├── models/      # struct ของแต่ละ collection
│   ├── routes/      # รวม route ทั้งหมด
│   ├── services/    # งานที่รันเบื้องหลัง (auto-cancel booking, event log worker)
│   └── ws/          # จัดการ WebSocket connection และการ broadcast
├── uploads/         # รูปโปสเตอร์หนังที่ admin อัพโหลด
└── main.go
```

---

### Environment Variables

สร้างไฟล์ `.env` ใน `/backend`

```env
JWT_SECRET=
PORT=8081

MONGODB_URI=
MONGODB_DB=cinema

REDIS_ADDR=

FIREBASE_PROJECT_ID=
GOOGLE_APPLICATION_CREDENTIALS=./serviceAccountKey.json
```

> ต้องเอาไฟล์ `serviceAccountKey.json` จาก Firebase Console (Project Settings > Service Accounts > Generate new private key) มาวางไว้ใน `/backend` ด้วย ไฟล์นี้และ `.env` ห้าม push ขึ้น git

---

### Authentication Flow

1. Frontend login ด้วย Firebase (email/password) แล้วได้ Firebase ID Token มา
2. ส่ง token นี้ไปที่ `POST /auth/login` ผ่าน header `Authorization: Bearer <token>`
3. Backend verify token กับ Firebase ถ้าเป็น user ใหม่จะสร้างข้อมูลใน MongoDB ให้อัตโนมัติ (role = user)
4. Backend ส่งข้อมูล user กลับมา (รวม role) ให้ frontend เก็บไว้
5. ทุก request ที่ต้อง login ต้องแนบ header `Authorization: Bearer <token>` เสมอ

> Firebase ID Token มีอายุ 1 ชั่วโมง frontend ใช้ `getIdToken()` ผ่าน interceptor ใน axios เพื่อขอ token ใหม่อัตโนมัติทุก request

---

### Role

มีสองแบบคือ `user` และ `admin` การจะให้ใครเป็น admin ต้องเรียก `PUT /admin/users/:id/role` โดยคนที่เรียกต้องเป็น admin อยู่แล้ว

> admin คนแรกของระบบต้องไปแก้ field `role` ใน MongoDB ตรงๆ

---

### API Endpoints

**Public (ไม่ต้อง login)**

| Method | Endpoint | คำอธิบาย |
|--------|----------|-----------|
| GET | `/movies` | ดูรายการหนังทั้งหมด รองรับ query `search`, `genre`, `status` |
| GET | `/movies/:id` | ดูรายละเอียดหนังเรื่องเดียว |
| GET | `/showtimes?movie_id=` | ดูรอบฉายของหนัง |
| GET | `/seats?showtime_id=` | ดูที่นั่งของรอบฉาย พร้อมสถานะ |
| WS | `/ws/showtimes/:showtime_id` | เชื่อม WebSocket รับสถานะที่นั่งและ countdown แบบ real-time |

**Auth**

| Method | Endpoint | คำอธิบาย |
|--------|----------|-----------|
| POST | `/auth/login` | verify firebase token แล้วสร้าง/ดึง user จาก MongoDB |

**User (ต้อง login)**

| Method | Endpoint | คำอธิบาย |
|--------|----------|-----------|
| GET | `/users/me` | ดูข้อมูลตัวเอง |
| PUT | `/users/me` | แก้ไขข้อมูลตัวเอง (ชื่อ เบอร์โทร วันเกิด) |
| POST | `/bookings` | จองที่นั่ง |
| GET | `/bookings/me` | ดูประวัติการจองของตัวเอง |
| GET | `/bookings/:id` | ดูรายละเอียดการจอง |
| PUT | `/bookings/:id/pay` | จ่ายเงิน |

**Admin (ต้อง login และเป็น admin)**

| Method | Endpoint | คำอธิบาย |
|--------|----------|-----------|
| POST | `/admin/movies` | เพิ่มหนัง (อัพโหลดรูปโปสเตอร์ผ่าน multipart/form-data field `poster`) |
| PUT | `/admin/movies/:id` | แก้ไขหนัง |
| DELETE | `/admin/movies/:id` | ลบหนัง |
| POST | `/admin/showtimes` | เพิ่มรอบฉาย |
| PUT | `/admin/showtimes/:id` | แก้ไขรอบฉาย |
| DELETE | `/admin/showtimes/:id` | ลบรอบฉาย |
| POST | `/admin/seats/generate` | สร้างที่นั่งทั้งหมดของรอบฉาย (ระบุ `rows` กับ `seats_per_row`) |
| GET | `/admin/users` | ดู user ทั้งหมด |
| PUT | `/admin/users/:id/role` | เปลี่ยน role ของ user |
| GET | `/admin/bookings` | ดูการจองทั้งหมด |
| PUT | `/admin/bookings/:id` | แก้ไขสถานะการจอง |
| DELETE | `/admin/bookings/:id` | ยกเลิกการจอง |

---

### Booking Flow

1. user เลือกที่นั่งแล้วกดจอง (`POST /bookings`)
   - backend lock ที่นั่งใน Redis ไว้ 5 นาที ถ้ามีคนล็อกไว้ก่อนแล้วจะได้ error 409 ทันที
   - สถานะที่นั่งใน MongoDB เปลี่ยนเป็น `locked`
   - broadcast ผ่าน WebSocket ให้ทุกคนที่ดูรอบฉายนี้เห็นว่าที่นั่งถูกล็อกแล้ว
   - ส่ง event `booking_success` เข้า RabbitMQ เพื่อบันทึก audit log

2. user จ่ายเงิน (`PUT /bookings/:id/pay`)
   - สถานะที่นั่งเปลี่ยนเป็น `booked` แบบถาวร
   - ปลด lock ใน Redis
   - broadcast ผ่าน WebSocket

3. ถ้า user ไม่จ่ายเงินภายใน 5 นาที มี background job รันทุก 1 นาทีคอยเช็คและยกเลิก booking ให้อัตโนมัติ คืนที่นั่งเป็น `available` ปลด lock และส่ง event `booking_timeout` เข้า RabbitMQ

4. ตอนที่ booking ยังรอจ่ายเงินอยู่ ระบบจะ broadcast countdown (เวลาที่เหลือ) ทุก 10 วินาทีผ่าน WebSocket

---

### สถานะที่นั่ง

| สถานะ | ความหมาย |
|-------|----------|
| `available` | ว่าง จองได้ |
| `locked` | มีคนกำลังจองอยู่ รอจ่ายเงิน (สูงสุด 5 นาที) |
| `booked` | จ่ายเงินแล้ว ถาวร |

---

### Database Schema (MongoDB Collections)

| Collection | Fields |
|-----------|--------|
| users | id, uid, email, role, name, phone, date_of_birth |
| movies | id, title, description, genre, poster_url, duration, status, created_at |
| showtimes | id, movie_id, room, start_time, end_time, seat_count, price, created_at |
| seats | id, showtime_id, label, row, number, status, locked_by, locked_until, booking_id |
| bookings | id, user_id, showtime_id, seat_ids, status, total_price, booking_number, created_at, paid_at |
| event_logs | id, event, user_id, booking_id, seat_ids, message, created_at |

---

### สิ่งที่ยังไม่ได้ทำ

- ยังไม่มี pagination ใน endpoint ที่ return list เช่น `/movies`, `/admin/bookings`, `/admin/event-logs`
- `/seats` ยังไม่ sort ตามแถวและเลขที่นั่ง
- รูปโปสเตอร์เก็บไว้ที่ local disk ของ server ถ้า deploy แบบ container แล้ว restart รูปจะหาย ควรย้ายไป cloud storage ในอนาคต
- การกรองหนังตามวันที่ฉายยังไม่มี endpoint รองรับ