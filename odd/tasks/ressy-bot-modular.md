# Ressy Bot modular — Discord, panel admin y música

Bot de Discord en Go con módulos instalables en caliente, panel administrativo web y primer módulo funcional de música desde YouTube. Todo el sistema se entrega en una sola imagen Docker.

## Specs

S1. "quiero crear un bot de discord completo modular hecho en go siguiendo las reglas de agents md, controlable desde un panel administrativo, cada nueva funcion sera un modulo instalable que puedes usar, quiero comenzar primero con un sistema para escuchar musica, a ver que me recomiendas como podemos comenzar"

S2. "Plugins .so en caliente"

S3. "Go web + OAuth de Discord (recomendado)"

S4. "todo este proyecto debe ir en un solo docker, entonces hay que ver la manera de que funcione bien, quiero primero traer musica desde youtube"

S5. "https://github.com/loonbac/Ressy-Bot.git usa este repo borra todo lo que estaba ahi, ya no lo necesito, vamos a empezar desde 0 no revises nada de lo que estaba ahi"

## Tasks

- T1 — S5 — inline — Crear commit raíz vacío, force push a `main` y borrar las demás ramas del repo `loonbac/Ressy-Bot`; sin leer el contenido previo. Commit: —
- T2 — S1,S2 — inline — Rellenar `AGENTS.md` §0 con el mapa Go; `go.mod`; `src/app/` con arranque; `src/modules/core/` con contrato de módulo y cargador de `.so` inyectable; pruebas del cargador. Commit: —
- T3 — S2 — inline — Registro de comandos del bot y activación/desactivación de módulos por guild, con pruebas. Commit: —
- T4 — S3 — inline — Módulo `admin`: servidor `net/http`, OAuth2 de Discord, cookie firmada, endpoint de carga de plugins y UI de administración; pruebas. Commit: —
- T5 — S1,S4 — worker `gentle-ai-worker` — Módulo `music`: rebanadas `play`, `skip`, `queue`, `stop`, `now-playing`, interfaz `AudioNode` con nodo falso y pruebas. Commit: —
- T6 — S4 — worker `gentle-ai-worker` — Adaptador Lavalink 4 con plugin de YouTube, `application.yml`, imagen Docker única con Go + Java + ffmpeg + toolchain de plugins y entrypoint; verificación de arranque. Commit: —
- T7 — S1-S5 — inline — Suite completa, differential checks, revisión nativa RDD y cierre del feature. Commit: —

## Log

L1. "quiero crear un bot de discord completo modular hecho en go siguiendo las reglas de agents md, controlable desde un panel administrativo, cada nueva funcion sera un modulo instalable que puedes usar, quiero comenzar primero con un sistema para escuchar musica, a ver que me recomiendas como podemos comenzar"

L2. Decisiones del usuario: "Plugins .so en caliente"; "Go web + OAuth de Discord (recomendado)"; "todo este proyecto debe ir en un solo docker, entonces hay que ver la manera de que funcione bien, quiero primero traer musica desde youtube"; "https://github.com/loonbac/Ressy-Bot.git usa este repo borra todo lo que estaba ahi, ya no lo necesito, vamos a empezar desde 0 no revises nada de lo que estaba ahi".

L3. Evidencia de estado: `ressy` contenía solo `AGENTS.md` y `.atl/`; `git rev-parse --show-toplevel` devolvía `/home/loonbac/Proyectos`. El remoto `loonbac/Ressy-Bot` tenía `main` y 7 ramas; acceso GitHub verificado con `gh auth status` (scopes `repo`, `workflow`). Go 1.26.8 y Docker 29.8.1 disponibles; sesión instabox activa.
