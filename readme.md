一、软件功能用例集

1、普通用户用例
1.1 用户张三注册、登录。
1.2 他创建标题为“耗散结构迷思”的文章，该文章内容包含有富文本形式的文字与图片。
1.3 他退出登录。
1.4 他重新登录并修改文章“耗散结构迷思”。
1.5 他发布文章“耗散结构迷思”。
1.6 他感觉不满意从而删除文章“耗散结构迷思”。
1.7 他浏览、查找他之前发表的文章。
1.8 他修改了自己的密码。

2 管理员用例
2.1 李明用内置管理员账号admin登录了系统。
2.2 他创建了一个新的管理员账号root。
2.3 他修改了自己的密码。
2.4 他因为用户李四说他忘记了密码，所以修改李四账号的密码为666666.
2.5 他浏览最近用户创建、发布的文章。
2.7 他发现用户张华的文章“民主的本质”违反了规定，决定删除该文章。
2.8 他发现张华多次违反规定，决定删除张华的账号。

前端svelte:
cst-2025-examine-svelte/
├── src/
│   ├── lib/
│   │   ├── api.js        # 封装 API 请求
│   │   ├── stores.js     # 全局状态管理 (例如用户信息)
│   │   └── request.js    # 响应拦截器函数
│   ├── routes/
│   │   ├── +layout.svelte # 全局布局
│   │   ├── +page.svelte   # 登录,注册页面
│   │   ├── user/
│   │   │   ├── profile/
│   │   │   │   └── +page.svelte # 个人中心页面
│   │   │   └── articles/
│   │   │   │   ├── +page.svelte # 文章列表
│   │   │   │   └── [id]/
│   │   │   │       └── +page.svelte # 文章详情
│   │   └── admin/
│   │   │   ├── +layout.svelte # 管理员后台布局
│   │   │   ├── users/
│   │   │   │   └── +page.svelte # 用户管理
│   │   │   └── articles/
│   │   │   │   ├── +page.svelte # 文章列表
│   │   │   │   └── [id]/
│   │   │   │       └── +page.svelte # 文章详情
│   │   └── +page.svelte   # 主页
│   ├── app.html          # 根 HTML 模板
│   └── static/           # 静态文件 (如图片)
├── svelte.config.js    # SvelteKit 配置
├── package.json        # 依赖管理
├── pnpm-lock.yaml      # pnpm 依赖锁定
└── README.md           # 说明文档

后端go:
cst-2025-examine-go/
├── API/
│   └── API.go      # 封装 API 接口
├── SQL/
│   └── SQL.go      # 封装 SQL 操作
├── go.mod          # Go 模块文件
├── go.sum          # Go 模块校验文件
└── main.go         # 项目入口文件

数据库
CREATE DATABASE testdb;
\c testdb;

CREATE SCHEMA testuser;

-- 创建 t_user 表
CREATE TABLE testuser.t_user (
    userid SERIAL PRIMARY KEY,                   -- 自增的用户ID
    username VARCHAR(255) UNIQUE NOT NULL,        -- 用户名，唯一且非空
    password VARCHAR(255) NOT NULL,               -- 存储加密后的密码
    mobilephone VARCHAR(15)                       -- 可选的手机号
);

-- 重置 t_user 表的自增序列，起始值为 20000
ALTER SEQUENCE testuser.t_user_userid_seq RESTART WITH 20000;

-- 创建 t_articles 表
CREATE TABLE testuser.t_articles (
    article_id SERIAL PRIMARY KEY,                -- 自增的文章ID
    userid INT NOT NULL,                          -- 外键，指定用户ID
    title VARCHAR(255) NOT NULL,                   -- 文章标题，非空
    content TEXT,                                  -- 富文本内容
    FOREIGN KEY (userid) REFERENCES testuser.t_user(userid) -- 外键约束，确保文章和用户的关联
);

-- 重置 t_articles 表的自增序列，起始值为 20000
ALTER SEQUENCE testuser.t_articles_article_id_seq RESTART WITH 20000;

-- 创建 t_admin 表
CREATE TABLE testuser.t_admin (
    admin_id SERIAL PRIMARY KEY,                  -- 自增的管理员ID
    username VARCHAR(255) UNIQUE NOT NULL,         -- 管理员用户名，唯一且非空
    password VARCHAR(255) NOT NULL,                -- 存储加密后的密码
    mobilephone VARCHAR(255) UNIQUE               -- 可选的管理员手机号
);

-- 重置 t_admin 表的自增序列，起始值为 20000
ALTER SEQUENCE testuser.t_admin_admin_id_seq RESTART WITH 20000;


管理员查看文章列表页面： /admin 
管理员创建管理员账号页面： /admin/create-admin 
管理员查看user的页面:/admin/users
管理员修改自己的密码:/admin/profile 
管理员修改users的密码:/admin/profile-user

user查看文章列表页面： /user
用户修改个人密码页面： /user/profile 
用户写文章页面： /user/articles 
编辑文章页面： /user/articles/edit 
查看单个文章页面： /user/articles/view 
文章详情页面： /articles/[id]


![image](https://github.com/user-attachments/assets/38c90a14-fb1b-48f7-9aa8-8119f3da975a)

![image](https://github.com/user-attachments/assets/e7dc878e-6503-4262-80e2-bdc0dd2f2624)

![image](https://github.com/user-attachments/assets/d23ce4a9-e5e6-4989-a92d-81d290e11c98)


