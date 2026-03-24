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
    settingsView: document.getElementById('settings-view'),
    profileView: document.getElementById('profile-view'),
    chatMessages: document.getElementById('chat-messages'),
    chatInput: document.getElementById('chatInput'),
    sendBtn: document.getElementById('sendMessage'),
    navSettingsBtn: document.getElementById('nav-settings'),
    profileBtn: document.getElementById('user-profile-btn'),
};

const state = {
    category: 'anime',
};

const ANI_DEFAULT_REDIRECT_URI = 'https://anilist.co/api/v2/oauth/pin';
const ANI_CLIENT_ID_STORAGE_KEY = 'anilist_client_id';

const isModalOpen = () => UI.modal && !UI.modal.classList.contains('hidden');

function openDetailsModal() {
    if (!UI.modal) return;
    UI.modal.classList.remove('hidden');
    UI.modal.style.display = 'flex';
    if (UI.reflection) UI.reflection.style.opacity = "1";
}

function closeDetailsModal() {
    if (!UI.modal) return;
    UI.modal.classList.add('hidden');
    UI.modal.style.display = '';
    if (UI.reflection) UI.reflection.style.opacity = "0";
}

window.closeDetailsModal = closeDetailsModal;

function escapeHtml(value) {
    return String(value ?? '')
        .replaceAll('&', '&amp;')
        .replaceAll('<', '&lt;')
        .replaceAll('>', '&gt;')
        .replaceAll('"', '&quot;')
        .replaceAll("'", '&#39;');
}

function formatAniListDate(isoDate) {
    if (!isoDate) return 'Sem data';
    const date = new Date(isoDate);
    if (Number.isNaN(date.getTime())) return 'Sem data';
    return date.toLocaleString('pt-BR', {
        dateStyle: 'short',
        timeStyle: 'short',
    });
}

function normalizeTitleForMatch(value) {
    return String(value || '')
        .toLowerCase()
        .normalize('NFD')
        .replace(/[\u0300-\u036f]/g, '')
        .replace(/^\s*\[[^\]]+\]\s*/g, '')
        .replace(/\s*[-:–—]?\s*(season|temporada|part|parte)\s*\d+\s*$/gi, '')
        .replace(/&/g, ' and ')
        .replace(/[^a-z0-9]+/g, ' ')
        .replace(/\s+/g, ' ')
        .trim();
}

function pickBestSearchResult(query, results) {
    if (!Array.isArray(results) || results.length === 0) return null;
    const q = normalizeTitleForMatch(query);
    if (!q) return null;
    return results.find((item) => {
        if (!item || !item.Name || !Array.isArray(item.Sources) || item.Sources.length === 0) return false;
        return normalizeTitleForMatch(item.Name) === q;
    }) || null;
}

if (UI.modal) {
    UI.modal.addEventListener('click', (e) => {
        if (e.target === UI.modal) closeDetailsModal();
    });
}

document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && isModalOpen()) closeDetailsModal();
});

// Search Logic
if (UI.search) {
    UI.search.addEventListener('keyup', (e) => {
        if (e.key === 'Enter') performSearch(UI.search.value);
    });
}

async function performSearch(query) {
    const normalizedQuery = String(query || '').trim();
    if (!normalizedQuery) {
        UI.dashboard.classList.remove('hidden');
        UI.searchContent.classList.add('hidden');
        return;
    }

    UI.status.classList.remove('hidden');
    UI.status.style.display = '';
    UI.dashboard.classList.add('hidden');
    UI.searchContent.classList.remove('hidden');
    UI.results.innerHTML = Array(8).fill('<div class="card" style="height: 300px; background: rgba(255,255,255,0.05); border-radius: 12px;"></div>').join('');

    try {
        const res = await fetch(`/api/search?q=${encodeURIComponent(normalizedQuery)}&type=${state.category}`);
        const contentType = res.headers.get('content-type') || '';
        const payload = contentType.includes('application/json') ? await res.json() : await res.text();
        if (!res.ok) {
            const errMsg = typeof payload === 'string' ? payload : (payload.error || 'Falha ao buscar conteudo.');
            throw new Error(errMsg);
        }
        const data = Array.isArray(payload) ? payload : [];
        renderSelection(UI.results, data);
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
    const safeImg = String(img).replace(/'/g, '%27');
    
    return `
        <article class="card ${isFeatured}" data-reflection="url('${safeImg}')">
            <div class="poster" style="background: url('${safeImg}');">
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
    if (!container) return;
    const normalizedItems = Array.isArray(items) ? items : [];
    container.innerHTML = normalizedItems.length ? normalizedItems.map(i => createCard(i, featured)).join("") : '<div style="padding: 20px;">Nenhum resultado encontrado.</div>';
    
    const cards = container.querySelectorAll('.card');
    cards.forEach((card, idx) => {
        const item = normalizedItems[idx];
        
        card.onmouseenter = () => {
            const img = card.dataset.reflection;
            if (UI.reflection) {
                UI.reflection.style.backgroundImage = img;
                UI.reflection.style.opacity = "1"; /* Vivid Engine: Full Opacity */
            }
        };
        
        card.onmouseleave = () => {
            if (UI.reflection && !isModalOpen()) UI.reflection.style.opacity = "0";
        };

        card.onclick = () => {
            if (item) showDetails(item);
        };
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
            <button id="modal-watchlist-btn"
                style="background: rgba(255,255,255,0.05); border: 1px solid var(--line); color: white; padding: 8px 16px; border-radius: 12px; cursor: pointer; display: flex; align-items: center; gap: 8px; font-weight: 600; font-size: 0.85rem; transition: 0.3s;">
                <i data-feather="plus"></i> Minha Lista
            </button>
        </div>
    `;
    updateWatchlistBtn(item.Name);
    const watchlistBtn = document.getElementById('modal-watchlist-btn');
    if (watchlistBtn) {
        watchlistBtn.onclick = () => {
            toggleWatchlist(item.Name, item.ImageURL, item.MediaType, item.TotalEpisodes);
        };
    }
    UI.modalImg.style.backgroundImage = `url('${item.ImageURL}')`;
    UI.modalEps.innerHTML = 'Aguardando Servidor...';
    UI.sourcesBox.innerHTML = '';
    openDetailsModal();

    try {
        let sources = item.Sources || [];
        
        // Se for sincronizado do Anilist, buscar fontes reais
        if (sources.length === 1 && sources[0].Name.includes("AniList")) {
            UI.modalEps.innerHTML = 'Buscando servidores disponíveis...';
            const res = await fetch(`/api/search?q=${encodeURIComponent(item.Name)}`);
            const results = await res.json().catch(() => []);
            if (!res.ok) {
                throw new Error('Falha ao buscar fontes para item sincronizado.');
            }
            const matched = pickBestSearchResult(item.Name, results);
            const fallback = Array.isArray(results)
                ? results.find((entry) => Array.isArray(entry.Sources) && entry.Sources.length > 0)
                : null;
            const picked = matched || fallback;
            if (picked && Array.isArray(picked.Sources)) {
                sources = picked.Sources;
            }
        }

        if (!sources.length) {
            UI.modalEps.innerHTML = 'Nenhum servidor encontrado para esta obra.';
            return;
        }

        // Renderizar botões de servidores
        const sourceEntries = [];
        sources.forEach((src) => {
            const btn = document.createElement('button');
            // USE ANIMENAME IF AVAILABLE TO DIFFERENTIATE DUB/LEG
            const label = src.AnimeName ? `${src.Name} - ${src.AnimeName.split(' ').slice(-1)}` : src.Name;
            btn.innerText = label;
            btn.title = src.AnimeName || src.Name;
            btn.className = "server-btn";
            btn.style.cssText = "padding: 8px 14px; background: rgba(255,255,255,0.05); border: 1px solid var(--line); color: white; border-radius: 8px; cursor: pointer; font-weight: 500; transition: 0.3s; font-size: 0.75rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 150px;";
            btn.onclick = () => { void loadEpisodes(src, item, btn); };
            UI.sourcesBox.appendChild(btn);
            sourceEntries.push({ source: src, button: btn });
        });

        // Auto-seleciona o primeiro servidor que retornar episódios.
        let loaded = false;
        for (const entry of sourceEntries) {
            const ok = await loadEpisodes(entry.source, item, entry.button, { silent: true });
            if (ok) {
                loaded = true;
                break;
            }
        }
        if (!loaded) {
            UI.modalEps.innerHTML = 'Nenhum episódio disponível nos servidores encontrados.';
        }

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

async function loadEpisodes(source, item, btn, options = {}) {
    const { silent = false } = options;
    if (!source || !source.URL || !source.Name) {
        if (!silent) UI.modalEps.innerHTML = 'Fonte inválida para carregar episódios.';
        return false;
    }

    // Estilizar botão selecionado
    const all = UI.sourcesBox.querySelectorAll('button');
    all.forEach(b => { b.style.background = "rgba(255,255,255,0.05)"; b.style.color = "white"; b.style.borderColor = "var(--line)"; });
    btn.style.background = "var(--cyan)";
    btn.style.color = "black";
    btn.style.borderColor = "var(--cyan)";

    if (!silent) UI.modalEps.innerHTML = 'Carregando episódios...';
    try {
        const res = await fetch(`/api/episodes?url=${encodeURIComponent(source.URL)}&source=${encodeURIComponent(source.Name)}`);
        const payload = await res.json().catch(() => ({}));
        if (!res.ok) {
            const message = payload && payload.error ? payload.error : 'Falha ao buscar episodios para a fonte selecionada.';
            throw new Error(message);
        }
        const eps = Array.isArray(payload) ? payload : [];
        if (!eps.length) {
            if (!silent) UI.modalEps.innerHTML = 'Sem episódios neste servidor.';
            return false;
        }
        UI.modalEps.innerHTML = '';
        
        eps.forEach(ep => {
            const epBtn = document.createElement('div');
            epBtn.style.cssText = "background: rgba(255,255,255,0.03); text-align: center; padding: 18px; border-radius: 12px; cursor: pointer; font-weight: 700; border: 1px solid var(--line); transition: 0.3s;";
            epBtn.innerText = ep.Number;
            epBtn.onclick = (e) => playEpisode(ep, item, source, e.currentTarget);
            UI.modalEps.appendChild(epBtn);
        });
        return true;
    } catch (e) {
        if (!silent) UI.modalEps.innerHTML = `Erro ao carregar episódios: ${e.message || 'falha desconhecida'}.`;
        return false;
    }
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
            const playData = await playRes.json().catch(() => ({}));
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
    if (isModalOpen()) closeDetailsModal();
    UI.dashboard.classList.add('hidden');
    UI.searchContent.classList.add('hidden');
    if(UI.watchlistView) UI.watchlistView.classList.add('hidden');
    if(UI.historyView) UI.historyView.classList.add('hidden');
    if(UI.settingsView) UI.settingsView.classList.add('hidden');
    if(UI.profileView) UI.profileView.classList.add('hidden');
}

// Navigation
const navDashboardBtn = document.getElementById('nav-dashboard');
const navWatchlistBtn = document.getElementById('nav-watchlist');
const navHistoryBtn = document.getElementById('nav-history');

if (navDashboardBtn) {
    navDashboardBtn.onclick = (e) => {
        hideAllViews();
        UI.dashboard.classList.remove('hidden');
        setActiveNav(e.currentTarget);
    };
}

if (navWatchlistBtn) {
    navWatchlistBtn.onclick = (e) => {
        hideAllViews();
        if(UI.watchlistView) {
            UI.watchlistView.classList.remove('hidden');
            renderWatchlist();
        }
        setActiveNav(e.currentTarget);
    };
}

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

if (navHistoryBtn) {
    navHistoryBtn.onclick = (e) => {
        hideAllViews();
        if(UI.historyView) {
            UI.historyView.classList.remove('hidden');
            renderHistory();
        }
        setActiveNav(e.currentTarget);
    };
}

if (UI.navSettingsBtn) {
    UI.navSettingsBtn.onclick = (e) => {
        hideAllViews();
        if (UI.settingsView) {
            UI.settingsView.classList.remove('hidden');
            renderSettings();
        }
        setActiveNav(e.currentTarget);
    };
}

if (UI.profileBtn) {
    UI.profileBtn.onclick = () => {
        hideAllViews();
        if (UI.profileView) {
            UI.profileView.classList.remove('hidden');
            renderProfile();
        }
        setActiveNav(null);
    };
}

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
            if (!isModalOpen()) UI.reflection.style.opacity = "0";
        };
        
        grid.appendChild(card);
    });
}

function renderSettings() {
    if (!UI.settingsView) return;

    const options = [
        { value: 'anime', label: 'Anime' },
        { value: 'movie', label: 'Filmes' },
        { value: 'tv', label: 'Series' },
    ];

    UI.settingsView.innerHTML = `
        <div class="section-head">
            <h2 style="font-family: 'Space Grotesk', sans-serif;">Configuracoes</h2>
            <button id="settings-back-dashboard" style="background: rgba(255,255,255,0.05); border: 1px solid var(--line); color: var(--muted); padding: 6px 12px; border-radius: 8px; cursor: pointer; font-size: 0.8rem;">
                Voltar ao Dashboard
            </button>
        </div>
        <div style="display: grid; gap: 18px; max-width: 900px;">
            <div style="background: rgba(255,255,255,0.03); border: 1px solid var(--line); border-radius: 14px; padding: 18px;">
                <h3 style="margin-bottom: 12px; font-size: 1rem;">Tipo de Conteudo da Busca</h3>
                <div id="settings-category-buttons" style="display: flex; gap: 10px; flex-wrap: wrap;">
                    ${options.map(opt => `
                        <button data-category-btn="${opt.value}" style="padding: 8px 12px; border-radius: 10px; border: 1px solid var(--line); background: ${state.category === opt.value ? 'rgba(57,226,255,0.18)' : 'rgba(255,255,255,0.03)'}; color: ${state.category === opt.value ? 'var(--cyan)' : 'var(--muted)'}; cursor: pointer;">
                            ${opt.label}
                        </button>
                    `).join('')}
                </div>
                <p style="margin-top: 12px; color: var(--muted); font-size: 0.85rem;">Define o filtro do campo de busca principal.</p>
            </div>
            <div style="background: rgba(255,255,255,0.03); border: 1px solid var(--line); border-radius: 14px; padding: 18px;">
                <h3 style="margin-bottom: 12px; font-size: 1rem;">Integracao AniList (Web)</h3>
                <p style="color: var(--muted); font-size: 0.9rem; line-height: 1.4;">
                    A visualizacao de status, sincronizacao, conexao e logout AniList esta no Perfil.
                </p>
            </div>
        </div>
    `;

    UI.settingsView.querySelectorAll('[data-category-btn]').forEach((btn) => {
        btn.onclick = () => {
            state.category = btn.getAttribute('data-category-btn') || 'anime';
            renderSettings();
        };
    });

    const backBtn = document.getElementById('settings-back-dashboard');
    if (backBtn) {
        backBtn.onclick = () => {
            if (navDashboardBtn) navDashboardBtn.click();
        };
    }
}

async function renderProfile() {
    if (!UI.profileView) return;

    const watchlist = JSON.parse(localStorage.getItem('anime_watchlist') || '[]');
    const localHistory = JSON.parse(localStorage.getItem('anime_history') || '[]');

    UI.profileView.innerHTML = `
        <div class="section-head">
            <h2 style="font-family: 'Space Grotesk', sans-serif;">Perfil</h2>
            <div style="display: flex; gap: 8px;">
                <button id="profile-refresh-btn" style="background: rgba(255,255,255,0.05); border: 1px solid var(--line); color: var(--muted); padding: 6px 12px; border-radius: 8px; cursor: pointer; font-size: 0.8rem;">Atualizar</button>
                <button id="profile-back-dashboard" style="background: rgba(255,255,255,0.05); border: 1px solid var(--line); color: var(--muted); padding: 6px 12px; border-radius: 8px; cursor: pointer; font-size: 0.8rem;">Dashboard</button>
            </div>
        </div>
        <div style="display: grid; gap: 18px; max-width: 900px;">
            <div style="background: rgba(255,255,255,0.03); border: 1px solid var(--line); border-radius: 14px; padding: 18px;">
                <h3 style="margin-bottom: 12px; font-size: 1rem;">Resumo Local</h3>
                <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 10px;">
                    <div style="background: rgba(0,0,0,0.2); border: 1px solid var(--line); border-radius: 10px; padding: 12px;">
                        <div style="color: var(--muted); font-size: 0.8rem;">Minha Lista</div>
                        <div style="font-size: 1.3rem; font-weight: 700;">${watchlist.length}</div>
                    </div>
                    <div style="background: rgba(0,0,0,0.2); border: 1px solid var(--line); border-radius: 10px; padding: 12px;">
                        <div style="color: var(--muted); font-size: 0.8rem;">Historico</div>
                        <div style="font-size: 1.3rem; font-weight: 700;">${localHistory.length}</div>
                    </div>
                </div>
            </div>
            <div style="background: rgba(255,255,255,0.03); border: 1px solid var(--line); border-radius: 14px; padding: 18px;">
                <h3 style="margin-bottom: 10px; font-size: 1rem;">AniList</h3>
                <div id="profile-anilist-status" style="color: var(--muted); font-size: 0.9rem;">Carregando status AniList...</div>
                <div id="profile-anilist-actions" style="display: flex; gap: 8px; margin-top: 12px;"></div>
                <div id="profile-anilist-connect" style="margin-top: 12px;"></div>
                <div id="profile-anilist-counts" style="display: grid; grid-template-columns: repeat(auto-fit, minmax(130px, 1fr)); gap: 8px; margin-top: 12px;"></div>
            </div>
            <div style="background: rgba(255,255,255,0.03); border: 1px solid var(--line); border-radius: 14px; padding: 18px;">
                <h3 style="margin-bottom: 12px; font-size: 1rem;">Sugestoes (AniList Sync)</h3>
                <div id="profile-anilist-suggestions" style="display: grid; grid-template-columns: repeat(auto-fill, minmax(210px, 1fr)); gap: 20px;"></div>
            </div>
            <div style="background: rgba(255,255,255,0.03); border: 1px solid var(--line); border-radius: 14px; padding: 18px;">
                <h3 style="margin-bottom: 12px; font-size: 1rem;">Historico AniList</h3>
                <div id="profile-anilist-history" style="display: grid; gap: 10px;"></div>
            </div>
        </div>
    `;

    const refreshBtn = document.getElementById('profile-refresh-btn');
    if (refreshBtn) {
        refreshBtn.onclick = () => renderProfile();
    }

    const backBtn = document.getElementById('profile-back-dashboard');
    if (backBtn) {
        backBtn.onclick = () => {
            if (navDashboardBtn) navDashboardBtn.click();
        };
    }

    const statusEl = document.getElementById('profile-anilist-status');
    const actionsEl = document.getElementById('profile-anilist-actions');
    const connectEl = document.getElementById('profile-anilist-connect');
    const countsEl = document.getElementById('profile-anilist-counts');
    const suggestionsEl = document.getElementById('profile-anilist-suggestions');
    const historyEl = document.getElementById('profile-anilist-history');

    let statusData = null;
    let syncData = null;
    let syncError = '';

    try {
        const statusRes = await fetch('/api/anilist/status', { cache: 'no-store' });
        statusData = await statusRes.json();

        if (statusData && statusData.loggedIn) {
            const syncRes = await fetch('/api/anilist/sync', { cache: 'no-store' });
            const syncBody = await syncRes.json();
            if (syncRes.ok) {
                syncData = syncBody;
            } else {
                syncError = syncBody.error || 'Falha ao sincronizar dados AniList.';
                if (syncRes.status === 401) {
                    statusData.loggedIn = false;
                }
            }
        }
    } catch (err) {
        syncError = err.message || 'Erro de conexao com AniList.';
    }

    const loggedIn = !!(statusData && statusData.loggedIn);
    const userName = statusData && statusData.user ? statusData.user.name : '';

    if (statusEl) {
        if (loggedIn) {
            statusEl.innerHTML = `
                <div style="color: #9df7bc; font-weight: 600;">Conectado como ${escapeHtml(userName || 'usuario AniList')}</div>
                <div style="color: var(--muted); margin-top: 5px;">Sincronizacao ativa para sugestoes e historico.</div>
            `;
        } else {
            const msg = syncError || (statusData && statusData.error) || 'Nao conectado no AniList.';
            statusEl.innerHTML = `
                <div style="color: #ffb5b5; font-weight: 600;">Nao conectado</div>
                <div style="color: var(--muted); margin-top: 5px;">${escapeHtml(msg)}</div>
                <div style="color: var(--muted); margin-top: 5px;">Use o fluxo abaixo para conectar com Client ID e codigo de autorizacao.</div>
            `;
        }
    }

    if (actionsEl) {
        actionsEl.innerHTML = loggedIn
            ? `<button id="profile-anilist-logout" style="background: rgba(255,90,90,0.15); border: 1px solid rgba(255,90,90,0.35); color: #ff9a9a; padding: 6px 12px; border-radius: 8px; cursor: pointer; font-size: 0.8rem;">Desconectar AniList</button>`
            : '';

        const logoutBtn = document.getElementById('profile-anilist-logout');
        if (logoutBtn) {
            logoutBtn.onclick = async () => {
                logoutBtn.disabled = true;
                logoutBtn.textContent = 'Desconectando...';
                try {
                    const resp = await fetch('/api/anilist/logout', { method: 'POST' });
                    if (!resp.ok) {
                        const payload = await resp.json().catch(() => ({}));
                        throw new Error(payload.error || 'Falha ao desconectar AniList.');
                    }
                    await renderProfile();
                } catch (err) {
                    alert(`Erro no logout AniList: ${err.message}`);
                    logoutBtn.disabled = false;
                    logoutBtn.textContent = 'Desconectar AniList';
                }
            };
        }
    }

    if (connectEl) {
        if (loggedIn) {
            connectEl.innerHTML = '';
        } else {
            const savedClientID = localStorage.getItem(ANI_CLIENT_ID_STORAGE_KEY) || '';
            connectEl.innerHTML = `
                <div style="display: grid; gap: 10px; background: rgba(0,0,0,0.2); border: 1px solid var(--line); border-radius: 10px; padding: 12px;">
                    <label style="font-size: 0.75rem; color: var(--muted);">Client ID AniList</label>
                    <input id="profile-anilist-client-id" type="text" value="${escapeHtml(savedClientID)}" placeholder="Ex: 12345" style="background: rgba(255,255,255,0.04); color: #fff; border: 1px solid var(--line); border-radius: 8px; padding: 9px 10px; outline: none;">

                    <label style="font-size: 0.75rem; color: var(--muted);">Client Secret (opcional)</label>
                    <input id="profile-anilist-client-secret" type="password" placeholder="Use se seu app AniList exigir secret" style="background: rgba(255,255,255,0.04); color: #fff; border: 1px solid var(--line); border-radius: 8px; padding: 9px 10px; outline: none;">

                    <label style="font-size: 0.75rem; color: var(--muted);">Codigo de autorizacao</label>
                    <input id="profile-anilist-code" type="text" placeholder="Cole o code retornado pelo AniList" style="background: rgba(255,255,255,0.04); color: #fff; border: 1px solid var(--line); border-radius: 8px; padding: 9px 10px; outline: none;">

                    <div style="display: flex; gap: 8px; flex-wrap: wrap;">
                        <button id="profile-anilist-open-auth" style="background: rgba(57,226,255,0.12); border: 1px solid rgba(57,226,255,0.4); color: var(--cyan); padding: 7px 12px; border-radius: 8px; cursor: pointer; font-size: 0.8rem;">1) Abrir autorizacao</button>
                        <button id="profile-anilist-connect-btn" style="background: rgba(157,247,188,0.12); border: 1px solid rgba(157,247,188,0.4); color: #9df7bc; padding: 7px 12px; border-radius: 8px; cursor: pointer; font-size: 0.8rem;">2) Conectar conta</button>
                    </div>

                    <div style="color: var(--muted); font-size: 0.78rem;">
                        Redirect URI padrao: <code>${ANI_DEFAULT_REDIRECT_URI}</code>
                    </div>
                </div>
            `;

            const clientIDEl = document.getElementById('profile-anilist-client-id');
            const clientSecretEl = document.getElementById('profile-anilist-client-secret');
            const codeEl = document.getElementById('profile-anilist-code');
            const openAuthBtn = document.getElementById('profile-anilist-open-auth');
            const connectBtn = document.getElementById('profile-anilist-connect-btn');

            if (openAuthBtn) {
                openAuthBtn.onclick = async () => {
                    const clientID = (clientIDEl ? clientIDEl.value : '').trim();
                    if (!clientID) {
                        alert('Informe o Client ID antes de abrir a autorizacao.');
                        return;
                    }

                    localStorage.setItem(ANI_CLIENT_ID_STORAGE_KEY, clientID);
                    openAuthBtn.disabled = true;
                    openAuthBtn.textContent = 'Abrindo...';
                    try {
                        const authRes = await fetch(`/api/anilist/auth-url?client_id=${encodeURIComponent(clientID)}&redirect_uri=${encodeURIComponent(ANI_DEFAULT_REDIRECT_URI)}`, { cache: 'no-store' });
                        const payload = await authRes.json().catch(() => ({}));
                        if (!authRes.ok || !payload.authUrl) {
                            throw new Error(payload.error || 'Falha ao gerar URL de autorizacao.');
                        }

                        const popup = window.open(payload.authUrl, '_blank', 'noopener,noreferrer');
                        if (!popup) {
                            alert(`Nao foi possivel abrir nova aba automaticamente.\n\nAbra manualmente:\n${payload.authUrl}`);
                        }
                        if (codeEl) codeEl.focus();
                    } catch (err) {
                        alert(`Erro ao abrir autorizacao AniList: ${err.message}`);
                    } finally {
                        openAuthBtn.disabled = false;
                        openAuthBtn.textContent = '1) Abrir autorizacao';
                    }
                };
            }

            if (connectBtn) {
                connectBtn.onclick = async () => {
                    const clientID = (clientIDEl ? clientIDEl.value : '').trim();
                    const clientSecret = (clientSecretEl ? clientSecretEl.value : '').trim();
                    const code = (codeEl ? codeEl.value : '').trim();

                    if (!clientID || !code) {
                        alert('Preencha Client ID e codigo de autorizacao.');
                        return;
                    }

                    localStorage.setItem(ANI_CLIENT_ID_STORAGE_KEY, clientID);
                    connectBtn.disabled = true;
                    connectBtn.textContent = 'Conectando...';
                    try {
                        const resp = await fetch('/api/anilist/login', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({
                                clientId: clientID,
                                clientSecret,
                                code,
                                redirectUri: ANI_DEFAULT_REDIRECT_URI,
                            }),
                        });

                        const payload = await resp.json().catch(() => ({}));
                        if (!resp.ok) {
                            throw new Error(payload.error || 'Falha ao conectar AniList.');
                        }

                        await renderProfile();
                    } catch (err) {
                        alert(`Erro ao conectar AniList: ${err.message}`);
                        connectBtn.disabled = false;
                        connectBtn.textContent = '2) Conectar conta';
                    }
                };
            }
        }
    }

    const counts = syncData && syncData.counts ? syncData.counts : null;
    if (countsEl) {
        if (loggedIn && counts) {
            const metricOrder = ['total', 'CURRENT', 'PLANNING', 'COMPLETED', 'PAUSED', 'DROPPED'];
            countsEl.innerHTML = metricOrder.map((key) => `
                <div style="background: rgba(0,0,0,0.2); border: 1px solid var(--line); border-radius: 10px; padding: 10px;">
                    <div style="color: var(--muted); font-size: 0.75rem;">${key}</div>
                    <div style="font-size: 1.1rem; font-weight: 700;">${counts[key] ?? 0}</div>
                </div>
            `).join('');
        } else {
            countsEl.innerHTML = '';
        }
    }

    const suggestions = syncData && Array.isArray(syncData.suggestions) ? syncData.suggestions : [];
    if (suggestionsEl) {
        if (!loggedIn) {
            suggestionsEl.innerHTML = `<p style="color: var(--muted);">Conecte sua conta AniList para carregar sugestoes personalizadas.</p>`;
        } else if (suggestions.length === 0) {
            suggestionsEl.innerHTML = `<p style="color: var(--muted);">Sem sugestoes no momento.</p>`;
        } else {
            renderSelection(suggestionsEl, suggestions.slice(0, 10));
        }
    }

    const aniHistory = syncData && Array.isArray(syncData.history) ? syncData.history : [];
    if (historyEl) {
        if (!loggedIn) {
            historyEl.innerHTML = `<p style="color: var(--muted);">Conecte sua conta AniList para visualizar o historico sincronizado.</p>`;
        } else if (aniHistory.length === 0) {
            historyEl.innerHTML = `<p style="color: var(--muted);">Nenhum item retornado pelo AniList.</p>`;
        } else {
            historyEl.innerHTML = aniHistory.slice(0, 12).map((item, idx) => `
                <button data-anilist-history-open="${idx}" style="text-align: left; background: rgba(255,255,255,0.03); border: 1px solid var(--line); color: #fff; border-radius: 10px; padding: 10px 12px; cursor: pointer;">
                    <div style="font-weight: 600;">${escapeHtml(item.name || 'Sem titulo')}</div>
                    <div style="color: var(--muted); font-size: 0.82rem; margin-top: 2px;">
                        ${escapeHtml(item.status || 'UNKNOWN')} • ${item.progress || 0}/${item.totalEpisodes || '?'} eps • ${escapeHtml(formatAniListDate(item.updatedAt))}
                    </div>
                </button>
            `).join('');

            historyEl.querySelectorAll('[data-anilist-history-open]').forEach((btn) => {
                btn.onclick = () => {
                    const idx = Number(btn.getAttribute('data-anilist-history-open'));
                    const item = aniHistory[idx];
                    if (!item) return;
                    const detailMedia = {
                        Name: item.name || 'Sem titulo',
                        ImageURL: item.imageUrl || '',
                        MediaType: item.mediaType || 'anime',
                        TotalEpisodes: item.totalEpisodes || 0,
                        Sources: item.sources && item.sources.length ? item.sources : [{ Name: 'AniList (Sync)', URL: item.name || '' }],
                    };
                    showDetails(detailMedia);
                };
            });
        }
    }
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
    document.querySelectorAll('.sidebar button').forEach(b => b.classList.remove('active'));
    if (btn) btn.classList.add('active');
}


if (UI.sendBtn && UI.chatInput && UI.chatMessages) {
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

        try {
            await fetch('/api/chat', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(newMsg)
            });
        } catch (err) {
            console.error('Chat send failed', err);
        }
    };
}

// Global bootstrap
window.addEventListener('load', () => {
    loadHome();
    if (window.feather) window.feather.replace();
}, { once: true });
