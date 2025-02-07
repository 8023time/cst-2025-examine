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

<style>
	@import 'notyf/notyf.min.css';
	@import 'bulma/css/bulma.css';
    .top{
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100px;
        background-color: #f5f5f5;
        color: #333;
        font-size: 24px;
        font-weight: 500;
        border-bottom: 1px solid #dcdcdc;
    }
   .title{
        margin-left: 20px;
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
