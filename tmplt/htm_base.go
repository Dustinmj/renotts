package tmplt

//BaseHTML - Base HTML Content
var baseHTML = `<!DOCTYPE HTML>
<html lang="en-US">
<head>
  <title>{{.Common.Title}}</title>
  <style type="text/css">
    /*! RenoTTS Tester */
    body{background-color:#FFF;padding-left:40px;color:#009fdb;padding-bottom:50px;}
    div{clear:both;padding-top:10px}
    div#output,div#request{display:none}
    p.txt{color:#666}
    p.txt span{float:left;clear:both;padding:.25em 0}
    h1{height:60px;line-height:60px}
    h1 img{float:left;vertical-align:middle;width:60px;height:60px;margin-right:1em}
    p.txt span strong{padding-right:10px}
    input,select{background:#efefef;padding:.25em;border:1px solid #acacac;width:30em}
    select{width:30.5em}
    input[type=button]{width:10em;display:block;margin-top:2em;background-color:#009fdb;color:#FFF}
    label{display:inline-block;width:6em;padding:1em 0 0}
    label.pad{width:auto;padding:1em 1em 0 0}
    label.pad+input[type=checkbox]{width:auto;display:inline-block;margin-right:3em}
    h2{color:#333;font-size:1.3em}
    div#output{min-height:300px}
    ul li{list-style-type:none;padding:0 0 1em;color:#333}
    a,a:visited{color:#369}
    a:hover,a:active{color:#333}
    ul li a{padding-left:10px}
    ul{padding:0}
    pre>code{display:inline-block;padding: 1rem 1.5rem;white-space: pre;}
    code{color:#333;font-family: monospace,monospace;padding: .2rem .5rem;margin: 0 .2rem;font-size: 90%;white-space: nowrap;background: #F1F1F1;border: 1px solid #E1E1E1;border-radius: 4px;}
    p{color:#333;max-width:600px;}
    ol{max-width:600px;color:#333}
    ol li{padding:.5em 0 0}
    input[type=button]{cursor:pointer}
  </style>
</head>
<body>
<div><h1><a href="{{.Common.BaseURI}}"><img alt="renotts.bl" src="data:image/png;base64,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"></a>
{{.Common.Title}}</h1>
{{ template "content" .}}
</div>
{{ template "javascript" .}}
</body>
</html>`
