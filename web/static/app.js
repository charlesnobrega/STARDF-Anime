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
        UI.dashboard.style.display = 'block';
        UI.searchContent.style.display = 'none';
        return;
    }

    UI.status.style.display = 'block';
    UI.dashboard.style.display = 'none';
    UI.searchContent.style.display = 'block';
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
    const source = item.Source || 'Desconhecido';
    const img = item.ImageURL || 'https://via.placeholder.com/400x600?text=No+Poster';
    
    return `
        <article class="card ${isFeatured}" data-reflection="url('${img}')">
            <div class="poster" style="background: url('${img}');">
                <span class="badge">${source}</span>
            </div>
            <div class="card-body">
                <h3>${name}</h3>
                <div class="card-meta">
                    <span>${item.MediaType || 'Série'} • ${item.TotalEpisodes || '?'} EPs</span>
                    <span class="rating">4.8</span>
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
        const res = await fetch(`/api/search?q=Solo Leveling&type=anime`);
        const data = await res.json();
        renderSelection(UI.trending, (data || []).slice(0, 4), true);
        renderSelection(UI.recommended, (data || []).slice(4, 11));
        renderSelection(UI.recent, (data || []).slice(11, 18));
    } catch (e) { console.error("Home loading error", e); }
}

async function showDetails(item) {
    UI.modalTitle.innerText = item.Name;
    UI.modalMeta.innerText = `${item.MediaType} | ${item.Source} | ⭐ 4.8`;
    UI.modalImg.style.backgroundImage = `url('${item.ImageURL}')`;
    UI.modalEps.innerHTML = 'Carregando episódios...';
    UI.modal.style.display = 'flex';
    UI.reflection.style.opacity = "1"; 

    try {
        const res = await fetch(`/api/episodes?url=${encodeURIComponent(item.URL)}&source=${encodeURIComponent(item.Source)}`);
        const eps = await res.json();
        UI.modalEps.innerHTML = eps.length ? '' : 'Sem episódios.';
        
        eps.forEach(ep => {
            const btn = document.createElement('div');
            btn.style.cssText = "background: rgba(255,255,255,0.05); text-align: center; padding: 18px; border-radius: 12px; cursor: pointer; font-weight: 700; border: 1px solid var(--line); transition: 0.3s;";
            btn.innerText = ep.Number;
            btn.onclick = () => window.open(`/api/stream?url=${encodeURIComponent(ep.URL)}&source=${encodeURIComponent(item.Source)}`, '_blank');
            UI.modalEps.appendChild(btn);
        });
    } catch (err) { UI.modalEps.innerHTML = `Erro: ${err.message}`; }
}

// Global bootstrap
window.addEventListener('load', () => {
    loadHome();
    if (window.feather) window.feather.replace();
}, { once: true });
