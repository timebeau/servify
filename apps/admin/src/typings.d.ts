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
    source?: string;
    session_id?: string;
    customer_id: number;
    customer_name?: string;
    agent_id?: number;
    agent_name?: string;
    created_at: string;
    updated_at: string;
    resolved_at?: string;
    closed_at?: string;
    custom_fields?: Record<string, any>;
    tags?: string[] | string;
    tag_list?: string[];
    satisfaction_score?: number;
  }

  interface TicketListParams {
    page?: number;
    page_size?: number;
    status?: string;
    priority?: string;
    category?: string;
    source?: string;
    tag?: string;
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

  interface ConversationMessage {
    id: string;
    conversation_id: string;
    sender: string;
    kind: string;
    content: string;
    metadata?: Record<string, string>;
    created_at: string;
  }

  // ---- 知识库文档 ----
  interface KnowledgeDoc {
    id: number;
    title: string;
    content?: string;
    category?: string;
    status: string;
    tags?: string[];
    is_public?: boolean;
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
    total_sessions: number;
    total_tickets: number;
    total_customers: number;
    total_agents: number;
    today_tickets: number;
    today_sessions: number;
    today_messages: number;
    open_tickets: number;
    assigned_tickets: number;
    resolved_tickets: number;
    closed_tickets: number;
    online_agents: number;
    busy_agents: number;
    active_sessions: number;
    avg_response_time: number;
    avg_resolution_time: number;
    customer_satisfaction: number;
    ai_usage_today: number;
    weknora_usage_today: number;
  }

  interface RemoteAssistTicketStats {
    total: number;
    open: number;
    resolved: number;
    closed: number;
    resolved_rate: number;
    closed_rate: number;
    avg_close_hours: number;
  }

  interface WorkspaceSession {
    id: string;
    platform: string;
    status: string;
    agent_id?: number;
    agent_name?: string;
    customer_id?: number;
    customer_name?: string;
    started_at: string;
  }

  interface WorkspaceOverview {
    total_active_sessions: number;
    waiting_queue: number;
    online_agents: number;
    busy_agents: number;
    channels: Array<{
      platform: string;
      active_sessions: number;
      waiting_sessions: number;
      avg_response_time: number;
    }>;
    recent_sessions: WorkspaceSession[];
    agent_stats?: {
      available_agents: Array<{ id: number; name?: string }>;
    };
  }

  /** 会话（来自 ticket 关联查询） */
  interface Session {
    id: string;
    platform?: string;
    status?: string;
    customer_id?: number;
    customer_name?: string;
    agent_id?: number;
    agent_name?: string;
    ticket_id?: number;
    started_at?: string;
    ended_at?: string;
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

  // ---- 安全管理 ----
  interface UserSecurityPreview {
    user_id: number;
    username: string;
    name: string;
    role: string;
    status: string;
    token_version: number;
    next_token_version: number;
    last_login?: string;
    token_valid_after?: string;
  }

  interface UserSecurityDetail {
    user_id: number;
    role: string;
    status: string;
    token_version: number;
    token_valid_after?: string;
    last_login?: string;
  }

  interface UserSecuritySession {
    session_id: string;
    status: string;
    token_version: number;
    device_fingerprint?: string;
    network_label?: string;
    location_label?: string;
    risk_score?: number;
    risk_level?: string;
    risk_reasons?: string[];
    family_public_ip_count?: number;
    family_device_count?: number;
    active_session_count?: number;
    family_hot_refresh_count?: number;
    reference_session_id?: string;
    ip_drift?: boolean;
    device_drift?: boolean;
    rapid_ip_change?: boolean;
    rapid_device_change?: boolean;
    refresh_recency?: string;
    rapid_refresh_activity?: boolean;
    user_agent?: string;
    client_ip?: string;
    last_seen_at?: string;
    last_refreshed_at?: string;
    revoked_at?: string;
    created_at?: string;
    updated_at?: string;
    is_current?: boolean;
  }

  interface RevokedToken {
    jti: string;
    user_id: number;
    session_id?: string;
    token_use?: string;
    reason?: string;
    expires_at?: string;
    revoked_at?: string;
  }
}
