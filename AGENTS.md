# Arquitectura modular para edición localizada

Objetivo único: encontrar el archivo exacto, cambiar solo eso y **no romper otra funcionalidad**.
Este archivo contiene solo las convenciones propias de estos proyectos. Las reglas de git,
revisión nativa, entrega e instabox están en el `AGENTS.md` global y no se repiten aquí.

## 0. Mapa de este repositorio (rellenar por repo)

- **Stack:** Go 1.26 con `go mod`; Discord `bwmarrin/discordgo` y audio Lavalink 4 se incorporan con el primer módulo funcional; una sola imagen Docker con Go, Java y ffmpeg.
- **Comandos:**
  - Prueba de un archivo: `go test ./src/modules/<modulo> -run Test<Nombre>`
  - Pruebas completas: `go test ./...`
  - Lint / tipos: `go vet ./...`
  - Build: `go build ./...`; imagen: `docker build -t ressy-bot .`
- **Checks:** `go vet ./...` y `go test ./...` desde la raíz; las reglas de arquitectura viven en este archivo. No hay linter de nombres ni CI todavía: se añaden en el paso 2 de la sección 7 cuando aparezca la primera corrección.
- **Módulos:**

  | Ruta | Responsabilidad | API pública |
  | :--- | :--- | :--- |
  | `src/modules/core/` | Contrato de módulos, catálogo de comandos y cargador de plugins `.so` | `core.Module`, `core.LoadModule`, `core.OpenPlugin`, `core.NewCommandCatalog` |
  | `src/modules/music/` | Cola y reproducción de música desde YouTube vía Lavalink | `music.New` (lo exporta el plugin) |
  | `src/modules/admin/` | Panel web con OAuth de Discord y control de módulos | `admin.New` |

- **No editar a mano:** `.atl/` (metadatos del harness), los `*.so` generados por `go build -buildmode=plugin` y `data/` (SQLite en ejecución).

Si el mapa está incompleto o es incorrecto, corregirlo en el mismo cambio que lo descubre.

## 1. Dónde está cada cosa

- Un módulo por capacidad del producto: `src/modules/<modulo>/`. Arranque, rutas y registro
  de dependencias en `src/app/`.
- **Un archivo por unidad exportada con lógica propia** (caso de uso, componente, hook,
  transformación, adaptador). El archivo se llama como su símbolo: `search-catalog.*` exporta
  `searchCatalog`; `CatalogItemCard.*` exporta `CatalogItemCard`.
- Una exportación principal por archivo. Los auxiliares privados viven en el mismo archivo,
  debajo de la exportación. Si dos archivos del mismo módulo los necesitan, van a
  `<modulo>/internal/<concepto>.*`.
- La prueba va al lado y se llama igual: `search-catalog.test.*` (o la convención del stack).
- **Dentro de un módulo se organiza por caso de uso, no por capa técnica.** Cada caso de uso
  es una rebanada: `<modulo>/<caso-de-uso>.*`, o la carpeta `<modulo>/<caso-de-uso>/` cuando
  tiene su propio adaptador o componente, con su lógica, su E/S y su prueba juntas. Prohibidas
  las carpetas horizontales `services/`, `repositories/`, `controllers/`, `handlers/`. Lo que
  comparten varias rebanadas (tipos y reglas del módulo) va en `<modulo>/domain.*`.
- Desde fuera de un módulo solo se importa su `public.*` (o `index`/`mod.rs`/equivalente).
- `shared/` solo con dos o más módulos que ya lo consumen. `shared/` no importa módulos.
- **Nombres únicos en el repo.** Buscar un símbolo debe devolver su definición y sus usos,
  nada más. Prohibidos `utils`, `helpers`, `common`, `data`, `handle`, `process`, `manager`
  como nombres de archivo o símbolo exportado. Prohibido construir dinámicamente nombres de
  funciones, eventos, rutas o claves: si el nombre no aparece literal, no se puede buscar.

**Cómo localizar:** mapa (sección 0) → módulo → búsqueda del símbolo → archivo. Si esos
pasos no bastan, el defecto está en el mapa o en el nombre: corregirlo.

### Interfaces que revelan intención

El nombre y la firma de lo exportado deben bastar para decidir si usarlo o cambiarlo **sin
leer su cuerpo**. Si hay que leer la implementación para saber qué hace, la interfaz está mal.

- **El nombre dice el efecto o el propósito, en el lenguaje del dominio, no el mecanismo:**
  `reserveSeat`, no `updateSeatRow`; `findOverdueInvoices`, no `filterInvoices`.
- **Consultas y comandos separados.** Las consultas (`get`, `find`, `list`, `count`, `is`,
  `has`, `can`) no tienen efectos observables. Los comandos usan un verbo de efecto
  (`reserve`, `cancel`, `publish`) y su resultado no es su propósito principal. Una función
  que consulta y modifica a la vez lo dice en su nombre (`takeNextJob`) o se divide.
- **Si cambia lo que hace, cambia su nombre.** Un nombre que ya no describe el comportamiento
  es un contrato roto (sección 2). Renombrar hace que el compilador y la búsqueda listen a
  todos los consumidores que hay que revisar.
- **La firma también revela intención:** tipos del dominio en vez de primitivos cuando llevan
  una regla (`Email`, `Money`, `SeatId`); un objeto de opciones con nombre en vez de
  booleanos posicionales; errores con el nombre de su causa (`SeatAlreadyReserved`), no
  `Error("failed")`.
- **Los efectos que no se pueden evitar se declaran:** precondiciones, postcondiciones e
  invariantes como aserciones en el código o como pruebas con nombre propio.
- **Diseñar desde el consumidor:** escribir primero la llamada en la prueba, como la leería
  quien la usa, y después la implementación.

Lo que de esta sección y de la regla de nombres únicos se pueda comprobar con lint (nombres
prohibidos, prefijos de booleanos, formato por tipo de símbolo) se lleva a lint según el
paso 2 de la sección 7. La prosa queda solo para lo que requiere criterio.

## 2. Antes de editar: medir el radio de impacto

1. **Listar los archivos que se van a tocar.**
2. **Buscar los usos de cada símbolo exportado que cambie** y clasificar el cambio:
   - **Local:** cambia la implementación, no el comportamiento observable.
   - **De contrato:** cambia la firma, el resultado, los errores o los efectos. Todos los
     consumidores encontrados entran en el alcance de la tarea.
3. **Si la unidad no tiene prueba,** escribir primero una que fije su comportamiento actual.
   Después cambiarla.
4. **Si lo que necesita un consumidor rompería a otro,** no añadir un parámetro booleano ni
   un `if` para ese caso en la función compartida. En este orden:
   1. Resolverlo en el lado del consumidor.
   2. Crear una función nueva con nombre propio.
   3. Cambiar el contrato y actualizar a todos los consumidores. Es obligatorio si la regla
      de negocio debe ser la misma para todos.

## 3. Durante la edición: contención

- Editar solo los archivos de la lista. Si aparece otro, añadirlo y justificarlo en el reporte.
- No tocar los internos de otro módulo para arreglar el propio. Solo su `public.*`, y si
  hace falta cambiarlo, es un cambio de contrato (paso 2.2).
- No mezclar el arreglo con renombres, reformateos, movimientos de archivos ni limpiezas
  «de paso».
- **Funcionalidad nueva = archivos nuevos + una línea de registro** en `src/app/`. No añadir
  ramas a funciones existentes que ya sirven a otros consumidores.
- Sin estado global mutable ni efectos al importar un archivo. Las dependencias entran por
  parámetro o constructor.
- No borrar, saltar ni debilitar pruebas para que pasen. Si falla una prueba ajena, o se
  rompió algo o el contrato cambió: decidir cuál y declararlo.

## 4. Después de editar: comprobar que no se rompió nada más

1. `git diff --stat` contra la lista del paso 2.1. Cada archivo extra se justifica o se revierte.
2. Ejecutar, según las reglas globales de compilación:
   - las pruebas de los archivos tocados;
   - las pruebas de los consumidores encontrados en el paso 2.2;
   - lint y tipos del repo completo (los tipos detectan consumidores rotos).
3. Si cambió un contrato o algo en `shared/`, ejecutar la suite completa.
4. Reportar: archivos tocados, contratos cambiados, consumidores verificados y qué no se
   pudo verificar.

## 5. Mantener el repo localizable

- Cada módulo no trivial tiene un `README.md` de 10–30 líneas: responsabilidad, API pública,
  qué módulos puede importar, invariantes y **dónde añadir** algo nuevo.
- Al tocar un archivo que tiene varias exportaciones independientes, o que mezcla UI, red y
  almacenamiento, separarlo como parte de esa tarea. No hacer barridos de reorganización
  fuera del alcance.

## 6. Correcciones bajo revisión nativa

Cuando la revisión nativa devuelve una corrección, el orden es obligatorio:

1. Registrar el forecast de líneas de corrección **antes de editar**.
2. Aplicar únicamente el fix aceptado, dentro de ese presupuesto.
3. Ejecutar y capturar la validación dirigida del fix.
4. Finalizar la revisión; solo entonces validar y commitear el árbol aprobado.

No editar entre la recepción del hallazgo y el forecast: invalida el binding del candidato
congelado y bloquea la revisión.

## 7. Cuando el agente se equivoca: árbol de corrección

Basado en *All code is slop (except mine)*, de Lauren Tan.
**Cada corrección se convierte en una comprobación, no en otra línea de este archivo.**

Cuando el usuario corrige un error del agente (un bug, un mal patrón, una suposición
equivocada), arreglar el caso concreto y además responder estas preguntas **en orden**.
Aplicar la primera que diga «sí»:

1. **¿El código puede hacerlo imposible?** → **Rediseñar.** Una API que no permita el uso
   incorrecto, tipos que no admitan el estado inválido, o prohibir el patrón en el código
   del producto.
2. **¿Lo puede detectar un lint, los tipos, el compilador, una prueba o CI?** →
   **Añadir el check** con un mensaje de error que diga qué hacer en lugar de eso, y
   eliminar el patrón viejo en todo el repo.
   - **Si el repo no tiene linter o no admite esa regla, instalarlo y configurarlo es parte
     de este paso**, no un motivo para bajar al paso 4. Usar el linter estándar del stack
     (y su plugin o regla personalizada cuando la regla no exista), con su comando en la
     sección 0 y en CI si el repo tiene CI.
   - Una regla nueva nace en modo error, no advertencia. Si hay demasiados casos viejos
     para limpiarlos en el mismo cambio, limitarla a los archivos nuevos o tocados y
     declarar la deuda pendiente en el reporte.
   - Solo bajar al paso 4 si ninguna herramienta del stack puede expresar el check.
     Decirlo en una frase.
3. **¿Se detecta ejecutando la app?** → **Verificar.** Añadir la prueba e2e o el paso de
   verificación (ejecutar la app, captura, salida observable) a la comprobación de la
   sección 4, con evidencia en el reporte.
4. **¿Es una decisión de criterio que nada de lo anterior detecta?** → **Escribir una regla**
   en la sección 8: una línea, no un párrafo.
5. **¿Nada de lo anterior?** → **Revisión humana.** Declararlo en el reporte como riesgo
   que solo el usuario puede revisar.

Reglas de aplicación:

- Los pasos 1 y 2 tocan código fuera de la tarea. Hacerlos como **cambio aparte** (commit o
  PR propio), no mezclados con el arreglo (sección 3). Si el usuario no lo pidió, proponerlo
  con la capa elegida y el check concreto, y esperar su respuesta. Instalar un linter o
  añadir una dependencia de desarrollo entra en esta misma propuesta.
- Antes de escribir una regla (paso 4), justificar en una frase por qué los pasos 1–3 no
  sirven.
- Si una regla de la sección 8 pasa a estar cubierta por un check, se borra de la sección 8.

## 8. Reglas aprendidas (paso 4 del árbol)

Una línea por regla: `- <regla> — <incidente que la originó>`.

<!-- vacío hasta que una corrección no pueda resolverse con los pasos 1–3 -->
