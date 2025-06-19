<script>
	import { Notyf } from 'notyf';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { UserAPI } from '$lib/api';
	import { AdminAPI } from '$lib/api';
	import { information } from '$lib/stores';

	let notyf;

	let mode = $state('login');
	let role = $state('user');

	let username = $state('');
	let password = $state('');
	let confirmPwd = $state('');
	let mobile_phone = $state('');

	const login = async () => {
		try {
			if (role === 'user') {
				const response = await UserAPI.login(username, password);
				information.setUserInformation(response.data.data);
				notyf.success('登录成功！');
				goto('/user');
			} else {
				const response = await AdminAPI.login(username, password);
				console.log(response.data);
				information.setUserInformation(response.data.data);
				notyf.success('登录成功！');
				goto('/admin');
			}
		} catch (err) {
			console.log(err);
			notyf.error('登录失败!');
		}
	};

	const register = async () => {
		if (password !== confirmPwd) {
			notyf.error('两次输入的密码不一致!');
			return;
		}
		try {
			await UserAPI.register(username, password, mobile_phone);
			notyf.success('注册成功，请登录!');
			mode = 'login';
		} catch (err) {
			notyf.error('注册失败!');
			console.log(err);
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

<div class="main">
	<div class="title h1">{mode === 'login' ? 'login' : 'register'}</div>
	<!-- 角色切换 -->
	<div class="role-switch">
		<button class="button is-small" on:click={() => (role = 'user')} class:active={role === 'user'}>
			普通用户
		</button>
		<button
			class="button is-small"
			on:click={() => (role = 'admin')}
			class:active={role === 'admin'}
			disabled={mode === 'register'}
		>
			管理员
		</button>
	</div>
	<div>
		name: <input class="input" bind:value={username} type="text" />
	</div>
	<div>
		password: <input class="input" bind:value={password} type="password" />
	</div>

	{#if mode === 'register'}
		<div>
			confirm pwd: <input class="input" bind:value={confirmPwd} type="password" />
		</div>
		<div>
			mobile_phone: <input class="input" bind:value={mobile_phone} type="text" />
		</div>
	{/if}

	<div>
		{#if mode === 'login'}
			<button class="button is-primary" on:click={login}>登录</button>
		{:else}
			<button class="button is-primary" on:click={register}>注册</button>
		{/if}
	</div>
	<!-- 模式切换 -->
	<div class="toggle-mode">
		{#if mode === 'login'}
			<p>还没有账号？ <a href="#" on:click={() => (mode = 'register')}>注册</a></p>
		{:else}
			<p>已有账号？ <a href="#" on:click={() => (mode = 'login')}>去登录</a></p>
		{/if}
	</div>
</div>

<style lang="scss">
	@import 'notyf/notyf.min.css';
	@import 'bulma/css/bulma.css';

	.main {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center; /* 垂直居中内容 */
		padding: 60px 45px; /* 增加卡片内边距 */
		box-sizing: border-box;
		background-color: #ffffff; /* 卡片背景色 */
		border-radius: 25px; /* 更大的圆角 */
		box-shadow: 0 25px 50px rgba(0, 0, 0, 0.15); /* 更深邃的阴影 */
		width: 100%;
		max-width: 550px; /* 略微增加卡片最大宽度 */
		margin: 80px auto; /* 卡片在页面上水平居中并添加顶部/底部外边距 */

		.title {
			font-size: 3.8em; /* 标题字号更大 */
			margin-bottom: 55px; /* 增加标题下方间距 */
			color: #2c3e50; /* 深色标题 */
			font-weight: 700; /* 更粗的字重 */
			text-align: center;
			letter-spacing: -0.05em; /* 调整字间距 */
			line-height: 1.1; /* 调整行高 */
		}

		/* 输入字段容器的样式 (.main 的直接子 div) */
		& > div {
			padding: 1.4em 0; /* 调整垂直内边距 */
			width: 100%;
			max-width: 420px; /* 限制内部元素宽度 */
			margin: 0 auto; /* 在 .main 内部居中 */
			text-align: left; /* 文本左对齐 */
		}
	}

	.input {
		width: 100%; /* 输入框占据其父 div 的全部宽度 */
		max-width: 100%; /* 确保不超过容器 */
		border-radius: 16px; /* 更多圆角 */
		padding: 18px 25px; /* 更多内边距 */
		font-size: 1.3em; /* 略微增大字号 */
		transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1); /* 更平滑的过渡效果 */
		border: 1px solid #c0c0c0; /* 默认边框，更柔和 */
		background-color: #ffffff; /* 更纯白的背景色 */
		color: #333; /* 输入文本颜色 */

		&:focus {
			border-color: #007bff; /* 聚焦时边框颜色 */
			box-shadow: 0 0 0 0.35em rgba(0, 123, 255, 0.5); /* 聚焦时阴影，更明显 */
			outline: none; /* 移除默认轮廓 */
		}

		&::placeholder {
			color: #b0b0b0; /* 占位符颜色 */
			opacity: 0.95; /* 略微透明 */
		}
	}

	.role-switch {
		display: flex;
		justify-content: center; /* 按钮居中 */
		gap: 1.5em; /* 略大的间距 */
		margin-bottom: 55px; /* 角色切换下方更多空间 */

		.button {
			border: 1px solid #c0c0c0; /* 添加边框 */
			cursor: pointer;
			padding: 1.2em 2.6em; /* 更多内边距 */
			border-radius: 35px; /* 更大的椭圆形按钮 */
			background-color: #ebebeb; /* 非激活状态的浅色背景 */
			color: #666; /* 默认文本颜色 */
			transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1); /* 更平滑的过渡 */
			font-weight: 600;
			box-shadow: 0 5px 10px rgba(0, 0, 0, 0.1);

			&:hover {
				background-color: #d8d8d8;
				color: #333;
				box-shadow: 0 10px 20px rgba(0, 0, 0, 0.15);
				transform: translateY(-5px);
			}
		}
		.active {
			background-color: #007bff; /* 蓝色主色 */
			border-color: #007bff;
			color: white;
			font-weight: 700;
			box-shadow: 0 12px 30px rgba(0, 123, 255, 0.7); /* 更明显的阴影 */
			transform: translateY(-5px);
		}
	}

	/* 主要登录/注册按钮 (Bulma 的 .button.is-primary) */
	.button.is-primary {
		width: 100%; /* 全宽按钮 */
		max-width: 420px; /* 与输入字段宽度一致 */
		padding: 22px 45px; /* 更多内边距 */
		font-size: 1.5em;
		font-weight: 700;
		border-radius: 16px;
		transition:
			background-color 0.4s cubic-bezier(0.25, 0.8, 0.25, 1),
			transform 0.3s cubic-bezier(0.25, 0.8, 0.25, 1),
			box-shadow 0.4s cubic-bezier(0.25, 0.8, 0.25, 1); /* 更平滑的动画效果 */
		background-color: #007bff; /* 确保主色 */
		border-color: #007bff;
		color: white;
		box-shadow: 0 15px 30px rgba(0, 123, 255, 0.5); /* 初始阴影 */

		&:hover {
			background-color: darken(#007bff, 15%); /* 悬停时颜色变深 */
			border-color: darken(#007bff, 15%);
			transform: translateY(-7px); /* 轻微上浮效果 */
			box-shadow: 0 20px 40px rgba(0, 123, 255, 0.7); /* 悬停时阴影更明显 */
		}

		&:active {
			transform: translateY(0); /* 按下效果 */
			box-shadow: 0 10px 20px rgba(0, 123, 255, 0.4); /* 按下时阴影变小 */
		}
	}

	.toggle-mode {
		margin-top: 50px; /* 切换模式上方更多空间 */
		font-size: 1.25em; /* 略微更大的字号 */
		color: #888;

		p {
			margin-bottom: 0.8em; /* 段落间距 */
			line-height: 1.7;
		}

		a {
			color: #007bff;
			cursor: pointer;
			text-decoration: none; /* 移除下划线 */
			font-weight: 600; /* 更粗的字重 */
			transition: color 0.3s ease; /* 颜色过渡 */

			&:hover {
				color: darken(#007bff, 25%); /* 悬停时颜色变深 */
				text-decoration: underline; /* 悬停时显示下划线 */
			}
		}
	}

	/* 针对小屏幕的响应式调整 */
	@media (max-width: 600px) {
		.main {
			margin: 40px auto; /* 小屏幕上减少顶部/底部外边距 */
			padding: 50px 35px; /* 减少内边距 */
			max-width: 98%; /* 允许卡片占据更多宽度 */
			box-shadow: 0 18px 36px rgba(0, 0, 0, 0.1); /* 小屏幕阴影更小 */
			border-radius: 20px;
		}

		.main > div {
			max-width: 95%; /* 内部元素占据更多宽度 */
			padding: 1.2em 0;
		}

		.title {
			font-size: 3em; /* 小屏幕标题字号更小 */
			margin-bottom: 45px;
		}

		.role-switch {
			flex-direction: column; /* 按钮垂直堆叠 */
			gap: 1.2em;
			margin-bottom: 45px;
			.button {
				width: 100%;
				padding: 1.1em 2.2em; /* 调整内边距 */
				font-size: 1.15em;
			}
		}
		.button.is-primary {
			max-width: 95%; /* 小屏幕按钮宽度调整 */
			padding: 20px 35px;
			font-size: 1.35em;
		}

		.input {
			padding: 16px 22px; /* 调整输入框内边距 */
			font-size: 1.2em;
		}

		.toggle-mode {
			font-size: 1.15em;
			margin-top: 40px;
		}
	}
</style>
