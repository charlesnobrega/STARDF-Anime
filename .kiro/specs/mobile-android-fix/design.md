# Design Document: Mobile Android Fix

## Overview

O StarDF-Anime Mobile Android Fix é uma solução completa para corrigir e otimizar a versão mobile em Flutter do agregador de anime StarDF-Anime. Este design aborda seis problemas críticos: interface sem glassmorphism, falhas na inicialização do servidor web local, bridge Go-Flutter quebrada, falta de sincronização com desktop, tratamento de erros inadequado e ausência de testes de integração.

A solução implementa uma arquitetura em camadas com separação clara entre UI (Flutter), lógica de negócio (Go), sincronização e tratamento de erros, garantindo paridade visual com o desktop, comunicação robusta entre plataformas e sincronização bidirecional de dados.

## Architecture

### Visão Geral da Arquitetura

```
┌─────────────────────────────────────────────────────────────────┐
│                     Flutter Mobile App (Android)                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │              UI Layer (Flutter Widgets)                   │   │
│  │  - Glassmorphism Components                              │   │
│  │  - Theme Management (Light/Dark)                         │   │
│  │  - Responsive Layout                                     │   │
│  └──────────────────────────────────────────────────────────┘   │
│                            ▲                                      │
│                            │                                      │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │         Business Logic Layer (Dart/Flutter)              │   │
│  │  - State Management (Provider/Riverpod)                  │   │
│  │  - Error Handling                                        │   │
│  │  - Logging                                               │   │
│  │  - Offline Queue Management                              │   │
│  └──────────────────────────────────────────────────────────┘   │
│                            ▲                                      │
│                            │                                      │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │         Bridge Layer (GoMobile Interface)                │   │
│  │  - Method Channel Communication                          │   │
│  │  - JSON Serialization/Deserialization                    │   │
│  │  - Error Mapping                                         │   │
│  │  - Concurrent Call Management                            │   │
│  └──────────────────────────────────────────────────────────┘   │
│                            ▲                                      │
│                            │ (Platform Channel)                   │
└────────────────────────────┼──────────────────────────────────────┘
                             │
                             │ (JNI/FFI)
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Go Backend (GoMobile)                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │         Web Server Layer (HTTP)                          │   │
│  │  - Local HTTP Server (localhost:8080)                    │   │
│  │  - Graceful Startup/Shutdown                             │   │
│  │  - Lifecycle Management                                  │   │
│  └──────────────────────────────────────────────────────────┘   │
│                            ▲                                      │
│                            │                                      │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │         Business Logic Layer (Go)                        │   │
│  │  - Anime Data Processing                                 │   │
│  │  - Watchlist Management                                  │   │
│  │  - History Tracking                                      │   │
│  │  - VIP Status Management                                 │   │
│  └──────────────────────────────────────────────────────────┘   │
│                            ▲                                      │
│                            │                                      │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │         Data Layer (SQLite + Drift)                      │   │
│  │  - Local Database                                        │   │
│  │  - Sync Queue Storage                                    │   │
│  │  - Offline Data Cache                                    │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
                             ▲
                             │ (HTTP/REST)
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│              Desktop Backend (Windows/Linux)                     │
│  - Sync Server                                                   │
│  - Data Persistence                                              │
│  - User State Management                                         │
└─────────────────────────────────────────────────────────────────┘
```

### Fluxo de Dados

1. **UI → Business Logic**: Usuário interage com UI Flutter
2. **Business Logic → Bridge**: Lógica prepara chamada para Go
3. **Bridge → Go Backend**: Chamada serializada em JSON via GoMobile
4. **Go Backend → Database**: Processamento e persistência
5. **Go Backend → Sync Engine**: Mudanças enfileiradas para sincronização
6. **Sync Engine → Desktop**: Sincronização bidirecional com desktop
7. **Desktop → Mobile**: Mudanças do desktop sincronizadas de volta

## Components and Interfaces

### 1. UI Layer (Flutter)

#### 1.1 Glassmorphism Design System

**Componentes Principais:**
- `GlassmorphicContainer`: Widget base com efeito de vidro fosco
- `GlassmorphicButton`: Botão com glassmorphism
- `GlassmorphicCard`: Card com transparência e blur
- `ThemeProvider`: Gerenciador de temas claro/escuro

**Implementação:**
```dart
class GlassmorphicContainer extends StatelessWidget {
  final Widget child;
  final double blur;
  final double opacity;
  final Color color;
  
  @override
  Widget build(BuildContext context) {
    return BackdropFilter(
      filter: ImageFilter.blur(sigmaX: blur, sigmaY: blur),
      child: Container(
        decoration: BoxDecoration(
          color: color.withOpacity(opacity),
          borderRadius: BorderRadius.circular(20),
          border: Border.all(
            color: Colors.white.withOpacity(0.2),
          ),
        ),
        child: child,
      ),
    );
  }
}
```

**Responsividade:**
- Layouts adaptativos para phones e tablets
- Breakpoints: mobile (<600dp), tablet (≥600dp)
- Orientação: portrait e landscape

#### 1.2 Screen Components

- **HomeScreen**: Lista de animes com lazy loading
- **WatchlistScreen**: Animes salvos com sincronização
- **HistoryScreen**: Episódios assistidos
- **SettingsScreen**: Configurações e temas
- **ErrorScreen**: Exibição de erros com retry

### 2. Business Logic Layer (Dart/Flutter)

#### 2.1 State Management

**Provider/Riverpod Architecture:**
```dart
// Anime Provider
final animeProvider = FutureProvider<List<Anime>>((ref) async {
  final bridge = ref.watch(bridgeProvider);
  return bridge.getAnimes();
});

// Sync State Provider
final syncStateProvider = StateNotifierProvider<SyncNotifier, SyncState>((ref) {
  return SyncNotifier(ref);
});

// Error Handler Provider
final errorHandlerProvider = Provider<ErrorHandler>((ref) {
  return ErrorHandler(ref);
});
```

#### 2.2 Error Handling

**ErrorHandler Class:**
- Captura exceções em toda a aplicação
- Mapeia erros para mensagens amigáveis
- Registra em log com contexto
- Oferece opções de retry

#### 2.3 Logging System

**Logger Configuration:**
- Níveis: DEBUG, INFO, WARNING, ERROR
- Armazenamento em arquivo local
- Rotação automática de logs
- Timestamps e thread IDs

#### 2.4 Offline Queue Management

**SyncQueue:**
- Armazena mudanças localmente quando offline
- Processa fila ao reconectar
- Resolve conflitos por timestamp

### 3. Bridge Layer (GoMobile)

#### 3.1 Method Channel Interface

**Métodos Expostos:**

```go
// GetAnimes retorna lista de animes
func GetAnimes(ctx context.Context) ([]byte, error)

// AddToWatchlist adiciona anime à watchlist
func AddToWatchlist(ctx context.Context, animeID string) ([]byte, error)

// MarkAsWatched marca episódio como assistido
func MarkAsWatched(ctx context.Context, animeID, episodeID string) ([]byte, error)

// GetSyncStatus retorna status de sincronização
func GetSyncStatus(ctx context.Context) ([]byte, error)

// SyncWithDesktop sincroniza com desktop
func SyncWithDesktop(ctx context.Context) ([]byte, error)
```

#### 3.2 Serialization

**JSON Schema:**
```json
{
  "status": "success|error",
  "data": {},
  "error": {
    "code": "ERROR_CODE",
    "message": "Descrição do erro",
    "details": {}
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 3.3 Error Handling

**Error Codes:**
- `BRIDGE_TIMEOUT`: Chamada excedeu 500ms
- `BRIDGE_SERIALIZATION_ERROR`: Erro ao serializar/desserializar
- `BRIDGE_EXECUTION_ERROR`: Erro na execução da função Go
- `BRIDGE_CONCURRENT_LIMIT`: Limite de chamadas simultâneas

### 4. Go Backend Layer

#### 4.1 Web Server

**Inicialização:**
```go
func StartWebServer(port int) error {
  server := &http.Server{
    Addr:         fmt.Sprintf(":%d", port),
    Handler:      setupRoutes(),
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 5 * time.Second,
  }
  
  go func() {
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
      log.Printf("Server error: %v", err)
    }
  }()
  
  return nil
}
```

**Lifecycle:**
- Startup: Inicializa em <2 segundos
- Pause: Pausa gracefully sem perder estado
- Resume: Reinicia e reconecta
- Shutdown: Fecha conexões gracefully

#### 4.2 Business Logic

**Anime Service:**
- Busca e processamento de dados
- Gerenciamento de watchlist
- Rastreamento de histórico
- Gerenciamento de status VIP

**Sync Service:**
- Sincronização bidirecional com desktop
- Resolução de conflitos por timestamp
- Fila de mudanças offline
- Processamento de fila ao reconectar

#### 4.3 Data Layer

**SQLite Database (Drift ORM):**
```dart
@DataClassName('AnimeData')
class Animes extends Table {
  IntColumn get id => integer().autoIncrement()();
  TextColumn get title => text()();
  TextColumn get description => text()();
  IntColumn get episodes => integer()();
  DateTimeColumn get addedAt => dateTime()();
  DateTimeColumn get syncedAt => dateTime().nullable()();
}

@DataClassName('WatchlistItem')
class Watchlist extends Table {
  IntColumn get id => integer().autoIncrement()();
  IntColumn get animeId => integer()();
  DateTimeColumn get addedAt => dateTime()();
  DateTimeColumn get syncedAt => dateTime().nullable()();
}

@DataClassName('SyncQueueItem')
class SyncQueue extends Table {
  IntColumn get id => integer().autoIncrement()();
  TextColumn get operation => text()(); // 'add', 'remove', 'update'
  TextColumn get entityType => text()(); // 'watchlist', 'history', 'vip'
  TextColumn get data => text()(); // JSON
  DateTimeColumn get createdAt => dateTime()();
  BoolColumn get synced => boolean().withDefault(const Constant(false))();
}
```

### 5. Sync Engine

#### 5.1 Sincronização Bidirecional

**Fluxo:**
1. Mudança local → Enfileirada em SyncQueue
2. Conexão com desktop → Envia mudanças pendentes
3. Desktop processa → Retorna confirmação
4. Mobile marca como sincronizado
5. Desktop envia mudanças → Mobile processa

#### 5.2 Resolução de Conflitos

**Estratégia:**
- Timestamp mais recente vence
- Logs de conflito para debugging
- Notificação ao usuário se necessário

#### 5.3 Offline Support

**Comportamento:**
- Mudanças armazenadas localmente
- Fila processada ao reconectar
- Sincronização automática em background

## Data Models

### Anime
```dart
class Anime {
  final String id;
  final String title;
  final String description;
  final int episodes;
  final String imageUrl;
  final DateTime addedAt;
  final DateTime? syncedAt;
}
```

### WatchlistItem
```dart
class WatchlistItem {
  final String id;
  final String animeId;
  final DateTime addedAt;
  final DateTime? syncedAt;
}
```

### HistoryEntry
```dart
class HistoryEntry {
  final String id;
  final String animeId;
  final int episodeNumber;
  final DateTime watchedAt;
  final DateTime? syncedAt;
}
```

### SyncQueueItem
```dart
class SyncQueueItem {
  final String id;
  final String operation; // 'add', 'remove', 'update'
  final String entityType; // 'watchlist', 'history', 'vip'
  final Map<String, dynamic> data;
  final DateTime createdAt;
  final bool synced;
}
```

### SyncConflict
```dart
class SyncConflict {
  final String entityId;
  final DateTime mobileTimestamp;
  final DateTime desktopTimestamp;
  final Map<String, dynamic> mobileData;
  final Map<String, dynamic> desktopData;
  final Map<String, dynamic> resolvedData;
}
```

## Error Handling

### Error Categories

1. **Network Errors**
   - Sem conexão
   - Timeout
   - Erro de servidor

2. **Bridge Errors**
   - Serialização falhou
   - Timeout na chamada
   - Tipo de dado não suportado

3. **Sync Errors**
   - Conflito de dados
   - Falha ao enviar mudanças
   - Falha ao receber mudanças

4. **Resource Errors**
   - Recurso não disponível
   - Permissão negada
   - Espaço em disco insuficiente

### Error Recovery

**Estratégia:**
- Retry automático com backoff exponencial
- Notificação ao usuário
- Logging completo com stack trace
- Manutenção de estado funcional

### Error UI

**Componentes:**
- `ErrorSnackBar`: Notificação rápida
- `ErrorDialog`: Erro crítico com ação
- `ErrorScreen`: Tela de erro com retry
- `ErrorLog`: Visualizador de logs

## Testing Strategy

### Unit Tests

**Cobertura:**
- Componentes Flutter (widgets)
- Lógica de negócio (providers, services)
- Serialização/desserialização
- Tratamento de erros
- Logging

**Exemplo:**
```dart
test('GlassmorphicContainer renders with blur effect', () {
  final widget = GlassmorphicContainer(
    blur: 10,
    opacity: 0.5,
    child: Text('Test'),
  );
  
  expect(find.byType(BackdropFilter), findsOneWidget);
});
```

### Integration Tests

**Cobertura:**
- Comunicação Go-Flutter
- Sincronização desktop-mobile
- Ciclo de vida da aplicação
- Transições offline/online
- Operações simultâneas

**Exemplo:**
```dart
testWidgets('Bridge call completes within 500ms', (WidgetTester tester) async {
  final bridge = BridgeService();
  final stopwatch = Stopwatch()..start();
  
  final result = await bridge.getAnimes();
  
  stopwatch.stop();
  expect(stopwatch.elapsedMilliseconds, lessThan(500));
  expect(result, isNotEmpty);
});
```

### Performance Tests

**Métricas:**
- Tempo de renderização de tela (<1s)
- Tempo de transição entre telas (<500ms)
- Tempo de carregamento de lista (<2s)
- Uso de memória em background (<50MB)
- Frame rate em foreground (60 FPS)

### Test Configuration

**Mínimo de iterações:** 100 por teste
**Timeout:** 60 segundos para suite completa
**Plataformas:** Android 8.0, 10, 12, 14

## Lifecycle Management

### Application Startup

```
1. App Launch
   ↓
2. Initialize Bridge (GoMobile)
   ↓
3. Start Web Server (localhost:8080)
   ↓
4. Initialize Database (SQLite)
   ↓
5. Load Cached Data
   ↓
6. Initialize Sync Engine
   ↓
7. Render UI
   ↓
8. Start Background Sync
```

### Pause/Resume

**Pause:**
- Pausa sincronização
- Pausa servidor web
- Salva estado

**Resume:**
- Restaura estado
- Reinicia servidor web
- Reinicia sincronização
- Reconecta com backend

### Shutdown

- Pausa sincronização
- Fecha servidor web
- Fecha conexões
- Salva estado para recuperação

## Logging and Debugging

### Log Levels

- **DEBUG**: Informações detalhadas para debugging
- **INFO**: Eventos importantes
- **WARNING**: Situações anormais
- **ERROR**: Erros que afetam funcionalidade

### Log Format

```
[2024-01-01 12:00:00.123] [INFO] [Thread-1] [BridgeService] Method call: getAnimes() -> 150ms
[2024-01-01 12:00:01.456] [ERROR] [Thread-2] [SyncEngine] Sync failed: Connection timeout
  Stack trace: ...
```

### Log Storage

- Arquivo local: `/data/data/com.stardf.anime/logs/`
- Rotação automática: 10MB por arquivo
- Retenção: 7 dias

### Debugging Tools

- Inspetor de estado da bridge
- Visualizador de fila de sincronização
- Monitor de performance
- Visualizador de logs em tempo real

## Implementation Roadmap

### Phase 1: Foundation (Week 1-2)
- [ ] Setup GoMobile bridge
- [ ] Implementar glassmorphism design system
- [ ] Configurar state management (Provider)
- [ ] Implementar error handling base

### Phase 2: Core Features (Week 3-4)
- [ ] Implementar web server local
- [ ] Implementar comunicação Go-Flutter
- [ ] Implementar sincronização básica
- [ ] Implementar offline queue

### Phase 3: Polish (Week 5-6)
- [ ] Otimizar performance
- [ ] Implementar logging completo
- [ ] Adicionar testes de integração
- [ ] Testar em múltiplas versões Android

### Phase 4: Release (Week 7)
- [ ] Testes finais
- [ ] Build release
- [ ] Deploy

## Success Criteria

1. ✓ Interface mobile com glassmorphism idêntica ao desktop
2. ✓ Servidor web local inicia em <2 segundos
3. ✓ Bridge Go-Flutter funciona com <500ms latência
4. ✓ Sincronização desktop-mobile em <5 segundos
5. ✓ Tratamento de erros completo com logging
6. ✓ Testes de integração com 100% de sucesso
7. ✓ Compatibilidade com Android 8.0+
8. ✓ Performance otimizada (60 FPS, <50MB background)
9. ✓ Ciclo de vida gerenciado corretamente
10. ✓ Logging e debugging completos
