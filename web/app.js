const cards = document.querySelector("#cards");
const statusLine = document.querySelector("#status");

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
    statusLine.textContent = "Could not create expedition";
    return;
  }
  event.currentTarget.reset();
  await loadExpeditions();
});

loadExpeditions();
