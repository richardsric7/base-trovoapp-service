// Minimal HashiCorp Vault KV v2 reader (no dependencies), matching the
// VAULT_ADDR / VAULT_TOKEN convention used by tm-api and the quote service.

async function readKV2({ addr, token, namespace, mount, path }) {
  const url = `${addr.replace(/\/+$/, "")}/v1/${mount.replace(/^\/+|\/+$/g, "")}/data/${path.replace(/^\/+/, "")}`;
  const headers = { "X-Vault-Token": token };
  if (namespace) headers["X-Vault-Namespace"] = namespace;
  const res = await fetch(url, { headers });
  if (res.status === 404) throw new Error(`Vault secret ${mount}/${path} not found`);
  if (!res.ok) throw new Error(`Vault read ${mount}/${path} failed: HTTP ${res.status}`);
  const body = await res.json();
  const data = body && body.data && body.data.data;
  if (!data || typeof data !== "object") throw new Error(`Vault secret ${mount}/${path} has no data`);
  return data;
}

// loadConfig returns the key/value config for this tool: from the Vault
// secret at secretPathEnv (default source), or from process.env when
// RECOVERY_CONFIG_SOURCE=env (local development only).
async function loadConfig(secretPathEnv, defaultPath) {
  const source = (process.env.RECOVERY_CONFIG_SOURCE || "vault").toLowerCase();
  if (source === "env") return { get: (k) => process.env[k], source: "environment" };
  if (source !== "vault") throw new Error(`RECOVERY_CONFIG_SOURCE must be "vault" or "env", got "${source}"`);

  const addr = process.env.VAULT_ADDR;
  const token = process.env.VAULT_TOKEN;
  if (!addr) throw new Error("VAULT_ADDR is not set");
  if (!token) throw new Error("VAULT_TOKEN is not set");
  const mount = process.env.VAULT_KV_MOUNT || "secret";
  const path = process.env[secretPathEnv] || defaultPath;
  const data = await readKV2({ addr, token, namespace: process.env.VAULT_NAMESPACE, mount, path });
  return { get: (k) => (data[k] === undefined || data[k] === null ? undefined : String(data[k])), source: `vault ${mount}/${path}` };
}

module.exports = { readKV2, loadConfig };
