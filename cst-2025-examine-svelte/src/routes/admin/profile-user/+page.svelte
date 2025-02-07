<script>
	import { onMount } from 'svelte';
	import { Notyf } from 'notyf';

	let notyf;

	let oldPassword = '';
	let newPassword = '';
	let confirmPassword = '';
	let message = '';

	const submitPasswordChange = async () => {
		if (newPassword !== confirmPassword) {
			notyf.error('新密码和确认密码不一致');
			return;
		}
		// 发送请求到后端修改密码
		const response = await fetch('/api/change-password', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				oldPassword,
				newPassword
			})
		});

		const result = await response.json();

		if (response.ok) {
			notyf.success('密码修改成功');
		} else {
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

<form on:submit={submitPasswordChange} class="form-container">
	<div class="form-group">
		<label for="oldPassword">旧密码:</label>
		<input type="password" id="oldPassword" bind:value={oldPassword} required class="input-field" />
	</div>

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

	<button type="submit" class="button is-primary">修改密码</button>
</form>

<style>
	@import 'notyf/notyf.min.css';
	@import 'bulma/css/bulma.css';
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
