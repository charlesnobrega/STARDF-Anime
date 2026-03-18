/* JavaScript: StarDF-Anime Web Premium Controller */
const UI = {
    search: document.getElementById('search-input'),
    grid: document.getElementById('media-grid'),
    title: document.getElementById('section-title'),
    indicators: document.getElementById('search-indicators'),
    navItems: document.querySelectorAll('.nav-item'),
};

const state = {
    currentCategory: 'anime', // Default
    results: [],
};

// Event Listeners
UI.search.addEventListener('keyup', (e) => {
    if (e.key === 'Enter') performSearch(UI.search.value);
});

UI.navItems.forEach(item => {
    item.addEventListener('click', () => {
        UI.navItems.forEach(el => el.classList.remove('active'));
        item.classList.add('active');
        handleNavigation(item.id);
    });
});

// Search Logic
async function performSearch(query) {
    if (!query) return;

    UI.indicators.style.display = 'block';
    UI.grid.innerHTML = ''; // Clear and show skeletons
    for (let i = 0; i < 6; i++) {
        const skeleton = document.createElement('div');
        skeleton.className = 'media-card skeleton';
        UI.grid.appendChild(skeleton);
    }

    try {
        const response = await fetch(`/api/search?q=${encodeURIComponent(query)}&type=${state.currentCategory}`);
        const data = await response.json();
        
        state.results = data || [];
        renderMediaGrid(state.results);
        UI.title.innerText = `Resultados para "${query}" (${state.results.length})`;
    } catch (err) {
        console.error("Erro na busca:", err);
        UI.grid.innerHTML = `<div style="grid-column: 1/-1; padding: 40px; text-align: center; color: #ff4d4d; background: rgba(255, 77, 77, 0.1); border-radius: 12px; border: 1px solid #ff4d4d;">❌ Erro na busca: ${err.message}</div>`;
    } finally {
        UI.indicators.style.display = 'none';
    }
}

// Rendering System
function renderMediaGrid(results) {
    UI.grid.innerHTML = '';
    
    if (results.length === 0) {
        UI.grid.innerHTML = '<div style="grid-column: 1/-1; text-align: center; color: var(--text-secondary); padding: 50px;">Nenhum título encontrado 😢</div>';
        return;
    }

    results.forEach((item, index) => {
        const card = document.createElement('div');
        card.className = 'media-card glass-panel';
        card.style.animation = `fadeInUp 0.6s ease forwards ${index * 0.05}s`;
        card.style.opacity = '0';
        
        // Placeholder image fallback
        const image = item.ImageURL || 'https://via.placeholder.com/400x600?text=Indispon%C3%ADvel';
        
        card.innerHTML = `
            <img src="${image}" alt="${item.Name}" onerror="this.src='https://via.placeholder.com/400x600?text=Anime'">
            <div class="overlay">
                <div class="title">${item.Name}</div>
                <div class="info">${item.Source} • ${(item.TotalEpisodes || '?')} eps</div>
            </div>
            <div class="badge">${item.MediaType === 1 ? 'Filme' : 'Anime'}</div>
        `;
        
        card.addEventListener('click', () => loadDetails(item));
        UI.grid.appendChild(card);
    });
}

// Navigation Handler
function handleNavigation(id) {
    switch (id) {
        case 'nav-anime':
            state.currentCategory = 'anime';
            UI.title.innerText = "Buscar Animes (PT-BR)";
            break;
        case 'nav-movies':
            state.currentCategory = 'movie';
            UI.title.innerText = "Buscar Filmes & Séries";
            break;
        case 'nav-watchlist':
            UI.title.innerText = "Minha Lista (Soon)";
            break;
    }
}

// Details system (Episodes etc.)
async function loadDetails(item) {
    console.log("Selecionado:", item);
    // TODO: Implementar modal de episódios
    alert(`Carregando: ${item.Name}\nFonte: ${item.Source}`);
}

// Animation (defined in CSS via JS)
const style = document.createElement('style');
style.textContent = `
    @keyframes fadeInUp {
        from { transform: translateY(20px); opacity: 0; }
        to { transform: translateY(0); opacity: 1; }
    }
`;
document.head.appendChild(style);
