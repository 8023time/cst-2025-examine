import axios from "axios";
import { get } from "svelte/store";
import { information } from "$lib/stores";
import { goto } from "$app/navigation";
import { getNotyf } from "$lib/component/notyf";

const baseURL = "http://localhost:8081";

let notyf = getNotyf()

const request = axios.create({
    baseURL,
    timeout: 5000,
    withCredentials: true,
});

// 请求拦截器
request.interceptors.request.use(
    (config) => {
        config.headers["Content-Type"] = "application/json";
        return config;
    },  
    (error) => Promise.reject(error)
);

export default request;
export { baseURL };
