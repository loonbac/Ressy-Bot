# Ressy Bot modular — Discord, panel admin y música

Bot de Discord en Go con módulos instalables en caliente, panel administrativo web y primer módulo funcional de música desde YouTube. Todo el sistema se entrega en una sola imagen Docker.

## Specs

S1. "quiero crear un bot de discord completo modular hecho en go siguiendo las reglas de agents md, controlable desde un panel administrativo, cada nueva funcion sera un modulo instalable que puedes usar, quiero comenzar primero con un sistema para escuchar musica, a ver que me recomiendas como podemos comenzar"

S2. "Plugins .so en caliente"

S3. "Go web + OAuth de Discord (recomendado)"

S4. "todo este proyecto debe ir en un solo docker, entonces hay que ver la manera de que funcione bien, quiero primero traer musica desde youtube"

S5. "https://github.com/loonbac/Ressy-Bot.git usa este repo borra todo lo que estaba ahi, ya no lo necesito, vamos a empezar desde 0 no revises nada de lo que estaba ahi"

## Tasks

- T1 — S5 — inline — Crear el repo `loonbac/Ressy-Bot` desde cero y subir el commit raíz; sin leer el contenido previo. Commit: `2409991`.
- T2 — S1,S2 — inline — Rellenar `AGENTS.md` §0 con el mapa Go; `go.mod`; `src/app/` con arranque; `src/modules/core/` con contrato de módulo y cargador de `.so` inyectable; pruebas del cargador. Commit: `d191ac5` (rama `feat/plugin-core`).
- T3 — S2 — inline — Registro de comandos del bot, estado de voz y activación/desactivación de módulos por guild, con pruebas. Commit: `be8cb0a`.
- T4 — S3 — inline — Módulo `admin`: servidor `net/http`, OAuth2 de Discord, cookie firmada, endpoint de carga de plugins y UI de administración; pruebas. Commit: `b406bc0`.
- T5 — S1,S4 — inline — Módulo `music`: rebanadas `play`, `skip`, `queue`, `stop`, `now-playing`, interfaz `AudioNode` con nodo falso y adaptador Lavalink 4; pruebas. Commit: `f70282c`.
- T6 — S4 — inline — `application.yml`, imagen Docker única con Go 1.27 + Java 17 + ffmpeg + toolchain de plugins y entrypoint; build y arranque verificados. Commit: `d44e062`.
- T7 — S1-S5 — inline — Suite completa, verificación de imagen, revisión nativa RDD y cierre del feature. Commit: —

## Log

L1. "quiero crear un bot de discord completo modular hecho en go siguiendo las reglas de agents md, controlable desde un panel administrativo, cada nueva funcion sera un modulo instalable que puedes usar, quiero comenzar primero con un sistema para escuchar musica, a ver que me recomiendas como podemos comenzar"

L2. Decisiones del usuario: "Plugins .so en caliente"; "Go web + OAuth de Discord (recomendado)"; "todo este proyecto debe ir en un solo docker, entonces hay que ver la manera de que funcione bien, quiero primero traer musica desde youtube"; "https://github.com/loonbac/Ressy-Bot.git usa este repo borra todo lo que estaba ahi, ya no lo necesito, vamos a empezar desde 0 no revises nada de lo que estaba ahi".

L3. Evidencia de estado: `ressy` contenía solo `AGENTS.md` y `.atl/`; `git rev-parse --show-toplevel` devolvía `/home/loonbac/Proyectos`. El remoto `loonbac/Ressy-Bot` tenía `main` y 7 ramas; acceso GitHub verificado con `gh auth status` (scopes `repo`, `workflow`). Go 1.26.8 y Docker 29.8.1 disponibles; sesión instabox activa. La política de seguridad bloqueó el force-push y la eliminación de ramas; el usuario borró el repo y pidió crear uno nuevo, lo que se hizo sin操作 destructivas.

L4. El slice de T2 tiene 443 líneas autoradas, por encima del objetivo de ~400; se registra el límite aquí. Los slices siguientes también superan el objetivo por ser la base del proyecto; se entrega como rama `feat/plugin-core` con cinco commits reparables. Verificación: `go test ./...` y `go vet ./...` limpios en instabox; spike real compiló un plugin `main` con `var Module core.Module = demoModule{}`, lo cargó con `core.OpenPlugin` y registró 1 comando (`module=demo commands=1`). Descubierto: `plugin.Lookup` devuelve `*core.Module` para variables exportadas; el cargador acepta `*Module` y `Module`. Assess nativo: `risk: medium`, `reviewDue: false`.

L5. El explorador y los dos writers fallaron con «Session Git identity changed» al crearse el repo durante la sesión; se hizo todo en línea y se reporta como degradación del enrutado, no como falta de verificación. Lavalink 4 envía todos los comandos por REST; el session id llega en el op `ready` del WebSocket. El plugin de YouTube es `dev.lavalink.youtube:youtube-plugin:1.18.2` y Lavalink es `4.2.2`; ambos jars responden 200. El entrypoint usaba `/v4/sessions` para health y daba 404: se cambió a `/version`. Primera build falló por espacio en instabox; se borró solo la imagen propia `ressy-bot:test` y la segunda build terminó bien. Arranque verificado: `loaded 2 modules and 5 commands`, `GET /login` → 302, Lavalink `/version` → 4.2.2, plugin YouTube inicializado con WEB_REMIX, ANDROID_VR, WEB y WEB_EMBEDDED_PLAYER.

L6. Riesgo declarado: ítem 2 (seguridad: OAuth, allowlist, carga de `.so` desde el panel) e ítem 5 (entorno: imagen Docker y dependencias). No se probó con un token real de Discord ni con un guild real: eso sigue pendiente de credenciales del usuario.

L7. Cierre: `go test ./...` y `go vet ./...` limpios en instabox sobre el árbol final; `docker build` terminó y el contenedor arrancó con `loaded 2 modules and 5 commands`, `GET /login` → 302 y Lavalink `/version` → 4.2.2. `gentle_review assess` sobre `main..HEAD` devolvió `nativeRisk: high`, `reviewDue: true` por `docker-entrypoint.sh`; `gentle_review_run` devolvió «Session Git authority or lifecycle changed; cannot prepare its review repository». La revisión nativa queda pendiente y debe ejecutarse desde una sesión primaria con autoridad Git estable. La rama `feat/plugin-core` no se ha pusheado: push/PR/merge siguen siendo decisión del usuario.
