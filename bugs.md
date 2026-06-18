# Sentinel — Análise de Bugs e Melhorias

Diagnóstico do projeto (backend Go, frontend React, infra Docker/k8s/CI),
organizado por severidade, com o motivo e a direção de correção de cada item.

---

## Bugs / Correção (prioridade máxima)

### 1. JWT_SECRET lido em momento errado — bug real
`internal/services/auth_service.go:15`:

```go
var secret = []byte(os.Getenv("JWT_SECRET"))
```

Essa var de pacote é inicializada **antes** do `main()` rodar `config.LoadDotEnv()`.
Já o middleware (`internal/http/middleware/auth.go:25`) lê `os.Getenv("JWT_SECRET")`
em **tempo de request**. Resultado: se o segredo vier só do `.env`, o token é
**assinado** com segredo vazio e **validado** com o segredo certo → todo token fica
inválido. No docker-compose passa porque `JWT_SECRET` é env var real, mas é uma
bomba-relógio.

**Correção:** carregar o segredo via config e passar por injeção de dependência
tanto para o serviço quanto para o middleware. Nunca ler `os.Getenv` espalhado.

### 2. Versão do Go incompatível entre go.mod e CI
`go.mod` exige `go 1.25.1`, mas `.github/workflows/ci.yml:21` usa
`go-version: "1.22"`. O build do CI quebra (ou comporta diferente).

**Correção:** alinhar a versão (subir o CI, ou abaixar o `go.mod` para uma versão real).

### 3. Possível panic em `ListContainers`
`internal/docker/containers.go:64` faz `c.Names[0]` sem checar tamanho.
`ListAllContainers` (linha 187) já protege com `len > 0`, mas essa não.

**Correção:** mesma guarda de `len(c.Names) > 0`.

### 4. Config de CPU/MEM nunca é usada
`ALERT_CPU_THRESHOLD` e `ALERT_MEM_THRESHOLD_MB` existem no config e no `.env`,
mas o `checker` só alerta sobre queda/restart — nunca lê stats nem compara
threshold. Config morta.

**Correção:** ou implementar o alerta por CPU/memória, ou remover a config para
não enganar.

---

## Segurança

### 5. CORS liberado para todo mundo
`router.go:13` usa `cors.Default()` (qualquer origem). Para uma ferramenta que
controla containers do host, é arriscado.

**Correção:** restringir à origem do frontend.

### 6. `/register` público
Qualquer um cria conta e ganha acesso ao painel de containers.

**Correção:** proteger, fazer seed de usuário, ou exigir convite.

### 7. Sem validação de input
`models/user.go` não tem tags `binding:"required,email"`; senha sem tamanho mínimo.

**Correção:** adicionar validação no DTO.

### 8. Sem rate limit no `/login`
Brute force livre.

**Correção:** middleware de rate limiting.

### 9. Vazamento de erro interno para o cliente
`utils/response_error.go:18` devolve `err.Error()` cru no 500 (inclui erros de
banco encapsulados).

**Correção:** logar internamente, responder mensagem genérica.

### 10. Segredo fraco como default
`docker-compose.yml:14` usa `JWT_SECRET:-secret`.

**Correção:** remover default, falhar se vazio.

---

## Arquitetura / Qualidade

### 11. Camadas inconsistentes
`auth_service` usa struct com DI; `container_service` são funções soltas chamando
`docker` direto; e o `StreamLogs` (controller) chama `docker` pulando o service.

**Correção:** padronizar tudo numa abordagem (DI com struct é o ideal).

### 12. `ListContainers() (interface{}, error)`
`container_service.go:5` retorna `interface{}`. Perde tipagem.

**Correção:** retornar `[]docker.ContainerInfo`.

### 13. `context.Background()` em todo lugar
Nenhuma chamada Docker propaga o contexto do request; não dá para cancelar/timeout.

**Correção:** passar `c.Request.Context()`.

### 14. Sem timeouts
`http.Post` ao Telegram (`alerter.go:164`) usa client default sem timeout (pode
travar); o servidor Gin sobe sem `ReadTimeout`/`WriteTimeout` (slowloris).

**Correção:** `http.Client{Timeout}` + `http.Server` configurado.

### 15. Sem graceful shutdown do HTTP
`router.Run` bloqueia e ignora o `ctx`; só checker/alerter respeitam o cancel.

**Correção:** usar `http.Server` com `Shutdown(ctx)`.

### 16. Pool de conexão do banco não configurado
Falta `SetMaxOpenConns`/`SetMaxIdleConns`/`SetConnMaxLifetime`.

### 17. Stream de logs vaza headers do Docker
O controller manda o stream multiplexado cru e o frontend "limpa" com regex
(`api.js:102`).

**Correção:** demuxar no servidor com `stdcopy.StdCopy`.

### 18. Idioma misturado
Erros/logs em PT e EN trocados (ex: "body invalid" vs "credenciais inválidas").

**Correção:** escolher um.

### 19. Detalhes menores
Alias feio `handler2` em `router.go:6`; `UserResponse` DTO nunca usado.

---

## Testes & Infra

### 20. Zero testes
CI só roda `go build`. Ponto que mais salta numa review.

**Correção:** começar com testes de unidade no `services` (auth/login) e
`alerting` (`buildMessage`/`stateKey`, que são puros e fáceis).

### 21. Dockerfile
`golang:alpine`/`alpine:latest` sem pin (não reproduzível), roda como **root**,
sem `HEALTHCHECK`.

**Correção:** pinar versões, criar usuário não-root.

### 22. CI contraditório no frontend
`no-cache: true` junto com `cache-from/to: gha` (`ci.yml:60-71`).

---

## Ordem sugerida de ataque

1. Bugs (1–4)
2. Segurança (5–10)
3. Testes (20)
4. O resto

Os itens **1**, **2** e **20** são os que um tech lead nota em 5 minutos.
