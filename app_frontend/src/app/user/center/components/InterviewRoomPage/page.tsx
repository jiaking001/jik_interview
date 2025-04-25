import { Card, message, Tag } from "antd";
import type { ActionType, ProColumns } from "@ant-design/pro-components";
import { ProTable } from "@ant-design/pro-components";
import React, { useRef, useState } from "react";
import { listMockInterviewVoByPageUsingPost } from "@/api/mockInterviewController"; // 假设这是获取模拟面试记录的 API
import Link from "next/link"; // 确保导入 Link 组件
import "./index.css";

interface MockInterview {
    id: string;
    jobPosition: string;
    workExperience: string;
    difficulty: string;
    status: number; // 0: 未开始, 1: 进行中, 2: 已结束
    createdAt: string;
}

interface Props {
    // 默认值（用于展示服务端渲染的数据）
    defaultInterviewList?: MockInterview[];
    defaultTotal?: number;
    // 默认搜索条件
    defaultSearchParams?: Record<string, any>;
}

const InterviewRoomPage: React.FC<Props> = (props) => {
    const { defaultInterviewList, defaultTotal, defaultSearchParams = {} } = props;
    const actionRef = useRef<ActionType>();
    const [interviewList, setInterviewList] = useState<MockInterview[]>(defaultInterviewList || []);
    const [total, setTotal] = useState<number>(defaultTotal || 0);
    const [init, setInit] = useState<boolean>(true);

    const columns: ProColumns<MockInterview>[] = [
        {
            title: "工作岗位",
            dataIndex: "jobPosition",
            valueType: "text",
            // hideInSearch: true,
            render: (_, record) => <Link href={`/mockInterview/chat/${record.id}`}>{record.jobPosition}</Link>,
        },
        {
            title: "工作年限",
            dataIndex: "workExperience",
            valueType: "text",
            hideInSearch: true,
        },
        {
            title: "难度",
            dataIndex: "difficulty",
            valueType: "text",
            hideInSearch: true,
        },
        {
            title: "状态",
            dataIndex: "status",
            valueType: "select",
            valueEnum: {
                0: { text: "未开始", status: "Default" },
                1: { text: "进行中", status: "Processing" },
                2: { text: "已结束", status: "Success" },
            },
            render: (_, record) => {
                // 设置默认值
                const statusValue = record.status || 0; // 如果 record.status 未定义，则默认为 0
                const statusText = statusValue === 0 ? "未开始" :
                    statusValue === 1 ? "进行中" : "已结束";
                const statusColor = statusValue === 0 ? "orange" :
                    statusValue === 1 ? "green" : "red";
                return (
                    <Tag color={statusColor}>
                        {statusText}
                    </Tag>
                );
            },
        },
        {
            title: "创建时间",
            dataIndex: "createTime",
            valueType: "dateTime",
            hideInSearch: true,
        },
    ];

    return (
        <div className="interview-table">
            <ProTable<MockInterview>
                actionRef={actionRef}
                size="large"
                search={{
                    labelWidth: "auto",
                    // 添加搜索表单
                    formItemProps: (form, { rowIndex }) => ({
                        style: { marginBottom: 0 },
                    }),
                }}
                form={{
                    initialValues: defaultSearchParams,
                    // 添加搜索表单项
                    layout: "inline",
                    fields: [
                        {
                            name: "jobPosition",
                            label: "工作岗位",
                            type: "text",
                            placeholder: "请输入工作岗位",
                        },
                        {
                            name: "workExperience",
                            label: "工作年限",
                            type: "text",
                            placeholder: "请输入工作年限",
                        },
                        {
                            name: "difficulty",
                            label: "难度",
                            type: "text",
                            placeholder: "请输入难度",
                        },
                    ],
                }}
                dataSource={interviewList}
                pagination={{
                    pageSize: 12,
                    showTotal: (total) => `总共 ${total} 条`,
                    showSizeChanger: false,
                    total,
                }}
                request={async (params, sort, filter) => {
                    if (init) {
                        setInit(false);
                        if (defaultInterviewList && defaultTotal) {
                            return;
                        }
                    }

                    const sortField = Object.keys(sort)?.[0] || "createdAt";
                    const sortOrder = sort?.[sortField] || "descend";

                    const { data, code } = await listMockInterviewVoByPageUsingPost({
                        ...params,
                        sortField,
                        sortOrder,
                        ...filter,
                    });

                    const newData = data?.records || [];
                    const newTotal = data?.total || 0;
                    setInterviewList(newData);
                    setTotal(newTotal);

                    return {
                        success: code === 0,
                        data: newData,
                        total: newTotal,
                    };
                }}
                columns={columns}
            />
        </div>
    );
};

export default InterviewRoomPage;