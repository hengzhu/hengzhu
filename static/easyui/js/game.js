$(function() {
	var str_game = $(".all_games").children("div").children("div").children("div:first-child").children("img");
	console.log(str_game.length)
	var str_odiv = $('#oDiv');
	str_game.mouseover(function () {
		str_odiv.html('<div class=row><div class="col-md-6 col-xs-6"><p>大唐游仙记</p><p><img src=../../img/content/android_icon.png><span>1.2.2</span></p><p><span>角色扮演</span><span>125MB</span></p></div><div class="col-md-6 col-xs-6"><img src=../../img/content/qrcode_for_gh_9eaaf35d2dd9_860.jpg><p>装进手机</p></div></div><h5>【游戏简介】</h5><p>××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××</p>')
		str_odiv.css("display","block");
		str_odiv.css('position','absolute')
        str_odiv.css("left",window.event.x + document.body.scrollLeft);
        str_odiv.css("top",window.event.y + document.body.scrollTop);
	})
	str_game.mouseout(function(){ 
			str_odiv.css("display","none");
	})
	window.onload = function(){
		$('.navbar-nav>li').removeClass("active");
			$('.navbar-nav>li:nth-child(2)').addClass("active");
		}
});