$(function() {
	$('.navbar-nav>li').click(function() {
		var str_li = $('.navbar-nav>li'),
		str_up = $(this).index();
		$('.navbar-nav>li').removeClass("active");
		$(this).addClass("active");
		if(str_up == 0){
			window.location.href="../../page/index.html"
		}else if(str_up == 1){
			window.location.href="../../page/game/game_center.html"
		}else if(str_up == 2){
			window.location.href="../../page/news/news_center.html"
		}else if(str_up == 3){
			window.location.href="../../page/account/user.html"
		}else if(str_up == 4){
			window.location.href="../../page/cs/customer_service_center.html"
		}else if(str_up == 5){
			window.location.href="../../page/index.html"
		}
	})	
})