package tmplt

//ListHTML - Baic List Template
var ListHTML = `
{{define "content"}}
<ul>
    {{ range $k, $v := .Data }}
       <li>{{$v}}</li>
    {{ end }}
</ul>{{end}}`
