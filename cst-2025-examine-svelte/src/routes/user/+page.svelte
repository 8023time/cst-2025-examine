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
    .article-container {
        max-width: 800px;
        margin: 50px auto;
        padding: 30px;
        background-color: #fff;
        border-radius: 8px;
        box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
        transition: transform 0.2s; 
        &:hover {
            transform: translateY(-5px); 	
        }
    }
    .article-actions {
        display: flex;
        justify-content: flex-end;
        margin-bottom: 20px;
        .edit-button,
        .delete-button {
            margin-left: 10px;
            padding: 8px 16px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            transition: background-color 0.2s;

            &:hover {
                opacity: 0.8;
            }
        }
        .edit-button {
            background-color: #3498db;
            color: #fff;
        }
        .delete-button {
            background-color: #e74c3c;
            color: #fff;
        }
    }
    .article-meta {
        display: flex;
        justify-content: space-between;
        margin-bottom: 20px;
        font-size: 14px;
        color: #7f8c8d;
    }
    .article-content {
        font-size: 16px;
        line-height: 1.8;
        color: #34495e;
    }
    h1 {
        font-size: 2.5rem;
        color: #2c3e50;
        margin-bottom: 20px;
        font-weight: bold;
    }
    p {
        text-align: center;
        color: #3498db;
        font-size: 1.2rem;
        font-weight: bold;
    }
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
            margin: 20px;
            padding: 20px;
        }
        h1 {
            font-size: 2rem;
        }
    }
</style>