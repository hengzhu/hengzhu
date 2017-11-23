function Login(e) {
	this.conf = e;
	this.init();
}
Login.prototype = {
	init: function() {
		$("#news_word").blur(function(tu2) {
			var email = $(this).val();
			var reg = /^(?=.*?[a-zA-Z])(?=.*?[0-9])([A-Za-z]|[0-9]|[^\w\s]*|[_]*){6,20}$/;
			if(reg.test(email)) {
				$("._span").css("visibility", "hidden");
				Login.conf.tu2=true;

			} else {
				$("._span").css("visibility", "visible");
				Login.conf.tu2=false;
			}
		});
		$("#wordie").blur(function(tu2) {
			var word = $(this).val(),
				_newsword = $('#news_word').val();
			if(word === _newsword) {
				$("._spani").css("visibility", "hidden");
				Login.conf.tu2=true;

			} else {
				$("._spani").css("visibility", "visible");
				Login.conf.tu2=false;
			}
		});
		$("#_buttona").click(function(event) {
			if(Login.conf.tu1 && Login.conf.tu2) {
				window.location.href = "../../page/index.html";
			} else {
				return false
			}
		});
		
	}
}
