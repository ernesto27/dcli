## Resumen

(Resume el trabajo realizado)

## Videos/imagenes/logs que muestren el resultado

(Si es un bug, subir un antes y despues del comportamiento corregido. Si se incluye output de consola o código, usar (```) para aplicar formato)

## Jira ID

(Id del ticket, o link al mismo)

## Tener en cuenta antes de crear el MR

- No crear el MR sin antes hacer un merge inverso para actualizar tu rama con lo último de la rama destino (en este caso, develop).
- No crear el MR sin antes realizar una autorevisión de los cambios usando el tab Changes.
- No realizar cambios desde el editor online de gitlab, ni tampoco realizar cambios en general sin testearlos de forma local.
- Los cambios de un MR solo deben estar relacionados al feature/bug en desarrollo.
- Hacer uso del plugin SonarLint (de Intellij) para mantener un código aceptable. En la medida de lo posible, evitar generar issues críticos. Si no es posible, dejar un comentario en el código.
- Para que un MR sea aprobado por el assignee es necesario que resuelvas todos los threads abiertos. Resolver un thread implica realizar acciones correctivas que vayan en sintonia con lo marcado en el thread.
- En general, el assignee abrirá solo un thread por cada tipo de incidencia que encuentre. Por ej, si se declaran varias constantes en minuscula, el assignee solo marcará una vez el issue. Es responsabilidad del autor del MR analizar el resto de los cambios para que cumplan con todas las directivas/buenas prácticas marcadas. Asimismo, tratar de retener en la memoria las incidencias marcadas para no repetirlas en futuros MR's.
- No incluir loggeos informales, del estilo "pasó por aca". Lo mismo aplica para los comentarios.
- Para marcar tareas pendientes (ya sea de desarrollo o de analisis), el formato del comentario debe ser: // TODO (TU_NOMBRE): DESCRIPCION (LINK_A_TICKET_EN_JIRA). No olvidar crear el ticket con la descripción.
- Ser detallista a la hora de testear; no limitarse a que el feature funcione al nivel más básico. Probar casos extremos, y probar al menos si los features existentes relacionados siguen funcionando correctamente.
- Finalmente, producción puede parecer lejos, pero algún día llegará. Pensar que todo lo desarrollado algún día tendrá impacto en un ambiente real.
