package routes

import (
	"erp/controllers"
	"erp/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/register", controllers.Register)
		}

		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			authenticated.GET("/auth/user", controllers.GetCurrentUser)

			users := authenticated.Group("/users")
			{
				users.GET("", controllers.GetUsers)
				users.GET("/:id", controllers.GetUser)
				users.POST("", controllers.CreateUser)
				users.PUT("/:id", controllers.UpdateUser)
				users.DELETE("/:id", controllers.DeleteUser)
			}

			roles := authenticated.Group("/roles")
			{
				roles.GET("", controllers.GetRoles)
				roles.GET("/:id", controllers.GetRole)
				roles.POST("", controllers.CreateRole)
				roles.PUT("/:id", controllers.UpdateRole)
				roles.DELETE("/:id", controllers.DeleteRole)
				roles.PUT("/:id/menus", controllers.AssignMenus)
			}

			menus := authenticated.Group("/menus")
			{
				menus.GET("", controllers.GetMenus)
				menus.GET("/:id", controllers.GetMenu)
				menus.POST("", controllers.CreateMenu)
				menus.PUT("/:id", controllers.UpdateMenu)
				menus.DELETE("/:id", controllers.DeleteMenu)
			}

			products := authenticated.Group("/products")
			{
				products.GET("", controllers.GetProducts)
				products.GET("/:id", controllers.GetProduct)
				products.POST("", controllers.CreateProduct)
				products.PUT("/:id", controllers.UpdateProduct)
				products.DELETE("/:id", controllers.DeleteProduct)
			}

			customers := authenticated.Group("/customers")
			{
				customers.GET("", controllers.GetCustomers)
				customers.GET("/:id", controllers.GetCustomer)
				customers.POST("", controllers.CreateCustomer)
				customers.PUT("/:id", controllers.UpdateCustomer)
				customers.DELETE("/:id", controllers.DeleteCustomer)
			}

			suppliers := authenticated.Group("/suppliers")
		{
			suppliers.GET("", controllers.GetSuppliers)
			suppliers.GET("/:id", controllers.GetSupplier)
			suppliers.POST("", controllers.CreateSupplier)
			suppliers.PUT("/:id", controllers.UpdateSupplier)
			suppliers.DELETE("/:id", controllers.DeleteSupplier)
		}

		purchaseOrders := authenticated.Group("/purchase-orders")
		{
			purchaseOrders.GET("", controllers.GetPurchaseOrders)
			purchaseOrders.GET("/:id", controllers.GetPurchaseOrder)
			purchaseOrders.POST("", controllers.CreatePurchaseOrder)
			purchaseOrders.PUT("/:id", controllers.UpdatePurchaseOrder)
			purchaseOrders.PUT("/:id/status", controllers.UpdatePurchaseOrderStatus)
			purchaseOrders.DELETE("/:id", controllers.DeletePurchaseOrder)
		}

		salesOrders := authenticated.Group("/sales-orders")
		{
			salesOrders.GET("", controllers.GetSalesOrders)
			salesOrders.GET("/:id", controllers.GetSalesOrder)
			salesOrders.POST("", controllers.CreateSalesOrder)
			salesOrders.PUT("/:id", controllers.UpdateSalesOrder)
			salesOrders.PUT("/:id/status", controllers.UpdateSalesOrderStatus)
			salesOrders.DELETE("/:id", controllers.DeleteSalesOrder)
		}
	}
}
}
