<script>
	import { Notyf } from 'notyf';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { UserAPI } from '$lib/api';
	import { AdminAPI } from '$lib/api';
	import {information}  from '$lib/stores';

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
	.input {
		width: unset;
		max-width: unset;
	}
	.main {
		display: flex;
		flex-direction: column;
		align-items: center;

		div {
			padding: 1em;
		}
	}
	.role-switch {
		display: flex;
		gap: 1em;

		.button {
			border: none;
			cursor: pointer;
			padding: 0.5em 1em;
		}
		.active {
			background-color: #007bff;
			color: white;
		}
	}
	.toggle-mode {
		margin-top: 1em;
		font-size: 14px;
		a {
			color: #007bff;
			cursor: pointer;
		}
	}
</style>
