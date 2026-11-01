# Aura Ciclos

Rastreador de ciclo menstrual em [Amarra](https://github.com/puppe1990/amarra-cais): previsão da menstruação e da janela fértil, registro de 70+ sintomas, três modos (ciclo, tentantes, gravidez), conta anônima e a Maya (assistente por regras).

Live: https://aura-ciclos.apps.gestaobem.com

Os registros ficam em SQLite no servidor. Não há criptografia ponta a ponta.

## O que o app faz

- **Ciclo** (`/ciclo`) — anel do ciclo, dia atual, próxima menstruação, precisão da previsão e registro de fluxo.
- **Registro** (`/registro`) — catálogo de sintomas, humor, energia, sono, água, exercício; **Vida sexual** (`/registro/vida-sexual`) com preservativo, orgasmo e vezes. No modo tentantes, BBT e testes.
- **Relatórios** (`/relatorios`) — regularidade, duração das fases, sintomas mais registrados e troca de modo.
- **Apoio** (`/apoio`) — Maya por palavras-chave, artigos revisados e lembretes (pílula, água, menstruação).
- **Anônimo** (`/anonimo`) — conta sem nome nem e-mail; a sessão fica neste aparelho.

A previsão usa a média dos últimos até 6 intervalos, fase lútea de 14 dias, janela fértil de ovulação−5 até +1. Ciclo irregular quando o desvio-padrão da amostra ≥ 4 dias. Teste de ovulação positivo vale mais que BBT, que vale mais que o calendário.

## Stack

- Go 1.26 (net/http stdlib) + Amarra Views + Drive
- HTML (`web/templates/`) + Tailwind 3 (paleta Stitch: alabaster, coral `#a83223`, lavanda)
- SQLite (`modernc.org/sqlite`, sem CGO)
- Testes com `gofakeit` em `internal/testdata` (`Seed = 20261101`)

## Comandos

```bash
amarra-cais install   # npm install + go mod tidy
amarra-cais dev       # http://localhost:8080
amarra-cais test      # go test ./...
amarra-cais css       # Tailwind → web/static/css/styles.css
amarra-cais build     # bin/server
amarra-cais doctor    # amarra.js, PWA, paleta, fontes
make ci               # test + lint + format-check
```

Demo em desenvolvimento: `demo@example.com` / `password` (Sofia, 7 ciclos de 28 dias).

## Ambiente

| Variável          | Padrão                  | Uso                                                |
| ----------------- | ----------------------- | -------------------------------------------------- |
| `PORT`            | `:8080`                 | Porta HTTP                                         |
| `DB_PATH`         | `./data/app.db`         | SQLite                                             |
| `ENV`             | `development`           | `production` exige `ADMIN_TOKEN`                   |
| `APP_URL`         | `http://localhost:8080` | URL canônica                                       |
| `LOCALE`          | `pt`                    | Catálogo i18n (`en` / `pt`)                        |
| `STATIC_DIR`      |                         | Em produção: `/opt/aura-ciclos/current/web/static` |
| `TRUSTED_PROXIES` |                         | `127.0.0.1` atrás do Caddy                         |

Health: `GET /health` → `{"status":"ok"}`.

## Amarra

O scaffold veio do CLI. Marcadores que não se remove: `registerRoutes`, `Close() error`, `<!-- cais:nav -->`, `// cais:live-views`.

Drive morphs `#amarra-main`. Formulários de sessão (login, signup, logout, anônimo) levam `data-amarra-skip`. Depois do login o destino é `/ciclo`.

Ver [AGENTS.md](AGENTS.md) para TDD, views, flash/CSRF e geradores.
