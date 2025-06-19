<script>
	import { information } from '$lib/stores';
	import { invalidate } from '$app/navigation';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	let inforname = information.userInformation.username;

	const logout = async () => {
		information.removeUserInformation();
		await goto('/');
		invalidate();
	};
</script>

<div class="container_DCM">
	<!-- 侧边栏 -->
	<aside class="sidebar">
		<div class="logo">
			<img
				src="https://cdn.jsdelivr.net/gh/8023time/image-storage-address/basic-img/avatar.jpg"
				alt="Logo"
			/>
		</div>
		<nav>
			<ul>
				<li class:active={$page.pathname === '/user'}>
					<a href="/user" on:click={() => goto('/user')}>首页</a>
				</li>
				<li class:active={$page.pathname === '/user/articles'}>
					<a href="/user/articles" on:click={() => goto('/user/articles')}>发表文章</a>
				</li>
				<li class:active={$page.pathname === '/user/profile'}>
					<a href="/user/profile" on:click={() => goto('/user/profile')}>修改信息</a>
				</li>
			</ul>
		</nav>
	</aside>
	<!-- 主体内容区域 -->
	<main class="content">
		<header class="header">
			<div class="header-right">
				<div class="user-info">
					<span class="username"> 你好! </span>
					<button class="button is-primary" on:click={logout}> 退出登录 </button>
				</div>
			</div>
		</header>
		<div class="main-content">
			<slot />
		</div>
	</main>
</div>

<style>
	@import 'notyf/notyf.min.css';
	@import 'bulma/css/bulma.css';

	:root {
		--sidebar-width: 250px;
		--header-height: 60px;
		--primary-color: #007bff;
		--secondary-color: #f5f5f5;
		--text-color-dark: #333;
		--text-color-light: #fff;
		--sidebar-bg: #2c3e50; /* 深蓝色调 */
		--sidebar-active: #34495e; /* 激活背景色 */
		--border-color: #e0e0e0;
	}

	html,
	body {
		height: 100%;
		margin: 0;
		padding: 0;
		overflow: hidden; /* 防止双滚动条 */
		font-family: 'Roboto', sans-serif; /* 假设 Roboto 字体可用 */
		-webkit-font-smoothing: antialiased;
		-moz-osx-font-smoothing: grayscale;
	}

	.container_DCM {
		display: flex;
		height: 100vh;
		width: 100%;
		overflow: hidden; /* 防止内部内容溢出导致滚动条 */
	}

	.sidebar {
		width: var(--sidebar-width);
		background-color: var(--sidebar-bg);
		color: var(--text-color-light);
		padding: 20px 0;
		display: flex;
		flex-direction: column;
		box-shadow: 2px 0 10px rgba(0, 0, 0, 0.1);
		z-index: 10;
		transition: all 0.3s ease; /* 平滑过渡 */
	}

	.logo {
		margin-bottom: 30px;
		text-align: center;
		padding: 0 20px;
	}
	.logo img {
		max-width: 120px; /* 调整 logo 大小 */
		border-radius: 50%;
		border: 3px solid #fff; /* 白色边框 */
		box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
	}
	.sidebar nav ul {
		list-style: none;
		padding: 0;
		margin: 0;
		flex-grow: 1; /* 占据剩余空间 */
	}
	.sidebar nav li {
		margin-bottom: 5px;
	}
	.sidebar nav a {
		display: flex;
		align-items: center;
		gap: 15px; /* 图标和文字的间距 */
		color: var(--text-color-light);
		text-decoration: none;
		padding: 12px 25px;
		border-radius: 0 25px 25px 0; /* 右侧圆角胶囊状 */
		transition: all 0.3s ease;
		font-weight: 500;
		position: relative;
		overflow: hidden;
	}

	.sidebar nav a:before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		width: 5px;
		height: 100%;
		background-color: transparent;
		transition: background-color 0.3s ease;
	}

	.sidebar nav a:hover:before,
	.sidebar nav li.active a:before {
		background-color: var(--primary-color); /* 激活时显示蓝色条 */
	}

	.sidebar nav a:hover {
		background-color: var(--sidebar-active);
		padding-left: 30px; /* 悬停时稍微扩大 */
	}
	.sidebar nav li.active a {
		background-color: var(--sidebar-active);
		font-weight: 700;
		box-shadow: inset 3px 0 0 var(--primary-color);
	}

	.content {
		flex: 1;
		display: flex;
		flex-direction: column;
		overflow-y: auto;
		background-color: var(--secondary-color);
	}
	.header {
		display: flex;
		justify-content: flex-end;
		align-items: center;
		padding: 0 30px;
		height: var(--header-height);
		border-bottom: 1px solid var(--border-color);
		background-color: #fff;
		box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
		z-index: 5;
	}
	.header-right {
		display: flex;
		align-items: center;
		gap: 20px;
	}
	.username {
		color: var(--text-color-dark);
		font-weight: 600;
		padding: 8px 15px;
		border-radius: 5px;
		transition: background-color 0.3s ease;
	}

	.username:hover {
		background-color: #f0f0f0;
	}

	.button.is-primary {
		background-color: var(--primary-color);
		border-color: var(--primary-color);
		color: white;
		font-weight: 600;
		padding: 8px 18px;
		border-radius: 8px;
		transition: all 0.3s ease;
		box-shadow: 0 4px 8px rgba(0, 123, 255, 0.2);

		&:hover {
			background-color: darken(var(--primary-color), 8%);
			border-color: darken(var(--primary-color), 8%);
			transform: translateY(-2px);
			box-shadow: 0 6px 12px rgba(0, 123, 255, 0.3);
		}

		&:active {
			transform: translateY(0);
			box-shadow: 0 2px 4px rgba(0, 123, 255, 0.1);
		}
	}

	.main-content {
		flex-grow: 1;
		padding: 30px;
		overflow-y: auto; /* 允许内容区域自身滚动 */
	}

	/* 响应式设计 */
	@media (max-width: 992px) {
		.sidebar {
			width: 60px; /* 缩小侧边栏 */
			overflow: hidden;
			padding: 20px 0;
		}
		.sidebar .logo img {
			max-width: 40px;
		}
		.sidebar nav a {
			justify-content: center;
			padding: 12px 0;
			border-radius: 5px; /* 调整为方形 */
		}
		.sidebar nav a span {
			display: none; /* 隐藏文字 */
		}
		.sidebar nav a:before {
			width: 100%;
			height: 5px;
			bottom: 0;
			left: 0;
			top: unset;
		}
		.sidebar nav a:hover {
			padding-left: 0; /* 移除悬停时的padding调整 */
		}
		.sidebar nav li.active a {
			box-shadow: 0 3px 0 var(--primary-color) inset; /* 底部条 */
		}
	}

	@media (max-width: 768px) {
		.container {
			flex-direction: column;
		}
		.sidebar {
			width: 100%;
			height: auto;
			padding: 10px 0;
			flex-direction: row; /* 横向排列 */
			justify-content: center; /* 居中导航 */
			box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
		}
		.sidebar .logo {
			display: none;
		}
		.sidebar nav ul {
			display: flex;
			flex-wrap: wrap;
			justify-content: center;
		}
		.sidebar nav li {
			margin: 0 5px;
		}
		.sidebar nav a {
			padding: 8px 15px;
			border-radius: 5px;
		}
		.sidebar nav a span {
			display: inline; /* 显示文字 */
		}
		.sidebar nav li.active a {
			box-shadow: inset 0 -3px 0 var(--primary-color); /* 底部条 */
		}
		.content {
			padding-top: 0;
		}
		.header {
			padding: 10px 20px;
			height: auto;
			flex-direction: column;
			align-items: flex-end;
			gap: 10px;
		}
		.header-right {
			width: 100%;
			justify-content: flex-end;
		}
		.main-content {
			padding: 20px;
		}
	}
</style>
