# Quick Reference - Orders API

## 🚀 Quick Start

### 1. Login & Get Token
```bash
POST /login
{
  "email": "admin@shop.com",
  "password": "123456"
}
```

### 2. Use Token in Headers
```
Authorization: Bearer {your_token}
```

---

## 📋 Endpoints Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/login` | ❌ | Login và lấy JWT token |
| GET | `/orders` | ✅ | List orders + pagination + filters |
| GET | `/orders/{id}` | ✅ | Chi tiết order + items |
| POST | `/orders` | ✅ | Tạo order mới |
| PUT | `/orders/{id}/pay` | ✅ | Thanh toán order |
| PUT | `/orders/{id}/cancel` | ✅ | Hủy order |
| PUT | `/orders/{id}/refund` | ✅ | Hoàn tiền order |

---

## 🔍 Filters Available

### List Orders (`GET /orders`)

```
?page=1              # Trang số
?limit=10            # Số items/trang
?status=paid         # pending | paid | cancelled
?date_from=2026-01-01
?date_to=2026-01-31
```

### Examples

```bash
# Orders đã thanh toán
GET /orders?status=paid

# Orders trong tháng 1
GET /orders?date_from=2026-01-01&date_to=2026-01-31

# Kết hợp: trang 2, pending orders, từ 10/01
GET /orders?page=2&status=pending&date_from=2026-01-10
```

---

## 💡 PowerShell Examples

```powershell
# 1. Login
$login = @{ email = "admin@shop.com"; password = "123456" } | ConvertTo-Json
$res = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method POST -Body $login -ContentType "application/json"
$token = $res.token

# 2. Set headers
$headers = @{ Authorization = "Bearer $token" }

# 3. List orders với filter
$orders = Invoke-RestMethod -Uri "http://localhost:8080/orders?status=paid" -Headers $headers
$orders.data
$orders.pagination

# 4. Get order detail
$detail = Invoke-RestMethod -Uri "http://localhost:8080/orders/1" -Headers $headers
$detail.data.items

# 5. Create order
$newOrder = @{
    customer_name = "Test"
    customer_phone = "0909123456"
    items = @(@{ product_id = 1; quantity = 2 })
} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/orders" -Method POST -Headers $headers -Body $newOrder -ContentType "application/json"
```

---

## 🔑 Key Features

✅ **JWT Authentication** - Tất cả endpoints orders protected  
✅ **Pagination** - page + limit  
✅ **Filter by Status** - pending, paid, cancelled  
✅ **Filter by Date** - date_from, date_to  
✅ **Order Detail** - Full items với product_name, price, quantity  
✅ **Error Handling** - 401 (Unauthorized), 404 (Not Found), 500 (Server Error)

---

## ⚠️ Common Errors

### 401 Unauthorized
```
authorization header required
```
**Fix:** Thêm `Authorization: Bearer {token}` vào header

### 404 Not Found
```
order not found
```
**Fix:** Kiểm tra order ID có tồn tại không

### Invalid Token
```
invalid token
```
**Fix:** Login lại để lấy token mới
