<script>
	import { onMount } from 'svelte';
	import { ArticleAPI } from '$lib/api';
	import { writable } from 'svelte/store';
	import { Notyf } from 'notyf';
	import { information } from '$lib/stores';
	import { goto } from '$app/navigation';

	let article = writable([]);
	let notyf;

	let name = JSON.parse(localStorage.getItem('user')).username;

	const getarticle = async () => {
		try {
			const response = await ArticleAPI.getuserAll();
			notyf.success('获取文章成功');
			article.set(response.data.data);
		} catch (error) {
			notyf.error('获取文章失败');
		}
	};

	const deleteArticle = async (artictle_id) => {
		try {
			console.log(artictle_id);
			const response = await ArticleAPI.delete(artictle_id);
			notyf.success('删除文章成功');
			getarticle();
		} catch (error) {
			notyf.error('删除文章失败');
		}
	};

	onMount(() => {
		notyf = new Notyf({
			duration: 3000,
			className: 'x-notification',
			position: { x: 'right', y: 'top' }
		});
		getarticle();
	});
</script>

{#if $article && $article.length > 0}
	{#each $article as articles}
		<div class="article-container">
			<div class="article-actions">
				<a href={`/user/articles/${articles.id}`} class="edit-button">编辑文章</a>
				<button class="delete-button" on:click={deleteArticle(articles.id)}>删除文章</button>
			</div>
			<h1>{articles.title}</h1>
			<div class="article-meta">
				<span class="author">作者: {name || '匿名'}</span>
			</div>
			<div class="article-content">
				{@html articles.content}
			</div>
		</div>
	{/each}
{:else}
	<p>无文章,快去发布吧...</p>
{/if}

<style lang="scss">
	@import 'notyf/notyf.min.css';
	@import 'bulma/css/bulma.css';

	/* 定义颜色变量 */
	$primary-blue: #007bff;
	$danger-red: #e74c3c;
	$info-blue: #3498db;
	$text-dark: #2c3e50;
	$text-medium: #34495e;
	$text-light: #7f8c8d;
	$bg-light: #f5f5f5;
	$white: #fff;
	$border-light: #e0e0e0;

	.article-container {
		max-width: 100%; /* 允许铺满左右 */
		margin: 40px 30px; /* 调整左右外边距，保留上下 */
		padding: 40px; /* 增加内边距 */
		background-color: $white;
		border-radius: 12px; /* 更大的圆角 */
		box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1); /* 更柔和的阴影 */
		transition:
			transform 0.3s cubic-bezier(0.25, 0.8, 0.25, 1),
			box-shadow 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);

		&:hover {
			transform: translateY(-8px); /* 更明显的上浮效果 */
			box-shadow: 0 12px 25px rgba(0, 0, 0, 0.15); /* 悬停时阴影更明显 */
		}
	}
	h1 {
		font-size: 2.8rem; /* 增大标题字号 */
		color: $text-dark;
		margin-bottom: 25px; /* 调整底部间距 */
		font-weight: 700; /* 更粗的字重 */
		line-height: 1.3; /* 调整行高 */
	}
	.article-meta {
		display: flex;
		justify-content: space-between;
		align-items: center; /* 垂直居中 */
		margin-bottom: 25px; /* 调整底部间距 */
		font-size: 15px; /* 调整字号 */
		color: $text-light;
		padding-bottom: 15px;
		border-bottom: 1px solid $border-light; /* 添加底部边框 */

		.author {
			font-weight: 500;
			color: $text-medium;
		}
	}
	.article-content {
		font-size: 17px; /* 增大字号 */
		line-height: 1.7; /* 调整行高 */
		color: $text-medium;
		word-wrap: break-word; /* 确保长单词换行 */
		img {
			max-width: 100%;
			height: auto;
			border-radius: 8px;
			margin: 15px 0;
		}
		pre {
			background-color: #f4f4f4;
			padding: 15px;
			border-radius: 8px;
			overflow-x: auto;
			font-size: 0.95em;
			line-height: 1.4;
		}
		code {
			background-color: #eee;
			padding: 2px 5px;
			border-radius: 4px;
			font-family: 'Fira Code', 'Consolas', monospace; /* 假设等宽字体可用 */
		}
	}
	.article-actions {
		display: flex;
		justify-content: flex-end;
		margin-bottom: 30px; /* 调整底部间距 */

		.edit-button,
		.delete-button {
			margin-left: 15px; /* 增加左侧间距 */
			padding: 10px 20px; /* 调整内边距 */
			border: none;
			border-radius: 6px; /* 调整圆角 */
			cursor: pointer;
			font-weight: 600;
			transition: all 0.3s ease; /* 更平滑的过渡 */
			text-decoration: none; /* 移除链接下划线 */
			display: inline-flex;
			align-items: center;
			justify-content: center;

			&:hover {
				transform: translateY(-2px); /* 轻微上浮 */
				box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
			}
		}
		.edit-button {
			background-color: $info-blue; /* 使用信息蓝 */
			color: $white;
			&:hover {
				background-color: darken($info-blue, 8%);
			}
		}
		.delete-button {
			background-color: $danger-red; /* 使用危险红 */
			color: $white;
			&:hover {
				background-color: darken($danger-red, 8%);
			}
		}
	}
	p {
		text-align: center;
		color: $primary-blue; /* 使用主蓝色 */
		font-size: 1.5rem; /* 增大字号 */
		font-weight: 700; /* 更粗的字重 */
		padding: 50px 20px; /* 增加内边距 */
		background-color: $white;
		border-radius: 8px;
		box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
		margin: 50px auto; /* 居中显示 */
		max-width: 600px;
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
	@media (max-width: 768px) {
		.article-container {
			margin: 20px 15px; /* 调整外边距 */
			padding: 25px; /* 调整内边距 */
			box-shadow: 0 4px 10px rgba(0, 0, 0, 0.08);
			&:hover {
				transform: translateY(-3px);
				box-shadow: 0 6px 12px rgba(0, 0, 0, 0.1);
			}
		}
		h1 {
			font-size: 2.2rem; /* 调整字号 */
			margin-bottom: 20px;
		}
		.article-meta {
			flex-direction: column; /* 垂直堆叠 */
			align-items: flex-start;
			gap: 5px;
			font-size: 14px;
			margin-bottom: 20px;
		}
		.article-actions {
			flex-direction: column; /* 按钮垂直堆叠 */
			align-items: flex-end;
			gap: 10px; /* 按钮间距 */
			margin-bottom: 20px;
			.edit-button,
			.delete-button {
				width: auto; /* 按钮宽度自适应 */
				margin-left: 0;
				padding: 10px 18px;
				font-size: 0.95em;
			}
		}
		.article-content {
			font-size: 16px;
		}
		p {
			margin: 30px 15px;
			font-size: 1.3rem;
			padding: 30px 15px;
		}
	}

	@media (max-width: 480px) {
		.article-container {
			padding: 20px;
		}
		h1 {
			font-size: 1.8rem;
		}
		.article-actions {
			.edit-button,
			.delete-button {
				padding: 8px 15px;
				font-size: 0.9em;
			}
		}
	}
</style>
