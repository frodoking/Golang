<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>List</title>
</head>
<body>
<ol>
	{{range $.images}}
 		<li><a href=\"/view?id={{.|urlquery}}">{{.|html}}</a></li>
 	{{end}}
</ol>
</body>
</html>