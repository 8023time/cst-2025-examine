package main

import (
	"cst-2025-examine-go/API"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	// 注册路由处理函数
	mux.HandleFunc("/articles", API.GetAllArticlesHandler)
	mux.HandleFunc("/articles/delete/", API.DeleteArticlesHandler)

	mux.HandleFunc("/user/login", API.LoginHandler)
	mux.HandleFunc("/user/loginout", API.LogoutHandler)
	mux.HandleFunc("/user/register", API.RegisterHandler)
	mux.HandleFunc("/user/articles", API.UserArticlesHandler)
	mux.HandleFunc("/user/articles/create", API.CreateArticleHandler)
	mux.HandleFunc("/user/articles/edit", API.EditArticleHandler)
	mux.HandleFunc("/user/profile", API.ProfileHandler)

	mux.HandleFunc("/admin/login", API.AdminLoginHandler)
	mux.HandleFunc("/admin/create-admin", API.CreateAdminHandler)
	mux.HandleFunc("/admin/users", API.GetUsersHandler)
	mux.HandleFunc("/admin/profile", API.AdminProfileHandler)
	mux.HandleFunc("/admin/users/delete/", API.DeleteUserHandler)
	mux.HandleFunc("/admin/users/profile", API.UsersProfileHandler)

	// 配置 CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"}, // 替换为你的前端应用的实际地址
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})

	// 应用 CORS 中间件
	handler := c.Handler(mux)

	port := ":8081"
	fmt.Println("Server is running on port: " + port)
	http.ListenAndServe(port, handler)
}
