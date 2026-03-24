/* ─── Vernon License — Shared API Utilities ─────────────────────────────── */

const API_BASE = '';

// ── Auth ────────────────────────────────────────────────────────────────────

function getToken() { return localStorage.getItem('vl_token'); }
function setToken(t) { localStorage.setItem('vl_token', t); }
function removeToken() { localStorage.removeItem('vl_token'); }

function getUser() {
  try { return JSON.parse(localStorage.getItem('vl_user') || 'null'); } catch { return null; }
}
function setUser(u) { localStorage.setItem('vl_user', JSON.stringify(u)); }
function removeUser() { localStorage.removeItem('vl_user'); }

function logout() {
  removeToken();
  removeUser();
  window.location.replace('/');
}

function requireAuth() {
  if (!getToken()) { window.location.replace('/'); }
}

// ── HTTP Client ─────────────────────────────────────────────────────────────

async function api(method, path, body) {
  const headers = { 'Content-Type': 'application/json' };
  const token = getToken();
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const res = await fetch(API_BASE + path, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  if (res.status === 401) {
    logout();
    throw new Error('Sesi berakhir');
  }

  const json = await res.json().catch(() => null);

  if (!res.ok) {
    const msg = json?.error?.message || `HTTP ${res.status}`;
    throw new Error(msg);
  }

  return json;
}

// ── Formatters ───────────────────────────────────────────────────────────────

function rupiah(n) {
  if (n == null) return '-';
  return new Intl.NumberFormat('id-ID', {
    style: 'currency', currency: 'IDR', minimumFractionDigits: 0,
  }).format(n);
}

function tanggal(s) {
  if (!s) return '-';
  return new Date(s).toLocaleDateString('id-ID', {
    day: 'numeric', month: 'short', year: 'numeric',
  });
}

function timeAgo(s) {
  if (!s) return '';
  const diff = Date.now() - new Date(s).getTime();
  const m = Math.floor(diff / 60000);
  if (m < 1) return 'baru saja';
  if (m < 60) return `${m} mnt lalu`;
  const h = Math.floor(m / 60);
  if (h < 24) return `${h} jam lalu`;
  const d = Math.floor(h / 24);
  return `${d} hari lalu`;
}

// ── Register service worker ──────────────────────────────────────────────────

if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js').catch(() => {});
  });
}
