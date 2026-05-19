# Implementation Plan: Mobile Android Fix

## Overview

Este plano implementa a correção completa da versão mobile em Flutter do StarDF-Anime, abordando interface com glassmorphism, inicialização do servidor web local, bridge Go-Flutter funcional, sincronização com desktop, tratamento de erros robusto, testes de integração, suporte a Android 8.0+, performance otimizada, ciclo de vida gerenciado e logging completo.

A implementação segue uma arquitetura em camadas com separação clara entre UI (Flutter), lógica de negócio (Dart), bridge (GoMobile) e backend (Go), garantindo comunicação robusta e sincronização bidirecional de dados.

## Tasks

- [ ] 1. Setup e Configuração Inicial
  - [x] 1.1 Configurar estrutura de projeto Flutter com GoMobile
    - Criar diretórios: lib/, go/, test/, integration_test/
    - Configurar pubspec.yaml com dependências (provider, drift, http, etc)
    - Configurar build.gradle para GoMobile
    - _Requirements: 3, 7_

  - [x] 1.2 Implementar sistema de logging base
    - Criar Logger com níveis DEBUG, INFO, WARNING, ERROR
    - Configurar armazenamento em arquivo local
    - Implementar rotação automática de logs
    - _Requirements: 10_

  - [x] 1.3 Configurar state management com Provider
    - Criar providers base para aplicação
    - Implementar AppState provider
    - Configurar error handling provider
    - _Requirements: 5_

  - [x] 1.4 Configurar banco de dados SQLite com Drift
    - Definir tabelas: Animes, Watchlist, History, SyncQueue
    - Gerar código Drift
    - Implementar migrations
    - _Requirements: 4, 9_

- [ ] 2. Implementar Design System Glassmorphism
  - [x] 2.1 Criar componentes base com glassmorphism
    - Implementar GlassmorphicContainer com BackdropFilter
    - Implementar GlassmorphicButton
    - Implementar GlassmorphicCard
    - _Requirements: 1_

  - [x] 2.2 Implementar sistema de temas claro/escuro
    - Criar ThemeProvider com paleta de cores
    - Implementar suporte a tema claro e escuro
    - Aplicar glassmorphism em ambos os temas
    - _Requirements: 1_

  - [x] 2.3 Implementar layouts responsivos
    - Criar breakpoints para mobile (<600dp) e tablet (≥600dp)
    - Implementar suporte a orientação portrait e landscape
    - Testar em diferentes tamanhos de tela
    - _Requirements: 1, 7, 8_

  - [x] 2.4 Escrever testes unitários para componentes glassmorphism
    - Testar renderização de GlassmorphicContainer
    - Testar aplicação de blur e opacity
    - Testar responsividade de layouts
    - _Requirements: 1_

- [ ] 3. Implementar Bridge Go-Flutter
  - [x] 3.1 Criar interface de método channel
    - Definir métodos: GetAnimes, AddToWatchlist, MarkAsWatched, GetSyncStatus, SyncWithDesktop
    - Implementar serialização JSON
    - Configurar timeout de 500ms
    - _Requirements: 3_

  - [x] 3.2 Implementar BridgeService em Dart
    - Criar classe BridgeService com method channel
    - Implementar chamadas assíncronas para Go
    - Implementar tratamento de erros com códigos estruturados
    - _Requirements: 3, 5_

  - [x] 3.3 Implementar funções Go expostas via GoMobile
    - Implementar GetAnimes() em Go
    - Implementar AddToWatchlist() em Go
    - Implementar MarkAsWatched() em Go
    - Implementar GetSyncStatus() em Go
    - Implementar SyncWithDesktop() em Go
    - _Requirements: 3_

  - [x] 3.4 Implementar serialização/desserialização JSON
    - Criar estruturas Go para dados
    - Implementar marshaling/unmarshaling
    - Testar com tipos: string, int, float, bool, array, object
    - _Requirements: 3_

  - [x] 3.5 Escrever testes unitários para bridge
    - Testar serialização JSON
    - Testar desserialização JSON
    - Testar tratamento de erros
    - _Requirements: 3_

- [ ] 4. Implementar Servidor Web Local
  - [x] 4.1 Implementar inicialização do servidor web em Go
    - Criar função StartWebServer(port int)
    - Configurar rotas HTTP básicas
    - Implementar graceful startup em <2 segundos
    - _Requirements: 2_

  - [x] 4.2 Implementar gerenciamento de ciclo de vida do servidor
    - Implementar PauseWebServer()
    - Implementar ResumeWebServer()
    - Implementar ShutdownWebServer()
    - Manter estado durante pause/resume
    - _Requirements: 2, 9_

  - [x] 4.3 Integrar servidor web com bridge Go-Flutter
    - Expor função StartWebServer via GoMobile
    - Chamar StartWebServer na inicialização da app
    - Gerenciar ciclo de vida com app lifecycle
    - _Requirements: 2, 3, 9_

  - [x] 4.4 Escrever testes de integração para servidor web
    - Testar inicialização em <2 segundos
    - Testar acessibilidade via localhost
    - Testar pause/resume sem perda de estado
    - _Requirements: 2_

- [ ] 5. Implementar Lógica de Negócio em Go
  - [~] 5.1 Implementar AnimeService em Go
    - Implementar GetAnimes() com busca de dados
    - Implementar AddToWatchlist(animeID)
    - Implementar MarkAsWatched(animeID, episodeID)
    - Implementar GetWatchlist()
    - _Requirements: 3, 4_

  - [~] 5.2 Implementar VIPService em Go
    - Implementar GetVIPStatus()
    - Implementar UpdateVIPStatus()
    - Implementar sincronização de status VIP
    - _Requirements: 4_

  - [~] 5.3 Implementar HistoryService em Go
    - Implementar GetHistory()
    - Implementar AddHistoryEntry(animeID, episodeID)
    - Implementar sincronização de histórico
    - _Requirements: 4_

  - [~] 5.4 Escrever testes unitários para serviços Go
    - Testar GetAnimes()
    - Testar AddToWatchlist()
    - Testar MarkAsWatched()
    - _Requirements: 3, 4_

- [ ] 6. Implementar Sincronização com Desktop
  - [~] 6.1 Implementar SyncEngine em Dart
    - Criar classe SyncEngine com gerenciamento de fila
    - Implementar enfileiramento de mudanças locais
    - Implementar processamento de fila ao reconectar
    - _Requirements: 4_

  - [~] 6.2 Implementar SyncService em Go
    - Implementar SyncWithDesktop() com envio de mudanças
    - Implementar recebimento de mudanças do desktop
    - Implementar resolução de conflitos por timestamp
    - _Requirements: 4_

  - [~] 6.3 Implementar offline queue management
    - Armazenar mudanças em SyncQueue quando offline
    - Processar fila ao reconectar
    - Marcar itens como sincronizados
    - _Requirements: 4_

  - [~] 6.4 Implementar detecção de conectividade
    - Usar connectivity_plus para monitorar rede
    - Disparar sincronização ao reconectar
    - Notificar usuário de status de conexão
    - _Requirements: 4_

  - [~] 6.5 Escrever testes de integração para sincronização
    - Testar sincronização em <5 segundos
    - Testar offline queue com dados de exemplo
    - Testar resolução de conflitos por timestamp
    - Testar transições offline/online
    - _Requirements: 4, 6_

- [ ] 7. Implementar Tratamento de Erros Robusto
  - [~] 7.1 Criar ErrorHandler centralizado
    - Implementar captura de exceções globais
    - Mapear erros para mensagens amigáveis
    - Registrar em log com contexto completo
    - _Requirements: 5_

  - [~] 7.2 Implementar categorias de erro
    - Definir NetworkError, BridgeError, SyncError, ResourceError
    - Implementar códigos de erro estruturados
    - Criar mapeamento de erro para mensagem
    - _Requirements: 5_

  - [~] 7.3 Implementar UI de erro
    - Criar ErrorSnackBar para notificações rápidas
    - Criar ErrorDialog para erros críticos
    - Criar ErrorScreen com opção de retry
    - _Requirements: 5_

  - [~] 7.4 Implementar retry automático com backoff exponencial
    - Implementar backoff exponencial (1s, 2s, 4s, 8s)
    - Limitar tentativas a 5
    - Notificar usuário após falha final
    - _Requirements: 5_

  - [~] 7.5 Escrever testes unitários para tratamento de erros
    - Testar captura de exceções
    - Testar mapeamento de erros
    - Testar retry automático
    - _Requirements: 5_

- [ ] 8. Implementar Screens e Navegação
  - [~] 8.1 Implementar HomeScreen
    - Exibir lista de animes com lazy loading
    - Implementar busca e filtros
    - Integrar com AnimeProvider
    - _Requirements: 1, 8_

  - [~] 8.2 Implementar WatchlistScreen
    - Exibir animes salvos
    - Permitir remover da watchlist
    - Sincronizar com desktop
    - _Requirements: 1, 4_

  - [~] 8.3 Implementar HistoryScreen
    - Exibir episódios assistidos
    - Mostrar progresso de cada anime
    - Sincronizar com desktop
    - _Requirements: 1, 4_

  - [~] 8.4 Implementar SettingsScreen
    - Configurar tema claro/escuro
    - Configurar porta do servidor web
    - Visualizar logs
    - _Requirements: 1, 10_

  - [~] 8.5 Implementar navegação entre screens
    - Configurar rotas com GoRouter
    - Implementar transições suaves
    - Manter estado de navegação
    - _Requirements: 1, 8_

  - [~] 8.6 Escrever testes de integração para screens
    - Testar renderização de HomeScreen
    - Testar navegação entre screens
    - Testar transições em <500ms
    - _Requirements: 1, 8_

- [ ] 9. Implementar Ciclo de Vida da Aplicação
  - [~] 9.1 Implementar AppLifecycleObserver
    - Criar listener para eventos de ciclo de vida
    - Implementar handlers para resumed, paused, detached
    - Integrar com bridge e sincronização
    - _Requirements: 9_

  - [~] 9.2 Implementar inicialização da aplicação
    - Inicializar bridge Go-Flutter
    - Iniciar servidor web
    - Inicializar banco de dados
    - Carregar dados em cache
    - Inicializar sync engine
    - _Requirements: 9_

  - [~] 9.3 Implementar pause/resume
    - Pausar sincronização
    - Pausar servidor web
    - Salvar estado
    - Restaurar estado ao resumir
    - _Requirements: 9_

  - [~] 9.4 Implementar shutdown graceful
    - Pausar sincronização
    - Fechar servidor web
    - Fechar conexões
    - Salvar estado para recuperação
    - _Requirements: 9_

  - [~] 9.5 Escrever testes de integração para ciclo de vida
    - Testar inicialização completa
    - Testar pause/resume sem perda de estado
    - Testar shutdown graceful
    - _Requirements: 9_

- [ ] 10. Implementar Logging e Debugging
  - [~] 10.1 Expandir sistema de logging
    - Implementar níveis DEBUG, INFO, WARNING, ERROR
    - Adicionar timestamps e thread IDs
    - Implementar contexto de execução
    - _Requirements: 10_

  - [~] 10.2 Implementar logging de operações
    - Registrar entrada/saída de funções
    - Registrar parâmetros relevantes
    - Registrar tempo de execução
    - _Requirements: 10_

  - [~] 10.3 Implementar logging de erros
    - Registrar stack trace completo
    - Registrar contexto de execução
    - Registrar estado da aplicação
    - _Requirements: 10_

  - [~] 10.4 Implementar ferramentas de debugging
    - Criar inspetor de estado da bridge
    - Criar visualizador de fila de sincronização
    - Criar monitor de performance
    - Criar visualizador de logs em tempo real
    - _Requirements: 10_

  - [~] 10.5 Implementar visualizador de logs na UI
    - Criar DebugScreen com logs em tempo real
    - Permitir filtrar por nível
    - Permitir exportar logs
    - _Requirements: 10_

  - [~] 10.6 Escrever testes unitários para logging
    - Testar formatação de logs
    - Testar rotação de arquivos
    - Testar níveis de log
    - _Requirements: 10_

- [ ] 11. Implementar Testes de Integração Completos
  - [~] 11.1 Escrever testes de integração Go-Flutter
    - Testar comunicação bridge com sucesso
    - Testar timeout de 500ms
    - Testar múltiplas chamadas simultâneas
    - _Requirements: 3, 6_

  - [~] 11.2 Escrever testes de integração de sincronização
    - Testar sincronização desktop-mobile
    - Testar resolução de conflitos
    - Testar offline/online transitions
    - _Requirements: 4, 6_

  - [~] 11.3 Escrever testes de integração de ciclo de vida
    - Testar inicialização completa
    - Testar pause/resume
    - Testar shutdown graceful
    - _Requirements: 9_

  - [~] 11.4 Escrever testes de integração de tratamento de erros
    - Testar cenários de falha de rede
    - Testar cenários de falha de bridge
    - Testar cenários de falha de sincronização
    - _Requirements: 5, 6_

  - [~] 11.5 Executar suite de testes de integração
    - Executar todos os testes
    - Validar 100% de sucesso
    - Completar em <60 segundos
    - _Requirements: 6_

- [ ] 12. Otimizar Performance para Mobile
  - [~] 12.1 Implementar lazy loading para listas
    - Usar ListView.builder para renderização eficiente
    - Carregar 20 itens inicialmente
    - Carregar mais ao scroll
    - _Requirements: 8_

  - [~] 12.2 Implementar cache de imagens
    - Usar CachedNetworkImage
    - Implementar cache local
    - Limpar cache periodicamente
    - _Requirements: 8_

  - [~] 12.3 Otimizar renderização de UI
    - Usar const widgets quando possível
    - Evitar rebuilds desnecessários
    - Usar RepaintBoundary para otimizar
    - _Requirements: 8_

  - [~] 12.4 Otimizar uso de memória
    - Monitorar uso de memória em background
    - Implementar limpeza de cache
    - Limitar tamanho de logs
    - _Requirements: 8_

  - [~] 12.5 Escrever testes de performance
    - Testar renderização em <1 segundo
    - Testar transições em <500ms
    - Testar carregamento de lista em <2 segundos
    - Testar uso de memória em background <50MB
    - Testar frame rate em foreground 60 FPS
    - _Requirements: 8_

- [ ] 13. Testar Compatibilidade Android 8.0+
  - [~] 13.1 Configurar testes em múltiplas versões Android
    - Configurar emuladores para Android 8.0, 10, 12, 14
    - Configurar CI/CD para testes em múltiplas versões
    - _Requirements: 7_

  - [~] 13.2 Testar permissões Android Runtime
    - Testar solicitação de permissões
    - Testar comportamento com permissões negadas
    - Testar em Android 8.0+
    - _Requirements: 7_

  - [~] 13.3 Testar suporte a diferentes tamanhos de tela
    - Testar em phones pequenos (4.5")
    - Testar em phones grandes (6.5")
    - Testar em tablets (7", 10")
    - _Requirements: 7_

  - [~] 13.4 Executar testes em múltiplas versões
    - Executar suite completa em Android 8.0
    - Executar suite completa em Android 10
    - Executar suite completa em Android 12
    - Executar suite completa em Android 14
    - Validar sem regressões
    - _Requirements: 7_

- [ ] 14. Checkpoint - Validar Implementação Completa
  - Executar todos os testes (unitários, integração, performance)
  - Validar 100% de sucesso
  - Verificar cobertura de todos os requisitos
  - Validar performance em múltiplas versões Android
  - Perguntar ao usuário se há dúvidas ou ajustes necessários
  - _Requirements: 1-10_

- [ ] 15. Build e Release
  - [~] 15.1 Preparar build release
    - Configurar versão do app
    - Gerar keystore para assinatura
    - Configurar build.gradle para release
    - _Requirements: 7_

  - [~] 15.2 Gerar APK/AAB release
    - Executar flutter build apk --release
    - Executar flutter build appbundle --release
    - Validar tamanho do build
    - _Requirements: 7_

  - [~] 15.3 Testar build release
    - Instalar APK em dispositivo real
    - Executar smoke tests
    - Validar funcionalidade completa
    - _Requirements: 7_

  - [~] 15.4 Documentar release
    - Criar changelog
    - Documentar correções e melhorias
    - Preparar notas de release
    - _Requirements: 1-10_

## Notes

- Tasks marcadas com `*` são opcionais e podem ser puladas para MVP mais rápido
- Cada task referencia requisitos específicos para rastreabilidade
- Checkpoints garantem validação incremental
- Testes unitários validam exemplos específicos e casos extremos
- Testes de integração validam fluxos completos
- Performance é validada em múltiplas versões Android
- Logging e debugging são integrados em toda a implementação
