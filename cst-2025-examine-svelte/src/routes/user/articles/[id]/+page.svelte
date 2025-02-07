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

<style>
	@import 'bulma/css/bulma.css';

	.top {
		background-color: #f9f9f9;
		padding: 20px;
		border-bottom: 1px solid #ddd;
		text-align: center;
	}
	.title {
		font-size: 24px;
		font-weight: bold;
	}
	.container {
		max-width: 800px;
		margin: 0 auto;
		padding: 40px 20px;
		background-color: #f9f9f9;
		border-radius: 10px;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
	}
	#editor {
		width: 100%;
		min-height: 400px;
		border: 1px solid #ddd;
		border-radius: 5px;
		padding: 15px;
		background-color: #fafafa;
	}
	#toolbar {
		width: 100%;
		border-bottom: 1px solid #ddd;
		margin-bottom: 15px;
	}
	.button.is-primary {
		width: 100%;
		font-size: 16px;
		padding: 12px;
		border-radius: 5px;
	}
	.input.is-large {
		font-size: 18px;
		padding: 12px;
		border-radius: 5px;
	}
	.label {
		font-size: 16px;
		font-weight: bold;
	}
	.field {
		margin-bottom: 20px;
	}
</style>
