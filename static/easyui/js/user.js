$(function ($, win, obj, ipt, ipt1) {
    $('.user_btn button').click(function () {
        var btn_1 = $('.user_btn button'),
            str_up = $(this).index(),
            str_security = $('.user_security>#_user_security'),
            str_security1 = $('.user_security>#_user_information'),
            str_security2 = $('.user_security>#_user_password'),
            str_security3 = $('.user_security>#_user_Recharge'),
            str_security5 = $('.user_security>._user_up'),
            str_security4 = $('.user_security>#_user_Package');
        btn_1.css("background-position", " 2px -194px");
        $(this).css("background-position", "2px -137px");
        if (str_up === 0) {
            str_security5.css("display", "none");
            str_security.css("display", "block");
        } else if (str_up === 2) {
            str_security5.css("display", "none");
            str_security1.css("display", "block");
        }
        $.ms_DatePicker({
            YearSelector: ".sel_year",
            MonthSelector: ".sel_month",
            DaySelector: ".sel_day"
        });
        $.ms_DatePicker();
        if (str_up === 1) {
            str_security5.css("display", "none");
            str_security2.css("display", "block");
            $("#paword").blur(function () {
                var email = $(this).val();
                var reg = /^(?=.*?[a-zA-Z])(?=.*?[0-9])([A-Za-z]|[0-9]|[^\w\s]*|[_]*){6,20}$/;
                if (reg.test(email)) {
                    $("._word").css("visibility", "hidden");
                } else {
                    $("._word").css("visibility", "visible");
                }
            });
            $("#pasword").blur(function () {
                var email = $(this).val();
                if ($("#paword").val() === email) {
                    $("._wordi").css("visibility", "hidden");
                } else {
                    $("._wordi").css("visibility", "visible");
                }
            });
        } else if (str_up === 4) {
            str_security5.css("display", "none");
            str_security3.css("display", "block");
        } else if (str_up === 5) {
            str_security5.css("display", "none");
            str_security4.css("display", "block");
        } else if (str_up === 3) {
            window.open("/service");
        }
        $("#skinblue").jeDate({
            isinitVal: true,
            skinCell: "jedateblue",
            format: 'YYYY-MM-DD hh:mm:ss'
        });
        $("#skinblue1").jeDate({
            isinitVal: true,
            skinCell: "jedateblue",
            format: 'YYYY-MM-DD hh:mm:ss'
        });
        $("#skinblue2").jeDate({
            isinitVal: true,
            skinCell: "jedateblue",
            format: 'YYYY-MM-DD hh:mm:ss'
        });
        $("#skinblue3").jeDate({
            isinitVal: true,
            skinCell: "jedateblue",
            format: 'YYYY-MM-DD hh:mm:ss'
        });
    });

    $('.news_nu>a:last-child>img').mouseover(function () {
        var btn_1 = $('.news_nu>a:last-child>img');
        btn_1.attr("src", "../../img/content/content_lqnormal_btn.png");
        $(this).attr("src", "../../img/content/content_lqclick_btn.png");
    }).mouseout(function () {
        var btn_1 = $('.news_nu>a:last-child>img');
        btn_1.attr("src", "../../img/content/content_lqnormal_btn.png");
    });

    var str_game = $(".user_wd").children("div").children("div").children("div:first-child").children("img");
    var str_odiv = $('#oDiv');
    str_game.mouseover(function () {
        str_odiv.html('<div class=row><div class=col-md-6><p>大唐游仙记</p><p><img src=../../img/content/android_icon.png><span>1.2.2</span></p><p><span>角色扮演</span><span>125MB</span></p></div><div class=col-md-6><img src=../../img/content/qrcode_for_gh_9eaaf35d2dd9_860.jpg><p>装进手机</p></div></div><h5>【游戏简介】</h5><p>××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××</p>')
        str_odiv.css("display", "block");
        str_odiv.css('position', 'absolute')
        str_odiv.css("left", window.event.x + document.body.scrollLeft);
        str_odiv.css("top", window.event.y + document.body.scrollTop);
    })
    str_game.mouseout(function () {
        str_odiv.css("display", "none");
    })
});

var countdown = 60;
var isf;
function settime(obj) {
    if (countdown == 0) {
        obj.removeAttribute("disabled");
        obj.value = "免费获取验证码";
        countdown = 60;
        return;
    } else {
        obj.setAttribute("disabled", true);
        this.ste = obj.value;
        obj.value = "重新发送(" + countdown + ")";
        countdown--;
    }
    isf = setTimeout(function () {
            settime(obj)
        }
        , 1000)
}
function ipt_str1(ipt) {
    var ipt_1 = $('#str_input>input:first-child'),
        ipt_2 = $('#str_input>input:nth-child(2)'),
        str_up = $(this).index(),
        str_h3 = $('.user_security>div:first-child>h3:nth-child(3)'),
        str_p1 = $('.user_security>div:first-child>p:nth-child(4)'),
        str_div1 = $('#div_phones'),
        str_div2 = $('#div_emil'),
        str_p2 = $('.user_security>div:first-child>p:nth-child(5)');
    ipt_2.css("background", "white");
    ipt_2.css("border", "1px solid #e5e5e5");
    ipt_2.css("color", "#237db2");
    ipt_1.css("background", "#237db2");
    ipt_1.css("border", "1px solid #237db2");
    ipt_1.css("color", "#fff");
    str_div2.css("display", "none");
    str_div1.css("display", "block");
}
function ipt_str2(ipt1) {
    var ipt_1 = $('#str_input>input:first-child'),
        ipt_2 = $('#str_input>input:nth-child(2)'),
        str_up = $(this).index(),
        str_h3 = $('.user_security>div:first-child>h3:nth-child(3)'),
        str_p1 = $('.user_security>div:first-child>p:nth-child(4)'),
        str_div1 = $('#div_phones'),
        str_div2 = $('#div_emil'),
        str_p2 = $('.user_security>div:first-child>p:nth-child(5)');
    ipt_1.css("background", "white");
    ipt_1.css("border", "1px solid #e5e5e5");
    ipt_1.css("color", "#237db2");
    ipt_2.css("background", "#237db2");
    ipt_2.css("border", "1px solid #237db2");
    ipt_2.css("color", "#fff");
    str_div1.css("display", "none");
    str_div2.css("display", "block");
}

//发送验证码
$(".sendvcode").click(function () {
    var sendtype = $(this).attr("name"),
        mobile = $("#_phones").val(),
        time = (Date.parse(new Date())) / 1e3,
        email = $("#_email").val(),
        bind, url;
    switch (sendtype) {
        case "mobile":
            bind = mobile;
            url = "/sendvcode?type=sms";
            break;
        case "email":
            bind = email;
            url = url = "/sendvcode?type=email";
            break;
        default:
            return false
    }
    $.post(url,
        {"time": time, "bind": bind},
        function (data) {
            alert(data.msg);
            if (data.status != 0) {
            } else {
                // window.location.href = "/account"
            }
        },
        "json")
});

//提交绑定信息
$(".submit_bind").click(function () {
    var sendtype = $(this).attr("name"),
        mobile = $("#_phones").val(),
        smscode = $("#_smscode").val(),
        email = $("#_email").val(),
        emailcode = $("#_emailcode").val(),
        bind, url, vcode;
    switch (sendtype) {
        case "mobile":
            bind = mobile;
            url = "/account/bind/?type=mobile";
            vcode = smscode;
            break;
        case "email":
            bind = email;
            url = url = "/account/bind/?type=email";
            vcode = emailcode;
            break;
        default:
            return false
    }
    $.post(url,
        {"vcode": vcode, "bind": bind},
        function (data) {
            alert(data.msg);
            if (data.status != 0) {

            } else {
                window.location.href = "/account"
            }
        },
        "json")
});

//修改密码
$("#submit_change_pwd").click(function () {
    $.post('/account/changepwd',
        $("#changepwd_form").serialize(),
        function (data) {
            alert(data.msg);
            if (data.status != 0) {

            } else if (data.status != 11) {
                window.location.href = "/login"
            }
        },
        "json")
});

//提交个人信息
$("#submit_detail").click(function () {
    $.post("/account/updatedetail",
        $('#update_detail_form').serialize(),
        function (data) {
            alert(data.msg);
            if (data.status != 0) {

            } else {
                // window.location.href = "/account"
            }
        },
        'json')
});

//按照时间条件查询订单
$("#submit_queryorder").click(function () {
    var h = $("#head").prop("outerHTML"),
        nickname = $("#user_nickname").text();
    var stime = Date.parse(new Date($("#skinblue").val())) / 1e3,
        etime = Date.parse(new Date($("#skinblue1").val())) / 1e3;
    $.get("/account",
        {"stime": stime, "etime": etime},
        function (json) {
            if (json.status == 0 && json.data != null) {
                $('#order_table').html("");
                var s = "", status;
                for (var i = 0; i < json.data.length; i++) {
                    if (json.data[i].io.order_status > 6) {
                        status = "成功"
                    } else if (json.data[i].io.order_status = 6) {
                        status = '<span style="color: red">失败</span>'
                    } else {
                        status = "未付款"
                    }
                    s += '<tr><td>' + json.data[i].order_time + '</td><td>' + nickname + '</td><td>' + json.data[i].game_name + '</td><td>' + status + '</td><td>' + json.data[i].payment_channel + '</td><td>¥' + json.data[i].io.order_real_amount.toFixed(2) + '</td>';
                }
                $('#order_table').html(h + s);
            } else {
                alert("没有数据");
            }
        },
        "json");
});