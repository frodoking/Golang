<!DOCTYPE html>
<html>
<head>
<link rel="stylesheet" href="/public/static/css/bootstrap.min.css" type="text/css">
<script type="text/javascript" src="public/static/js/bootstrap.min.js"></script>
<meta charset="utf-8">
<title>List</title>
</head>
<body>
<ol class="breadcrumb">
	{{range $.images}}
 		<li><a href="/view?id={{.|urlquery}}">{{.|html}}</a></li>
 	{{end}}
</ol>
</body>
</html>