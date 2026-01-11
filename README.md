# Clothing Store API

## 📖 Giới thiệu
Backend API cho hệ thống cửa hàng quần áo (Clothing Store), được viết bằng Go (Golang). Hệ thống cung cấp các chức năng quản lý sản phẩm, đơn hàng, thanh toán, vận chuyển và mã giảm giá.

## 🚀 Tính năng chính
- **Xác thực người dùng (Auth)**: Đăng ký, Đăng nhập (JWT).
- **Quản lý sản phẩm (Products)**: Thêm, sửa, xóa, xem danh sách sản phẩm.
- **Quản lý đơn hàng (Orders)**: Tạo đơn hàng, xem lịch sử đơn hàng, hủy đơn, thanh toán.
- **Mã giảm giá (Vouchers)**: Quản lý mã giảm giá, tính toán giảm giá, giới hạn sử dụng, xác thực voucher.
- **Thanh toán (Payments)**: Tích hợp thanh toán, lịch sử giao dịch, webhook.
- **Vận chuyển (Shipping)**: Quản lý thông tin vận chuyển, tính phí, theo dõi đơn hàng (tracking), cập nhật trạng thái vận chuyển.

## 🛠 Công nghệ sử dụng
- **Ngôn ngữ**: Go 1.24
- **Database**: MySQL
- **Thư viện chính**:
    - `github.com/go-sql-driver/mysql`: MySQL driver
    - `github.com/golang-jwt/jwt/v5`: JSON Web Tokens
    - `golang.org/x/crypto`: Hashing mật khẩu

## ⚙️ Cài đặt và Chạy dự án

### 1. Yêu cầu
- Go 1.24 trở lên
- MySQL Server

### 2. Clone dự án
```bash
git clone https://github.com/yourusername/clothing-store.git
cd clothing-store
```

### 3. Cấu hình Database
Tạo database `clothing_store` trong MySQL.
Cấu hình các biến môi trường (Environment Variables). Bạn có thể set trực tiếp trong terminal hoặc tạo file script.

Ví dụ (PowerShell):
```powershell
$Env:DB_USER="clothing_app"
$Env:DB_PASSWORD="StrongPassword@123"
$Env:DB_HOST="127.0.0.1"
$Env:DB_PORT="3306"
$Env:DB_NAME="clothing_store"
```

Ví dụ (Bash):
```bash
export DB_USER="clothing_app"
export DB_PASSWORD="StrongPassword@123"
export DB_HOST="127.0.0.1"
export DB_PORT="3306"
export DB_NAME="clothing_store"
```

### 4. Chạy Migration và Seed Data
Dự án có sẵn các script để khởi tạo database và dữ liệu mẫu. Hãy chạy lần lượt các lệnh sau:

**Bước 1: Tạo các bảng (Migrations)**
```bash
go run cmd/migrate_features/main.go
# và
go run cmd/migrate_payment/main.go
```

**Bước 2: Tạo dữ liệu mẫu (Seed)**
```bash
go run cmd/seed_user/main.go
```

### 5. Chạy Server
```bash
go run main.go
```
Server sẽ chạy tại `http://localhost:8080` (mặc định).

## 📚 Tài liệu API
Chi tiết về các API endpoint và tính năng mới có thể tìm thấy trong thư mục `pkg/`:
- [Tổng quan tính năng mới](pkg/FEATURES_SUMMARY.md)
- [Chi tiết triển khai](pkg/IMPLEMENTATION_SUMMARY.md)
- [Tài liệu API Orders](pkg/API_ORDERS.md)
- [Tài liệu API Orders Complete](pkg/API_ORDERS_COMPLETE.md)
- [Quick Reference](pkg/QUICK_REFERENCE.md)

## 🤝 Đóng góp
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
