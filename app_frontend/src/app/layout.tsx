"use client";
import {AntdRegistry} from "@ant-design/nextjs-registry";
import "./globals.css";
import BasicLayout from "@/layouts/BasicLayout";
import React, {useCallback, useEffect} from "react";
import store, {AppDispatch} from "@/stores";
import {Provider, useDispatch} from "react-redux";
import {getLoginUserUsingGet} from "@/api/userController";
import {setLoginUser} from "@/stores/loginUser";

/**
 * 全局初始化逻辑
 * @param children
 * @constructor
 */
const InitLayout: React.FC<
    Readonly<{
        children: React.ReactNode;
    }>
> = ({children}) => {
    const dispatch = useDispatch<AppDispatch>();
    // 初始全局用户状态
    const doInitLoginUser = useCallback(async () => {
        setTimeout(() => {
            const testUser = {userName: "John", id: 1, userAvatar: "/assets/logo.png"};
            dispatch(setLoginUser(testUser));
        }, 3000);
        // const res = await getLoginUserUsingGet()
        // if (res.data) {
        //     // 更新全局用户状态
        // } else {
        //     setTimeout(() => {
        //         const testUser = {userName: "John", id: 1, userAvatar: "/assets/notLoginUser.png"};
        //         dispatch(setLoginUser(testUser));
        //     }, 3000);
        // }
    }, []);

    // 只执行一次
    useEffect(() => {
        doInitLoginUser();
    }, []);
    return (
        children
    );
};

export default function RootLayout({
                                       children,
                                   }: Readonly<{
    children: React.ReactNode;
}>) {

    return (
        <html lang="zh">
        <body>
        <AntdRegistry>
            <Provider store={store}>
                <InitLayout>
                    <BasicLayout>
                        {children}
                    </BasicLayout>
                </InitLayout>
            </Provider>
        </AntdRegistry>
        </body>
        </html>
    );
}
