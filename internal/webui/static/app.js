/* ANIMEPRO Full-Screen Controller v9.1 */
const UI = {
    search: document.getElementById('searchInput'),
    dashboard: document.getElementById('dashboard-rails'),
    searchContent: document.getElementById('search-content'),
    trending: document.getElementById('trendingRail'),
    recommended: document.getElementById('recommendedRail'),
    recent: document.getElementById('recentRail'),
    results: document.getElementById('results-grid'),
    status: document.getElementById('search-status'),
    reflection: document.querySelector('.ambient-reflection'),
    // Modal
    modal: document.getElementById('modal-overlay'),
    modalTitle: document.getElementById('modal-title'),
    modalMeta: document.getElementById('modal-meta'),
    modalImg: document.getElementById('modal-poster'),
    modalEps: document.getElementById('episodes-box'),
    sourcesBox: document.getElementById('sources-buttons'),
    // Views
    watchlistView: document.getElementById('watchlist-view'),
    historyView: document.getElementById('history-view'),
    chatMessages: document.getElementById('chat-messages'),
    chatInput: document.getElementById('chatInput'),
    sendBtn: document.getElementById('sendMessage'),
};

const state = {
    category: 'anime',
};

// Search Logic
UI.search.addEventListener('keyup', (e) => {
    if (e.key === 'Enter') performSearch(UI.search.value);
});

async function performSearch(query) {
    if (!query) {
        UI.dashboard.classList.remove('hidden');
        UI.searchContent.classList.add('hidden');
        return;
    }

    UI.status.classList.remove('hidden');
    UI.dashboard.classList.add('hidden');
    UI.searchContent.classList.remove('hidden');
    UI.results.innerHTML = Array(8).fill('<div class="card" style="height: 300px; background: rgba(255,255,255,0.05); border-radius: 12px;"></div>').join('');

    try {
        const res = await fetch(`/api/search?q=${encodeURIComponent(query)}&type=${state.category}`);
        const data = await res.json();
        renderSelection(UI.results, data || []);
    } catch (err) {
        UI.results.innerHTML = `<div style="grid-column: 1/-1; padding: 50px; text-align: center; color: var(--pink);">Erro: ${err.message}</div>`;
    } finally {
        UI.status.style.display = 'none';
    }
}

// Clean Card Template (No Letters on Image, Titles below)
function createCard(item, featured = false) {
    const isFeatured = featured ? 'featured' : '';
    const name = item.Name || 'Sem título';
    const numSources = item.Sources ? item.Sources.length : 0;
    const sourceLabel = numSources > 1 ? `${numSources} Servidores` : (item.Sources && item.Sources[0] ? item.Sources[0].Name : 'AniList');
    const img = item.ImageURL || 'https://placehold.co/400x600/101525/fff?text=No+Poster';
    
    return `
        <article class="card ${isFeatured}" data-reflection="url('${img}')">
            <div class="poster" style="background: url('${img}');">
                <span class="badge">${sourceLabel}</span>
            </div>
            <div class="card-body">
                <h3>${name}</h3>
                <div class="card-meta">
                    <span>${item.MediaType || 'Série'} • ${item.TotalEpisodes || '?'} EPs</span>
                    <span class="rating">⭐ 4.8</span>
                </div>
                ${featured ? '<button class="watch-btn">Assistir agora</button>' : ''}
            </div>
        </article>
    `;
}

function renderSelection(container, items, featured = false) {
    container.innerHTML = (items && items.length) ? items.map(i => createCard(i, featured)).join("") : '<div style="padding: 20px;">Nenhum resultado encontrado.</div>';
    
    const cards = container.querySelectorAll('.card');
    cards.forEach((card, idx) => {
        const item = items[idx];
        
        card.onmouseenter = () => {
            const img = card.dataset.reflection;
            UI.reflection.style.backgroundImage = img;
            UI.reflection.style.opacity = "1"; /* Vivid Engine: Full Opacity */
        };
        
        card.onmouseleave = () => {
            if (UI.modal.style.display !== 'flex') UI.reflection.style.opacity = "0";
        };

        card.onclick = () => showDetails(item);
    });
}

// Start sequence
async function loadHome() {
    try {
        const res = await fetch(`/api/trending`);
        const data = await res.json();
        renderSelection(UI.trending, (data || []).slice(0, 4), true);
        renderSelection(UI.recommended, (data || []).slice(4, 8));
        renderSelection(UI.recent, (data || []).slice(8, 10)); // Adjust slice
    } catch (e) { console.error("Home loading error", e); }
    
    // Global Watchlist functions
    window.toggleWatchlist = (name, img, type, total) => {
        let list = JSON.parse(localStorage.getItem('anime_watchlist') || '[]');
        const exists = list.find(i => i.Name === name);
        if (exists) {
            list = list.filter(i => i.Name !== name);
        } else {
            list.push({ Name: name, ImageURL: img, MediaType: type, TotalEpisodes: total });
        }
        localStorage.setItem('anime_watchlist', JSON.stringify(list));
        updateWatchlistBtn(name);
    };
    
    window.updateWatchlistBtn = (name) => {
        const btn = document.getElementById('modal-watchlist-btn');
        if (!btn) return;
        const list = JSON.parse(localStorage.getItem('anime_watchlist') || '[]');
        const exists = list.find(i => i.Name === name);
        btn.innerHTML = exists ? '<i data-feather="check"></i> Na Lista' : '<i data-feather="plus"></i> Minha Lista';
        btn.style.background = exists ? 'rgba(57, 226, 255, 0.1)' : 'rgba(255,255,255,0.05)';
        btn.style.borderColor = exists ? 'var(--cyan)' : 'var(--line)';
        if (window.feather) window.feather.replace();
    };
}

async function showDetails(item) {
    UI.modalTitle.innerText = item.Name;
    UI.modalMeta.innerHTML = `
        <div style="display: flex; align-items: center; justify-content: space-between; width: 100%;">
            <span>${item.MediaType || 'Série'} | ⭐ 4.8</span>
            <button id="modal-watchlist-btn" onclick="toggleWatchlist('${item.Name.replace(/'/g, "\\'")}', '${item.ImageURL}', '${item.MediaType}', ${item.TotalEpisodes})" 
                style="background: rgba(255,255,255,0.05); border: 1px solid var(--line); color: white; padding: 8px 16px; border-radius: 12px; cursor: pointer; display: flex; align-items: center; gap: 8px; font-weight: 600; font-size: 0.85rem; transition: 0.3s;">
                <i data-feather="plus"></i> Minha Lista
            </button>
        </div>
    `;
    updateWatchlistBtn(item.Name);
    UI.modalImg.style.backgroundImage = `url('${item.ImageURL}')`;
    UI.modalEps.innerHTML = 'Aguardando Servidor...';
    UI.sourcesBox.innerHTML = '';
    UI.modal.style.display = 'flex';
    UI.reflection.style.opacity = "1"; 

    try {
        let sources = item.Sources || [];
        
        // Se for sincronizado do Anilist, buscar fontes reais
        if (sources.length === 1 && sources[0].Name.includes("AniList")) {
            UI.modalEps.innerHTML = 'Buscando servidores disponíveis...';
            const res = await fetch(`/api/search?q=${encodeURIComponent(item.Name)}`);
            const results = await res.json();
            const matched = results.find(r => r.Name.toLowerCase() === item.Name.toLowerCase()) || results[0];
            if (matched && matched.Sources) {
                sources = matched.Sources;
            }
        }

        if (!sources.length) {
            UI.modalEps.innerHTML = 'Nenhum servidor encontrado para esta obra.';
            return;
        }

        // Renderizar botões de servidores
        sources.forEach((src, idx) => {
            const btn = document.createElement('button');
            // USE ANIMENAME IF AVAILABLE TO DIFFERENTIATE DUB/LEG
            const label = src.AnimeName ? `${src.Name} - ${src.AnimeName.split(' ').slice(-1)}` : src.Name;
            btn.innerText = label;
            btn.title = src.AnimeName || src.Name;
            btn.className = "server-btn";
            btn.style.cssText = "padding: 8px 14px; background: rgba(255,255,255,0.05); border: 1px solid var(--line); color: white; border-radius: 8px; cursor: pointer; font-weight: 500; transition: 0.3s; font-size: 0.75rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 150px;";
            btn.onclick = () => loadEpisodes(src, item, btn);
            UI.sourcesBox.appendChild(btn);
            if (idx === 0) btn.click(); // Auto-selecionar primeiro
        });

        // Load Chat
        const resChat = await fetch('/api/chat');
        const msgs = await resChat.json();
        UI.chatMessages.innerHTML = msgs.map(m => `
            <div style="margin-bottom: 15px; padding: 12px 18px; background: rgba(255,255,255,0.05); border-radius: 15px; border-left: 4px solid var(--cyan);">
                <strong style="color: var(--pink); display: block; font-size: 0.8rem; margin-bottom: 5px;">${m.user}</strong>
                <span style="font-size: 0.95rem;">${m.text}</span>
            </div>
        `).join('');
        UI.chatMessages.scrollTop = UI.chatMessages.scrollHeight;

    } catch (err) { UI.modalEps.innerHTML = `Erro: ${err.message}`; }
}

async function loadEpisodes(source, item, btn) {
    // Estilizar botão selecionado
    const all = UI.sourcesBox.querySelectorAll('button');
    all.forEach(b => { b.style.background = "rgba(255,255,255,0.05)"; b.style.color = "white"; b.style.borderColor = "var(--line)"; });
    btn.style.background = "var(--cyan)";
    btn.style.color = "black";
    btn.style.borderColor = "var(--cyan)";

    UI.modalEps.innerHTML = 'Carregando episódios...';
    try {
        const res = await fetch(`/api/episodes?url=${encodeURIComponent(source.URL)}&source=${encodeURIComponent(source.Name)}`);
        const eps = await res.json();
        UI.modalEps.innerHTML = eps.length ? '' : 'Sem episódios neste servidor.';
        
        eps.forEach(ep => {
            const epBtn = document.createElement('div');
            epBtn.style.cssText = "background: rgba(255,255,255,0.03); text-align: center; padding: 18px; border-radius: 12px; cursor: pointer; font-weight: 700; border: 1px solid var(--line); transition: 0.3s;";
            epBtn.innerText = ep.Number;
            epBtn.onclick = (e) => playEpisode(ep, item, source, e.currentTarget);
            UI.modalEps.appendChild(epBtn);
        });
    } catch (e) { UI.modalEps.innerHTML = "Erro ao carregar episódios."; }
}

async function playEpisode(ep, item, source, btn) {
    const originalText = btn.innerText;
    btn.innerText = "⏳";
    btn.style.opacity = "0.5";
    btn.style.pointerEvents = "none";

    try {
        const res = await fetch(`/api/stream?url=${encodeURIComponent(ep.URL)}&source=${encodeURIComponent(source.Name)}`);
        
        if (!res.ok) {
            let errorMsg = "Erro no servidor ao buscar stream.";
            try {
                const data = await res.json();
                if (data.error) errorMsg = data.error;
            } catch (e) {
                errorMsg = await res.text() || errorMsg;
            }
            throw new Error(errorMsg);
        }

        const data = await res.json();
        
        if (!data.url) throw new Error("Stream URL não encontrada.");

        const referer = data.metadata ? (data.metadata.referer || data.metadata.Source || "") : "";
        const title = `${item.Name} - EP ${ep.Number}`;
        
        // Save to History
        saveToHistory(item, ep, source);
        
        const playRes = await fetch(`/api/play?url=${encodeURIComponent(data.url)}&referer=${encodeURIComponent(referer)}&title=${encodeURIComponent(title)}`);
        if (!playRes.ok) {
            const playData = await playRes.json();
            throw new Error(playData.error || "Erro desconhecido ao iniciar o player.");
        }
        
    } catch (err) {
        alert("Erro ao iniciar player: " + err.message);
    } finally {
        btn.innerText = originalText;
        btn.style.opacity = "1";
        btn.style.pointerEvents = "all";
    }
}

function hideAllViews() {
    UI.dashboard.classList.add('hidden');
    UI.searchContent.classList.add('hidden');
    if(UI.watchlistView) UI.watchlistView.classList.add('hidden');
    if(UI.historyView) UI.historyView.classList.add('hidden');
}

// Navigation
document.getElementById('nav-dashboard').onclick = (e) => {
    hideAllViews();
    UI.dashboard.classList.remove('hidden');
    setActiveNav(e.currentTarget);
};

document.getElementById('nav-watchlist').onclick = (e) => {
    hideAllViews();
    if(UI.watchlistView) {
        UI.watchlistView.classList.remove('hidden');
        renderWatchlist();
    }
    setActiveNav(e.currentTarget);
};

function renderWatchlist() {
    const list = JSON.parse(localStorage.getItem('anime_watchlist') || '[]');
    UI.watchlistView.innerHTML = `
        <div class="section-head">
            <h2 style="font-family: 'Space Grotesk', sans-serif;">Minha Lista</h2>
            <div style="color: var(--muted); font-size: 0.9rem;">${list.length} itens salvos</div>
        </div>
        <div id="watchlist-grid" style="display: grid; grid-template-columns: repeat(auto-fill, minmax(210px, 1fr)); gap: 30px; margin-top: 30px;"></div>
    `;
    
    const grid = document.getElementById('watchlist-grid');
    if (list.length === 0) {
        grid.innerHTML = `<p style="color:var(--muted); grid-column: 1/-1; padding: 40px; text-align: center;">Sua lista está vazia. Adicione animes para acompanhar.</p>`;
        return;
    }

    renderSelection(grid, list);
}

document.getElementById('nav-history').onclick = (e) => {
    hideAllViews();
    if(UI.historyView) {
        UI.historyView.classList.remove('hidden');
        renderHistory();
    }
    setActiveNav(e.currentTarget);
};

function saveToHistory(item, ep, source) {
    let history = JSON.parse(localStorage.getItem('anime_history') || '[]');
    // Remove if already exists to move to top
    history = history.filter(h => h.item.Name !== item.Name);
    history.unshift({
        item,
        ep,
        source,
        timestamp: new Date().getTime()
    });
    // Limit to 20 items
    if (history.length > 20) history.pop();
    localStorage.setItem('anime_history', JSON.stringify(history));
}

function renderHistory() {
    const history = JSON.parse(localStorage.getItem('anime_history') || '[]');
    UI.historyView.innerHTML = `
        <div class="section-head">
            <h2 style="font-family: 'Space Grotesk', sans-serif;">Continuar Assistindo</h2>
            <button onclick="clearHistory()" style="background: rgba(255,255,255,0.05); border: 1px solid var(--line); color: var(--muted); padding: 5px 12px; border-radius: 8px; cursor: pointer; font-size: 0.75rem;">Limpar Tudo</button>
        </div>
        <div id="history-grid" style="display: grid; grid-template-columns: repeat(auto-fill, minmax(210px, 1fr)); gap: 30px; margin-top: 30px;"></div>
    `;
    
    const grid = document.getElementById('history-grid');
    if (history.length === 0) {
        grid.innerHTML = `<p style="color:var(--muted); grid-column: 1/-1; padding: 40px; text-align: center;">Nenhum histórico recente encontrado.</p>`;
        return;
    }

    history.forEach(h => {
        const cardWrapper = document.createElement('div');
        const itemWithEp = {...h.item, Name: `${h.item.Name} (EP ${h.ep.Number})`};
        cardWrapper.innerHTML = createCard(itemWithEp);
        const card = cardWrapper.firstElementChild;
        
        card.onclick = () => showDetails(h.item);
        card.onmouseenter = () => {
            UI.reflection.style.backgroundImage = `url('${h.item.ImageURL}')`;
            UI.reflection.style.opacity = "1";
        };
        card.onmouseleave = () => {
            if (UI.modal.style.display !== 'flex') UI.reflection.style.opacity = "0";
        };
        
        grid.appendChild(card);
    });
}

function clearHistory() {
    if(confirm("Deseja limpar todo o histórico?")) {
        localStorage.removeItem('anime_history');
        renderHistory();
    }
}

// "Messages" Sidebar remover
const navMsgs = document.getElementById('nav-messages');
if(navMsgs) navMsgs.style.display = 'none';

function setActiveNav(btn) {
    document.querySelectorAll('.nav button').forEach(b => b.classList.remove('active'));
    btn.classList.add('active');
}


UI.sendBtn.onclick = async () => {
    const text = UI.chatInput.value.trim();
    if (!text) return;
    
    UI.chatInput.value = '';
    const newMsg = { user: "Hiroshi K.", text };
    
    // Optimistic update
    const div = document.createElement('div');
    div.style.cssText = "margin-bottom: 15px; padding: 12px 18px; background: rgba(255,255,255,0.05); border-radius: 15px; border-left: 4px solid var(--cyan); opacity: 0.7;";
    div.innerHTML = `<strong style="color: var(--pink); display: block; font-size: 0.8rem; margin-bottom: 5px;">${newMsg.user}</strong><span>${newMsg.text}</span>`;
    UI.chatMessages.appendChild(div);
    UI.chatMessages.scrollTop = UI.chatMessages.scrollHeight;

    await fetch('/api/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newMsg)
    });
};

// Global bootstrap
window.addEventListener('load', () => {
    loadHome();
    if (window.feather) window.feather.replace();
}, { once: true });
