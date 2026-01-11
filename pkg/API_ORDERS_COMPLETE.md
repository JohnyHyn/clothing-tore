# Orders API - Complete Documentation

## 🔐 Authentication

Tất cả các endpoints orders đều yêu cầu JWT authentication.

### Get JWT Token

**POST** `/login`

```json
{
  "email": "admin@shop.com",
  "password": "123456"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Sử dụng Token

Thêm token vào header của mọi request:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

---

## 📋 1. List Orders với Pagination & Filters

**GET** `/orders`

🔐 **Requires Authentication**

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Số trang hiện tại |
| `limit` | integer | 10 | Số lượng orders mỗi trang (max: 100) |
| `status` | string | - | Filter theo status: `pending`, `paid`, `cancelled` |
| `date_from` | string | - | Filter từ ngày (format: `YYYY-MM-DD HH:MM:SS`) |
| `date_to` | string | - | Filter đến ngày (format: `YYYY-MM-DD HH:MM:SS`) |

### Examples

```bash
# Lấy tất cả orders (trang 1, 10 items)
GET /orders
Authorization: Bearer {token}

# Filter theo status = paid
GET /orders?status=paid
Authorization: Bearer {token}

# Filter theo date range
GET /orders?date_from=2026-01-01&date_to=2026-01-31
Authorization: Bearer {token}

# Kết hợp: trang 2, status pending, từ ngày 2026-01-10
GET /orders?page=2&limit=20&status=pending&date_from=2026-01-10
Authorization: Bearer {token}
```

### Response

```json
{
  "data": [
    {
      "id": 1,
      "customer_name": "Nguyen Van A",
      "customer_phone": "0901234567",
      "total_price": 500000,
      "status": "paid",
      "created_at": "2026-01-11 14:30:00",
      "items": []
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  },
  "filters": {
    "status": "paid",
    "date_from": "2026-01-10",
    "date_to": ""
  }
}
```

---

## 📦 2. Get Order Detail (with Items)

**GET** `/orders/{id}`

🔐 **Requires Authentication**

Lấy thông tin chi tiết của một order, bao gồm danh sách items đầy đủ.

### Example

```bash
GET /orders/1
Authorization: Bearer {token}
```

### Response

```json
{
  "data": {
    "id": 1,
    "customer_name": "Nguyen Van A",
    "customer_phone": "0901234567",
    "total_price": 500000,
    "status": "paid",
    "created_at": "2026-01-11 14:30:00",
    "items": [
      {
        "product_id": 10,
        "product_name": "Áo thun nam",
        "quantity": 2,
        "price": 150000
      },
      {
        "product_id": 15,
        "product_name": "Quần jean",
        "quantity": 1,
        "price": 200000
      }
    ]
  }
}
```

### Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 401 | Unauthorized (missing or invalid token) |
| 404 | Order not found |
| 500 | Internal Server Error |

---

## ➕ 3. Create Order

**POST** `/orders`

🔐 **Requires Authentication**

### Request Body

```json
{
  "customer_name": "Nguyen Van A",
  "customer_phone": "0901234567",
  "items": [
    {
      "product_id": 10,
      "quantity": 2
    },
    {
      "product_id": 15,
      "quantity": 1
    }
  ]
}
```

### Response

```json
{
  "id": 1,
  "customer_name": "Nguyen Van A",
  "customer_phone": "0901234567",
  "total_price": 500000,
  "status": "pending",
  "items": [...]
}
```

---

## 💳 4. Pay Order

**PUT** `/orders/{id}/pay`

🔐 **Requires Authentication**

Đánh dấu order là đã thanh toán.

### Example

```bash
PUT /orders/1/pay
Authorization: Bearer {token}
```

### Response

```
Order paid successfully
```

---

## ❌ 5. Cancel Order

**PUT** `/orders/{id}/cancel`

🔐 **Requires Authentication**

Hủy order (chỉ hủy được order có status = `pending`).

### Example

```bash
PUT /orders/1/cancel
Authorization: Bearer {token}
```

### Response

```
Order cancelled successfully
```

---

## 💰 6. Refund Order

**PUT** `/orders/{id}/refund`

🔐 **Requires Authentication**

Hoàn tiền cho order (chỉ refund được order có status = `paid`).

### Request Body

```json
{
  "reason": "Sản phẩm bị lỗi"
}
```

### Example

```bash
PUT /orders/1/refund
Authorization: Bearer {token}
Content-Type: application/json

{
  "reason": "Sản phẩm bị lỗi"
}
```

### Response

```
Order refunded successfully
```

---

## 🧪 Testing với PowerShell

### 1. Login để lấy token

```powershell
$loginBody = @{
    email = "admin@shop.com"
    password = "123456"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/login" `
    -Method POST `
    -Body $loginBody `
    -ContentType "application/json"

$token = $response.token
Write-Host "Token: $token"
```

### 2. List orders với filters

```powershell
$headers = @{
    Authorization = "Bearer $token"
}

# List tất cả orders
$orders = Invoke-RestMethod -Uri "http://localhost:8080/orders" `
    -Method GET `
    -Headers $headers

# Filter theo status
$paidOrders = Invoke-RestMethod -Uri "http://localhost:8080/orders?status=paid&page=1&limit=5" `
    -Method GET `
    -Headers $headers

$paidOrders.data
$paidOrders.pagination
```

### 3. Get order detail

```powershell
$orderDetail = Invoke-RestMethod -Uri "http://localhost:8080/orders/1" `
    -Method GET `
    -Headers $headers

$orderDetail.data
$orderDetail.data.items
```

### 4. Create order

```powershell
$newOrderBody = @{
    customer_name = "Test Customer"
    customer_phone = "0909123456"
    items = @(
        @{
            product_id = 1
            quantity = 2
        }
    )
} | ConvertTo-Json

$newOrder = Invoke-RestMethod -Uri "http://localhost:8080/orders" `
    -Method POST `
    -Headers $headers `
    -Body $newOrderBody `
    -ContentType "application/json"

$newOrder
```

---

## 🧪 Testing với curl

### 1. Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@shop.com","password":"123456"}'

# Save token
export TOKEN="eyJhbGciOiJIUzI1NiIs..."
```

### 2. List orders với filters

```bash
# List all
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/orders

# Filter by status
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/orders?status=paid&page=1&limit=10"

# Filter by date
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/orders?date_from=2026-01-01&date_to=2026-01-31"
```

### 3. Get order detail

```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/orders/1
```

### 4. Create order

```bash
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "Test Customer",
    "customer_phone": "0909123456",
    "items": [
      {"product_id": 1, "quantity": 2}
    ]
  }'
```

---

## 🔒 Error Responses

### 401 Unauthorized

Khi không có token hoặc token không hợp lệ:

```json
authorization header required
```

hoặc

```json
invalid token
```

### Cách fix

- Đảm bảo header `Authorization: Bearer {token}` có trong request
- Token phải còn hạn (kiểm tra expiration)
- Token phải được lấy từ endpoint `/login`

---

## 📊 Order Status Flow

```
pending → paid → (có thể refund) → cancelled
   ↓
cancelled (từ pending)
```

- **pending**: Order mới tạo, chưa thanh toán
- **paid**: Đã thanh toán
- **cancelled**: Đã hủy (từ pending hoặc refund từ paid)

---

## ✅ Checklist Implementation

- ✅ **Filter orders theo status** - Query param `status`
- ✅ **Filter orders theo date** - Query params `date_from`, `date_to`
- ✅ **Order detail + items** - Endpoint `/orders/{id}` trả về đầy đủ items
- ✅ **JWT Auth protection** - Tất cả endpoints orders đều require JWT token
- ✅ **Pagination** - Hỗ trợ `page` và `limit`
- ✅ **Clear error messages** - 401, 404, 500 với messages rõ ràng
