# 🎬 OPHIM CLI - RẠP CHIẾU PHIM "VÙNG KÍN" 🍿

Chào cục dzàng! Đây là công cụ xem phim lậu đỉnh cao chạy ngay trong Terminal, được tối ưu cho các mấy cưng nào thích vừa code vừa cày phim. Giao diện Catppuccin siêu cute, có lịch sử xem phim và tính năng chào hỏi cực sến.



## ✨ Tính năng
- 🔍 **Tìm phim:** Search phát ra luôn.
- 🎞️ **Chọn tập:** Hỗ trợ cả phim bộ và phim lẻ.
- 📜 **Lịch sử:** Nhắc cưng xem tiếp tập đang dang dở (Cục dzàng có muốn coi tiếp hem?).
- 🎨 **Giao diện:** Hệ màu Catppuccin siêu mịn, hỗ trợ cuộn chuột xem mô tả phim.
- 🚀 **Tốc độ:** Chạy bằng Go + MPV, mượt hơn cả người yêu cũ trở mặt.

---

## 🛠 Yêu cầu hệ thống

Cái "động cơ" chính để phát phim là **MPV**. Cưng phải cài nó trước:

- **MacOS:** `brew install mpv`
- **Linux (Ubuntu/Debian):** `sudo apt update && sudo apt install mpv`
- **Windows:** Tải bản Zip tại [mpv.io](https://mpv.io/installation/), giải nén và **add vào PATH** hệ thống.

---

## 🏗 Hướng dẫn cài đặt (Dành cho anh em)

### 1. Clone Project
```bash
git clone [https://github.com/your-username/ophim-cli.git](https://github.com/your-username/ophim-cli.git)
cd ophim-cli
```
### 2. Cài con Go
```bash
go mod tidy
```

### 3. Cấu hình "vùng kín" (.env)
Tạo file .env ngay root nho. Xong cưng làm ơn gửi cho tui cái tin nhắn/email qua facebook https://www.facebook.com/panadolextra9103/ hoặc email anhthuhuynh9103@gmail.com nho
Sau khi cưng nhận được file .env từ tui, có 2 cách để xài:
- **Cách lười:** Luôn mở terminal trong đúng folder ophim-cli rồi mới gõ `go run main.go`.
- **Cách pro (Khuyên dùng):** Mở file main.go, tìm các hàm searchMoviesCmd và fetchDetailMsg, dán thẳng mấy cái link API vô code luôn rồi hãy `go build`. Làm vậy thì cưng đứng ở đâu trên máy gõ xemphim nó cũng chạy, không cần lôi cái file .env đi theo khắp nơi.

---

## 🚀 Cách Build & Chạy

### Cách 1: Chạy trực tiếp (test cho lẹ)
```bash
go run main.go
```

### Cách 2: Build thành lệnh hệ thống (khuyên mấy cưng xài)
Cưng vô cái folder ophim-cli nha
- **MacOS/Linux:** `go build -o xemphim`
`sudo mv xemphim /usr/local/bin/`
Vậy là từ giờ cưng chỉ cần gõ "xemphim" trên Terminal/iTerm/Kitty của cưng là ào ào liền
- **Windows (PowerShell):** `go build -o xemphim.exe`
Sau đó add folder chứa file này vào PATH hoặc copy vào C:\Windows

---

## ⌨️ Phím tắt khi xem (MPV)
| Phím | Tác dụng |
|------|----------|
| `Space` | Tạm dừng / Xem tiếp |
| `M` | Tắt/Mở tiếng |
| `F` | Bật/Tắt Fullscreen |
| `Q` | Thoát phim quay lại CLI |
| `Mũi tên Trái/Phải` | Tua phim (-5s / +5s) |

---

## LƯU Ý
Để hiện icon đẹp như trên terminal của cj thì mấy cưng nên xài terminal xịn như iTerm2 hoặc Kitty nho, và nhớ cài Nerd Fonts nèk!!!

---
# CHÚC MÍ CƯNG XEM PHIM ZUI