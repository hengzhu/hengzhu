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
	<form id="regform">
		<h3><b>注册账号</b></h3>
		<p>账号<input id="_user" type="text" name="username" required="true" /><img class="_user" style="visibility:hidden" src="static/easyui/img/login/error.png"/></p>
		<p>密码<input id="_word" type="password" name="password" required="true" /><img class="_word" style="visibility:hidden" src="static/easyui/img/login/1.png"/></p>
		<p>确认密码<input id="_wordi" type="password" name="repassword" required="true"/><img class="_wordi" style="visibility:hidden" src="static/easyui/img/login/2.png"/></p>
		<!--<p>邮箱<input id="_emil" type="text" /><img class="_emil" style="visibility:hidden" src="static/easyui/img/login/3.png"/></p>-->
		<p>手机号码<input id="_phone" name="mobile" type="text"/><img class="_phone" style="visibility:hidden" src="static/easyui/img/login/4.png"/></p>
		<p id="_phone22" style="display:none"><input type="text" name="smscode" placeholder="请输入短信验证码"/><input type="button"  onclick="settime(this)" name="" id="sendsmscode" value="免费获取验证码" /></p>
		<!--				<p>微信号码<input id="_weixin" type="text"/><img class="_weixin" style="visibility:hidden" src="static/easyui/img/login/5.png"/></p>-->
		<p>验证码<input id="_phone111" placeholder="请输入验证码" name="captcha" required="true" type="text"/>{{create_captcha}}<img class="_phone111" style="visibility:hidden" src="static/easyui/img/login/6.png"/></p>

		<img id="_regbuttona" src="static/easyui/img/login/zhuce_wanchengzhuce_btn.png"/>
	</form>
</article>
<div id="footer"></div>
<script type="text/javascript" src="static/easyui/js/jquery-2.1.0.js"></script>
<script type="text/javascript" src="static/easyui/js/bootstrap.min.js"></script>
<script>
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
<link rel="stylesheet" href="static/easyui/css/login/register.css" />
<script type="text/javascript" src="static/easyui/js/register.js" ></script>
</body>
</html>
{{template "easyui/website/header_or_footer/footer.tpl" .}}