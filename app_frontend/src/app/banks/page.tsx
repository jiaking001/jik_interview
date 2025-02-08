"use server";
import Title from "antd/es/typography/Title";
import "./index.css";
import {listQuestionBankVoByPageUsingPost} from "@/api/questionBankController";
import QuestionBankList from "@/components/QuestionBankList";

/**
 * 主页
 * @constructor
 */
export default async function HomePage() {

    let questionBankList = [];

    try {
        const questionBankRes = await listQuestionBankVoByPageUsingPost({
            pageSize: 200,
            sortField: 'createTime',
            sortOrder: 'descend',
        });
        questionBankList = questionBankRes.data.records ?? [];
    } catch (e) {
        console.error('获取题库列表失败，' + e.message);
    }

    return (
        <div id="homePage" className="max-width-content">
            <Title level={3}>题库大全</Title>
            <QuestionBankList questionBankList={questionBankList}/>
        </div>
    );
}
