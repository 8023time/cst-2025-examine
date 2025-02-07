<script>
	import { onMount } from 'svelte';
	import { AdminAPI } from '$lib/api';
	import { Notyf } from 'notyf';
	import { goto } from '$app/navigation';

	let notyf;

	let users = []

	const getusersAll = async () => {
		try {
			const response = await AdminAPI.getAllUsers();
			users = response.data.data
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

<style>
	.container {
		max-width: 600px;
		margin: 0 auto;
		padding: 20px;
		background-color: #f9f9f9;
		border-radius: 10px;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
	}
	h1 {
		font-size: 24px;
		margin-bottom: 20px;
		text-align: center;
	}
	.user-item {
		display: flex;
		justify-content: space-between;
		padding: 10px;
		margin-bottom: 10px;
		background-color: #fff;
		border: 1px solid #ddd;
		border-radius: 5px;
	}
	.button-group {
		display: flex;
		gap: 10px;
	}
	.change-button {
		background-color: #4caf50;
		color: white;
		border: none;
		padding: 5px 15px;
		border-radius: 5px;
		cursor: pointer;
		transition: background-color 0.3s ease;
	}
	.change-button:hover {
		background-color: #3e8e41;
	}
	.delete-button {
		background-color: #ff4d4f;
		color: white;
		border: none;
		padding: 5px 15px;
		border-radius: 5px;
		cursor: pointer;
		transition: background-color 0.3s ease;
	}
	.delete-button:hover {
		background-color: #f44336;
	}
</style>
