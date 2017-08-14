package tmplt

//URLListHTML - URL list html template
var URLListHTML = `
{{define "content"}}
<ul>
    {{ range $k, $v := .Data }}
       <li><strong>{{index $v 0}}</strong><a href="{{index $v 1}}">{{index $v 1}}</a></li>
    {{ end }}
</ul>{{end}}`
