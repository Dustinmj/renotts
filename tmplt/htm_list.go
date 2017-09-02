package tmplt

//ListHTML - Baic List Template
var ListHTML = `
{{define "css"}}{{end}}
{{define "content"}}
<ul class="surround">
    {{ range $k, $v := .Data }}
       <li>{{$v}}</li>
    {{ end }}
</ul>{{end}}{{define "javascript"}}{{end}}`
