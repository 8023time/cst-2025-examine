package API

import (
	"context"
	"cst-2025-examine-go/SQL"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Article struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Author   string `json:"author"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

var store = sessions.NewCookieStore([]byte("alongstory"))

func respondWithJSON(w http.ResponseWriter, statusCode int, status int, msg string, data interface{}) {
	log.Printf("响应: statusCode=%d, status=%d, msg=%s", statusCode, status, msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"status": status,
		"msg":    msg,
	}
	if data != nil {
		response["data"] = data
	}
	json.NewEncoder(w).Encode(response)
}

//	POST{
//		username,
//		password
//	}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondWithJSON(w, http.StatusMethodNotAllowed, -1, "不支持的请求方法", nil)
		return
	}

	var loginReq LoginRequest
	err := json.NewDecoder(req.Body).Decode(&loginReq)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, -1, "请求体解析失败: "+err.Error(), nil)
		fmt.Println("Request body parsing error:", err)
		return
	}

	if loginReq.Username == "" || loginReq.Password == "" {
		respondWithJSON(w, http.StatusBadRequest, -1, "用户名或密码不能为空", nil)
		return
	}

	// 获取数据库连接
	conn, err := SQL.ConnectSQL()
	if err != nil {
		http.Error(w, fmt.Sprintf("无法连接数据库, %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	// 执行 SQL 查询
	var userid int
	var stored_password string
	err = conn.QueryRow(
		context.Background(),
		"SELECT userid, password FROM t_user WHERE username=$1",
		loginReq.Username,
	).Scan(&userid, &stored_password)

	if err != nil {
		if err == pgx.ErrNoRows {
			respondWithJSON(w, http.StatusUnauthorized, -1, "登录失败: 用户名不存在", nil)
		} else {
			respondWithJSON(w, http.StatusInternalServerError, -1, "登录失败: 数据库查询错误", nil)
		}
		fmt.Println("Database query error:", err)
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(stored_password), []byte(loginReq.Password))
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, -1, "登录失败: 密码错误", nil)
		fmt.Println("Password mismatch error:", err)
		return
	}

	// 保存 session
	session, err := store.Get(req, "session-information")
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取 session: "+err.Error(), nil)
		fmt.Println("Session error:", err)
		return
	}
	session.Values["userid"] = userid
	session.Values["username"] = loginReq.Username

	err = session.Save(req, w)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法保存会话: "+err.Error(), nil)
		fmt.Println("Session save error:", err)
		return
	}

	data := map[string]interface{}{
		"username": loginReq.Username,
		"userid":   userid,
	}

	respondWithJSON(w, http.StatusOK, 0, "登录成功", data)
}

func LogoutHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session-information")
	session.Options.MaxAge = -1
	session.Save(req, w)

	respondWithJSON(w, http.StatusOK, 0, "success", nil)
}

//	POST{
//		username,
//		password,
//		mobilephone
//	}
type RegisterRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Mobilephone string `json:"mobilephone"`
}

func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondWithJSON(w, http.StatusMethodNotAllowed, -1, "不支持的请求方法", nil)
		return
	}

	var registerReq RegisterRequest
	err := json.NewDecoder(req.Body).Decode(&registerReq)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, -1, "请求体解析失败: "+err.Error(), nil)
		fmt.Println(err)
		return
	}

	if registerReq.Username == "" || registerReq.Password == "" || registerReq.Mobilephone == "" {
		respondWithJSON(w, http.StatusBadRequest, -1, "用户名、密码和手机号是必填项", nil)
		return
	}

	// 使用连接池获取数据库连接
	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法连接到数据库: "+err.Error(), nil)
		fmt.Println(err)
		return
	}
	defer conn.Close(context.Background())

	// 检查手机号是否已存在
	var exists bool
	err = conn.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM t_user WHERE mobilephone=$1)", registerReq.Mobilephone,
	).Scan(&exists)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "检查手机号时发生错误: "+err.Error(), nil)
		fmt.Println(err)
		return
	}

	if exists {
		respondWithJSON(w, http.StatusBadRequest, -1, "手机号已被注册", nil)
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "加密密码失败: "+err.Error(), nil)
		fmt.Println(err)
		return
	}

	// 插入用户数据
	_, err = conn.Exec(
		context.Background(),
		"INSERT INTO t_user (username, password, mobilephone) VALUES ($1, $2, $3)",
		registerReq.Username, string(hashedPassword), registerReq.Mobilephone,
	)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "注册失败: "+err.Error(), nil)
		fmt.Println(err)
		return
	}

	// 返回成功响应
	respondWithJSON(w, http.StatusOK, 0, "注册成功", nil)
}

//	post{
//		title,
//		content
//	}
type CreateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func CreateArticleHandler(w http.ResponseWriter, req *http.Request) {
	var articleReq CreateArticleRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&articleReq)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, -1, "请求数据解析失败: "+err.Error(), nil)
		fmt.Println(err)
		return
	}

	if articleReq.Title == "" || articleReq.Content == "" {
		respondWithJSON(w, http.StatusBadRequest, -1, "标题和内容是必填项", nil)
		return
	}

	// 获取 session 中的用户信息
	session, err := store.Get(req, "session-information")
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取 session: "+err.Error(), nil)
		fmt.Println(err)
		return
	}

	userID, ok := session.Values["userid"].(int)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, -1, "未登录或 session 过期", nil)
		return
	}

	// 获取数据库连接池
	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法连接到数据库: "+err.Error(), nil)
		fmt.Println(err)
		return
	}
	defer conn.Close(context.Background())

	fmt.Println("userID:", userID, "title:", articleReq.Title, "content:", articleReq.Content)

	// 插入文章到数据库
	_, err = conn.Exec(
		context.Background(),
		"INSERT INTO t_articles (userid, title, content) VALUES ($1, $2, $3)",
		userID, articleReq.Title, articleReq.Content,
	)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "创建文章失败: "+err.Error(), nil)
		fmt.Println(err)
		return
	}

	// 成功响应
	respondWithJSON(w, http.StatusOK, 0, "文章创建成功", nil)
}

func UserArticlesHandler(w http.ResponseWriter, req *http.Request) {
	// 获取 session 数据
	session, err := store.Get(req, "session-information")

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取 session: "+err.Error(), nil)
		fmt.Println("Session error:", err)
		return
	}

	// 从 session 中获取用户 ID
	userID, ok := session.Values["userid"].(int)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, -1, "未登录或 session 过期", nil)
		return
	}

	// 获取数据库连接池
	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法连接到数据库: "+err.Error(), nil)
		fmt.Println("Database connection error:", err)
		return
	}
	defer conn.Close(context.Background()) // 使用连接池，避免重复连接和关闭

	// 查询当前用户的文章
	rows, err := conn.Query(
		context.Background(),
		"SELECT articleid, title, content FROM t_articles WHERE userid=$1",
		userID,
	)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "获取文章失败: "+err.Error(), nil)
		fmt.Println("Query error:", err)
		return
	}
	defer rows.Close()

	// 存储文章列表
	articles := []Article{}
	for rows.Next() {
		var article Article
		// 扫描当前行到文章结构体
		err = rows.Scan(&article.ID, &article.Title, &article.Content)
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, -1, "文章数据解析失败: "+err.Error(), nil)
			fmt.Println("Row scan error:", err)
			return
		}
		// 将文章添加到列表中
		articles = append(articles, article)
	}

	// 检查查询是否有任何文章数据
	if len(articles) == 0 {
		respondWithJSON(w, http.StatusOK, 0, "没有找到相关文章", nil)
		return
	}

	// 返回文章列表
	respondWithJSON(w, http.StatusOK, 0, "success", articles)
}

//	put{
//		articleid,
//		title,
//		content
//	}
type EditArticleRequest struct {
	ArticleID int    `json:"articleid"`
	Title     string `json:"title"`
	Content   string `json:"content"`
}

func EditArticleHandler(w http.ResponseWriter, req *http.Request) {
	// 解析请求数据
	var articleReq EditArticleRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&articleReq)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, -1, "请求数据解析失败: "+err.Error(), nil)
		fmt.Println("Request body parsing error:", err)
		return
	}

	// 校验请求数据的有效性
	if articleReq.ArticleID == 0 || articleReq.Title == "" || articleReq.Content == "" {
		respondWithJSON(w, http.StatusBadRequest, -1, "文章 ID、标题和内容是必填项", nil)
		return
	}

	// 获取 session 信息
	session, err := store.Get(req, "session-information")
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取 session: "+err.Error(), nil)
		fmt.Println("Session error:", err)
		return
	}

	// 从 session 中获取用户 ID
	userID, ok := session.Values["userid"].(int)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, -1, "未登录或 session 过期", nil)
		return
	}

	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取数据库连接: "+err.Error(), nil)
		fmt.Println("Connection acquire error:", err)
		return
	}
	defer conn.Close(context.Background())

	// 检查文章是否存在，并验证文章的所有者
	var existingUserID int
	err = conn.QueryRow(
		context.Background(),
		"SELECT userid FROM t_articles WHERE articleid=$1", articleReq.ArticleID,
	).Scan(&existingUserID)

	if err != nil {
		if err == pgx.ErrNoRows {
			respondWithJSON(w, http.StatusNotFound, -1, "文章不存在", nil)
			fmt.Println("Article not found:", articleReq.ArticleID)
			return
		}
		respondWithJSON(w, http.StatusInternalServerError, -1, "验证文章所有者失败: "+err.Error(), nil)
		fmt.Println("Error querying article owner:", err)
		return
	}

	// 确认当前用户是文章的所有者
	if existingUserID != userID {
		respondWithJSON(w, http.StatusForbidden, -1, "您没有权限编辑该文章", nil)
		return
	}

	// 更新文章内容
	_, err = conn.Exec(
		context.Background(),
		"UPDATE t_articles SET title=$1, content=$2 WHERE articleid=$3 AND userid=$4",
		articleReq.Title, articleReq.Content, articleReq.ArticleID, userID,
	)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "更新文章失败: "+err.Error(), nil)
		fmt.Println("Update article error:", err)
		return
	}

	// 返回成功响应
	respondWithJSON(w, http.StatusOK, 0, "文章更新成功", nil)
}

//	DELETE{
//		articleid
//	}
func DeleteArticleHandler(w http.ResponseWriter, req *http.Request) {
	// 获取文章 ID
	path := req.URL.Path
	parts := strings.Split(path, "/")

	articleID := parts[len(parts)-1]
	if articleID == "" {
		respondWithJSON(w, http.StatusBadRequest, -1, "文章 ID 不能为空", nil)
		return
	}

	// 获取 session 信息
	session, err := store.Get(req, "session-information")
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取 session: "+err.Error(), nil)
		fmt.Println("Session error:", err)
		return
	}

	// 获取当前用户的 ID
	userID, ok := session.Values["userid"].(int)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, -1, "未登录或 session 过期", nil)
		return
	}

	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取数据库连接: "+err.Error(), nil)
		fmt.Println("Connection acquire error:", err)
		return
	}
	defer conn.Close(context.Background())

	// 查询文章的所有者 ID
	var existingUserID int
	err = conn.QueryRow(
		context.Background(),
		"SELECT userid FROM t_articles WHERE articleid=$1", articleID,
	).Scan(&existingUserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			respondWithJSON(w, http.StatusNotFound, -1, "文章不存在", nil)
			fmt.Println("Article not found:", articleID)
			return
		}
		respondWithJSON(w, http.StatusInternalServerError, -1, "验证文章所有者失败: "+err.Error(), nil)
		fmt.Println("Error querying article owner:", err)
		return
	}

	// 确认当前用户是文章的所有者
	if existingUserID != userID {
		respondWithJSON(w, http.StatusForbidden, -1, "您没有权限删除该文章", nil)
		return
	}

	// 删除文章
	sqlStatement := `DELETE FROM t_articles WHERE articleid=$1 AND userid=$2`
	_, err = conn.Exec(
		context.Background(),
		sqlStatement,
		articleID,
		userID,
	)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "删除文章失败: "+err.Error(), nil)
		fmt.Println("Error deleting article:", err)
		return
	}

	// 返回成功响应
	respondWithJSON(w, http.StatusOK, 0, "文章删除成功", nil)
}

//	post{
//		currentpassword,
//		newpassword
//	}
func ProfileHandler(w http.ResponseWriter, req *http.Request) {
	// 获取 session 信息
	session, err := store.Get(req, "session-information")
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取 session: "+err.Error(), nil)
		fmt.Println("Session error:", err)
		return
	}

	// 获取当前用户的 ID
	userID, ok := session.Values["userid"].(int)
	if !ok || userID == 0 {
		respondWithJSON(w, http.StatusForbidden, -1, "未登录或 session 过期", nil)
		return
	}

	// 解析请求数据
	var requestData struct {
		CurrentPassword string `json:"currentpassword"`
		NewPassword     string `json:"newpassword"`
	}
	err = json.NewDecoder(req.Body).Decode(&requestData)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, -1, "请求数据格式错误", nil)
		fmt.Println("JSON decode error:", err)
		return
	}

	// 校验必填项
	if requestData.CurrentPassword == "" || requestData.NewPassword == "" {
		respondWithJSON(w, http.StatusBadRequest, -1, "当前密码和新密码是必填项", nil)
		return
	}

	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取数据库连接: "+err.Error(), nil)
		fmt.Println("Connection acquire error:", err)
		return
	}
	defer conn.Close(context.Background())

	// 查询当前用户的密码
	var storedPassword string
	err = conn.QueryRow(
		context.Background(),
		"SELECT password FROM t_user WHERE userid=$1", userID,
	).Scan(&storedPassword)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "获取用户信息失败: "+err.Error(), nil)
		fmt.Println("Error querying stored password:", err)
		return
	}

	// 比对当前密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(requestData.CurrentPassword))
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, -1, "当前密码错误", nil)
		return
	}

	// 加密新密码
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "密码加密失败: "+err.Error(), nil)
		fmt.Println("Error hashing new password:", err)
		return
	}

	// 更新密码
	_, err = conn.Exec(
		context.Background(),
		"UPDATE t_user SET password=$1 WHERE userid=$2",
		string(hashedNewPassword),
		userID,
	)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "更新密码失败: "+err.Error(), nil)
		fmt.Println("Error updating password:", err)
		return
	}

	// 返回成功响应
	respondWithJSON(w, http.StatusOK, 0, "密码修改成功", nil)
}

//	post{
//		username,
//		password
//	}
func AdminLoginHandler(w http.ResponseWriter, req *http.Request) {
	var requestData struct {
		AdminName string `json:"username"`
		Password  string `json:"password"`
	}

	// 解码请求体
	err := json.NewDecoder(req.Body).Decode(&requestData)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, -1, "请求数据格式错误", nil)
		fmt.Println("Error decoding request data:", err)
		return
	}

	// 校验管理员账号和密码不能为空
	if requestData.AdminName == "" || requestData.Password == "" {
		respondWithJSON(w, http.StatusBadRequest, -1, "管理员用户名和密码不能为空", nil)
		return
	}

	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取数据库连接: "+err.Error(), nil)
		fmt.Println("Connection acquire error:", err)
		return
	}
	defer conn.Close(context.Background())

	// 查询存储的密码以及管理员ID
	var storedPassword string
	var userID int // 新增的字段，用于存储 adminid
	err = conn.QueryRow(
		context.Background(),
		"SELECT adminid, password FROM t_admin WHERE username=$1", requestData.AdminName,
	).Scan(&userID, &storedPassword)

	if err != nil {
		// 不要在登录失败时泄露具体原因，避免攻击者获知是用户名还是密码错误
		respondWithJSON(w, http.StatusUnauthorized, -1, "用户名或密码错误", nil)
		fmt.Println("Error querying stored password:", err)
		return
	}

	// 比对密码
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(requestData.Password))
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, -1, "用户名或密码错误", nil)
		fmt.Println("Password mismatch:", err)
		return
	}

	// 创建 session 并保存管理员信息
	session, _ := store.Get(req, "session-information")
	session.Values["username"] = requestData.AdminName
	err = session.Save(req, w)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法保存 session", nil)
		fmt.Println("Error saving session:", err)
		return
	}

	// 将 userid 和 username 返回给前端
	data := map[string]interface{}{
		"userid":   userID, // 添加 userid
		"username": requestData.AdminName,
	}

	// 返回成功响应
	respondWithJSON(w, http.StatusOK, 0, "登录成功", data)
}

//	post{
//		username,
//		password,
//		mobilephone
//	}
func CreateAdminHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondWithJSON(w, http.StatusMethodNotAllowed, -1, "不支持的请求方法", nil)
		return
	}

	// 检查是否具有管理员权限
	session, _ := store.Get(req, "session-information")
	if session.Values["username"] == nil {
		respondWithJSON(w, http.StatusForbidden, -1, "没有权限访问", nil)
		return
	}

	var newAdmin struct {
		AdminName string `json:"username"`
		Password  string `json:"password"`
		Mobile    string `json:"mobilephone"`
	}
	defer req.Body.Close() // 确保请求体关闭

	// 解析请求体
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&newAdmin); err != nil {
		respondWithJSON(w, http.StatusBadRequest, -1, "请求数据格式不正确", nil)
		fmt.Println("Error decoding request data:", err)
		return
	}

	// 校验输入
	if newAdmin.AdminName == "" || newAdmin.Password == "" {
		respondWithJSON(w, http.StatusBadRequest, -1, "管理员名称和密码是必填项", nil)
		return
	}

	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取数据库连接: "+err.Error(), nil)
		fmt.Println("Connection acquire error:", err)
		return
	}
	defer conn.Close(context.Background())

	// 查询是否已有相同用户名的管理员
	var existingUsername string
	err = conn.QueryRow(
		context.Background(),
		"SELECT username FROM t_admin WHERE username=$1", newAdmin.AdminName,
	).Scan(&existingUsername)

	// 如果有重复的用户名，则返回冲突错误
	if err == nil {
		respondWithJSON(w, http.StatusConflict, -1, "用户名已存在", nil)
		return
	}
	if err != pgx.ErrNoRows {
		respondWithJSON(w, http.StatusInternalServerError, -1, "查询管理员失败: "+err.Error(), nil)
		fmt.Println("Error querying admin username:", err)
		return
	}

	// 对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newAdmin.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "密码加密失败: "+err.Error(), nil)
		fmt.Println("Error hashing password:", err)
		return
	}

	// 插入新管理员到数据库
	_, err = conn.Exec(
		context.Background(),
		"INSERT INTO t_admin (username, password, mobilephone) VALUES ($1, $2, $3)",
		newAdmin.AdminName,
		string(hashedPassword),
		newAdmin.Mobile, // 添加手机号字段
	)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "创建新管理员失败: "+err.Error(), nil)
		fmt.Println("Error inserting new admin:", err)
		return
	}

	// 返回成功响应
	respondWithJSON(w, http.StatusOK, 0, "新管理员创建成功", nil)
}

// User 用于表示用户信息
func GetUsersHandler(w http.ResponseWriter, req *http.Request) {
	// 从封装好的连接池获取连接
	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法连接到数据库: "+err.Error(), nil)
		fmt.Println("Error connecting to database:", err)
		return
	}
	// 连接池是长时间有效的，所以这里不需要每次请求都显式关闭连接

	// 查询所有用户信息
	rows, err := conn.Query(context.Background(), "SELECT userid, username, mobilephone FROM t_user")
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "查询用户失败: "+err.Error(), nil)
		fmt.Println("Error querying users:", err)
		return
	}
	defer rows.Close() // 确保 rows 被关闭

	var users []map[string]interface{} // 用于存储查询到的用户信息

	// 遍历查询结果
	for rows.Next() {
		var userId int
		var username, mobilephone string

		// 扫描当前行数据
		err := rows.Scan(&userId, &username, &mobilephone)
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, -1, "处理用户数据失败: "+err.Error(), nil)
			fmt.Println("Error scanning user data:", err)
			return
		}

		// 创建用户信息的 map，并添加到 users 列表中
		user := map[string]interface{}{
			"userid":      userId,
			"username":    username,
			"mobilephone": mobilephone,
		}
		users = append(users, user)
	}

	// 如果没有查询到用户，返回 404 错误
	if len(users) == 0 {
		respondWithJSON(w, http.StatusNotFound, -1, "没有找到任何用户", nil)
		return
	}

	// 返回成功的响应
	respondWithJSON(w, http.StatusOK, 0, "成功获取用户列表", users)
}

//	post{
//		userid
//	}
func DeleteUserHandler(w http.ResponseWriter, req *http.Request) {
	// 验证管理员是否已登录
	session, _ := store.Get(req, "session-information")
	if session.Values["username"] == nil {
		respondWithJSON(w, http.StatusForbidden, -1, "没有权限访问", nil)
		return
	}

	// 获取请求中的用户ID
	path := req.URL.Path
	parts := strings.Split(path, "/")

	userID := parts[len(parts)-1]
	if userID == "" {
		respondWithJSON(w, http.StatusBadRequest, -1, "缺少用户ID", nil)
		return
	}

	// 从封装好的连接池获取连接
	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法连接到数据库: "+err.Error(), nil)
		fmt.Println("Error connecting to database:", err)
		return
	}

	// 确保要删除的用户存在
	var existingUserID int
	err = conn.QueryRow(
		context.Background(),
		"SELECT userid FROM t_user WHERE userid=$1", userID,
	).Scan(&existingUserID)

	// 如果用户不存在
	if err != nil {
		if err == pgx.ErrNoRows {
			respondWithJSON(w, http.StatusNotFound, -1, "用户不存在", nil)
		} else {
			respondWithJSON(w, http.StatusInternalServerError, -1, "验证用户失败: "+err.Error(), nil)
		}
		fmt.Println("Error verifying user:", err)
		return
	}

	// 删除用户
	_, err = conn.Exec(
		context.Background(),
		"DELETE FROM t_user WHERE userid=$1", userID,
	)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "删除用户失败: "+err.Error(), nil)
		fmt.Println("Error deleting user:", err)
		return
	}

	// 返回成功响应
	respondWithJSON(w, http.StatusOK, 0, "用户删除成功", nil)
}

//	post{
//		articleid
//	}
func DeleteArticlesHandler(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	parts := strings.Split(path, "/")

	articleID := parts[len(parts)-1]
	if articleID == "" {
		respondWithJSON(w, http.StatusBadRequest, -1, "缺少文章ID", nil)
		return
	}

	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法连接到数据库: "+err.Error(), nil)
		fmt.Println("Error connecting to database:", err)
		return
	}

	// 确保文章存在
	var existingArticleID int
	err = conn.QueryRow(
		context.Background(),
		"SELECT articleid FROM t_articles WHERE articleid=$1", articleID,
	).Scan(&existingArticleID)

	if err != nil {
		if err == pgx.ErrNoRows {
			respondWithJSON(w, http.StatusNotFound, -1, "文章不存在", nil)
		} else {
			respondWithJSON(w, http.StatusInternalServerError, -1, "验证文章失败: "+err.Error(), nil)
		}
		fmt.Println("Error verifying article:", err)
		return
	}

	// 删除文章
	_, err = conn.Exec(
		context.Background(),
		"DELETE FROM t_articles WHERE articleid=$1", articleID,
	)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "删除文章失败: "+err.Error(), nil)
		fmt.Println("Error deleting article:", err)
		return
	}

	// 返回成功响应
	respondWithJSON(w, http.StatusOK, 0, "文章删除成功", nil)
}

//	post{
//		oldpassword,
//		newpassword
//	}
func AdminProfileHandler(w http.ResponseWriter, req *http.Request) {
	// 验证管理员是否已登录
	session, _ := store.Get(req, "session-information")
	if session.Values["username"] == nil {
		respondWithJSON(w, http.StatusForbidden, -1, "没有权限访问", nil)
		return
	}

	// 解析请求体中的旧密码和新密码
	var requestData struct {
		OldPassword string `json:"oldpassword"`
		NewPassword string `json:"newpassword"`
	}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&requestData); err != nil {
		respondWithJSON(w, http.StatusBadRequest, -1, "请求数据格式不正确", nil)
		return
	}

	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取数据库连接: "+err.Error(), nil)
		fmt.Println("Connection acquire error:", err)
		return
	}
	defer conn.Close(context.Background())

	// 获取当前管理员的密码
	var storedPassword string
	adminUsername := session.Values["username"].(string)
	err = conn.QueryRow(
		context.Background(),
		"SELECT password FROM t_admin WHERE username=$1",
		adminUsername,
	).Scan(&storedPassword)

	if err != nil {
		if err == pgx.ErrNoRows {
			respondWithJSON(w, http.StatusNotFound, -1, "管理员不存在", nil)
		} else {
			respondWithJSON(w, http.StatusInternalServerError, -1, "验证旧密码失败: "+err.Error(), nil)
		}
		fmt.Println("Error querying password:", err)
		return
	}

	// 使用 bcrypt 来比较密码（存储的密码是加密后的）
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(requestData.OldPassword))
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, -1, "旧密码不正确", nil)
		return
	}

	// 新密码加密
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "加密新密码失败: "+err.Error(), nil)
		fmt.Println("Error hashing new password:", err)
		return
	}

	// 更新密码
	_, err = conn.Exec(
		context.Background(),
		"UPDATE t_admin SET password=$1 WHERE username=$2",
		string(hashedNewPassword), // 保存加密后的新密码
		adminUsername,
	)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "更新密码失败: "+err.Error(), nil)
		fmt.Println("Error updating password:", err)
		return
	}

	// 返回成功响应
	respondWithJSON(w, http.StatusOK, 0, "密码修改成功", nil)
}

//	post{
//		userid,
//		newpassword
//	}
func UsersProfileHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session-information")
	if session.Values["username"] == nil {
		respondWithJSON(w, http.StatusForbidden, -1, "没有权限访问", nil)
		return
	}

	var requestData struct {
		UserID      int    `json:"userid"`
		NewPassword string `json:"newpassword"`
	}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&requestData); err != nil {
		respondWithJSON(w, http.StatusBadRequest, -1, "请求数据格式不正确", nil)
		return
	}

	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取数据库连接: "+err.Error(), nil)
		fmt.Println("Connection acquire error:", err)
		return
	}
	defer conn.Close(context.Background())

	// 获取用户的密码
	var storedPassword string
	err = conn.QueryRow(
		context.Background(),
		"SELECT password FROM t_user WHERE userid=$1",
		requestData.UserID,
	).Scan(&storedPassword)

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "获取用户信息失败: "+err.Error(), nil)
		fmt.Println("Error getting user information:", err)
		return
	}

	// 加密新密码
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "密码加密失败: "+err.Error(), nil)
		fmt.Println("Error hashing new password:", err)
		return
	}

	// 更新密码
	_, err = conn.Exec(
		context.Background(),
		"UPDATE t_user SET password=$1 WHERE userid=$2",
		string(hashedNewPassword),
		requestData.UserID,
	)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "更新密码失败: "+err.Error(), nil)
		fmt.Println("Error updating password:", err)
		return
	}

	respondWithJSON(w, http.StatusOK, 0, "密码修改成功", nil)
}

func GetAllArticlesHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := SQL.ConnectSQL()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "无法获取数据库连接: "+err.Error(), nil)
		fmt.Println("Connection acquire error:", err)
		return
	}
	defer conn.Close(context.Background())

	// 执行 SQL 查询操作
	rows, err := conn.Query(context.Background(), "SELECT userid, articleid, title, content FROM t_articles")
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, -1, "查询文章失败: "+err.Error(), nil)
		fmt.Println("Error querying articles:", err)
		return
	}
	defer rows.Close() // 确保 rows 被关闭

	// 存储查询结果
	var articles []map[string]interface{}

	// 遍历查询结果
	for rows.Next() {
		var userId, articleId int
		var title, content string

		// 扫描当前行数据
		err := rows.Scan(&userId, &articleId, &title, &content)
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, -1, "处理文章数据失败: "+err.Error(), nil)
			fmt.Println("Error scanning article data:", err)
			return
		}

		// 将当前文章信息存入列表
		article := map[string]interface{}{
			"userid":    userId,
			"articleid": articleId,
			"title":     title,
			"content":   content,
		}
		articles = append(articles, article)
	}

	// 如果没有查询到文章，返回 404 错误
	if len(articles) == 0 {
		respondWithJSON(w, http.StatusNotFound, -1, "没有找到任何文章", nil)
		return
	}

	// 返回成功的响应
	respondWithJSON(w, http.StatusOK, 0, "成功获取文章列表", articles)
}
