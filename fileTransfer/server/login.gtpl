<html>

<head>
<title>The first web page</title>
</head>

<body>
<form action="http://172.42.9.188:9090/login" method="post">
用户名：<input type="text" name="username">
密码：<input type="passwrod" name="password">
<input type="hidden" name="token" value="{{.}}">
<input type="submit" value="登录">
</form>
</body>

</html>