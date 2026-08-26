const cards = document.querySelector("#cards");
const statusLine = document.querySelector("#status");

async function clientErrorMessage(response) {
  let payload = null;
  try {
    payload = await response.json();
  } catch {
    payload = null;
  }
  const code = payload && payload.code;
  const detail = payload && payload.message;
  if (code === "invalid_input") {
    return detail ? `无法登记：${detail}` : "输入不完整或数值不合规，请检查后再提交";
  }
  if (code === "invalid_state") {
    return detail ? `状态冲突：${detail}` : "该队伍当前状态不允许此操作";
  }
  if (code === "not_found") {
    return "未找到对应的记录或队伍";
  }
  if (response.status >= 500) {
    return "服务器暂时无法处理，请稍后重试";
  }
  return detail || "提交失败，请检查输入";
}

function render(items) {
  cards.innerHTML = "";
  if (!items.length) {
    cards.innerHTML = "<p class='meta'>No expeditions yet. Start with a planned route.</p>";
    return;
  }
  for (const item of items) {
    const card = document.createElement("article");
    card.className = "card";
    card.innerHTML = `<h3>${item.name}</h3><p class="meta">${item.region} · ${item.lead}</p><p class="meta">Status: ${item.status}</p>`;
    cards.append(card);
  }
}

async function loadExpeditions() {
  statusLine.textContent = "Refreshing...";
  const response = await fetch("/api/expeditions");
  const payload = await response.json();
  render(payload.items || []);
  statusLine.textContent = `${(payload.items || []).length} expedition(s)`;
}

document.querySelector("#refresh").addEventListener("click", loadExpeditions);
document.querySelector("#expedition-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = new FormData(event.currentTarget);
  const response = await fetch("/api/expeditions", {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({
      name: form.get("name"),
      region: form.get("region"),
      lead: form.get("lead"),
      start_date: new Date().toISOString()
    })
  });
  if (!response.ok) {
    statusLine.textContent = await clientErrorMessage(response);
    return;
  }
  event.currentTarget.reset();
  await loadExpeditions();
});

loadExpeditions();
