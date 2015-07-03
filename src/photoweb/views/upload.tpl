<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<link rel="stylesheet" href="/public/static/css/bootstrap.min.css" type="text/css">
<script type="text/javascript" src="public/static/js/bootstrap.min.js"></script>
<title>Upload</title>
</head>
<body>
<form role="form" method="POST" action="/upload" enctype="multipart/form-data">
	<div class="form-group">
 		Choose an image to upload : <input name="image" type="file"/> 
 		<input class="form-control" type="submit" value="Upload"/>
 	</div>
</form> 
</body>
</html>