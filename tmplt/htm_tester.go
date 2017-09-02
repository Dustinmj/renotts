package tmplt

//TestHTML - renotts tester html
var TestHTML = `
{{define "css"}}
div#tester div#output,div#tester div#request{display:none;background-color:#EFEFEF;color:#333;border:1px solid #DDD;border-radius:4px;margin:2rem 0;padding:0 1rem;max-width:50rem;font-size:80%}div#tester div#output p.txt,div#tester div#request p.txt{color:#666}div#tester div#output p.txt span,div#tester div#request p.txt span{display:block;padding:.25em 0}div#tester div#output p.txt span strong,div#tester div#request p.txt span strong{padding-right:10px}div#tester div#output h2,div#tester div#request h2{color:#009fdb}div#tester div#output{min-height:10rem}div#tester input,div#tester select{background:#EFEFEF;padding:.25em;border:1px solid #ACACAC;width:30em}div#tester select{width:30.5em}div#tester input[type=button]{width:10em;display:block;margin-top:1em;background-color:#009fdb;color:#FFF;font-size:105%;font-weight:700;cursor:pointer}div#tester input[type=checkbox]{width:auto;display:inline-block;margin-right:3em}div#tester label{display:inline-block;width:7em;padding:0 0 1em;color:#333}div#tester label .pad{width:auto;padding:1em 1em 0 0}
{{end}}
{{define "content"}}
<!-- start copy -->
<div id="tester">
    <form action="" method="post">
        <label for="host">Host</label>
        <input type="text" id="host" name="host" value="192.168.1.70:8080" />
        <br />
        <label for="path">Path</label>
        <input type="text" id="path" name="path" value="/tts/polly" />
        <br />
        <label for="voice">Voice</label>
        <select id="voice" name="voice">
            <option>Joanna</option><option>Ivy</option><option>Joey</option><option>Justin</option><option>Kendra</option><option>Kimberly</option>
            <option>Salli</option> <option>Astrid</option> <option>Amy</option> <option>Brian</option> <option>Carla</option> <option>Carmen</option> <option>Celine</option> <option>Chantal</option>
            <option>Conchita</option> <option>Cristiano</option> <option>Dora</option> <option>Emma</option> <option>Enrique</option> <option>Ewa</option> <option>Filiz</option>
            <option>Geraint</option> <option>Giorgio</option> <option>Gwyneth</option> <option>Hans</option> <option>Ines</option> <option>Jacek</option> <option>Jan</option> <option>Karl</option>
            <option>Liv</option> <option>Lotte</option> <option>Mads</option> <option>Maja</option> <option>Marlene</option> <option>Mathieu</option> <option>Maxim</option> <option>Miguel</option>
            <option>Mizuki</option> <option>Naja</option> <option>Nicole</option> <option>Raveena</option> <option>Ricardo</option> <option>Penelope</option> <option>Ruben</option>
            <option>Russell</option> <option>Tatyana</option> <option>Vitoria</option>
        </select>
        <br />
        <label for="sampleRate">Sample Rate</label>
        <select id="sampleRate" name="sampleRate">
            <option>22050</option>
            <option>16000</option>
            <option>8000</option>
        </select>
        <br />
        <label for="tts">Text</label>
        <input type="text" id="tts" name="tts" value="This is a test." />
        <br />
        <label class="pad" for="padBefore">Silence Before</label>
        <input type="checkbox" id="padBefore" name="padBefore" />
        <label class="pad" for="padAfter">Silence After</label>
        <input type="checkbox" id="padAfter" name="padAfter" />
        <input type="button" id="talk" value="Send Request" />
    </form>
    <div id="request">
        <h2>Request:</h2>
        <p class="txt">
        </p>
    </div>
    <div id="output">
        <h2>Response:</h2>
        <img class="loading" src="data:image/gif;base64,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" alt="Loading" />
        <p class="txt">
        </p>
    </div>
</div> <!--end copy-->
{{end}}
{{define "javascript"}}
<script type="text/javascript">
var renotts={init:function(){var e=this;document.getElementsByTagName("form")[0].querySelector("input[type=button]").addEventListener("click",function(){e.talk()})},talk:function(){var e,t=this.inputData();t.padA&&t.padB?e="Both":t.padA?e="After":t.padB&&(e="Before");var n=this.requestParams(t,e),a="http://"+t.host+t.path;this.showData("request",{URI:a,Method:"POST",Data:n}),this.showData("output",{}),this.showLoading(),this.sendRequest(a,n)},sendRequest:function(e,t){var n=this,a=new XMLHttpRequest;a.onreadystatechange=function(){a.readyState==XMLHttpRequest.DONE&&n.ttsCallback(a)},a.open("POST",e,!0),a.setRequestHeader("Content-type","application/json"),a.send(t)},showLoading:function(){document.getElementById("output").querySelector(".loading").style.display="block"},hideLoading:function(){document.getElementById("output").querySelector(".loading").style.display="none"},showData:function(e,t){var n=document.getElementById(e);n.style.display="block";var e=n.querySelector(".txt");e.innerHTML="";for(k in t)e.appendChild(this.label(k,t[k]))},requestParams:function(e,t){return JSON.stringify({text:e.text,sampleRate:e.sampleRate,voice:e.voice,padding:t})},inputData:function(){return{host:document.getElementById("host").value,path:document.getElementById("path").value,text:document.getElementById("tts").value,padB:this.isChecked("padBefore"),padA:this.isChecked("padAfter"),sampleRate:this.selectedVal("sampleRate"),voice:this.selectedVal("voice")}},isChecked:function(e){var e=document.getElementById(e);return e.checked},selectedVal:function(e){var e=document.getElementById(e);return e.options[e.selectedIndex].text},ttsCallback:function(e){if("undefined"!=e){var t=e.getAllResponseHeaders()||"Not Available";"undefined"!=t&&(t=t.replace(/[^\x20-\x7E]/gim," | "));var n=[e.status,e.statusText].join(" - ");this.output(t,n,e.response)}},output:function(e,t,n){n=n||"No Data",t=t||"Error",e=e||"Not Available",this.showData("output",{Status:t,"Response Headers":e,"Response Body":n}),this.hideLoading()},label:function(e,t){"object"==typeof t&&(t=JSON.stringify(t));var n=document.createElement("span"),a=document.createElement("strong");a.innerText=e;var o=document.createElement("em");return o.innerText=t,n.appendChild(a),n.appendChild(o),n}};window.addEventListener("DOMContentLoaded",function(){renotts.init()},!1);
</script>
{{end}}`
