# 🎉 New Features Documentation

## 🎫 VOUCHER (Mã Giảm Giá)
## 💳 PAYMENT TRANSACTIONS (Giao Dịch Thanh Toán)
## 🚚 SHIPPING (Vận Chuyển)

---

# 🎫 1. VOUCHER API

## 📋 Endpoints Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/vouchers` | ❌ | List all vouchers |
| POST | `/vouchers` | ✅ | Create voucher (admin) |
| GET | `/vouchers/{code}` | ❌ | Get voucher by code |
| PUT | `/vouchers/{id}` | ✅ | Update voucher (admin) |
| DELETE | `/vouchers/{id}` | ✅ | Delete voucher (admin) |
| POST | `/vouchers/validate` | ❌ | Validate voucher |

---

## 1.1 Create Voucher (Admin)

**POST** `/vouchers`  
🔐 **Requires Authentication**

### Request Body

```json
{
  "code": "SUMMER2026",
  "description": "Giảm giá mùa hè 2026",
  "type": "percentage",
  "value": 20,
  "min_order": 100000,
  "max_discount": 50000,
  "usage_limit": 100,
  "start_date": "2026-06-01T00:00:00Z",
  "end_date": "2026-08-31T23:59:59Z",
  "is_active": true
}
```

### Field Descriptions

- `code`: Mã voucher (unique, tự động uppercase)
- `type`: `"percentage"` hoặc `"fixed"`
- `value`: Giá trị giảm (% nếu type=percentage, số tiền nếu type=fixed)
- `min_order`: Giá trị đơn hàng tối thiểu
- `max_discount`: Giảm tối đa (chỉ cho percentage)
- `usage_limit`: Số lần dùng tối đa
- `start_date`, `end_date`: Thời gian hiệu lực

---

## 1.2 Validate Voucher

**POST** `/vouchers/validate`  
❌ **No Authentication Required**

### Request

```json
{
  "code": "SUMMER2026",
  "order_amount": 150000
}
```

### Response (Success)

```json
{
  "valid": true,
  "voucher": {
    "id": 1,
    "code": "SUMMER2026",
    "type": "percentage",
    "value": 20,
    ...
  },
  "discount": 30000,
  "message": "Voucher is valid"
}
```

### Response (Error)

```
400 Bad Request
"voucher not found" / "voucher has expired" / "order amount below minimum"
```

---

## 1.3 List Vouchers

**GET** `/vouchers?page=1&limit=10&active_only=true`

### Query Parameters

- `page`: Số trang (default: 1)
- `limit`: Số items/trang (default: 10)
- `active_only`: `true` để chỉ lấy vouchers đang active

---

## 1.4 Get Voucher by Code

**GET** `/vouchers/SUMMER2026`

### Response

```json
{
  "data": {
    "id": 1,
    "code": "SUMMER2026",
    "description": "Giảm giá mùa hè 2026",
    "type": "percentage",
    "value": 20,
    "used_count": 45,
    "usage_limit": 100,
    ...
  }
}
```

---

# 💳 2. PAYMENT TRANSACTIONS API

## 📋 Endpoints Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/payments/create` | ✅ | Tạo payment request |
| POST | `/payments/webhook` | ❌ | Webhook từ payment gateway |
| GET | `/payments/history?order_id={id}` | ✅ | Lịch sử thanh toán |

---

## 2.1 Create Payment

**POST** `/payments/create`  
🔐 **Requires Authentication**

### Request

```json
{
  "order_id": 1,
  "provider": "momo"
}
```

### Response

```json
{
  "payment_url": "https://payment-gateway.com/pay?order=1&amount=500000",
  "transaction_id": "TRANS_1_1736551234567890"
}
```

---

## 2.2 Payment History

**GET** `/payments/history?order_id=1`  
🔐 **Requires Authentication**

### Response

```json
{
  "data": [
    {
      "id": 1,
      "order_id": 1,
      "transaction_id": "TRANS_1_...",
      "provider": "momo",
      "amount": 500000,
      "status": "success",
      "created_at": "2026-01-11T15:30:00Z"
    }
  ]
}
```

---

# 🚚 3. SHIPPING API

## 📋 Endpoints Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/shippings` | ✅ | Tạo shipping cho order |
| GET | `/shippings?order_id={id}` | ✅ | Get shipping by order |
| GET | `/shippings/track?tracking_code={code}` | ❌ | Track đơn hàng |
| PUT | `/shippings/status` | ✅ | Cập nhật trạng thái |
| GET | `/shippings/history?shipping_id={id}` | ✅ | Lịch sử vận chuyển |
| GET | `/shippings/calculate-fee?method={m}&city={c}` | ❌ | Tính phí ship |
| PUT | `/shippings/{id}/tracking` | ✅ | Cập nhật tracking code |

---

## 3.1 Create Shipping

**POST** `/shippings`  
🔐 **Requires Authentication**

### Request

```json
{
  "order_id": 1,
  "method": "express",
  "provider": "ghn",
  "tracking_code": "GHN123456789",
  "address": "123 Nguyễn Huệ",
  "city": "Hồ Chí Minh",
  "district": "Quận 1",
  "ward": "Phường Bến Nghé",
  "receiver_name": "Nguyễn Văn A",
  "receiver_phone": "0901234567",
  "note": "Giao giờ hành chính",
  "estimated_date": "2026-01-15T00:00:00Z"
}
```

### Response

```json
{
  "message": "Shipping created successfully",
  "data": {
    "id": 1,
    "order_id": 1,
    "fee": 50000,
    "status": "pending",
    ...
  }
}
```

---

## 3.2 Track Shipping

**GET** `/shippings/track?tracking_code=GHN123456789`  
❌ **No Authentication Required**

### Response

```json
{
  "shipping": {
    "id": 1,
    "tracking_code": "GHN123456789",
    "status": "in_transit",
    "receiver_name": "Nguyễn Văn A",
    ...
  },
  "history": [
    {
      "status": "in_transit",
      "location": "Đà Nẵng",
      "note": "Đang vận chuyển",
      "created_at": "2026-01-12T10:00:00Z"
    },
    {
      "status": "picked_up",
      "location": "Hà Nội",
      "note": "Đã lấy hàng",
      "created_at": "2026-01-11T16:00:00Z"
    }
  ]
}
```

---

## 3.3 Update Shipping Status

**PUT** `/shippings/status`  
🔐 **Requires Authentication**

### Request

```json
{
  "shipping_id": 1,
  "status": "delivered",
  "location": "Hồ Chí Minh",
  "note": "Giao hàng thành công"
}
```

### Shipping Statuses

- `pending`: Chờ xử lý
- `picked_up`: Đã lấy hàng
- `in_transit`: Đang vận chuyển
- `delivered`: Đã giao
- `failed`: Giao thất bại

---

## 3.4 Calculate Shipping Fee

**GET** `/shippings/calculate-fee?method=express&city=Hà Nội`  
❌ **No Authentication Required**

### Response

```json
{
  "method": "express",
  "city": "Hà Nội",
  "fee": 50000
}
```

### Shipping Methods

- `standard`: 30,000đ (3-5 ngày)
- `express`: 50,000đ (1-2 ngày)
- `same_day`: 80,000đ (trong ngày)

Phí tăng thêm 20,000đ cho các tỉnh xa (Hà Giang, Cao Bằng, Lai Châu, Lào Cai)

---

# 🧪 Testing Examples

## PowerShell Examples

### 1. Vouchers

```powershell
# Login first
$login = @{ email = "admin@shop.com"; password = "123456" } | ConvertTo-Json
$res = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method POST -Body $login -ContentType "application/json"
$token = $res.token
$headers = @{ Authorization = "Bearer $token" }

# Create voucher
$voucher = @{
    code = "SUMMER2026"
    description = "Giảm 20% mùa hè"
    type = "percentage"
    value = 20
    min_order = 100000
    max_discount = 50000
    usage_limit = 100
    start_date = "2026-06-01T00:00:00Z"
    end_date = "2026-08-31T23:59:59Z"
    is_active = $true
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/vouchers" -Method POST -Headers $headers -Body $voucher -ContentType "application/json"

# Validate voucher
$validate = @{
    code = "SUMMER2026"
    order_amount = 150000
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/vouchers/validate" -Method POST -Body $validate -ContentType "application/json"

# List vouchers
Invoke-RestMethod -Uri "http://localhost:8080/vouchers?active_only=true"
```

### 2. Shipping

```powershell
# Create shipping
$shipping = @{
    order_id = 1
    method = "express"
    provider = "ghn"
    tracking_code = "GHN123456789"
    address = "123 Nguyễn Huệ"
    city = "Hồ Chí Minh"
    district = "Quận 1"
    receiver_name = "Nguyễn Văn A"
    receiver_phone = "0901234567"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/shippings" -Method POST -Headers $headers -Body $shipping -ContentType "application/json"

# Track shipping
Invoke-RestMethod -Uri "http://localhost:8080/shippings/track?tracking_code=GHN123456789"

# Calculate shipping fee
Invoke-RestMethod -Uri "http://localhost:8080/shippings/calculate-fee?method=express&city=Hà Nội"

# Update status
$status = @{
    shipping_id = 1
    status = "in_transit"
    location = "Đà Nẵng"
    note = "Đang vận chuyển"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/shippings/status" -Method PUT -Headers $headers -Body $status -ContentType "application/json"
```

### 3. Payment History

```powershell
# Get payment history
Invoke-RestMethod -Uri "http://localhost:8080/payments/history?order_id=1" -Headers $headers
```

---

# 💾 Database Migration

Chạy migration để tạo tables:

```powershell
# Windows
go run ./cmd/migrate_features/main.go
```

Tables được tạo:
- `vouchers` - Danh sách vouchers
- `voucher_usage` - Lịch sử sử dụng vouchers
- `shippings` - Thông tin vận chuyển
- `shipping_history` - Lịch sử trạng thái vận chuyển

---

# 🔑 Key Features Summary

## Vouchers
✅ Hỗ trợ 2 loại: percentage & fixed  
✅ Usage limit (giới hạn số lần dùng)  
✅ Date range validation  
✅ Minimum order amount  
✅ Maximum discount (cho percentage)  
✅ Real-time validation

## Shipping
✅ Multiple providers (GHN, GHTK, Viettel Post)  
✅ Real-time tracking  
✅ Status history  
✅ Automatic fee calculation  
✅ Estimated delivery date  
✅ Receiver information

## Payment Transactions
✅ Multiple providers support  
✅ Transaction history  
✅ Webhook integration  
✅ Payment status tracking

---

# ⚠️ Next Steps

1. **Run Migration**
   ```
   go run ./cmd/migrate_features/main.go
   ```

2. **Start Server**
   ```
   go run main.go
   ```

3. **Test Endpoints** using PowerShell examples above

4. **Seed Data** (optional) - Create some test vouchers and shipping data
