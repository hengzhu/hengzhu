{{template "easyui/website/header_or_footer/header.tpl" .}}
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>有才互娱</title>
	<meta name="keywords" content="" />
	<meta name="description" content="" />
</head>

<body>
<div id="html"></div>
<article class="container page_center">
	<div>
		<h3><b>登陆账号</b></h3>
		<p>账号</p><input type="text" placeholder="请输入账号" name="" id="user"  /><br><span class="_span" style="color: red;visibility:hidden">用户名错误</span>
		<p>密码</p><input type="password" placeholder="请输入密码" name="" id="password"  /><br><span class="_spani" style="color: red;visibility:hidden">密码错误</span>
		<p>
			<label><input name="Remember" type="checkbox" value="" />记住密码 </label>
			<label><input name="Auto" type="checkbox" value=""  />自动登陆 </label>
			<span><a href="../../page/login/Forget .html">忘记密码?</a></span><span><a href="/reg">注册新账号?</a></span>
		</p>
		<button id="_buttona"><img src="static/easyui/img/login/denglu_queren_btn.png"/></button>
	</div>
</article>
<div id="footer"></div>
<script>
//	$("#html").load("../header_or_footer/header.html");
//	$("#footer").load("easyui/website/header_or_footer/footer.tpl");
	$(function() {
		window.onload = function(){
			$('.navbar-nav>li').removeClass("active");
			$('.navbar-nav>li:nth-child(4)').addClass("active");
		}
	});
	var $ = jQuery.noConflict();
	$(function() {
		window.Login = new Login({ });
	});
</script>

<link rel="stylesheet" href="static/easyui/css/header_or_footer/footer.css" />
<link rel="stylesheet" href="static/easyui/css/header_or_footer/header.css" />
<link rel="stylesheet" href="static/easyui/css/login/login.css" />
<script type="text/javascript" src="static/easyui/js/login.js" ></script>
</body>
</html>
{{template "easyui/website/header_or_footer/footer.tpl" .}}