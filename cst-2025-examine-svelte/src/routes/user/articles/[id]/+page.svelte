<script>
	import { onMount, onDestroy } from 'svelte';
	import { ArticleAPI } from '$lib/api';
	import { Notyf } from 'notyf';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	let id = page.params.id;
	let editorInstance = null;
	$: title = '';
	$: content = '';
	let notyf = new Notyf({
		duration: 3000,
		className: 'x-notification',
		position: { x: 'right', y: 'top' }
	});

	$: id_article = {};
	$: users_article = {};

	const get_id_article = async () => {
		try {
			const response = await ArticleAPI.getuserAll();
			users_article = response.data.data;
			id_article = users_article.find((item) => item.id == id);
			console.log(id_article);
			if (id_article) {
				title = id_article.title;
				content = id_article.content;
			}
		} catch (error) {
			notyf.error('文章获取失败!');
		}
	};

	onMount(async () => {
		// 加载编辑器
		if (typeof window !== 'undefined') {
			const { createEditor, createToolbar } = await import('@wangeditor/editor');
			import('@wangeditor/core/dist/css/style.css');
			const editor = createEditor({
				selector: '#editor',
				html: '<p><br></p>',
				config: {
					placeholder: '请输入内容...',
					MENU_CONF: {
						uploadImage: {
							fieldName: 'image',
							base64LimitSize: 10 * 1024 * 1024,
							server: '/upload',
							maxFileSize: 5 * 1024 * 1024,
							onSuccess(file, res) {
								editor.insertHtml(`<img src="${res.url}" alt="image"/>`);
							},
							onError(err) {
								console.error('上传失败:', err);
							}
						}
					},
					onChange(editor) {
						content = editor.getHtml(); // 获取富文本内容
					}
				},
				mode: 'default'
			});
			const toolbar = createToolbar({
				editor,
				selector: '#toolbar',
				config: {},
				mode: 'default'
			});
			editorInstance = editor;
		}

		await get_id_article(); // 获取文章数据

		// 如果获取到了文章内容，将其填充到编辑器中
		if (id_article && editorInstance) {
			editorInstance.setHtml(id_article.content); // 填充富文本编辑器内容
		}
	});

	onDestroy(() => {
		if (editorInstance) {
			editorInstance.destroy();
		}
	});

	const handleSubmit = async () => {
		try {
			const response = await ArticleAPI.update(parseInt(id), title, content);
			notyf.success('文章编辑成功!');
			goto('/user');
		} catch (error) {
			notyf.error('文章编辑失败!');
		}
	};
</script>

<div class="top">
	<h1 class="title">编辑文章</h1>
</div>

<div class="container">
	<div class="field">
		<label for="title" class="label">文章标题:</label>
		<div class="control">
			<input
				id="title"
				type="text"
				bind:value={title}
				placeholder="请输入文章标题"
				class="input is-large"
			/>
		</div>
	</div>
	<div id="toolbar"></div>
	<div id="editor"></div>
	<div class="field">
		<div class="control">
			<button class="button is-primary is-large" on:click={handleSubmit}>提交</button>
		</div>
	</div>
</div>

<style lang="scss">
	@import 'bulma/css/bulma.css';

	/* 定义颜色变量 */
	$primary-blue: #007bff;
	$success-green: #28a745;
	$text-dark: #333;
	$text-medium: #555;
	$bg-light: #f8f9fa;
	$white: #fff;
	$border-light: #e0e0e0;

	.top {
		background-color: $bg-light; /* 使用浅色背景 */
		padding: 25px;
		border-bottom: 1px solid $border-light; /* 更柔和的底部边框 */
		text-align: center;
		margin-bottom: 30px; /* 增加底部间距 */
		box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05); /* 轻微阴影 */
	}
	.title {
		font-size: 2.5em; /* 增大字号 */
		font-weight: 700; /* 更粗的字重 */
		color: $text-dark;
		line-height: 1.2;
		margin-bottom: 0; /* 移除 Bulma 默认 margin */
	}
	.container {
		max-width: 100%; /* 允许铺满左右 */
		margin: 0 30px; /* 调整左右外边距，保留上下 */
		padding: 40px; /* 增加内边距 */
		background-color: $white;
		border-radius: 12px; /* 更大的圆角 */
		box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1); /* 柔和阴影 */
		margin-bottom: 50px; /* 底部间距 */
	}
	#editor {
		width: 100%;
		min-height: 400px;
		border: 1px solid $border-light; /* 使用变量边框 */
		border-radius: 8px; /* 统一圆角 */
		padding: 15px; /* 确保内容有内边距 */
		background-color: $white; /* 白色背景 */
		box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.05); /* 内部阴影 */
		font-size: 1em; /* 调整字体大小 */
		line-height: 1.6;
		color: $text-medium;
		/* 覆盖 wangeditor 默认样式 */
		:global(.w-e-text-container) {
			line-height: 1.6 !important;
			font-size: 1em !important;
			color: $text-medium !important;
		}
	}
	#toolbar {
		width: 100%;
		border: 1px solid $border-light; /* 边框 */
		border-bottom: none; /* 底部无边框，与编辑器连接 */
		border-top-left-radius: 8px;
		border-top-right-radius: 8px;
		margin-bottom: 0; /* 移除默认间距 */
		background-color: #f5f5f5;
		padding: 8px 0; /* 调整内边距 */
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05); /* 轻微阴影 */
		z-index: 100; /* 确保在最上层 */
		position: relative; /* 确保 z-index 有效 */
	}
	.button.is-primary {
		width: 100%;
		max-width: 250px; /* 限制按钮最大宽度 */
		font-size: 1.1em; /* 调整字号 */
		padding: 14px 25px; /* 调整内边距 */
		border-radius: 8px; /* 统一圆角 */
		background-color: $primary-blue;
		border-color: $primary-blue;
		color: $white;
		font-weight: 600;
		transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1); /* 平滑过渡 */
		box-shadow: 0 6px 12px rgba(0, 123, 255, 0.25);

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
	.input.is-large {
		font-size: 1.15em; /* 调整字号 */
		padding: 14px 18px; /* 调整内边距 */
		border-radius: 8px; /* 统一圆角 */
		border: 1px solid $border-light; /* 统一边框 */
		transition: all 0.3s ease;

		&:focus {
			border-color: $primary-blue;
			box-shadow: 0 0 0 0.2em rgba(0, 123, 255, 0.3);
			outline: none;
		}
	}
	.label {
		font-size: 1.05em; /* 调整字号 */
		font-weight: 600;
		color: $text-dark;
		margin-bottom: 8px; /* 标签下方间距 */
	}
	.field {
		margin-bottom: 25px; /* 调整字段间距 */
		/* 确保提交按钮的field没有过多的底部间距 */
		&:last-child {
			margin-bottom: 0;
			display: flex;
			justify-content: center;
			margin-top: 30px; /* 按钮上方间距 */
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
			padding: 20px 15px;
			margin-bottom: 20px;
		}
		.title {
			font-size: 2em;
		}
		.container {
			margin: 0 15px 30px 15px; /* 调整小屏幕下的左右外边距 */
			padding: 25px;
			box-shadow: 0 4px 10px rgba(0, 0, 0, 0.08);
		}
		#editor {
			min-height: 300px;
		}
		.button.is-primary {
			max-width: 100%;
			font-size: 1em;
			padding: 12px 20px;
		}
		.input.is-large {
			font-size: 1em;
			padding: 12px 15px;
		}
		.label {
			font-size: 0.95em;
		}
		.field {
			margin-bottom: 20px;
			&:last-child {
				margin-top: 20px;
			}
		}
	}

	@media (max-width: 480px) {
		.container {
			padding: 20px;
			margin: 0 10px 30px 10px; /* 进一步调整小屏幕下的左右外边距 */
		}
		.title {
			font-size: 1.8em;
		}
		.button.is-primary {
			padding: 10px 18px;
			font-size: 0.95em;
		}
		.input.is-large {
			padding: 10px 12px;
		}
	}
</style>
