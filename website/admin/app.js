const tokenKey = "jlg_admin_token";

const loginView = document.querySelector("#loginView");
const dashboardView = document.querySelector("#dashboardView");
const loginForm = document.querySelector("#loginForm");
const loginError = document.querySelector("#loginError");
const usernameInput = document.querySelector("#usernameInput");
const passwordInput = document.querySelector("#passwordInput");
const dateInput = document.querySelector("#dateInput");
const keywordInput = document.querySelector("#keywordInput");
const levelInput = document.querySelector("#levelInput");
const statusText = document.querySelector("#statusText");

function today() {
    const now = new Date();
    const offset = now.getTimezoneOffset() * 60000;
    return new Date(now.getTime() - offset).toISOString().slice(0, 10);
}

function getToken() {
    return localStorage.getItem(tokenKey) || "";
}

function setToken(token) {
    localStorage.setItem(tokenKey, token);
}

function clearToken() {
    localStorage.removeItem(tokenKey);
}

function showLogin(message = "") {
    loginView.classList.remove("hidden");
    dashboardView.classList.add("hidden");
    loginError.textContent = message;
}

function showDashboard() {
    loginView.classList.add("hidden");
    dashboardView.classList.remove("hidden");
}

async function apiFetch(path, options = {}) {
    const headers = {
        ...(options.headers || {}),
    };
    if (getToken()) {
        headers.Authorization = `Bearer ${getToken()}`;
    }

    const response = await fetch(path, {
        ...options,
        headers,
    });

    if (response.status === 401 && !path.startsWith("/api/admin/login")) {
        clearToken();
        showLogin("登录已过期，请重新登录。");
        throw new Error("unauthorized");
    }

    const body = await response.json();
    if (!response.ok || body.code !== 200) {
        throw new Error(body.data || body.msg || "请求失败");
    }
    return body.data;
}

function formatDateTime(value) {
    if (!value) return "";
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return value;
    return date.toLocaleString("zh-CN", {
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
        hour12: false,
    });
}

function escapeHTML(value) {
    return String(value ?? "")
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;")
        .replaceAll("'", "&#039;");
}

function endpoint(path, params = {}) {
    const url = new URL(path, window.location.origin);
    Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== "") {
            url.searchParams.set(key, value);
        }
    });
    return url.pathname + url.search;
}

function renderRank(container, items, labelKey, valueKey) {
    if (!items || items.length === 0) {
        container.innerHTML = '<div class="empty">暂无数据</div>';
        return;
    }
    container.innerHTML = items.map((item) => `
        <div class="rank-row">
            <strong title="${escapeHTML(item[labelKey])}">${escapeHTML(item[labelKey])}</strong>
            <span>${escapeHTML(item[valueKey])}</span>
        </div>
    `).join("");
}

function statusClass(status) {
    if (status >= 500) return "error";
    if (status >= 400) return "warn";
    return "";
}

function renderAccessLogs(data) {
    document.querySelector("#accessSource").textContent = data.source_exists ? "" : "日志文件不存在";
    const body = document.querySelector("#accessLogBody");
    if (!data.items || data.items.length === 0) {
        body.innerHTML = '<tr><td colspan="6" class="empty">暂无访问日志</td></tr>';
        return;
    }

    body.innerHTML = data.items.map((item) => `
        <tr>
            <td class="mono">${escapeHTML(formatDateTime(item.time))}</td>
            <td class="mono">${escapeHTML(item.ip)}</td>
            <td>${escapeHTML(item.method)}</td>
            <td class="path-cell mono" title="${escapeHTML(item.path)}">${escapeHTML(item.path)}</td>
            <td><span class="status-pill ${statusClass(item.status)}">${escapeHTML(item.status)}</span></td>
            <td class="ua-cell" title="${escapeHTML(item.user_agent)}">${escapeHTML(item.user_agent)}</td>
        </tr>
    `).join("");
}

function renderTextLogs(container, sourceEl, data, emptyText) {
    sourceEl.textContent = data.source_exists ? "" : "日志文件不存在";
    if (!data.items || data.items.length === 0) {
        container.innerHTML = `<div class="empty">${emptyText}</div>`;
        return;
    }
    container.innerHTML = data.items.map((item) => `
        <pre class="log-line">${escapeHTML(item.raw)}</pre>
    `).join("");
}

function renderSummary(summary) {
    document.querySelector("#totalRequests").textContent = summary.total_requests || 0;
    document.querySelector("#uniqueIps").textContent = summary.unique_ips || 0;
    document.querySelector("#status4xx").textContent = summary.status_4xx || 0;
    document.querySelector("#status5xx").textContent = summary.status_5xx || 0;
    renderRank(document.querySelector("#topPaths"), summary.top_paths || [], "path", "count");
    renderRank(document.querySelector("#statusCounts"), summary.status_counts || [], "status", "count");
}

async function loadDashboard() {
    const params = {
        date: dateInput.value || today(),
        keyword: keywordInput.value.trim(),
        limit: 100,
    };
    statusText.textContent = "正在加载...";

    const [summary, accessLogs, errorLogs, appLogs] = await Promise.all([
        apiFetch(endpoint("/api/admin/ops/summary", { date: params.date })),
        apiFetch(endpoint("/api/admin/ops/access-logs", params)),
        apiFetch(endpoint("/api/admin/ops/error-logs", params)),
        apiFetch(endpoint("/api/admin/ops/app-logs", { ...params, level: levelInput.value })),
    ]);

    renderSummary(summary);
    renderAccessLogs(accessLogs);
    renderTextLogs(document.querySelector("#errorLogs"), document.querySelector("#errorSource"), errorLogs, "暂无 Nginx 错误日志");
    renderTextLogs(document.querySelector("#appLogs"), document.querySelector("#appSource"), appLogs, "暂无后端应用日志");
    statusText.textContent = `已更新：${new Date().toLocaleTimeString("zh-CN", { hour12: false })}`;
}

loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    loginError.textContent = "";
    try {
        const data = await apiFetch("/api/admin/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                username: usernameInput.value.trim(),
                password: passwordInput.value,
            }),
        });
        setToken(data.token);
        passwordInput.value = "";
        showDashboard();
        await loadDashboard();
    } catch (error) {
        if (error.message !== "unauthorized") {
            loginError.textContent = "登录失败，请检查账号和密码。";
        }
    }
});

document.querySelector("#logoutBtn").addEventListener("click", () => {
    clearToken();
    showLogin();
});

document.querySelector("#refreshBtn").addEventListener("click", () => {
    loadDashboard().catch((error) => {
        if (error.message !== "unauthorized") {
            statusText.textContent = error.message || "加载失败";
        }
    });
});

dateInput.value = today();

if (getToken()) {
    showDashboard();
    loadDashboard().catch((error) => {
        if (error.message !== "unauthorized") {
            statusText.textContent = error.message || "加载失败";
        }
    });
} else {
    showLogin();
}
