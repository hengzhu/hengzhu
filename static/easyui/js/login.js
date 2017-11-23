function Login(e) {
    this.conf = e;
    this.init();
}
Login.prototype = {
    init: function () {
        $("#user").blur(function (tu1) {
            var email = $(this).val();
            var reg = /^\w{6,20}$/;
            var emailreg = /^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$/;
            if (reg.test(email) || emailreg.test(email)) {
                $("._span").css("visibility", "hidden");
                Login.conf.tu1 = true;
            } else {
                $("._span").css("visibility", "visible");
                Login.conf.tu1 = false;

            }
        });

        $("#password").blur(function (tu2) {
            var email = $(this).val();
            var reg = /^(?=.*?[a-zA-Z])(?=.*?[0-9])([A-Za-z]|[0-9]|[^\w\s]*|[_]*){6,20}$/;
            if (reg.test(email)) {
                $("._spani").css("visibility", "hidden");
                Login.conf.tu2 = true;

            } else {
                $("._spani").css("visibility", "visible");
                Login.conf.tu2 = false;
            }
        });


        $('input[name="Auto"]').click(function () {
//			console.log($('input[name="Auto"]').is(':checked')+"+"+$('input[name="Remember"]').is(':checked'));
            if ($('input[name="Auto"]').is(':checked') == true) {
                $('input[name="Remember"]').prop("checked", true);
            }else {
                $('input[name="Remember"]').prop("checked", false);
            }
        });
        $('input[name="Remember"]').click(function () {
            console.log($('input[name="Auto"]').is(':checked') + "+" + $('input[name="Remember"]').is(':checked'));
            if ($('input[name="Remember"]').is(':checked') == false) {
                $('input[name="Auto"]').prop("checked", false);
            }
        });

        $("#_buttona").click(function (event) {
            var username = $("#user").val();
            var password = $("#password").val();
            var autologin = $('input[name="Auto"]')[0].checked;
            if (Login.conf.tu1 && Login.conf.tu2) {
                $.post("/login?isajax=1",
                    {"username": username, "password": password, "autologin":autologin},
                    function (data) {
                        if (data.status != 0) {
                            alert(data.msg)
                        } else {
                            $("#footer").after(data.data);
                            alert(data.msg);
                            // window.location.href = "/account"
                            window.setTimeout(function(){ location.href = "/account"; },500);
                        }
                    },
                    "json"
                )
            } else {
                return false
            }
        });

    }
};
