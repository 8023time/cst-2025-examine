<script>
	import { onMount } from 'svelte';
	import { AdminAPI } from '$lib/api';
	import { Notyf } from 'notyf';
	import { goto } from '$app/navigation';

	let notyf;

	let users = [];

	const getusersAll = async () => {
		try {
			const response = await AdminAPI.getAllUsers();
			users = response.data.data;
			notyf.success('获取用户数据成功');
		} catch (error) {
			notyf.error('获取用户数据失败');
		}
	};

	const deleteUser = async (id) => {
		users = users.filter((user) => user.id !== id);
		try {
			const response = await AdminAPI.deleteUser(id);
			notyf.success('删除用户成功');
			window.location.reload();
		} catch (error) {
			notyf.error('删除用户失败');
		}
	};

	const changePassword = (id) => {
		goto(`/admin/users/${id}`);
	};

	onMount(() => {
		notyf = new Notyf({
			duration: 3000,
			className: 'x-notification',
			position: { x: 'right', y: 'top' }
		});
		getusersAll();
	});
</script>

<div class="container">
	<h1>用户管理页面</h1>
	{#if users.length === 0}
		<p>暂无用户数据</p>
	{:else}
		<ul>
			{#each users as user}
				<li class="user-item">
					<span>{user.username} - {user.mobilephone || '暂无无联系方式'}</span>
					<div class="button-group">
						<button class="change-button" on:click={changePassword(user.userid)}>修改密码</button>
						<button class="delete-button" on:click={deleteUser(user.userid)}>删除</button>
					</div>
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style lang="scss">
	/* 定义颜色变量 */
	$primary-blue: #007bff;
	$success-green: #28a745;
	$danger-red: #e74c3c;
	$text-dark: #2c3e50;
	$text-medium: #555;
	$bg-light: #f8f9fa;
	$white: #fff;
	$border-light: #e0e0e0;

	.container {
		max-width: 100%; /* 允许铺满左右 */
		margin: 40px 30px; /* 调整左右外边距，保留上下 */
		padding: 35px; /* 增加内边距 */
		background-color: $white;
		border-radius: 12px; /* 更大的圆角 */
		box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1); /* 柔和阴影 */
	}
	h1 {
		font-size: 2.5em; /* 增大字号 */
		margin-bottom: 30px; /* 增加底部间距 */
		text-align: center;
		color: $text-dark;
		font-weight: 700;
	}
	ul {
		list-style: none;
		padding: 0;
		margin: 0;
	}
	.user-item {
		display: flex;
		justify-content: space-between;
		align-items: center; /* 垂直居中 */
		padding: 15px 20px; /* 增加内边距 */
		margin-bottom: 12px; /* 增加底部间距 */
		background-color: $white;
		border: 1px solid $border-light;
		border-radius: 8px; /* 更大的圆角 */
		box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05); /* 轻微阴影 */
		transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1); /* 平滑过渡 */

		&:hover {
			transform: translateY(-5px); /* 轻微上浮 */
			box-shadow: 0 6px 15px rgba(0, 0, 0, 0.1);
		}

		span {
			font-size: 1.1em; /* 调整字号 */
			color: $text-medium;
			font-weight: 500;
			flex-grow: 1; /* 占据更多空间 */
			word-break: break-word; /* 避免文字溢出 */
			margin-right: 15px; /* 与按钮组的间距 */
		}
	}
	.button-group {
		display: flex;
		gap: 10px; /* 按钮间距 */
	}
	.change-button,
	.delete-button {
		border: none;
		padding: 10px 18px; /* 调整内边距 */
		border-radius: 6px; /* 统一圆角 */
		cursor: pointer;
		font-weight: 600;
		color: $white;
		transition: all 0.3s ease;
		box-shadow: 0 2px 5px rgba(0, 0, 0, 0.08);

		&:hover {
			transform: translateY(-1px);
			box-shadow: 0 4px 8px rgba(0, 0, 0, 0.12);
		}

		&:active {
			transform: translateY(0);
			box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
		}
	}
	.change-button {
		background-color: $success-green; /* 绿色 */
		&:hover {
			background-color: darken($success-green, 8%);
		}
	}
	.delete-button {
		background-color: $danger-red; /* 红色 */
		&:hover {
			background-color: darken($danger-red, 8%);
		}
	}
	p {
		text-align: center;
		color: $text-medium;
		font-size: 1.2em;
		font-weight: 600;
		padding: 30px;
		background-color: $white;
		border-radius: 8px;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
		margin: 30px auto; /* 居中显示 */
		max-width: 500px;
	}

	/* Notyf 样式保持不变 */
	.x-notification {
		background-color: #2ecc71;
		color: white;
		border-radius: 4px;
		padding: 10px 20px;
		font-size: 16px;
	}
	.x-notification.error {
		background-color: #e74c3c;
	}

	/* 响应式设计 */
	@media (max-width: 768px) {
		.container {
			margin: 20px 15px; /* 调整外边距 */
			padding: 25px; /* 调整内边距 */
			box-shadow: 0 4px 10px rgba(0, 0, 0, 0.08);
		}
		h1 {
			font-size: 2em;
			margin-bottom: 25px;
		}
		.user-item {
			flex-direction: column; /* 垂直堆叠 */
			align-items: flex-start;
			padding: 15px;
			gap: 10px; /* 间距 */
			span {
				font-size: 1em;
				margin-bottom: 5px;
				margin-right: 0;
			}
			.button-group {
				width: 100%;
				justify-content: flex-end; /* 按钮靠右 */
				gap: 8px;
			}
			.change-button,
			.delete-button {
				padding: 8px 15px;
				font-size: 0.9em;
			}
		}
		p {
			margin: 20px 15px;
			padding: 25px;
			font-size: 1.1em;
		}
	}

	@media (max-width: 480px) {
		.container {
			padding: 20px;
			margin: 20px 10px; /* 进一步调整小屏幕下的左右外边距 */
		}
		h1 {
			font-size: 1.8em;
			margin-bottom: 20px;
		}
		.user-item {
			padding: 12px;
			span {
				font-size: 0.95em;
			}
			.button-group {
				gap: 5px;
			}
			.change-button,
			.delete-button {
				padding: 6px 12px;
				font-size: 0.8em;
			}
		}
		p {
			font-size: 1em;
			padding: 20px;
		}
	}
</style>
