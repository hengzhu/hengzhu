function Login(e) {
    this.conf = e;
    this.init();
}
Login.prototype = {
    init: function () {
        $("#_user").blur(function (tu1) {
            var email = $(this).val();
            var reg = /^\w{6,20}$/;
            if (reg.test(email)) {
                $("._user").css("visibility", "hidden");
                Login.conf.tu1 = true;
            } else {
                $("._user").css("visibility", "visible");
                Login.conf.tu1 = false;

            }
        });

        $("#_word").blur(function (tu2) {
            var email = $(this).val();
            var reg = /^(?=.*?[a-zA-Z])(?=.*?[0-9])([A-Za-z]|[0-9]|[^\w\s]*|[_]*){6,20}$/;
            if (reg.test(email)) {
                $("._word").css("visibility", "hidden");
                Login.conf.tu2 = true;

            } else {
                $("._word").css("visibility", "visible");
                Login.conf.tu2 = false;
            }
        });
        $("#_wordi").blur(function (tu3) {
            var email = $(this).val();
            if ($("#_word").val() === email) {
                $("._wordi").css("visibility", "hidden");
                Login.conf.tu3 = true;

            } else {
                $("._wordi").css("visibility", "visible");
                Login.conf.tu3 = false;
            }
        });
        $("#_emil").blur(function (tu4) {
            var email = $(this).val();
            var reg = /\w+[@]{1}\w+[.]\w+/;
            if (reg.test(email)) {
                $("._emil").css("visibility", "hidden");
                Login.conf.tu4 = true;

            } else {
                $("._emil").css("visibility", "visible");
                Login.conf.tu4 = false;
            }
        });
        $("#_phone").blur(function (tu5) {
            var email = $(this).val();
            var reg = /^1(3|4|5|7|8)\d{9}$/;
            if (reg.test(email)) {
                $("._phone").css("visibility", "hidden");
                $("#_phone22").css("display", "block");
                Login.conf.tu5 = true;

            } else {
                $("._phone").css("visibility", "visible");
                $("#_phone22").css("display", "none");
                Login.conf.tu5 = false;
            }
        });
        $("#_phone111").blur(function (tu6) {
            var email = $(this).val();
            var reg = /^[a-z0-9]{4}$/;
            if (reg.test(email)) {
                $("._phone111").css("visibility", "hidden");
                Login.conf.tu2 = true;

            } else {
                $("._phone111").css("visibility", "visible");
                Login.conf.tu2 = false;
            }
        });
        $("#_regbuttona").click(function (event) {
            if (Login.conf.tu1 && Login.conf.tu2 && Login.conf.tu3) {
                $.ajax({
                    type: "POST",
                    url: '/reg?isajax=1',
                    data: $('#regform').serialize(), // 要提交的表单,必须使用name属性
                    async: true,
                    success: function (data) {
                        // $("#common").html(data);//输出提交的表表单
                        if (data.status != 0) {
                            alert(data.msg)
                        } else {
                            alert(data.msg);
                            window.location.href = "/account"
                        }
                    },
                    dataType: "json"
                    // error: function (request) {
                    //     alert("Connection error");
                    // }
                });
                // window.location.href = "/../page/index.html";
            } else {
                return false
            }
        });
        $("#sendsmscode").click(function (event) {
            alert("send!");
            var mobile = $("#_phone").val();
            var time = (Date.parse(new Date()))/1e3;
            $.ajax({
                type:"POST",
                url:"/sendvcode?type=sms",
                data:{"time":time,"mobile":mobile},
                dataType:"json",
                success: function (data) {
                    // $("#common").html(data);//输出提交的表表单
                    if (data.status != 0) {
                        alert(data.msg)
                    } else {
                        alert(data.msg);
                        // window.location.href = "/account"
                    }
                }
            });
        });

    }
};
var countdown = 60;
function settime(obj) {
    // var mobile = $("#_phone").val();
    // alert(mobile);
    if (countdown == 0) {
        obj.removeAttribute("disabled");
        obj.value = "免费获取验证码";
        countdown = 60;
        return;
    } else {
        obj.setAttribute("disabled", true);
        obj.value = "重新发送(" + countdown + ")";
        countdown--;
    }
    setTimeout(function () {
            settime(obj)
        }
        , 1000)
}