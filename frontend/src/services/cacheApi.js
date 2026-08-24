const API_BASE = import.meta.env.VITE_API_URL || "http://localhost:8080";

export async function setItem(key, value, ttlSeconds) {
  const res = await fetch(`${API_BASE}/set`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ key, value, ttl_seconds: ttlSeconds }),
  });
  return res.ok;
}

export async function getItem(key) {
  const res = await fetch(`${API_BASE}/get?key=${encodeURIComponent(key)}`);
  return res.json();
}

export async function getAll() {
  const res = await fetch(`${API_BASE}/all`);
  return res.json();
}