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

	// CORS Middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	productService := &service.ProductService{DB: database}
	productHandler := &handler.ProductHandler{
		ProductService: productService,
	}

	// Helpers for RBAC
	adminStaff := func(next http.HandlerFunc) http.HandlerFunc {
		return middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			middleware.RoleMiddleware([]string{"admin", "staff"}, next)(w, r)
		})
	}
	authUser := middleware.AuthMiddleware

	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetProducts(w, r)
		case http.MethodPost:
			adminStaff(productHandler.CreateProduct)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetProductByID(w, r)
		case http.MethodPut:
			adminStaff(productHandler.UpdateProduct)(w, r)
		case http.MethodDelete:
			adminStaff(productHandler.DeleteProduct)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	orderService := &service.OrderService{DB: database}
	orderHandler := &handler.OrderHandler{OrderService: orderService}

	authService := &service.AuthService{DB: database}
	authHandler := &handler.AuthHandler{AuthService: authService}

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authHandler.Login(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authHandler.Register(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/forgot-password", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authHandler.ForgotPassword(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/reset-password", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authHandler.ResetPassword(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	cartService := &service.CartService{DB: database}
	cartHandler := &handler.CartHandler{CartService: cartService}

	mux.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authUser(cartHandler.GetCart)(w, r)
		case http.MethodPost:
			authUser(cartHandler.AddToCart)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/cart/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			authUser(cartHandler.UpdateQuantity)(w, r)
		case http.MethodDelete:
			authUser(cartHandler.RemoveItem)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	paymentService := &service.PaymentService{DB: database}
	paymentHandler := &handler.PaymentHandler{PaymentService: paymentService}

	voucherService := &service.VoucherService{DB: database}
	voucherHandler := &handler.VoucherHandler{VoucherService: voucherService}

	shippingService := &service.ShippingService{DB: database}
	shippingHandler := &handler.ShippingHandler{ShippingService: shippingService}

	mux.HandleFunc("/payments/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authUser(paymentHandler.CreatePayment)(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/payments/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			paymentHandler.Webhook(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	// Payment history endpoint
	mux.HandleFunc("/payments/history", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			adminStaff(paymentHandler.GetHistory)(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	// VOUCHER ROUTES
	mux.HandleFunc("/vouchers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			voucherHandler.ListVouchers(w, r)
		case http.MethodPost:
			adminStaff(voucherHandler.CreateVoucher)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/vouchers/validate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			voucherHandler.ValidateVoucher(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/vouchers/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			voucherHandler.GetVoucher(w, r)
		} else if r.Method == http.MethodPut {
			adminStaff(voucherHandler.UpdateVoucher)(w, r)
		} else if r.Method == http.MethodDelete {
			adminStaff(voucherHandler.DeleteVoucher)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// SHIPPING ROUTES
	mux.HandleFunc("/shippings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Get by order ID
			authUser(shippingHandler.GetShippingByOrder)(w, r)
		case http.MethodPost:
			authUser(shippingHandler.CreateShipping)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/shippings/track", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			shippingHandler.TrackShipping(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/shippings/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			adminStaff(shippingHandler.UpdateShippingStatus)(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/shippings/history", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			adminStaff(shippingHandler.GetShippingHistory)(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/shippings/calculate-fee", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			shippingHandler.CalculateShippingFee(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/shippings/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/tracking") && r.Method == http.MethodPut {
			adminStaff(shippingHandler.UpdateTrackingCode)(w, r)
			return
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/my/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authUser(orderHandler.ListMyOrders)(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			adminStaff(orderHandler.ListOrders)(w, r)
		case http.MethodPost:
			authUser(orderHandler.CreateOrder)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/cancel") && r.Method == http.MethodPut {
			authUser(orderHandler.CancelOrder)(w, r)
			return
		}

		if strings.HasSuffix(r.URL.Path, "/pay") && r.Method == http.MethodPut {
			authUser(orderHandler.PayOrder)(w, r)
			return
		}

		if r.Method == http.MethodGet {
			authUser(orderHandler.GetOrder)(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			authUser(orderHandler.DeleteOrder)(w, r)
			return
		}
		if r.Method == http.MethodPut && strings.HasSuffix(r.URL.Path, "/approve") {
			adminStaff(orderHandler.ApproveOrder)(w, r)
			return
		}
		if r.Method == http.MethodPut && strings.HasSuffix(r.URL.Path, "/refund") {
			adminStaff(orderHandler.RefundOrder)(w, r)
			return
		}

		http.NotFound(w, r)
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// You can add middleware here if needed
		handler.UploadHandler(w, r)
	})

	// Serve static files from uploads directory
	fs := http.FileServer(http.Dir("./uploads"))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(mux)))
}
