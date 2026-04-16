var __defProp = Object.defineProperty;
var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __publicField = (obj, key, value) => __defNormalProp(obj, typeof key !== "symbol" ? key + "" : key, value);
var __defProp2 = Object.defineProperty;
var __defNormalProp2 = (obj, key, value) => key in obj ? __defProp2(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __publicField2 = (obj, key, value) => __defNormalProp2(obj, typeof key !== "symbol" ? key + "" : key, value);
function getDefaultExportFromCjs(x) {
  return x && x.__esModule && Object.prototype.hasOwnProperty.call(x, "default") ? x["default"] : x;
}
var eventemitter3 = { exports: {} };
(function(module) {
  var has = Object.prototype.hasOwnProperty, prefix = "~";
  function Events() {
  }
  if (Object.create) {
    Events.prototype = /* @__PURE__ */ Object.create(null);
    if (!new Events().__proto__) prefix = false;
  }
  function EE(fn, context, once) {
    this.fn = fn;
    this.context = context;
    this.once = once || false;
  }
  function addListener(emitter, event, fn, context, once) {
    if (typeof fn !== "function") {
      throw new TypeError("The listener must be a function");
    }
    var listener = new EE(fn, context || emitter, once), evt = prefix ? prefix + event : event;
    if (!emitter._events[evt]) emitter._events[evt] = listener, emitter._eventsCount++;
    else if (!emitter._events[evt].fn) emitter._events[evt].push(listener);
    else emitter._events[evt] = [emitter._events[evt], listener];
    return emitter;
  }
  function clearEvent(emitter, evt) {
    if (--emitter._eventsCount === 0) emitter._events = new Events();
    else delete emitter._events[evt];
  }
  function EventEmitter2() {
    this._events = new Events();
    this._eventsCount = 0;
  }
  EventEmitter2.prototype.eventNames = function eventNames() {
    var names = [], events, name;
    if (this._eventsCount === 0) return names;
    for (name in events = this._events) {
      if (has.call(events, name)) names.push(prefix ? name.slice(1) : name);
    }
    if (Object.getOwnPropertySymbols) {
      return names.concat(Object.getOwnPropertySymbols(events));
    }
    return names;
  };
  EventEmitter2.prototype.listeners = function listeners(event) {
    var evt = prefix ? prefix + event : event, handlers = this._events[evt];
    if (!handlers) return [];
    if (handlers.fn) return [handlers.fn];
    for (var i = 0, l = handlers.length, ee = new Array(l); i < l; i++) {
      ee[i] = handlers[i].fn;
    }
    return ee;
  };
  EventEmitter2.prototype.listenerCount = function listenerCount(event) {
    var evt = prefix ? prefix + event : event, listeners = this._events[evt];
    if (!listeners) return 0;
    if (listeners.fn) return 1;
    return listeners.length;
  };
  EventEmitter2.prototype.emit = function emit(event, a1, a2, a3, a4, a5) {
    var evt = prefix ? prefix + event : event;
    if (!this._events[evt]) return false;
    var listeners = this._events[evt], len = arguments.length, args, i;
    if (listeners.fn) {
      if (listeners.once) this.removeListener(event, listeners.fn, void 0, true);
      switch (len) {
        case 1:
          return listeners.fn.call(listeners.context), true;
        case 2:
          return listeners.fn.call(listeners.context, a1), true;
        case 3:
          return listeners.fn.call(listeners.context, a1, a2), true;
        case 4:
          return listeners.fn.call(listeners.context, a1, a2, a3), true;
        case 5:
          return listeners.fn.call(listeners.context, a1, a2, a3, a4), true;
        case 6:
          return listeners.fn.call(listeners.context, a1, a2, a3, a4, a5), true;
      }
      for (i = 1, args = new Array(len - 1); i < len; i++) {
        args[i - 1] = arguments[i];
      }
      listeners.fn.apply(listeners.context, args);
    } else {
      var length = listeners.length, j;
      for (i = 0; i < length; i++) {
        if (listeners[i].once) this.removeListener(event, listeners[i].fn, void 0, true);
        switch (len) {
          case 1:
            listeners[i].fn.call(listeners[i].context);
            break;
          case 2:
            listeners[i].fn.call(listeners[i].context, a1);
            break;
          case 3:
            listeners[i].fn.call(listeners[i].context, a1, a2);
            break;
          case 4:
            listeners[i].fn.call(listeners[i].context, a1, a2, a3);
            break;
          default:
            if (!args) for (j = 1, args = new Array(len - 1); j < len; j++) {
              args[j - 1] = arguments[j];
            }
            listeners[i].fn.apply(listeners[i].context, args);
        }
      }
    }
    return true;
  };
  EventEmitter2.prototype.on = function on(event, fn, context) {
    return addListener(this, event, fn, context, false);
  };
  EventEmitter2.prototype.once = function once(event, fn, context) {
    return addListener(this, event, fn, context, true);
  };
  EventEmitter2.prototype.removeListener = function removeListener(event, fn, context, once) {
    var evt = prefix ? prefix + event : event;
    if (!this._events[evt]) return this;
    if (!fn) {
      clearEvent(this, evt);
      return this;
    }
    var listeners = this._events[evt];
    if (listeners.fn) {
      if (listeners.fn === fn && (!once || listeners.once) && (!context || listeners.context === context)) {
        clearEvent(this, evt);
      }
    } else {
      for (var i = 0, events = [], length = listeners.length; i < length; i++) {
        if (listeners[i].fn !== fn || once && !listeners[i].once || context && listeners[i].context !== context) {
          events.push(listeners[i]);
        }
      }
      if (events.length) this._events[evt] = events.length === 1 ? events[0] : events;
      else clearEvent(this, evt);
    }
    return this;
  };
  EventEmitter2.prototype.removeAllListeners = function removeAllListeners(event) {
    var evt;
    if (event) {
      evt = prefix ? prefix + event : event;
      if (this._events[evt]) clearEvent(this, evt);
    } else {
      this._events = new Events();
      this._eventsCount = 0;
    }
    return this;
  };
  EventEmitter2.prototype.off = EventEmitter2.prototype.removeListener;
  EventEmitter2.prototype.addListener = EventEmitter2.prototype.on;
  EventEmitter2.prefixed = prefix;
  EventEmitter2.EventEmitter = EventEmitter2;
  {
    module.exports = EventEmitter2;
  }
})(eventemitter3);
var eventemitter3Exports = eventemitter3.exports;
const EventEmitter = /* @__PURE__ */ getDefaultExportFromCjs(eventemitter3Exports);
function isStructuredApiResponse(value) {
  return typeof value === "object" && value !== null && "success" in value;
}
class ApiClient {
  constructor(options) {
    __publicField2(this, "options");
    this.options = {
      timeout: 1e4,
      headers: {
        "Content-Type": "application/json"
      },
      debug: false,
      ...options
    };
  }
  // 通用请求方法
  async request(method, endpoint, data, options) {
    const url = `${this.options.baseUrl}${endpoint}`;
    const headers = { ...this.options.headers, ...options == null ? void 0 : options.headers };
    this.log(`${method} ${url}`, data);
    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), this.options.timeout);
      const response = await fetch(url, {
        method,
        headers,
        body: data ? JSON.stringify(data) : void 0,
        signal: controller.signal
      });
      clearTimeout(timeoutId);
      const result = await response.json();
      const structuredResult = isStructuredApiResponse(result) ? result : void 0;
      this.log(`响应:`, result);
      if (!response.ok) {
        return {
          success: false,
          error: (structuredResult == null ? void 0 : structuredResult.error) || `HTTP ${response.status}: ${response.statusText}`
        };
      }
      return {
        success: true,
        data: (structuredResult == null ? void 0 : structuredResult.data) ?? result,
        message: structuredResult == null ? void 0 : structuredResult.message
      };
    } catch (error) {
      this.log(`请求失败:`, error);
      if (error instanceof Error) {
        if (error.name === "AbortError") {
          return { success: false, error: "Request timeout" };
        }
        return { success: false, error: error.message };
      }
      return { success: false, error: "Unknown error" };
    }
  }
  // 客户相关 API
  async getCustomer(customerId) {
    return this.request("GET", `/api/customers/${customerId}`);
  }
  async createCustomer(customerData) {
    return this.request("POST", "/api/customers", customerData);
  }
  async updateCustomer(customerId, customerData) {
    return this.request("PUT", `/api/customers/${customerId}`, customerData);
  }
  // 会话相关 API
  async createSession(sessionData) {
    return this.unsupported(
      "REST session creation is not exposed by the current server contract. Use the WebSocket chat flow instead.",
      sessionData
    );
  }
  async getSession(sessionId) {
    return this.request("GET", `/api/omni/sessions/${encodeURIComponent(String(sessionId))}`);
  }
  async endSession(sessionId) {
    return this.request("POST", `/api/omni/sessions/${encodeURIComponent(String(sessionId))}/close`);
  }
  async getCustomerSessions(customerId) {
    return this.request("GET", `/api/customers/${customerId}/sessions`);
  }
  // 消息相关 API
  async sendMessage(messageData) {
    const sessionId = messageData.session_id;
    if (sessionId === void 0 || sessionId === null || String(sessionId).trim() === "") {
      return { success: false, error: "session_id is required" };
    }
    return this.request(
      "POST",
      `/api/omni/sessions/${encodeURIComponent(String(sessionId))}/messages`,
      { content: messageData.content }
    );
  }
  async getSessionMessages(sessionId, options) {
    const params = new URLSearchParams();
    if (options == null ? void 0 : options.page) params.append("page", options.page.toString());
    if (options == null ? void 0 : options.limit) params.append("limit", options.limit.toString());
    const query = params.toString() ? `?${params.toString()}` : "";
    const response = await this.request(
      "GET",
      `/api/omni/sessions/${encodeURIComponent(String(sessionId))}/messages${query}`
    );
    if (!response.success) {
      return {
        success: false,
        error: response.error,
        message: response.message
      };
    }
    const messages = Array.isArray(response.data) ? response.data : [];
    return {
      success: true,
      data: {
        messages,
        total: messages.length,
        page: (options == null ? void 0 : options.page) ?? 1
      },
      message: response.message
    };
  }
  // AI 相关 API
  async askAI(question, sessionId) {
    var _a, _b, _c;
    const response = await this.request("POST", "/api/v1/ai/query", { query: question, session_id: sessionId ? String(sessionId) : "" });
    if (!response.success) {
      return response;
    }
    return {
      success: true,
      data: {
        answer: ((_a = response.data) == null ? void 0 : _a.answer) ?? ((_b = response.data) == null ? void 0 : _b.content) ?? "",
        confidence: ((_c = response.data) == null ? void 0 : _c.confidence) ?? 0
      },
      message: response.message
    };
  }
  async getAIStatus() {
    return this.request("GET", "/api/v1/ai/status");
  }
  // 工单相关 API
  async createTicket(ticketData) {
    return this.request("POST", "/api/tickets", ticketData);
  }
  async getTicket(ticketId) {
    return this.request("GET", `/api/tickets/${ticketId}`);
  }
  async updateTicket(ticketId, updates) {
    return this.request("PUT", `/api/tickets/${ticketId}`, updates);
  }
  async getCustomerTickets(customerId) {
    return this.request("GET", `/api/tickets?customer_id=${encodeURIComponent(String(customerId))}`);
  }
  // 满意度评价 API
  async submitSatisfaction(satisfactionData) {
    return this.request("POST", "/api/satisfactions", satisfactionData);
  }
  async getSatisfactionByTicket(ticketId) {
    return this.request("GET", `/api/tickets/${ticketId}/satisfaction`);
  }
  // 队列相关 API
  async joinQueue(queueData) {
    return this.unsupported("Queue join is not exposed by the current server contract.", queueData);
  }
  async getQueueStatus(customerId) {
    return this.unsupported("Queue status is not exposed by the current server contract.", customerId);
  }
  async leaveQueue(customerId) {
    return this.unsupported("Queue leave is not exposed by the current server contract.", customerId);
  }
  // 文件上传 API
  async uploadFile(file, sessionId) {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("session_id", sessionId.toString());
    const url = `${this.options.baseUrl}/api/v1/upload`;
    try {
      const response = await fetch(url, {
        method: "POST",
        body: formData,
        headers: {
          // 不设置 Content-Type，让浏览器自动设置 multipart/form-data 边界
          ...Object.fromEntries(
            Object.entries(this.options.headers).filter(
              ([key]) => key.toLowerCase() !== "content-type"
            )
          )
        }
      });
      const result = await response.json();
      const structuredResult = isStructuredApiResponse(result) ? result : void 0;
      if (!response.ok) {
        return { success: false, error: (structuredResult == null ? void 0 : structuredResult.error) || "Upload failed" };
      }
      return {
        success: true,
        data: (structuredResult == null ? void 0 : structuredResult.data) ?? result
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : "Upload failed"
      };
    }
  }
  // WebRTC 相关 API
  async startCall(sessionId, callType) {
    return this.unsupported(
      "WebRTC call REST endpoints are not exposed by the current server contract. Use WebSocket signaling instead.",
      { session_id: sessionId, type: callType }
    );
  }
  async endCall(callId) {
    return this.unsupported("WebRTC call REST endpoints are not exposed by the current server contract.", callId);
  }
  async getCallStatus(callId) {
    return this.unsupported("WebRTC call REST endpoints are not exposed by the current server contract.", callId);
  }
  // 设置认证头
  setAuthToken(token) {
    this.options.headers["Authorization"] = `Bearer ${token}`;
  }
  // 设置客户 ID
  setCustomerId(customerId) {
    this.options.headers["X-Customer-ID"] = customerId.toString();
  }
  log(...args) {
    if (this.options.debug) {
      console.warn("[ServifyAPI]", ...args);
    }
  }
  unsupported(error, context) {
    this.log(error, context);
    return Promise.resolve({ success: false, error });
  }
}
class ServifyError extends Error {
  constructor(message, options) {
    super(message);
    __publicField2(this, "code");
    __publicField2(this, "cause");
    __publicField2(this, "details");
    __publicField2(this, "retryable");
    this.name = "ServifyError";
    this.code = options.code;
    this.cause = options.cause;
    this.details = options.details;
    this.retryable = options.retryable ?? false;
  }
}
const DEFAULT_RECONNECT_POLICY = {
  maxAttempts: 5,
  baseDelayMs: 1e3,
  backoffFactor: 2,
  maxDelayMs: 3e4
};
function normalizeReconnectPolicy(policy, legacy) {
  return {
    maxAttempts: (policy == null ? void 0 : policy.maxAttempts) ?? (legacy == null ? void 0 : legacy.reconnectAttempts) ?? DEFAULT_RECONNECT_POLICY.maxAttempts,
    baseDelayMs: (policy == null ? void 0 : policy.baseDelayMs) ?? (legacy == null ? void 0 : legacy.reconnectDelay) ?? DEFAULT_RECONNECT_POLICY.baseDelayMs,
    backoffFactor: (policy == null ? void 0 : policy.backoffFactor) ?? DEFAULT_RECONNECT_POLICY.backoffFactor,
    maxDelayMs: (policy == null ? void 0 : policy.maxDelayMs) ?? DEFAULT_RECONNECT_POLICY.maxDelayMs
  };
}
function computeReconnectDelay(policy, attempt) {
  const exponent = Math.max(attempt - 1, 0);
  const delay = policy.baseDelayMs * Math.pow(policy.backoffFactor ?? 2, exponent);
  return Math.min(delay, policy.maxDelayMs ?? delay);
}
function shouldReconnect(decision, policy) {
  if (decision.isManualClose) {
    return false;
  }
  return decision.attempt < policy.maxAttempts;
}
class WebSocketManager extends EventEmitter {
  constructor(options) {
    super();
    __publicField2(this, "ws", null);
    __publicField2(this, "options");
    __publicField2(this, "reconnectAttempts", 0);
    __publicField2(this, "reconnectTimer", null);
    __publicField2(this, "heartbeatTimer", null);
    __publicField2(this, "isManualClose", false);
    __publicField2(this, "subscribers", /* @__PURE__ */ new Set());
    __publicField2(this, "kind", "websocket");
    __publicField2(this, "state", "idle");
    const reconnectPolicy = normalizeReconnectPolicy(options.reconnectPolicy, {
      reconnectAttempts: options.reconnectAttempts,
      reconnectDelay: options.reconnectDelay
    });
    this.options = {
      protocols: [],
      reconnectAttempts: 5,
      reconnectDelay: 1e3,
      heartbeatInterval: 3e4,
      debug: false,
      reconnectPolicy,
      authProvider: options.authProvider,
      onTokenRefreshRequired: options.onTokenRefreshRequired ?? (async () => void 0),
      ...options
    };
  }
  async connect(_options) {
    const connectionUrl = await this.resolveConnectionUrl();
    return new Promise((resolve, reject) => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.state = "connected";
        resolve();
        return;
      }
      this.isManualClose = false;
      this.state = "connecting";
      this.log("正在连接 WebSocket...", connectionUrl);
      try {
        this.ws = new WebSocket(connectionUrl, this.options.protocols);
      } catch (error) {
        this.state = "error";
        reject(error);
        return;
      }
      this.ws.onopen = () => {
        this.log("WebSocket 连接成功");
        this.reconnectAttempts = 0;
        this.state = "connected";
        this.startHeartbeat();
        this.emit("connected");
        resolve();
      };
      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          this.handleMessage(message);
        } catch (error) {
          this.log("解析消息失败:", error);
          this.emit("error", new Error("Invalid message format"));
        }
      };
      this.ws.onclose = (event) => {
        this.log("WebSocket 连接关闭:", event.code, event.reason);
        this.stopHeartbeat();
        this.state = this.isManualClose ? "closed" : "idle";
        this.emit("disconnected", event.reason || "连接关闭");
        if (shouldReconnect(
          {
            attempt: this.reconnectAttempts,
            isManualClose: this.isManualClose
          },
          this.options.reconnectPolicy
        )) {
          this.state = "reconnecting";
          this.scheduleReconnect();
        }
      };
      this.ws.onerror = (event) => {
        this.log("WebSocket 错误:", event);
        this.state = "error";
        const err = new ServifyError("WebSocket connection error", {
          code: "transport_unavailable",
          retryable: true,
          details: { url: this.options.url }
        });
        this.emit("error", err);
        reject(err);
      };
    });
  }
  async disconnect() {
    this.isManualClose = true;
    this.state = "closed";
    this.stopHeartbeat();
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
  async send(message, _options) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      this.log("WebSocket 未连接，无法发送消息");
      const err = new ServifyError("WebSocket not connected", {
        code: "transport_disconnected",
        retryable: true
      });
      this.emit("error", err);
      throw err;
    }
    try {
      this.ws.send(JSON.stringify(message));
      this.log("发送消息:", message);
    } catch (error) {
      this.log("发送消息失败:", error);
      const err = new ServifyError("Failed to send message", {
        code: "transport_unavailable",
        cause: error,
        retryable: true
      });
      this.emit("error", err);
      throw err;
    }
  }
  isConnected() {
    var _a;
    return ((_a = this.ws) == null ? void 0 : _a.readyState) === WebSocket.OPEN;
  }
  subscribe(handler) {
    this.subscribers.add(handler);
    return () => {
      this.subscribers.delete(handler);
    };
  }
  handleMessage(message) {
    this.log("收到消息:", message);
    for (const subscriber of this.subscribers) {
      subscriber(message);
    }
    switch (message.type) {
      case "message":
      case "text-message":
        this.emit("message", this.normalizeMessage(message, "customer"));
        break;
      case "agent-message":
        this.emit("message", this.normalizeMessage(message, "agent"));
        break;
      case "ai-response":
        this.emit("message", this.normalizeMessage(message, "system", true));
        break;
      case "session_update":
        this.emit("session_updated", message.data);
        break;
      case "agent_status":
        if (typeof message.data === "object" && message.data !== null) {
          const agentStatus = message.data;
          if (agentStatus.type === "assigned" && agentStatus.agent) {
            this.emit("agent_assigned", agentStatus.agent);
          } else if (agentStatus.type === "typing" && typeof agentStatus.typing === "boolean") {
            this.emit("agent_typing", agentStatus.typing);
          }
        }
        break;
      case "error":
        this.emit("error", new Error(this.extractMessageText(message.data) || "Unknown error"));
        break;
      case "system":
        if (typeof message.data === "object" && message.data !== null && "type" in message.data && message.data.type === "pong") {
          this.log("收到心跳响应");
        }
        break;
      case "webrtc-offer":
        this.emit("webrtc:offer", message.data);
        break;
      case "webrtc-answer":
        this.emit("webrtc:answer", message.data);
        break;
      case "webrtc-candidate":
        this.emit("webrtc:candidate", this.extractICECandidate(message.data));
        break;
      case "webrtc-state-change":
        this.emit("webrtc:state", this.extractRemoteAssistState(message.data));
        break;
      default:
        this.log("未知消息类型:", message.type);
    }
  }
  extractRemoteAssistState(data) {
    const runtimeState = this.extractRemoteAssistRuntimeState(data);
    switch (runtimeState.state) {
      case "new":
      case "checking":
      case "connecting":
        return "connecting";
      case "connected":
        return "connected";
      case "failed":
        return "failed";
      case "closed":
      case "disconnected":
        return "ended";
      default:
        return "connecting";
    }
  }
  extractRemoteAssistRuntimeState(data) {
    if (typeof data === "object" && data !== null) {
      return {
        connectionId: "connection_id" in data && typeof data.connection_id === "string" ? data.connection_id : void 0,
        state: "state" in data && typeof data.state === "string" ? data.state : "connecting"
      };
    }
    return { state: "connecting" };
  }
  scheduleReconnect() {
    this.reconnectAttempts++;
    this.emit("reconnecting", this.reconnectAttempts);
    const delay = computeReconnectDelay(this.options.reconnectPolicy, this.reconnectAttempts);
    this.log(`${delay}ms 后重连 (第 ${this.reconnectAttempts}/${this.options.reconnectPolicy.maxAttempts} 次)`);
    this.reconnectTimer = setTimeout(() => {
      this.connect().catch(() => {
      });
    }, delay);
  }
  startHeartbeat() {
    this.stopHeartbeat();
    this.heartbeatTimer = setInterval(() => {
      if (this.isConnected()) {
        void this.send({
          type: "system",
          data: { type: "ping", timestamp: (/* @__PURE__ */ new Date()).toISOString() }
        }).catch(() => void 0);
      }
    }, this.options.heartbeatInterval);
  }
  stopHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }
  log(...args) {
    if (this.options.debug) {
      console.warn("[ServifyWS]", ...args);
    }
  }
  async resolveConnectionUrl() {
    const token = await this.resolveAuthToken(this.options.authProvider);
    if (!token) {
      return this.options.url;
    }
    const url = new URL(this.options.url);
    url.searchParams.set("access_token", token);
    return url.toString();
  }
  async resolveAuthToken(authProvider) {
    if (!authProvider) {
      return null;
    }
    const currentToken = await authProvider.getToken();
    if (currentToken == null ? void 0 : currentToken.accessToken) {
      return currentToken.accessToken;
    }
    if (!authProvider.refreshToken) {
      return null;
    }
    await this.options.onTokenRefreshRequired();
    const refreshedToken = await authProvider.refreshToken();
    if (refreshedToken == null ? void 0 : refreshedToken.accessToken) {
      return refreshedToken.accessToken;
    }
    throw new ServifyError("Authentication refresh required", {
      code: "auth_refresh_required",
      retryable: false,
      details: { url: this.options.url }
    });
  }
  extractMessageText(data) {
    if (typeof data === "object" && data !== null && "message" in data && typeof data.message === "string") {
      return data.message;
    }
    return null;
  }
  normalizeMessage(message, senderType, isAIResponse = false) {
    const data = typeof message.data === "object" && message.data !== null ? message.data : { content: String(message.data ?? "") };
    const content = typeof data.content === "string" ? data.content : typeof data.message === "string" ? data.message : String(message.data ?? "");
    return {
      id: typeof data.id === "string" || typeof data.id === "number" ? data.id : `ws-${Date.now()}`,
      session_id: typeof data.session_id === "string" || typeof data.session_id === "number" ? data.session_id : message.session_id ?? "",
      sender_type: senderType,
      content,
      message_type: "text",
      is_ai_response: isAIResponse,
      metadata: data,
      created_at: (/* @__PURE__ */ new Date()).toISOString()
    };
  }
  extractICECandidate(data) {
    if (typeof data === "object" && data !== null && "candidate" in data && typeof data.candidate === "object" && data.candidate !== null) {
      return data.candidate;
    }
    return data;
  }
}
function negotiateCapabilities(available, requested) {
  const granted = [];
  const rejected = [];
  for (const request of requested) {
    const descriptor = available.find((entry) => entry.name === request.name);
    if (!descriptor) {
      rejected.push({ request, reason: "unsupported" });
      continue;
    }
    if (!descriptor.enabled) {
      rejected.push({ request, reason: "disabled", descriptor });
      continue;
    }
    granted.push({ ...descriptor });
  }
  return { granted, rejected };
}
class StaticCapabilitySet {
  constructor(entries) {
    __publicField2(this, "entries");
    this.entries = entries.slice();
  }
  all() {
    return this.entries.slice();
  }
  has(name) {
    return this.entries.some((entry) => entry.name === name && entry.enabled);
  }
  get(name) {
    return this.entries.find((entry) => entry.name === name);
  }
  negotiate(requested) {
    return negotiateCapabilities(this.entries, requested);
  }
}
function createWebCapabilitySet() {
  return new StaticCapabilitySet([
    { name: "chat", enabled: true, version: "1" },
    { name: "realtime", enabled: true, version: "1" },
    { name: "knowledge", enabled: true, version: "1" },
    { name: "remote_assist", enabled: true, version: "1" },
    { name: "voice", enabled: false, version: "reserved" }
  ]);
}
class ServifySDK extends EventEmitter {
  constructor(config) {
    super();
    __publicField2(this, "config");
    __publicField2(this, "api");
    __publicField2(this, "ws", null);
    __publicField2(this, "currentCustomer", null);
    __publicField2(this, "currentSession", null);
    __publicField2(this, "currentAgent", null);
    __publicField2(this, "messageQueue", []);
    __publicField2(this, "remoteAssistPeer", null);
    __publicField2(this, "remoteAssistStream", null);
    __publicField2(this, "isInitialized", false);
    __publicField2(this, "id");
    __publicField2(this, "capabilities");
    __publicField2(this, "events", this);
    __publicField2(this, "authProvider");
    __publicField2(this, "transport", {
      get kind() {
        return "session";
      },
      get state() {
        return "idle";
      },
      connect: async () => void 0,
      disconnect: async () => void 0,
      send: async () => void 0,
      isConnected: () => false,
      subscribe: () => () => void 0
    });
    this.id = config.sessionId || `web-${Date.now()}`;
    this.capabilities = createWebCapabilitySet();
    this.config = {
      autoConnect: true,
      reconnectAttempts: 5,
      reconnectDelay: 1e3,
      debug: false,
      ...config
    };
    this.api = new ApiClient({
      baseUrl: this.config.apiUrl,
      debug: this.config.debug
    });
    if (this.config.customerId) {
      this.api.setCustomerId(parseInt(this.config.customerId));
    }
    this.log("SDK 初始化完成", this.config);
  }
  // 初始化 SDK
  async initialize() {
    if (this.isInitialized) {
      this.log("SDK 已初始化");
      return;
    }
    try {
      this.log("正在初始化 SDK...");
      await this.initializeCustomer();
      if (this.config.autoConnect) {
        await this.connect();
      }
      this.isInitialized = true;
      this.log("SDK 初始化成功");
    } catch (error) {
      this.log("SDK 初始化失败:", error);
      throw error;
    }
  }
  // 连接到服务器
  async connect() {
    if (!this.currentCustomer) {
      throw new Error("Customer not initialized. Call initialize() first.");
    }
    const wsUrl = this.config.wsUrl || this.config.apiUrl.replace(/^http/, "ws") + "/api/v1/ws";
    const realtimeSessionID = this.resolveRealtimeSessionID();
    this.ws = new WebSocketManager({
      url: `${wsUrl}?session_id=${encodeURIComponent(realtimeSessionID)}`,
      reconnectAttempts: this.config.reconnectAttempts,
      reconnectDelay: this.config.reconnectDelay,
      reconnectPolicy: this.config.reconnectPolicy,
      authProvider: this.config.authProvider,
      onTokenRefreshRequired: this.config.onTokenRefreshRequired,
      debug: this.config.debug
    });
    this.ws.on("connected", () => this.emit("connected"));
    this.ws.on("disconnected", (reason) => this.emit("disconnected", reason));
    this.ws.on("reconnecting", (attempt) => this.emit("reconnecting", attempt));
    this.ws.on("message", (message) => this.handleIncomingMessage(message));
    this.ws.on("session_updated", (session) => {
      this.currentSession = session;
      this.emit("session_updated", session);
    });
    this.ws.on("agent_assigned", (agent) => {
      this.currentAgent = agent;
      this.emit("agent_assigned", agent);
    });
    this.ws.on("agent_typing", (isTyping) => this.emit("agent_typing", isTyping));
    this.ws.on("webrtc:offer", (offer) => this.emit("webrtc:offer", offer));
    this.ws.on("webrtc:answer", (answer) => this.emit("webrtc:answer", answer));
    this.ws.on("webrtc:candidate", (candidate) => this.emit("webrtc:candidate", candidate));
    this.ws.on("webrtc:state", (state) => this.updateRemoteAssistState(state));
    this.ws.on("error", (error) => this.emit("error", error));
    await this.ws.connect();
  }
  // 断开连接
  disconnect() {
    var _a;
    void this.endRemoteAssist();
    (_a = this.ws) == null ? void 0 : _a.disconnect();
    this.ws = null;
  }
  // 开始聊天会话
  async startChat(options) {
    var _a, _b;
    if (!this.currentCustomer) {
      throw new Error("Customer not initialized");
    }
    const sessionData = {
      id: ((_a = this.currentSession) == null ? void 0 : _a.id) ?? this.id,
      customer_id: this.currentCustomer.id,
      status: "active",
      channel: "web",
      priority: (options == null ? void 0 : options.priority) || "normal",
      started_at: (/* @__PURE__ */ new Date()).toISOString(),
      created_at: (/* @__PURE__ */ new Date()).toISOString(),
      updated_at: (/* @__PURE__ */ new Date()).toISOString()
    };
    this.currentSession = sessionData;
    if (!((_b = this.ws) == null ? void 0 : _b.isConnected())) {
      await this.connect();
    }
    this.emit("session_created", this.currentSession);
    if (options == null ? void 0 : options.message) {
      await this.sendMessage(options.message);
    }
    return this.currentSession;
  }
  // 发送消息
  async sendMessage(content, options) {
    var _a, _b, _c;
    if (!this.currentSession) {
      throw new Error("No active session. Start a chat first.");
    }
    const messageData = {
      id: `local-${Date.now()}`,
      session_id: this.currentSession.id,
      sender_type: "customer",
      sender_id: (_a = this.currentCustomer) == null ? void 0 : _a.id,
      content,
      message_type: (options == null ? void 0 : options.type) || "text",
      attachments: options == null ? void 0 : options.attachments,
      metadata: options == null ? void 0 : options.metadata,
      created_at: (/* @__PURE__ */ new Date()).toISOString()
    };
    if (!((_b = this.ws) == null ? void 0 : _b.isConnected())) {
      await this.connect();
    }
    await ((_c = this.ws) == null ? void 0 : _c.send({
      type: "text-message",
      data: {
        content,
        message_type: (options == null ? void 0 : options.type) || "text",
        attachments: options == null ? void 0 : options.attachments,
        metadata: options == null ? void 0 : options.metadata
      }
    }));
    return messageData;
  }
  // 结束会话
  async endSession() {
    if (!this.currentSession) {
      return;
    }
    await this.endRemoteAssist();
    const endedSession = { ...this.currentSession, status: "closed", ended_at: (/* @__PURE__ */ new Date()).toISOString() };
    this.currentSession = null;
    this.currentAgent = null;
    this.emit("session_ended", endedSession);
  }
  // 创建工单
  async createTicket(ticketData) {
    if (!this.currentCustomer) {
      throw new Error("Customer not initialized");
    }
    const data = {
      ...ticketData,
      customer_id: this.currentCustomer.id,
      status: "open"
    };
    const response = await this.api.createTicket(data);
    if (!response.success || !response.data) {
      throw new Error(response.error || "Failed to create ticket");
    }
    this.emit("ticket_created", response.data);
    return response.data;
  }
  async startRemoteAssist(options) {
    var _a, _b, _c, _d;
    if (!((_a = this.ws) == null ? void 0 : _a.isConnected())) {
      throw new Error("WebSocket not connected. Call connect() first.");
    }
    this.updateRemoteAssistState("starting");
    await this.endRemoteAssist();
    const peer = new RTCPeerConnection({
      iceServers: (options == null ? void 0 : options.iceServers) || ((_b = this.config.remoteAssist) == null ? void 0 : _b.iceServers) || []
    });
    this.remoteAssistPeer = peer;
    peer.createDataChannel(
      (options == null ? void 0 : options.dataChannelLabel) || ((_c = this.config.remoteAssist) == null ? void 0 : _c.dataChannelLabel) || "servify-remote-assist"
    );
    peer.onicecandidate = (event) => {
      var _a2;
      if (!event.candidate || !((_a2 = this.ws) == null ? void 0 : _a2.isConnected())) {
        return;
      }
      void this.ws.send({
        type: "webrtc-candidate",
        data: event.candidate.toJSON()
      });
    };
    peer.onconnectionstatechange = () => {
      const state = peer.connectionState;
      if (state === "connected") {
        this.updateRemoteAssistState("connected");
      } else if (state === "connecting") {
        this.updateRemoteAssistState("connecting");
      } else if (state === "failed") {
        this.updateRemoteAssistState("failed");
      } else if (state === "closed" || state === "disconnected") {
        this.updateRemoteAssistState("ended");
      }
    };
    peer.ontrack = (event) => {
      this.emit("webrtc:track", event);
    };
    const shouldCaptureScreen = (options == null ? void 0 : options.captureScreen) ?? ((_d = this.config.remoteAssist) == null ? void 0 : _d.captureScreen) ?? false;
    if (shouldCaptureScreen) {
      this.remoteAssistStream = await this.captureRemoteAssistStream(options);
      for (const track of this.remoteAssistStream.getTracks()) {
        peer.addTrack(track, this.remoteAssistStream);
      }
    }
    const offer = await peer.createOffer();
    await peer.setLocalDescription(offer);
    const localDescription = peer.localDescription;
    if (!localDescription) {
      throw new Error("Failed to create local WebRTC description");
    }
    await this.ws.send({
      type: "webrtc-offer",
      data: localDescription.toJSON()
    });
    this.emit("webrtc:offer", localDescription.toJSON());
    this.updateRemoteAssistState("offered");
    return peer;
  }
  async acceptRemoteAnswer(answer) {
    if (!this.remoteAssistPeer) {
      throw new Error("Remote assist has not started");
    }
    await this.remoteAssistPeer.setRemoteDescription(answer);
    this.updateRemoteAssistState("connecting");
  }
  async addRemoteIce(candidate) {
    if (!this.remoteAssistPeer) {
      throw new Error("Remote assist has not started");
    }
    await this.remoteAssistPeer.addIceCandidate(candidate);
  }
  async endRemoteAssist() {
    if (this.remoteAssistStream) {
      for (const track of this.remoteAssistStream.getTracks()) {
        track.stop();
      }
      this.remoteAssistStream = null;
    }
    if (this.remoteAssistPeer) {
      this.remoteAssistPeer.onicecandidate = null;
      this.remoteAssistPeer.onconnectionstatechange = null;
      this.remoteAssistPeer.ontrack = null;
      this.remoteAssistPeer.close();
      this.remoteAssistPeer = null;
    }
    this.updateRemoteAssistState("ended");
  }
  // 提交满意度评价
  async submitSatisfaction(satisfaction) {
    var _a;
    if (!this.currentCustomer) {
      throw new Error("Customer not initialized");
    }
    const data = {
      ...satisfaction,
      customer_id: this.currentCustomer.id,
      agent_id: (_a = this.currentAgent) == null ? void 0 : _a.id
    };
    const response = await this.api.submitSatisfaction(data);
    if (!response.success || !response.data) {
      throw new Error(response.error || "Failed to submit satisfaction");
    }
    return response.data;
  }
  // AI 问答
  async askAI(question) {
    var _a;
    const response = await this.api.askAI(question, (_a = this.currentSession) == null ? void 0 : _a.id);
    if (!response.success || !response.data) {
      throw new Error(response.error || "Failed to get AI response");
    }
    return response.data;
  }
  // 文件上传
  async uploadFile(file) {
    if (!this.currentSession) {
      throw new Error("No active session");
    }
    const response = await this.api.uploadFile(file, this.currentSession.id);
    if (!response.success || !response.data) {
      throw new Error(response.error || "Failed to upload file");
    }
    return response.data;
  }
  // 获取历史消息
  async getMessages(options) {
    if (!this.currentSession) {
      throw new Error("No active session");
    }
    const response = await this.api.getSessionMessages(this.currentSession.id, options);
    if (!response.success || !response.data) {
      throw new Error(response.error || "Failed to get messages");
    }
    return response.data;
  }
  // 获取客户信息
  getCustomer() {
    return this.currentCustomer;
  }
  // 获取当前会话
  getSession() {
    return this.currentSession;
  }
  // 获取当前客服代理
  getAgent() {
    return this.currentAgent;
  }
  getIdentity() {
    var _a, _b, _c, _d, _e, _f;
    return {
      customerId: (_b = (_a = this.currentCustomer) == null ? void 0 : _a.id) == null ? void 0 : _b.toString(),
      conversationId: (_d = (_c = this.currentSession) == null ? void 0 : _c.id) == null ? void 0 : _d.toString(),
      agentId: (_f = (_e = this.currentAgent) == null ? void 0 : _e.id) == null ? void 0 : _f.toString()
    };
  }
  getState() {
    return {
      initialized: this.isInitialized,
      connected: this.isConnected(),
      customer: this.currentCustomer,
      session: this.currentSession,
      agent: this.currentAgent
    };
  }
  updateState(patch) {
    return { ...this.getState(), ...patch };
  }
  // 检查连接状态
  isConnected() {
    var _a;
    return ((_a = this.ws) == null ? void 0 : _a.isConnected()) ?? false;
  }
  // 私有方法：初始化客户信息
  async initializeCustomer() {
    const customerId = Number.parseInt(this.config.customerId || "0", 10);
    this.currentCustomer = {
      id: Number.isFinite(customerId) && customerId > 0 ? customerId : 0,
      name: this.config.customerName || "Anonymous",
      email: this.config.customerEmail || "",
      status: "active",
      created_at: (/* @__PURE__ */ new Date()).toISOString(),
      updated_at: (/* @__PURE__ */ new Date()).toISOString()
    };
    if (this.currentCustomer.id > 0) {
      this.api.setCustomerId(this.currentCustomer.id);
    }
  }
  // 私有方法：处理收到的消息
  handleIncomingMessage(message) {
    this.messageQueue.push(message);
    this.emit("message", message);
  }
  resolveRealtimeSessionID() {
    var _a;
    const sessionID = (_a = this.currentSession) == null ? void 0 : _a.id;
    if (typeof sessionID === "string" || typeof sessionID === "number") {
      return String(sessionID);
    }
    return this.id;
  }
  async captureRemoteAssistStream(options) {
    var _a;
    if (typeof navigator === "undefined" || !navigator.mediaDevices || typeof navigator.mediaDevices.getDisplayMedia !== "function") {
      throw new Error("Screen capture is not supported in this browser");
    }
    return navigator.mediaDevices.getDisplayMedia({
      video: true,
      audio: (options == null ? void 0 : options.audio) ?? ((_a = this.config.remoteAssist) == null ? void 0 : _a.audio) ?? false
    });
  }
  updateRemoteAssistState(state) {
    this.emit("webrtc:state", state);
  }
  // 私有方法：日志输出
  log(...args) {
    if (this.config.debug) {
      console.warn("[ServifySDK]", ...args);
    }
  }
}
function createWebServifySDK(config) {
  return new ServifySDK(config);
}
class VanillaServifySDK {
  constructor(config) {
    __publicField(this, "sdk");
    __publicField(this, "eventCallbacks", /* @__PURE__ */ new Map());
    this.sdk = createWebServifySDK(config);
    this.sdk.on("connected", () => this.triggerCallback("connected"));
    this.sdk.on("disconnected", (reason) => this.triggerCallback("disconnected", reason));
    this.sdk.on("message", (message) => this.triggerCallback("message", message));
    this.sdk.on("session_created", (session) => this.triggerCallback("sessionCreated", session));
    this.sdk.on("session_updated", (session) => this.triggerCallback("sessionUpdated", session));
    this.sdk.on("session_ended", (session) => this.triggerCallback("sessionEnded", session));
    this.sdk.on("agent_assigned", (agent) => this.triggerCallback("agentAssigned", agent));
    this.sdk.on("agent_typing", (isTyping) => this.triggerCallback("agentTyping", isTyping));
    this.sdk.on("error", (error) => this.triggerCallback("error", error));
    this.sdk.on("ticket_created", (ticket) => this.triggerCallback("ticketCreated", ticket));
    this.sdk.on("webrtc:offer", (offer) => this.triggerCallback("webrtc:offer", offer));
    this.sdk.on("webrtc:answer", (answer) => this.triggerCallback("webrtc:answer", answer));
    this.sdk.on("webrtc:candidate", (candidate) => this.triggerCallback("webrtc:candidate", candidate));
    this.sdk.on("webrtc:track", (event) => this.triggerCallback("webrtc:track", event));
    this.sdk.on("webrtc:state", (state) => this.triggerCallback("webrtc:state", state));
  }
  normalizePriority(priority) {
    if (priority === "low" || priority === "normal" || priority === "high" || priority === "urgent") {
      return priority;
    }
    return void 0;
  }
  normalizeMessageType(type) {
    if (type === "image" || type === "file") {
      return type;
    }
    return "text";
  }
  /**
   * 初始化 SDK
   */
  async init() {
    return this.sdk.initialize();
  }
  /**
   * 连接到服务器
   */
  async connect() {
    return this.sdk.connect();
  }
  /**
   * 断开连接
   */
  disconnect() {
    this.sdk.disconnect();
  }
  /**
   * 开始聊天
   */
  async startChat(options) {
    return this.sdk.startChat({
      priority: this.normalizePriority(options == null ? void 0 : options.priority),
      message: options == null ? void 0 : options.message
    });
  }
  /**
   * 发送消息
   */
  async sendMessage(content, type = "text") {
    return this.sdk.sendMessage(content, { type: this.normalizeMessageType(type) });
  }
  /**
   * 结束会话
   */
  async endChat() {
    return this.sdk.endSession();
  }
  /**
   * 发起远程协助基础链路
   */
  async startRemoteAssist(options) {
    return this.sdk.startRemoteAssist(options);
  }
  /**
   * 接收对端 WebRTC answer
   */
  async acceptRemoteAnswer(answer) {
    return this.sdk.acceptRemoteAnswer(answer);
  }
  /**
   * 注入对端 ICE candidate
   */
  async addRemoteIce(candidate) {
    return this.sdk.addRemoteIce(candidate);
  }
  /**
   * 结束远程协助
   */
  async endRemoteAssist() {
    return this.sdk.endRemoteAssist();
  }
  /**
   * AI 问答
   */
  async askAI(question) {
    return this.sdk.askAI(question);
  }
  /**
   * 上传文件
   */
  async uploadFile(file) {
    const result = await this.sdk.uploadFile(file);
    return {
      fileUrl: result.file_url,
      fileName: result.file_name,
      fileSize: result.file_size
    };
  }
  /**
   * 创建工单
   */
  async createTicket(data) {
    return this.sdk.createTicket({
      ...data,
      priority: this.normalizePriority(data.priority)
    });
  }
  /**
   * 提交满意度评价
   */
  async submitRating(rating, comment) {
    return this.sdk.submitSatisfaction({
      rating,
      comment
    });
  }
  /**
   * 获取历史消息
   */
  async getMessages(page = 1, limit = 50) {
    return this.sdk.getMessages({ page, limit });
  }
  /**
   * 获取客户信息
   */
  getCustomer() {
    return this.sdk.getCustomer();
  }
  /**
   * 获取当前会话
   */
  getSession() {
    return this.sdk.getSession();
  }
  /**
   * 获取当前客服代理
   */
  getAgent() {
    return this.sdk.getAgent();
  }
  /**
   * 检查连接状态
   */
  isConnected() {
    return this.sdk.isConnected();
  }
  /**
   * 添加事件监听器（简化版）
   */
  on(event, callback) {
    if (!this.eventCallbacks.has(event)) {
      this.eventCallbacks.set(event, []);
    }
    this.eventCallbacks.get(event).push(callback);
  }
  /**
   * 移除事件监听器
   */
  off(event, callback) {
    if (!callback) {
      this.eventCallbacks.delete(event);
      return;
    }
    const callbacks = this.eventCallbacks.get(event);
    if (callbacks) {
      const index = callbacks.indexOf(callback);
      if (index > -1) {
        callbacks.splice(index, 1);
      }
    }
  }
  /**
   * 触发回调函数
   */
  triggerCallback(event, ...args) {
    const callbacks = this.eventCallbacks.get(event);
    if (callbacks) {
      callbacks.forEach((callback) => {
        try {
          callback(...args);
        } catch (error) {
          console.warn(`Error in ${event} callback:`, error);
        }
      });
    }
  }
}
if (typeof window !== "undefined") {
  window.Servify = VanillaServifySDK;
  window.createServify = (config) => new VanillaServifySDK(config);
}
export {
  VanillaServifySDK,
  VanillaServifySDK as default
};
//# sourceMappingURL=index.esm.js.map
