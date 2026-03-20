import { WebServifyConfig, Message, Customer, ChatSession, Agent } from '@servify/core';

/**
 * 为原生 JavaScript 提供更简单的 API 接口
 */
export declare class VanillaServifySDK {
    private sdk;
    private eventCallbacks;
    constructor(config: WebServifyConfig);
    /**
     * 初始化 SDK
     */
    init(): Promise<void>;
    /**
     * 连接到服务器
     */
    connect(): Promise<void>;
    /**
     * 断开连接
     */
    disconnect(): void;
    /**
     * 开始聊天
     */
    startChat(options?: {
        priority?: string;
        message?: string;
    }): Promise<ChatSession>;
    /**
     * 发送消息
     */
    sendMessage(content: string, type?: string): Promise<Message>;
    /**
     * 结束会话
     */
    endChat(): Promise<void>;
    /**
     * AI 问答
     */
    askAI(question: string): Promise<{
        answer: string;
        confidence: number;
    }>;
    /**
     * 上传文件
     */
    uploadFile(file: File): Promise<{
        fileUrl: string;
        fileName: string;
        fileSize: number;
    }>;
    /**
     * 创建工单
     */
    createTicket(data: {
        title: string;
        description: string;
        priority?: string;
        category: string;
    }): Promise<any>;
    /**
     * 提交满意度评价
     */
    submitRating(rating: number, comment?: string): Promise<any>;
    /**
     * 获取历史消息
     */
    getMessages(page?: number, limit?: number): Promise<{
        messages: Message[];
        total: number;
        page: number;
    }>;
    /**
     * 获取客户信息
     */
    getCustomer(): Customer | null;
    /**
     * 获取当前会话
     */
    getSession(): ChatSession | null;
    /**
     * 获取当前客服代理
     */
    getAgent(): Agent | null;
    /**
     * 检查连接状态
     */
    isConnected(): boolean;
    /**
     * 添加事件监听器（简化版）
     */
    on(event: string, callback: (...args: unknown[]) => void): void;
    /**
     * 移除事件监听器
     */
    off(event: string, callback?: (...args: unknown[]) => void): void;
    /**
     * 触发回调函数
     */
    private triggerCallback;
}
declare global {
    interface Window {
        Servify: typeof VanillaServifySDK;
        createServify: (config: WebServifyConfig) => VanillaServifySDK;
    }
}
export type { WebServifyConfig as ServifyConfig, Message, Customer, ChatSession, Agent } from '@servify/core';
export default VanillaServifySDK;
