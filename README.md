# 🎬 OPHIM CLI – TÔI YÊU PHÂU LỊM (OPEN ENGINE EDITION) 🍿

---

## ⚠️ TUYÊN BỐ MIỄN TRỪ TRÁCH NHIỆM (DISCLAIMER)

Dự án này được tạo ra với mục đích **nghiên cứu kỹ thuật** (*Go + MPV integration*) và vọc vạch CLI.

* Tác giả **KHÔNG cung cấp** bất kỳ nội dung, video, hay API lấy phim nào.
* Phần mềm này chỉ là một **Media Client (vỏ)**.
* Người dùng cuối chịu trách nhiệm hoàn toàn về việc tìm kiếm nguồn nội dung (Content) và tuân thủ bản quyền của nguồn đó.
* Tác giả không chịu trách nhiệm cho bất kỳ rắc rối pháp lý nào phát sinh từ phía người dùng.

> Dùng hay không là quyền của cưng, đi tù hay không là chuyện của cưng. 😌

---

## ✨ Tính năng "Ba trợn"

* 🔍 **Search Engine** – Tìm phim thông qua API cưng tự cấu hình.
* 🎞️ **Multi-Source** – Hỗ trợ phim bộ, phim lẻ, phim "vùng kín" (tùy thuộc vào nguồn cưng có).
* 📜 **Session Recovery** – Nhắc cưng xem tiếp tập đang dang dở (*Cục dzàng có muốn coi tiếp hem?*).
* 🎨 **Aesthetic UI** – Hệ màu Catppuccin siêu mịn, hỗ trợ cuộn chuột xem mô tả.
* 🚀 **High Performance** – Viết bằng Go, mượt hơn cả cách người yêu cũ trở mặt.

---

# 🏗 Cấu hình "Nguồn nước" (API Configuration)

Vì mục đích bảo mật và phủi bỏ trách nhiệm, project này **KHÔNG đi kèm API**.
Cưng cần chuẩn bị một Server trả về JSON theo đúng định dạng bên dưới.

---

## 1️⃣ Thiết lập file `.env`

Tạo file `.env` ngay tại thư mục root:

```env
# Link API gốc (Base URL)
API_BASE_URL="https://your-hidden-provider.com/api"

# Endpoint tìm kiếm (ví dụ: /v1/search?keyword=)
SEARCH_PATH="/v1/search?keyword="

# Endpoint chi tiết phim (ví dụ: /v1/movie/)
DETAIL_PATH="/v1/movie/"
```

---

## 2️⃣ JSON Schema yêu cầu

Để App có thể parse dữ liệu, API của cưng phải trả về đúng cấu trúc này:

### 🔎 Search Result

```json
{
  "status": true,
  "items": [
    {
      "name": "Tên phim cực căng",
      "slug": "ten-phim-cuc-cang",
      "origin_name": "Hardcore Movie Name",
      "year": 2024
    }
  ]
}
```

---

### 🎬 Movie Details

```json
{
  "movie": {
    "name": "Tên phim",
    "content": "Mô tả nội dung phim cực sến...",
    "episodes": [
      {
        "server_name": "Server Vietsub",
        "server_data": [
          { 
            "name": "Tập 1", 
            "link_m3u8": "https://stream.link/playlist.m3u8"
          }
        ]
      }
    ]
  }
}
```

---

> ⚠️ **Lưu ý:**
> Mẹ sẽ **KHÔNG trả lời** bất kỳ tin nhắn/email nào hỏi về việc *"xin link phim"*.
> Mọi gói tin hỏi về API lậu sẽ bị hốt lên C50 ngay lập tức. 🚓

---

# 🛠 Yêu cầu hệ thống

Cái "động cơ" chính để phát phim là **MPV**. Cưng phải cài nó trước.

### 🍎 MacOS

```bash
brew install mpv
```

### 🐧 Linux

```bash
sudo apt update && sudo apt install mpv
```

### 🪟 Windows

* Tải bản Zip tại: [https://mpv.io](https://mpv.io)
* Giải nén và add vào `PATH` hệ thống.

---

# 🚀 Cách Build & Chạy

## ▶️ Cách 1: Chạy trực tiếp

```bash
go run main.go
```

---

## 🏗 Cách 2: Build thành lệnh hệ thống (Khuyên dùng)

### MacOS / Linux

```bash
go build -o xemphim
sudo mv xemphim /usr/local/bin/
```

### Windows (PowerShell)

```powershell
go build -o xemphim.exe
# Sau đó add folder này vào PATH hệ thống
```

---

# ⌨️ Phím tắt khi xem (MPV)

| Phím    | Tác dụng                |
| ------- | ----------------------- |
| `Space` | Tạm dừng / Xem tiếp     |
| `M`     | Tắt/Mở tiếng            |
| `F`     | Bật/Tắt Fullscreen      |
| `Q`     | Thoát phim quay lại CLI |
| `← / →` | Tua phim (-5s / +5s)    |

---

# 📜 LƯU Ý

Để hiện icon đẹp như trên terminal của "chị":

* Nên dùng terminal xịn như **iTerm2** hoặc **Kitty**
* Cài **Nerd Fonts** (gợi ý: *JetBrainsMono Nerd Font*)

Không có font là nó hiện ô vuông "vô tri" ráng chịu nha cục dzàng 😌

---

# ❤️ CHÚC MÍ CƯNG HỌC CODE VUI VẺ

*(VÀ XEM PHIM TRONG SỰ KÍN ĐÁO!)* 🍿
