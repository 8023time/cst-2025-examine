import request from "$lib/request"

export const UserAPI = {
    register: (username, password, mobilephone) => request({
        url: "/user/register",
        method: "POST",
        data: {
            username, password, mobilephone
        }
    }),
    login: (username, password) => request({
        url: "/user/login",
        method: "POST",
        data: {
            username, password
        }
    }),
    logout: () => request({
        url: "/user/logout",
        method: "POST"
    }),
    updatePassword: (currentpassword, newpassword) => request({
        url: "/user/profile",
        method: "POST",
        data: {
            currentpassword, newpassword
        }
    }),
};


export const ArticleAPI = {
    getAll: () => request({
        url: "/articles",
        method: "GET"
    }),
    getuserAll: () => request({
        url: "/user/articles",
        method: "GET"
    }),
    create: (title, content) => request({
        url: "/user/articles/create",
        method: "POST",
        data: {
            title, content
        }
    }),
    update: (articleid, title, content) => request({
        url: `/user/articles/edit`,
        method: "PUT",
        data: {
            articleid,title, content
        }
    }),
    delete: (article_id) => request({
        url: `/articles/delete/${article_id}`,
        method: "DELETE"
    }),
};


export const AdminAPI = {
    login: (username, password) => request({
        url: "/admin/login",
        method: "POST",
        data: {
            username, password
        }
    }),
    createAdmin: (username, password, mobilephone) => request({
        url: "/admin/create-admin",
        method: "POST",
        data: {
            username, password, mobilephone
        }
    }),
    getAllUsers: () => request({
        url: "/admin/users",
        method: "GET"
    }),
    resetUserPassword: (userid, newpassword) => {
        return request({
            url: "/admin/users/profile",
            method: "POST",
            data: {
                userid, newpassword
            }
        });
    },
    resetMyPassword: (oldpassword, newpassword) => {
        return request({
            url: "/admin/profile",
            method: "POST",
            data: {
                oldpassword, newpassword
            }
        });
    },
    deleteUser: (userid) => request({
        url: `/admin/users/delete/${userid}`,
        method: "DELETE"
    }),
};
