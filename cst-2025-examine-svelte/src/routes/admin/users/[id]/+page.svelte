<script>
	import { onMount } from 'svelte';
	import { Notyf } from 'notyf';
	import { AdminAPI } from '$lib/api';
	import { writable } from 'svelte/store';
	import { page } from '$app/state';

	let notyf;

	$: id = page.params.id;

	let newPassword = '';
	let confirmPassword = '';

	let userlist = [];
	let user = writable({});

	const getuserinformation = async () => {
		try {
			const response = await AdminAPI.getAllUsers();
			if (response.status === 0) {
				userlist = response.data;
				user.set(userlist.find((item) => item.user_id === id));
			} else {
			}
		} catch (error) {
			notyf.error('获取用户信息失败');
		}
	};

	const submitPasswordChange = async () => {
		if (newPassword !== confirmPassword) {
			notyf.error('两次输入的密码不一致');
			return;
		}
		try {
			const response = await AdminAPI.resetUserPassword(parseInt(id), newPassword);
			if (response.data.status === 0) {
				notyf.success('密码修改成功');
				setTimeout(() => {
					window.location.href = '/admin/users';
				}, 1000);
			} else {
				notyf.error('密码修改失败');
			}
		} catch (error) {
			notyf.error('密码修改失败');
		}
	};

	onMount(() => {
		notyf = new Notyf({
			duration: 3000,
			className: 'x-notification',
			position: { x: 'right', y: 'top' }
		});
	});
</script>

<div class="top">
	<h1 class="title">修改用户密码</h1>
</div>

<form on:submit={submitPasswordChange} class="form-container">
	<div class="form-group">
		<label for="newPassword">新密码:</label>
		<input type="password" id="newPassword" bind:value={newPassword} required class="input-field" />
	</div>

	<div class="form-group">
		<label for="confirmPassword">确认新密码:</label>
		<input
			type="password"
			id="confirmPassword"
			bind:value={confirmPassword}
			required
			class="input-field"
		/>
	</div>

	<button type="submit" class="button is-primary">确定修改{user.username}用户的密码</button>
</form>

<style lang="scss">
	@import 'notyf/notyf.min.css';
	@import 'bulma/css/bulma.css';

	/* 定义颜色变量 */
	$primary-blue: #007bff;
	$text-dark: #2c3e50;
	$text-medium: #555;
	$bg-light: #f8f9fa;
	$white: #fff;
	$border-light: #dcdcdc;

	.top {
		display: flex;
		justify-content: center;
		align-items: center;
		height: 180px; /* 增加高度 */
		background-color: $bg-light;
		border-bottom: 1px solid $border-light; /* 更柔和的底部边框 */
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05); /* 轻微阴影 */
		margin-bottom: 40px; /* 增加底部间距 */
	}
	.title {
		font-size: 3em; /* 增大字号 */
		font-weight: 700; /* 更粗的字重 */
		color: $text-dark;
		text-align: center;
		line-height: 1.2;
		letter-spacing: -0.03em;
		margin-left: 0; /* 移除原来的左边距 */
	}
	.form-container {
		max-width: 100%; /* 允许铺满左右 */
		margin: 0 30px 50px 30px; /* 调整左右外边距，保留上下 */
		padding: 40px; /* 增加内边距 */
		background-color: $white;
		border-radius: 15px; /* 更大的圆角 */
		box-shadow: 0 10px 25px rgba(0, 0, 0, 0.12); /* 柔和且明显的阴影 */
		/* margin-bottom: 50px; 底部间距已包含在上方 margin 中 */
	}
	.form-group {
		margin-bottom: 25px; /* 增加间距 */
	}
	label {
		display: block;
		font-weight: 600; /* 更粗的字重 */
		margin-bottom: 10px; /* 增加下方间距 */
		color: $text-dark;
		font-size: 1.05em;
	}
	.input-field {
		width: 100%;
		padding: 14px 18px; /* 调整内边距 */
		border: 1px solid $border-light; /* 统一边框颜色 */
		border-radius: 8px; /* 统一圆角 */
		font-size: 1.1em; /* 调整字号 */
		transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1); /* 平滑过渡 */
		background-color: #fcfcfc; /* 略微偏白的背景色 */
		color: $text-medium;

		&:focus {
			border-color: $primary-blue; /* 聚焦时边框颜色 */
			box-shadow: 0 0 0 0.2em rgba(0, 123, 255, 0.3); /* 聚焦时阴影 */
			outline: none; /* 移除默认轮廓 */
		}
		&::placeholder {
			color: #a0a0a0;
			opacity: 0.9;
		}
	}
	.button.is-primary {
		width: 100%;
		max-width: 300px; /* 限制最大宽度 */
		padding: 16px 25px; /* 调整内边距 */
		font-size: 1.2em; /* 调整字号 */
		border-radius: 8px; /* 统一圆角 */
		cursor: pointer;
		background-color: $primary-blue;
		border-color: $primary-blue;
		color: $white;
		font-weight: 700;
		transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1); /* 平滑过渡 */
		box-shadow: 0 6px 12px rgba(0, 123, 255, 0.25);
		margin-top: 20px; /* 增加上方间距 */

		&:hover {
			background-color: darken($primary-blue, 8%);
			border-color: darken($primary-blue, 8%);
			transform: translateY(-3px);
			box-shadow: 0 8px 16px rgba(0, 123, 255, 0.35);
		}
		&:active {
			transform: translateY(0);
			box-shadow: 0 4px 8px rgba(0, 123, 255, 0.2);
		}
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
		.top {
			height: 120px;
			margin-bottom: 25px;
		}
		.title {
			font-size: 2.2em;
			margin-left: 0; /* 确保小屏幕下也无左边距 */
		}
		.form-container {
			margin: 0 15px 30px 15px; /* 调整小屏幕下的左右外边距 */
			padding: 25px;
			box-shadow: 0 4px 10px rgba(0, 0, 0, 0.08);
		}
		.form-group {
			margin-bottom: 20px;
		}
		label {
			font-size: 1em;
			margin-bottom: 5px;
		}
		.input-field {
			padding: 12px 15px;
			font-size: 1em;
		}
		.button.is-primary {
			max-width: 100%;
			padding: 14px 20px;
			font-size: 1.1em;
			margin-top: 15px;
		}
	}

	@media (max-width: 480px) {
		.top {
			height: 100px;
			margin-bottom: 20px;
		}
		.title {
			font-size: 1.8em;
		}
		.form-container {
			padding: 20px;
			margin: 0 10px 30px 10px; /* 进一步调整小屏幕下的左右外边距 */
		}
		.input-field {
			padding: 10px 12px;
			font-size: 0.95em;
		}
		.button.is-primary {
			padding: 12px 18px;
			font-size: 1em;
		}
	}
</style>
