const BASE = '/api/v1';

async function fetchJSON(url, opts = {}) {
  const res = await fetch(url, {
    headers: { 'Content-Type': 'application/json', ...opts.headers },
    ...opts,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || `HTTP ${res.status}`);
  }
  return res.json();
}

export function getDashboardStats() {
  return fetchJSON(`${BASE}/dashboard/stats`);
}

export function getMerchantNetwork() {
  return fetchJSON(`${BASE}/dashboard/network`);
}

export function checkMerchant(qrisMerchantId, merchantName = '') {
  return fetchJSON(`${BASE}/merchant/check`, {
    method: 'POST',
    body: JSON.stringify({ qris_merchant_id: qrisMerchantId, merchant_name: merchantName }),
  });
}

export function searchMerchants(q) {
  return fetchJSON(`${BASE}/merchant/search?q=${encodeURIComponent(q)}`);
}

export function getMerchantByID(id) {
  return fetchJSON(`${BASE}/merchant/${encodeURIComponent(id)}`);
}

export function submitReport(merchantId, reporterNote = '', evidenceUrl = '') {
  return fetchJSON(`${BASE}/report/submit`, {
    method: 'POST',
    body: JSON.stringify({ merchant_id: merchantId, reporter_note: reporterNote, evidence_url: evidenceUrl }),
  });
}

export function getReports() {
  return fetchJSON(`${BASE}/report/list`);
}
