const API_BASE = window.__API_BASE__ || 'http://localhost:8080';

function switchView(nextView) {
  document.querySelectorAll('.view').forEach((el) => el.classList.add('hidden'));
  const target = document.getElementById(`view-${nextView}`);
  if (target) {
    target.classList.remove('hidden');
  }
}

async function loadResources() {
  const body = document.getElementById('resources-table-body');
  const message = document.getElementById('resources-message');
  body.innerHTML = '';
  message.textContent = '読み込み中...';

  try {
    const res = await fetch(`${API_BASE}/api/resources`);
    if (!res.ok) {
      throw new Error(`status: ${res.status}`);
    }

    const payload = await res.json();
    const items = payload.data || [];

    if (items.length === 0) {
      message.textContent = '要員データはまだありません。';
      return;
    }

    for (const item of items) {
      const row = document.createElement('tr');
      row.innerHTML = `
        <td>${item.resource_id}</td>
        <td>${item.resource_name ?? ''}</td>
        <td>${item.department ?? ''}</td>
        <td>${item.note ?? ''}</td>
      `;
      body.appendChild(row);
    }

    message.textContent = `${items.length} 件を表示しています。`;
  } catch (error) {
    message.textContent = `読み込みに失敗しました: ${error.message}`;
  }
}

for (const button of document.querySelectorAll('nav button[data-view]')) {
  button.addEventListener('click', async () => {
    const view = button.getAttribute('data-view');
    switchView(view);
    if (view === 'resources') {
      await loadResources();
    }
  });
}

document.getElementById('reload-resources').addEventListener('click', loadResources);
