<!DOCTYPE html>
<html>

	<head>
		<meta charset="UTF-8">
		<title>有才互娱</title>
		<meta name="keywords" content="" />
		<meta name="description" content="" />
		<link rel="stylesheet" href="/static/easyui/css/bootstrap.min.css" />
		<link rel="stylesheet" href="/static/easyui/css/bootstrap-theme.min.css" />
		<link rel="stylesheet" href="/static/easyui/css/header_or_footer/header.css" />
	</head>

	<body>
		<nav class="navbar nav_bg" role="navigation">
			<div class="container">
				<div class="container-fluid">
					<div class="navbar-header">
						<button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#example-navbar-collapse">
		            <span class="sr-only">切换导航</span>
		            <span class="icon-bar"></span>
		            <span class="icon-bar"></span>
		            <span class="icon-bar"></span>
		        </button>
						<a class="navbar-brand" href="#"><img src="/static/easyui/img/nav/nav_logo.png"></a>
					</div>
					<div class="collapse navbar-collapse" id="example-navbar-collapse">
						<ul class="nav navbar-nav">
							<li class="active">
								<a href="#">首页</a>
							</li>
							<li>
								<a href="#">游戏中心</a>
							</li>
							<li>
								<a href="/news">新闻中心</a>
							</li>
							<li>
								<a href="/account">账号中心</a>
							</li>
							<li>
								<a href="/service">客服中心</a>
							</li>
							<li>
								<a href="#">玩家论坛</a>
							</li>
						</ul>
					</div>
					<div class="nav_login">
						{{if .islogin}}<p><span><a href="/logout"></a></span><span><a href="/logout">注销</a></span></p>{{else}}
						<p><span><a href="/login">登陆</a></span><span><a href="/reg">注册</a></span></p>{{end}}
						<img src="/static/easyui/img/nav/nav_dlzc_btn.png"/>
					</div>
				</div>
			</div>
		</nav>
		<script type="text/javascript" src="/static/easyui/js/jquery-2.1.0.js"></script>
		<script type="text/javascript" src="/static/easyui/js/header.js"></script>
	</body>

</html>