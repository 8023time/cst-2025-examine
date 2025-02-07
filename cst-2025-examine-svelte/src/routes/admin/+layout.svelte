<script>
	import { information } from '$lib/stores';
	import { invalidate } from '$app/navigation';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	const logout = async () => {
		information.removeUserInformation()
		await goto('/')
		invalidate()
	};
</script>

<div class="container">
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
				<li class:active={$page.pathname === '/admin'}>
					<a href="/admin" on:click={() => goto('/admin')}>首页</a>
				</li>
				<li class:active={$page.pathname === '/admin/articles'}>
					<a href="/admin/users" on:click={() => goto('/admin/users')}>管理用户</a>
				</li>
                <li class:active={$page.pathname === '/admin/profile'}>
					<a href="/admin/profile" on:click={() => goto('/admin/profile')}>修改个人密码</a>
				</li>
                <li class:active={$page.pathname === '/admin/profile'}>
					<a href="/admin/create-admin" on:click={() => goto('/admin/create-admin')}>创建管理员</a>
				</li>
			</ul>
		</nav>
	</aside>
	<!-- 主体内容区域 -->
	<main class="content">
		<header class="header">
			<div class="header-right">
				<div class="user-info">
					<span class="username">
						你好! 管理员
					</span>
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
	.container {
		/* weight: 100%; */
		display: flex;
		height: 100vh;
	}
	.sidebar {
		width: 200px;
		background-color: #333;
		color: #fff;
		padding: 20px;
	}
	.logo {
		margin-bottom: 20px;
		text-align: center;
	}
	.logo img {
		max-width: 100%;
		border-radius: 50%;
	}
	.sidebar nav ul {
		list-style: none;
		padding: 0;
	}
	.sidebar nav li {
		margin-bottom: 10px;
	}
	.sidebar nav a {
		display: block;
		color: #fff;
		text-decoration: none;
		padding: 10px;
		border-radius: 5px;
		transition: background-color 0.3s ease;
	}
	.sidebar nav a:hover,
	.sidebar nav li.active a {
		background-color: #555;
	}
	.content {
		flex: 1;
		padding: 20px;
		overflow-y: auto;
		background-color: #f5f5f5;
	}
	.header {
		display: flex;
		justify-content: flex-end;
		align-items: center;
		padding: 10px 20px;
		border-bottom: 1px solid #ccc;
		background-color: #fff;
	}
	.header-right {
		display: flex;
		align-items: center;
		gap: 15px;
	}
	.username {
		text-decoration: none;
		color: #333;
		align-items: center;
		font-weight: bold;
		padding: 5px 10px;
		border-radius: 5px;
		transition:
			background-color 0.3s ease,
			color 0.3s ease;
	}
	.user-info {
		position: relative;
	}
	.main-content {
		padding: 20px;
	}
</style>
