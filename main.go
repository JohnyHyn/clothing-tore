package main

import (
	_ "clothing-store/docs"
	"clothing-store/internal/db"
	"clothing-store/internal/handler"
	"clothing-store/internal/middleware"
	"clothing-store/internal/service"
	"log"
	"net/http"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Clothing Store API
// @version         1.0
// @description     This is a sample server for Clothing Store.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	database := db.Connect()
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	productService := &service.ProductService{DB: database}
	productHandler := &handler.ProductHandler{
		ProductService: productService,
	}
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetProducts(w, r)
		case http.MethodPost:
			// Protect CreateProduct
			middleware.AuthMiddleware(productHandler.CreateProduct)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetProductByID(w, r)
		case http.MethodPut:
			middleware.AuthMiddleware(productHandler.UpdateProduct)(w, r)
		case http.MethodDelete:
			middleware.AuthMiddleware(productHandler.DeleteProduct)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	orderService := &service.OrderService{DB: database}
	orderHandler := &handler.OrderHandler{OrderService: orderService}

	authService := &service.AuthService{DB: database}
	authHandler := &handler.AuthHandler{AuthService: authService}

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authHandler.Login(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	paymentService := &service.PaymentService{DB: database}
	paymentHandler := &handler.PaymentHandler{PaymentService: paymentService}

	voucherService := &service.VoucherService{DB: database}
	voucherHandler := &handler.VoucherHandler{VoucherService: voucherService}

	shippingService := &service.ShippingService{DB: database}
	shippingHandler := &handler.ShippingHandler{ShippingService: shippingService}

	http.HandleFunc("/payments/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.AuthMiddleware(paymentHandler.CreatePayment)(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/payments/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			paymentHandler.Webhook(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	// Payment history endpoint
	http.HandleFunc("/payments/history", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				paymentHandler.GetHistory(w, r)
			}))(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	// VOUCHER ROUTES
	http.HandleFunc("/vouchers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			voucherHandler.ListVouchers(w, r)
		case http.MethodPost:
			middleware.AuthMiddleware(voucherHandler.CreateVoucher)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/vouchers/validate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			voucherHandler.ValidateVoucher(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/vouchers/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			voucherHandler.GetVoucher(w, r)
		} else if r.Method == http.MethodPut {
			middleware.AuthMiddleware(voucherHandler.UpdateVoucher)(w, r)
		} else if r.Method == http.MethodDelete {
			middleware.AuthMiddleware(voucherHandler.DeleteVoucher)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// SHIPPING ROUTES
	http.HandleFunc("/shippings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Get by order ID
			middleware.AuthMiddleware(shippingHandler.GetShippingByOrder)(w, r)
		case http.MethodPost:
			middleware.AuthMiddleware(shippingHandler.CreateShipping)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/shippings/track", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			shippingHandler.TrackShipping(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/shippings/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			middleware.AuthMiddleware(shippingHandler.UpdateShippingStatus)(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/shippings/history", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			middleware.AuthMiddleware(shippingHandler.GetShippingHistory)(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/shippings/calculate-fee", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			shippingHandler.CalculateShippingFee(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/shippings/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/tracking") && r.Method == http.MethodPut {
			middleware.AuthMiddleware(shippingHandler.UpdateTrackingCode)(w, r)
			return
		}
		http.NotFound(w, r)
	})

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Protect ListOrders with JWT
			middleware.AuthMiddleware(orderHandler.ListOrders)(w, r)
		case http.MethodPost:
			// Protect CreateOrder with JWT
			middleware.AuthMiddleware(orderHandler.CreateOrder)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/cancel") && r.Method == http.MethodPut {
			// Protect CancelOrder with JWT
			middleware.AuthMiddleware(orderHandler.CancelOrder)(w, r)
			return
		}

		if strings.HasSuffix(r.URL.Path, "/pay") && r.Method == http.MethodPut {
			// Protect PayOrder with JWT
			middleware.AuthMiddleware(orderHandler.PayOrder)(w, r)
			return
		}

		if r.Method == http.MethodGet {
			// Protect GetOrder with JWT
			middleware.AuthMiddleware(orderHandler.GetOrder)(w, r)
			return
		}
		if r.Method == http.MethodPut && strings.HasSuffix(r.URL.Path, "/refund") {
			// Protect RefundOrder with JWT
			middleware.AuthMiddleware(orderHandler.RefundOrder)(w, r)
			return
		}

		http.NotFound(w, r)
	})

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
