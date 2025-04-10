"use client";
import {GithubFilled, LogoutOutlined, UserOutlined,} from '@ant-design/icons';
import {ProLayout,} from '@ant-design/pro-components';
import {Dropdown, message,} from 'antd';
import React from 'react';
import Image from 'next/image';
import {usePathname, useRouter} from "next/navigation";
import Link from "next/link";
import GlobalFooter from "@/components/GlobalFooter";
import "./index.css";
import {menus} from "../../../config/menus";
import {AppDispatch, RootState} from "@/stores";
import {useDispatch, useSelector} from "react-redux";
import getAccessibleMenus from "@/access/menuAccess";
import {userLogoutUsingPost} from "@/api/userController";
import {setLoginUser} from "@/stores/loginUser";
import {DEFAULT_USER} from "@/constants/user";
import SearchInput from "@/layouts/BasicLayout/components/SerchInput";

interface Props {
    children: React.ReactNode;
}

/**
 * 全局通用布局
 * @param children
 * @constructor
 */

export default function BasicLayout({children}: Props) {
    const pathname = usePathname();

    const loginUser = useSelector((state: RootState) => state.loginUser);
    const dispatch = useDispatch<AppDispatch>();
    const router = useRouter();

    /**
     * 用户注销
     */
    const userLogout = async () => {
        try {
            const res = await userLogoutUsingPost();
            message.success("已退出登录");
            dispatch(setLoginUser(DEFAULT_USER));
            router.push("/user/login");
        } catch (e: any) {
            message.error('操作失败，' + e.message);
        }
    }

    return (
        <div
            id="basic-layout"
            style={{
                height: '100vh',
                overflow: 'auto',
            }}
        >
            <ProLayout
                title="面试知识库"
                layout="top"
                logo={
                    <Image
                        src="/assets/logo.png"
                        height={32}
                        width={32}
                        alt="JIK"
                    />
                }
                location={{
                    pathname,
                }}
                avatarProps={{
                    src: loginUser.userAvatar || "/assets/notLoginUser.png",
                    size: 'small',
                    title: loginUser.userName || "未登录",
                    render: (props, dom) => {
                        if (!loginUser.id) {
                            return (
                                <div
                                    onClick={() => {
                                        router.push("/user/login");
                                    }}
                                >
                                    {dom}
                                </div>
                            );
                        }
                        return (
                            <Dropdown
                                menu={{
                                    items: [
                                        {
                                            key: 'userCenter',
                                            icon: <UserOutlined/>,
                                            label: '个人中心',
                                        },
                                        {
                                            key: 'logout',
                                            icon: <LogoutOutlined/>,
                                            label: '退出登录',
                                        },
                                    ],
                                    onClick: async (event: { key: React.Key }) => {
                                        const {key} = event;
                                        if (key === 'logout') {
                                            userLogout();
                                        } else if (key === 'userCenter') {
                                            router.push("/user/center");
                                        }
                                    },
                                }}
                            >
                                {dom}
                            </Dropdown>
                        );
                    },
                }}
                actionsRender={(props) => {
                    if (props.isMobile) return [];
                    return [
                        <a key="github" href="https://github.com/jiaking001/jik_interview" target="_blank">
                            <GithubFilled key="GithubFilled"/>,
                        </a>,
                    ];
                }}
                headerTitleRender={(logo, title, _) => {
                    return (
                        <a>
                            {logo}
                            {title}
                        </a>
                    );
                }}
                // 渲染底部栏
                footerRender={() => {
                    return <GlobalFooter/>;
                }}
                onMenuHeaderClick={(e) => console.log(e)}
                // 定义菜单
                menuDataRender={() => {
                    return getAccessibleMenus(loginUser, menus);
                }}
                // 定义了菜单项的渲染方式
                menuItemRender={(item, dom) => (
                    <Link
                        href={item.path || "/"}
                        target={item.target}
                    >
                        {dom}
                    </Link>
                )}
            >
                {children}

            </ProLayout>
        </div>
    );
};