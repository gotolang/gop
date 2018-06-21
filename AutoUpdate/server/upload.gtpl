<html>

<head>
<title>文件上传</title>
<link rel="stylesheet" href="up.css" />
</head>

<body align="center">
<h1>选择要上传的文件</h1>

<div id="up" >

<form enctype="multipart/form-data" action="http://172.42.9.188:9090/upload" method="post">
<input type="file" name="uploadfile">
<input type="hidden" name="token" value="{{.}}">
<input type="submit" value="Upload">
</form>
</div>
</body>

</html>