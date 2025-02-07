<script>
    import { onMount } from 'svelte';
    import { ArticleAPI } from '$lib/api';
    import { writable } from 'svelte/store';
    import { Notyf } from 'notyf';
    import { AdminAPI } from '$lib/api';

    let article = writable();
    let notyf;
    let users = {};
    let authorNames = writable({}); // 用来存储每篇文章的作者名字，按 userid 索引

    const getAuthor = async (userid) => {
        const response = await AdminAPI.getAllUsers();
        users = response.data.data;
        const userT = users.find(user => user.userid === userid);
        if (userT) {
            // 将作者名字存入 store，并且用 userid 作为键
            authorNames.update(names => ({ ...names, [userid]: userT.username }));
        }
    };

    const getarticle = async () => {
        try {
            const response = await ArticleAPI.getAll();
            notyf.success('获取文章成功');
            article.set(response.data.data);
            // 在加载文章后获取每篇文章的作者名字
            response.data.data.forEach(article => {
                getAuthor(article.userid); // 传入每篇文章的 userid
            });

        } catch (error) {
        }
    };

    const deleteArticle = async (article_id) => {
        try {		
            const response = await ArticleAPI.delete(article_id);
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
                <button class="delete-button" on:click={() => deleteArticle(articles.articleid)}>删除文章</button>
            </div>
            <h1>{articles.title}</h1>
            <div class="article-meta">
                <!-- 使用 $authorNames 来动态获取并显示作者名字 -->
                <span class="author">作者: {$authorNames[articles.userid] || '加载中...'}</span>
            </div>
            <div class="article-content">
                {@html articles.content}
            </div>
        </div>
    {/each}
{:else}
    <p>无文章,用户未发布文章...</p>
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
