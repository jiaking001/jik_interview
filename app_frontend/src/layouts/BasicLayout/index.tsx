"use client";
import {
    GithubFilled,
    LogoutOutlined,
    SearchOutlined,
} from '@ant-design/icons';
import type {ProSettings} from '@ant-design/pro-components';
import {ProLayout,} from '@ant-design/pro-components';
import {Dropdown, Input,} from 'antd';
import React, {useState} from 'react';
import Image from 'next/image';
import {usePathname} from "next/navigation";
import Link from "next/link";
import GlobalFooter from "@/components/GlobalFooter";
import "./index.css";
import menus from "../../../config/menu";

/**
 * 搜索框
 * @constructor
 */
const SearchInput = () => {
    return (
        <div
            key="SearchOutlined"
            aria-hidden
            style={{
                display: 'flex',
                alignItems: 'center',
                marginInlineEnd: 24,
            }}
            onMouseDown={(e) => {
                e.stopPropagation();
                e.preventDefault();
            }}
        >
            <Input
                style={{
                    borderRadius: 4,
                    marginInlineEnd: 12,
                }}
                prefix={
                    <SearchOutlined/>
                }
                placeholder="搜索题目"
                variant="borderless"
            />
        </div>
    );
};

interface Props {
    children: React.ReactNode;
}

export default function BasicLayout({children}: Props) {
    const pathname = usePathname();
    return (
        <div
            id="basic-layout"
            style={{
                height: '100vh',
                overflow: 'auto',
            }}
        >
            <ProLayout
                title="面试刷题平台"
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
                    src: 'https://gw.alipayobjects.com/zos/antfincdn/efFD%24IOql2/weixintupian_20170331104822.jpg',
                    size: 'small',
                    title: "jiaking",
                    render: (props, dom) => {
                        return (
                            <Dropdown
                                menu={{
                                    items: [
                                        {
                                            key: 'logout',
                                            icon: <LogoutOutlined/>,
                                            label: '退出登录',
                                        },
                                    ],
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
                        <SearchInput key="search"/>,
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
                    return menus;
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