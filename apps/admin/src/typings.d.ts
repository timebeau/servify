declare namespace API {
  /** 当前用户 */
  interface CurrentUser {
    id?: number;
    username?: string;
    email?: string;
    role?: string;
    principal_kind?: 'admin' | 'agent' | 'service' | 'end_user';
    permissions?: string[];
  }

  /** 通用分页响应 */
  interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    page_size: number;
    pages: number;
  }

  /** 错误响应 */
  interface ErrorResponse {
    error: string;
    message: string;
    code?: string;
  }

  /** 成功响应 */
  interface SuccessResponse {
    message: string;
    data?: any;
  }

  // ---- 工单 ----
  interface Ticket {
    id: number;
    title: string;
    description?: string;
    status: string;
    priority: string;
    category?: string;
    customer_id: number;
    customer_name?: string;
    agent_id?: number;
    agent_name?: string;
    created_at: string;
    updated_at: string;
    resolved_at?: string;
    closed_at?: string;
    custom_fields?: Record<string, any>;
    tags?: string[];
    satisfaction_score?: number;
  }

  interface TicketListParams {
    page?: number;
    page_size?: number;
    status?: string;
    priority?: string;
    category?: string;
    agent_id?: number;
    customer_id?: number;
    search?: string;
  }

  // ---- 客户 ----
  interface Customer {
    id: number;
    name: string;
    email?: string;
    phone?: string;
    company?: string;
    industry?: string;
    source?: string;
    tags?: string[];
    priority?: string;
    created_at: string;
    updated_at: string;
  }

  // ---- 客服 ----
  interface Agent {
    id: number;
    name: string;
    email: string;
    status: string;
    skills?: string[];
    current_sessions?: number;
    max_sessions?: number;
    created_at: string;
  }

  // ---- 会话 ----
  interface Conversation {
    id: string;
    customer_id: number;
    customer_name?: string;
    agent_id?: number;
    agent_name?: string;
    status: string;
    channel?: string;
    created_at: string;
    updated_at: string;
  }

  // ---- 知识库文档 ----
  interface KnowledgeDoc {
    id: number;
    title: string;
    content?: string;
    category?: string;
    status: string;
    created_at: string;
    updated_at: string;
  }

  // ---- 自动化规则 ----
  interface Automation {
    id: number;
    name: string;
    description?: string;
    trigger_type: string;
    conditions?: Record<string, any>;
    actions?: Record<string, any>;
    enabled: boolean;
    created_at: string;
  }

  // ---- 宏 ----
  interface Macro {
    id: number;
    name: string;
    description?: string;
    content: string;
    category?: string;
    created_at: string;
  }

  // ---- SLA 配置 ----
  interface SLAConfig {
    id: number;
    name: string;
    description?: string;
    priority: string;
    first_response_time: number;
    resolution_time: number;
    enabled: boolean;
  }

  // ---- 班次 ----
  interface Shift {
    id: number;
    agent_id: number;
    agent_name?: string;
    start_time: string;
    end_time: string;
    status: string;
  }

  // ---- 满意度 ----
  interface Satisfaction {
    id: number;
    ticket_id: number;
    customer_id: number;
    score: number;
    comment?: string;
    created_at: string;
  }

  // ---- 自定义字段 ----
  interface CustomField {
    id: number;
    name: string;
    key: string;
    field_type: 'string' | 'number' | 'boolean' | 'date' | 'select' | 'multi_select';
    options?: string[];
    required: boolean;
    entity_type: string;
  }

  // ---- 审计日志 ----
  interface AuditLog {
    id: number;
    action: string;
    resource_type: string;
    resource_id: string;
    actor_user_id?: number;
    principal_kind?: string;
    details?: Record<string, any>;
    success: boolean;
    created_at: string;
  }

  // ---- 统计 ----
  interface DashboardStats {
    total_conversations: number;
    total_tickets: number;
    total_customers: number;
    avg_satisfaction: number;
    open_tickets: number;
    online_agents: number;
  }

  // ---- AI ----
  interface AIStatus {
    provider: string;
    model: string;
    available: boolean;
    latency_ms?: number;
  }

  // ---- 集成 ----
  interface Integration {
    id: number;
    name: string;
    type: string;
    config?: Record<string, any>;
    enabled: boolean;
    created_at: string;
  }

  // ---- 语音 ----
  interface VoiceProtocol {
    id: string;
    type: string;
    status: string;
    started_at: string;
    duration?: number;
  }
}
