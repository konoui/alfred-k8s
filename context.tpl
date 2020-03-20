{{ range .contexts -}}
*   {{.name}}   {{.context.cluster}}   {{.context.user}}    {{.context.namespace}}
{{ end -}}
