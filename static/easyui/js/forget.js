$(function() {
	$('.page_center>div:first-child>button').click(function() {
		str_up = $(this).index();
		if(str_up == 1){
			$('.page_center').html('<div id=email_un><p>账号绑定邮箱<input id="email_ipt" type=text onblur="emailblur(\'ss\')"></p><br><span class="_spani" style="color: red;visibility:hidden;margin-left: 258px;">邮箱错误</span><button><a href="../../page/login/new_password.html"><img src=../../img/login/wjmm_queren_btn.png></a></button></div>')
		}else if(str_up == 2){
			$('.page_center').html('<div id=phone_un><p>手机号码<input id="phone_ipt" type=text onblur="emailblur(\'aa\')"></p><br><span class="_span" style="color: red;visibility:hidden;margin-left: 125px;display: block;margin-top: -14px;">手机号错误</span><p>验证码<input type=text><img src=../../img/login/wjmm_yzm_btn.png></p><button><a href="../../page/login/new_password.html"><img src=../../img/login/wjmm_queren_btn.png></a></button></div>')
		}
	});
	$('#imgbtn').click(function() {
		var str_time = $('#daojishi');
		if(str_time.text()<=0){
			$('.gift_unts').css("display","block");
			$('.gift_unt').css("display","none");
		}else{
			$('.gift_unt').css("display","block");
			$('.gift_unts').css("display","none");
		}
	});
	$('.gift_unt>p:last-child>button:first-child').click(function() {
		$('.gift_unt').css("display","none");
		$('.gift_unts').css("display","none");
	});
});
	
	function emailblur(e) {
		var email = $('#email_ipt').val();
		var reg = /\w+[@]{1}\w+[.]\w+/;
			if(reg.test(email)) {
				$("._spani").css("visibility", "hidden");
			} else {
				$("._spani").css("visibility", "visible");
			}
	}
function emailblur(w) {
		var email = $('#phone_ipt').val();
		var reg = /^1(3|4|5|7|8)\d{9}$/;
			if(reg.test(email)) {
				$("._span").css("visibility", "hidden");
			} else {
				$("._span").css("visibility", "visible");
			}
	}
