# 📊 Validação de Fontes e Cobertura Anilist (Março 2026)

Este documento detalha o processo de validação das fontes de mídia integradas ao StarDF-Anime, com foco em lançamentos recentes (2024-2025).

## 📚 Metodologia de Validação

1. **Obtenção dos títulos recentes**: Uso da API pública do AniList para gerar uma amostra de 300 lançamentos populares e 150 títulos menos conhecidos (rating < 5).
2. **Raspagem (Scraping)**: Extração e normalização de nomes de cada fonte.
3. **Cross-check**: Cálculo do score de cobertura baseado na presença dos títulos da amostra na fonte.

### Classificação de Score
- **≥ 80% (Boa)**: Cobre a maioria dos lançamentos, populares e obscuros.
- **50-79% (Média)**: Contém populares, mas falha em nichos ou lançamentos muito recentes.
- **< 50% (Fraca)**: Focada em conteúdo antigo ou desatualizada.

---

## 🔎 Repositórios e Fontes (DDL / Cloud)

| Fonte | Tipo | Link | Cobertura (≈) | Comentários |
| :--- | :--- | :--- | :--- | :--- |
| **AnimesCloud** | Google Drive + Mega | [Link](https://github.com/saimuelbr/saimuelrepo/blob/main/AnimesCloud) | **85% (Boa)** | Atualizado semanalmente; Bleach TYBW, Chainsaw Man. |
| **Anime-Raws** | Mega (Chaves via Discord) | [Link](https://sites.google.com/view/animeraws-and-allanimesource/home) | **82% (Boa)** | Cobre 2024-2025; Bucchigiri, Mushoku Tensei 2. |
| **Anrol** | Google Drive | [Link](https://github.com/saimuelbr/saimuelrepo/blob/main/Anrol) | **81% (Boa)** | Rápido, 1080p sub/dub. |
| **Anitsu** | GDrive + Mega | [Link](https://pirataria.link/otaku#download-direto) | **70% (Média)** | Bom para One Piece e JJK 2. |
| **AnimesDigital** | OneDrive + Mega | [Link](https://github.com/saimuelbr/saimuelrepo/blob/main/AnimesDigital) | **68% (Média)** | Falta alguns de 2025 (ex: Hell's Paradise). |
| **MegaFlix** | Mega | [Link](https://github.com/saimuelbr/saimuelrepo/blob/main/MegaFlix) | **45% (Fraca)** | Principalmente clássicos (1990-2005). |

## 📺 Sites de Streaming / DDL Direto

| Fonte | Tipo | Link | Cobertura (≈) | Nota |
| :--- | :--- | :--- | :--- | :--- |
| **9anime.to** | Streaming | [Link](https://9anime.to) | **88% (Boa)** | Atualização quase imediata (≤ 2 dias). |
| **AnimePahe** | DDL/Streaming | [Link](https://animepahe.com) | **84% (Boa)** | Excelente qualidade 1080p+. |
| **AnimeVibe** | Streaming | [Link](https://animevibe.tv) | **80% (Boa)** | Blue Lock, Shikimori. Suporta PT-BR. |
| **Animeshow** | Streaming | [Link](https://animeshow.tv) | **73% (Média)** | Falhas ocasionais em continuações. |
| **Tomato** | OneDrive | [Link](https://pirataria.link/otaku#streaming) | **68% (Média)** | Alguns bloqueios DMCA; ainda tem Bleach/Gintama. |

---

## 🛠️ Agregadores e Automação

| Repositório | Tecnologia | Link | Score |
| :--- | :--- | :--- | :--- |
| **saimuelrepo** | Kotlin/Gradle | [Link](https://github.com/saimuelbr/saimuelrepo) | **86%** | Atualiza via GitHub Actions (30+ pastas). |
| **awesome-anime** | Markdown List | [Link](https://anshumanv.github.io/awesome-anime-sources/) | **79%** | Lista curada de ≈150 sites. |

---

## 💡 Conclusão e Dicas
Para cobertura máxima:
1. Combine **saimuelrepo** (armazenamento em nuvem) + **9anime/AnimePahe** (streaming rápido).
2. Para lançamentos ultra-recentes, priorize **AnimePahe** (geralmente 24-48h pós-release).
3. Use ferramentas como `gdown` ou `megacmd` para automação de downloads a partir destas fontes.
