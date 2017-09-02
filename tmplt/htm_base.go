package tmplt

//BaseHTML - Base HTML Content
var baseHTML = `<!DOCTYPE HTML>
<html lang="en-US">
<head>
  <title>{{.Common.Title}}</title>
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <link href="https://fonts.googleapis.com/css?family=Roboto:400,700" rel="stylesheet">
  <style type="text/css">
  body{font-family:Roboto,HelveticaNeue,"Helvetica Neue",Helvetica,Arial,sans-serif;background-color:#FFF;color:#333;padding:0 0 3em;margin:0;font-size:16px;min-width:960px}body>div,body>header{padding:0 0 0 2.5em;margin:0;min-width:575px}body>header{height:90px;background-color:#009fdb;margin-bottom:2em}body h1{height:60px;line-height:60px;color:#FFF;margin:0 0 1em;padding:10px 0 0;font-weight:400;font-size:2em}body h1 img{margin-right:1em;vertical-align:middle}body h2{color:#333;font-size:1.3em}body footer{margin:40px 0 0 2.5em}body footer p{font-size:.75em;color:#66666}body a{color:#369;text-decoration:underline;cursor:pointer}body a:hover{color:#444}body ul{padding:0;margin:0}body ul li{list-style-type:circle;padding:0 0 1.25em;color:#333}body ul li a{font-size:90%;white-space:nowrap;background-color:#EFEFEF;border:1px solid #E1E1E1;border-radius:4px;padding:.2em .5em;margin-left:.5em}body ul li a:hover{background-color:#CACACA}body ul li strong{color:#333;font-weight:400}body ul li strong:after{content:":"}body ul li strong::first-letter{text-transform:uppercase}body ul.surround li{font-size:90%;background:#F1F1F1;border:1px solid #E1E1E1;border-radius:4px;padding:.2em .5em;margin:0 0 1em;max-width:35em}body ul.surround li::first-letter{text-transform:uppercase}body code{display:inline-block;padding:1em 1.5em;color:#333;font-family:monospace,monospace;margin:1em .2em;font-size:90%;white-space:nowrap;background:#F1F1F1;border:1px solid #E1E1E1;border-radius:4px;width:auto;max-width:80em;overflow-x:scroll;resize:horizontal}@media (max-width:600px){body{font-size:14px}}
  {{ template "css" .}}
  </style>
</head>
<body>
<header>
    <h1>
        <a href='/'><img alt="renotts.bl" src="data:image/png;base64,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"></a>
    {{.Common.Title}}</h1>
</header>
<div>
    {{ template "content" .}}
</div>
<footer>
    <p>RenoTTS &copy;2017 Dustin Jorge | Licensed under the terms of the <a href="https://opensource.org/licenses/MIT">MIT License</a>.</p>
</footer>
{{ template "javascript" .}}
</body>
</html>`
