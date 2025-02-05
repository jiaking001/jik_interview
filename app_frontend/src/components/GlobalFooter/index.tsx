import React from 'react';
import './index.css';

/**
 * 全局底部栏组件
 * @constructor
 */
export default function GlobalFooter() {
    const currentYear : number = new Date().getFullYear();

    return (
        <div
            className="global-footer"
        >
            <div>© {currentYear} 面试刷题平台</div>
            <div>
                <a href="https://github.com/jiaking001" target="_blank">
                    by JiaKing
                </a>
            </div>
            <div>
                <a href="https://github.com/jiaking001/jik_interview" target="_blank">
                    项目地址
                </a>
            </div>
        </div>
    );
};