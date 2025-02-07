<script>
	import { onMount } from 'svelte';
	import { Notyf } from 'notyf';
	import { AdminAPI } from '$lib/api';

	let notyf;

	let admin_name = '';
	let password = '';
	let mobile_phone = '';

	const submitPasswordChange = async () => {
		try {			
			const response = await AdminAPI.createAdmin(admin_name, password, mobile_phone);
			if (response.data.status === 0) {
				notyf.success('管理员账户创建成功！');
			} else {
				notyf.error('管理员账户创建失败！');
			}
		} catch (error) {
			notyf.error('管理员账户创建失败！');
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
	<h1 class="title">创建管理员账户</h1>
</div>

<form on:submit={submitPasswordChange} class="form-container">
	<div class="form-group">
		<label for="admin_name">名字:</label>
		<input type="text" bind:value={admin_name} required class="input-field" />
	</div>
	<div class="form-group">
		<label for="password">密码:</label>
		<input type="password" bind:value={password} required class="input-field" />
	</div>
	<div class="form-group">
		<label for="mobile_phone">电话号码:</label>
		<input type="text" bind:value={mobile_phone} required class="input-field" />
	</div>

	<button type="submit" class="button is-primary">创建管理员账户</button>
</form>

<style>
	@import 'notyf/notyf.min.css';
	@import 'bulma/css/bulma.css';
	.top {
		display: flex;
		align-items: center;
		justify-content: center;
		height: 150px;
		background-color: #f5f5f5;
	}
	.title {
		font-size: 36px;
		font-weight: 700;
		color: #333;
	}
	.form-container {
		max-width: 450px;
		margin: 40px auto;
		padding: 30px;
		background-color: #ffffff;
		border-radius: 12px;
		box-shadow: 0 6px 15px rgba(0, 0, 0, 0.1);
	}
	.form-group {
		margin-bottom: 20px;
	}
	label {
		display: block;
		font-weight: 500;
		margin-bottom: 8px;
		color: #333;
	}
	.input-field {
		width: 100%;
		padding: 12px;
		border: 1px solid #dcdcdc;
		border-radius: 6px;
		font-size: 14px;
		transition:
			border-color 0.3s ease,
			box-shadow 0.3s ease;
	}
	.input-field:focus {
		border-color: #4caf50;
		box-shadow: 0 0 8px rgba(76, 175, 80, 0.4);
		outline: none;
	}
	.button.is-primary {
		width: 100%;
		padding: 12px;
		font-size: 16px;
		border-radius: 6px;
		cursor: pointer;
		transition:
			background-color 0.3s ease,
			transform 0.3s ease;
	}
	.form-container .form-group input,
	.form-container button {
		margin-top: 8px;
	}
</style>
