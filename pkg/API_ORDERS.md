# Orders API - List Orders với Pagination

## Endpoint: List Orders

**GET** `/orders`

Lấy danh sách đơn hàng với hỗ trợ phân trang (pagination).

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Số trang hiện tại |
| `limit` | integer | 10 | Số lượng orders mỗi trang (max: 100) |

### Ví dụ Request

```bash
# Lấy trang đầu tiên với 10 orders (mặc định)
GET http://localhost:8080/orders

# Lấy trang 2, mỗi trang 20 orders
GET http://localhost:8080/orders?page=2&limit=20

# Lấy trang 3, mỗi trang 5 orders
GET http://localhost:8080/orders?page=3&limit=5
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
      "created_at": "2026-01-11T14:30:00Z",
      "items": []
    },
    {
      "id": 2,
      "customer_name": "Tran Thi B",
      "customer_phone": "0907654321",
      "total_price": 750000,
      "status": "pending",
      "created_at": "2026-01-11T15:00:00Z",
      "items": []
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

### Response Fields

#### Data
- `id`: ID của order
- `customer_name`: Tên khách hàng
- `customer_phone`: Số điện thoại khách hàng
- `total_price`: Tổng giá trị đơn hàng
- `status`: Trạng thái đơn hàng (`pending`, `paid`, `cancelled`)
- `created_at`: Thời gian tạo đơn hàng
- `items`: Danh sách sản phẩm trong đơn hàng (trống trong list view)

#### Pagination
- `page`: Trang hiện tại
- `limit`: Số lượng items mỗi trang
- `total`: Tổng số orders
- `total_pages`: Tổng số trang

### Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 500 | Internal Server Error |

## Test với curl

```bash
# Test pagination
curl "http://localhost:8080/orders?page=1&limit=5"

# Test với default parameters
curl "http://localhost:8080/orders"
```

## Test với PowerShell

```powershell
# Lấy danh sách orders
Invoke-WebRequest -Uri "http://localhost:8080/orders?page=1&limit=10" -Method GET

# Parse JSON response
$response = Invoke-RestMethod -Uri "http://localhost:8080/orders?page=1&limit=10" -Method GET
$response.data
$response.pagination
```
