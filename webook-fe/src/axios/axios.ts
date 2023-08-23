/*
 * @Author: hugo lee hugo2lee@gmail.com
 * @Date: 2023-08-07 09:37
 * @LastEditors: hugo lee hugo2lee@gmail.com
 * @LastEditTime: 2023-08-23 16:16
 * @FilePath: /geektime-basic-go/webook-fe/src/axios/axios.ts
 * @Description: 
 * 
 * Copyright (c) 2023 by hugo, All Rights Reserved. 
 */
import axios from "axios";
const instance = axios.create({
    // 这边记得修改你对应的配置文件
    baseURL:  "http://live.webook.com:81",
    withCredentials: true
})


instance.interceptors.response.use(function (resp) {
    const newToken = resp.headers["x-jwt-token"]
    console.log("resp headers", resp.headers)
    console.log("token" + newToken)
    if (newToken) {
        localStorage.setItem("token", newToken)
    }
    if (resp.status == 401) {
        window.location.href="/users/login"
    }
    return resp
}, (err) => {
    console.log(err)
    if (err.response.status == 401) {
        window.location.href="/users/login"
    }
    return err
})

// 在这里让每一个请求都加上 authorization 的头部
instance.interceptors.request.use((req) => {
    const token = localStorage.getItem("token")
    req.headers.setAuthorization("Bearer " + token, true)
    return req
}, (err) => {
    console.log(err)
})

export default instance