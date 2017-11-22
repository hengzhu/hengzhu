{{template "easyui/website/header_or_footer/header.tpl" .}}
<html xmlns="http://www.w3.org/1999/html">

<head>
    <meta charset="UTF-8">
    <title>有才互娱</title>
    <meta name="keywords" content=""/>
    <meta name="description" content=""/>
    <link rel="stylesheet" href="static/easyui/css/bootstrap.min.css"/>
    <link rel="stylesheet" href="static/easyui/css/bootstrap-theme.min.css"/>
</head>

<body>
<div id="html"></div>
<article class="container page_center">
    <div class="row">
        <div class="col-md-4 col-xs-4">
            <div class="user_nc">
                <div class="row">
                    <img class="col-md-5 col-xs-5" src="static/easyui/img/user/account_zhanghao_icon.png"/>
                    <div class="col-md-7 col-xs-7">
                        <h4 id="user_nickname">{{.member.NickName}}</h4>
                        <p>账号: <span>{{.member.UserName}}</span></p>
                    </div>
                </div>
                <div>
                    <p><span>安全等级:</span> <span {{if ge .star 1}}id="_user_span_xin"{{end}}></span><span
                                {{if ge .star 2}}id="_user_span_xin"{{end}} ></span><span
                                {{if ge .star 3}}id="_user_span_xin"{{end}}></span><span
                                {{if ge .star 4}}id="_user_span_xin"{{end}}></span><span
                                {{if ge .star 5}}id="_user_span_xin"{{end}}></span></p>
                    <!--<p style="display: none;">账户余额: <span> 1234</span></p>-->
                    <p><span>绑定信息:</span>
                        {{if .IsBindMobile}}
                        <span></span>{{else}}<span style="background-position: -47px -67px;"></span>{{end}}
                        {{if .IsBindEmail}}<span></span>{{else}}<span
                                style="background-position: -47px -6px;"></span>{{end}}</p>
                    <p style="clear: both;margin: 0;"></p>
                </div>
            </div>
            <div class="user_btn">
                <button>账户安全</button>
                <button>修改密码</button>
                <button>完善个人信息</button>
                <button>联系客服</button>
                <button>充值中心</button>
                <button>我的礼包</button>
            </div>
            <div class="user_lb">
                <h3><b>最新礼包</b>
                    <small><a href="../../page/gift/gift.html"><img src="static/easyui/img/content/content_gd_btn.png"/></a>
                    </small>
                </h3>
                <div class="news_content news_lb">
                    <div class="row news_nu">
                        <a href="../../page/gift/gifts.html"><img class="col-xs-3 col-md-3"
                                                                  src="static/easyui/img/content/images/zuixinlibao_icon.jpg"/></a>
                        <div class="col-xs-9 col-md-9">
                            <p>圣斗士3D</p>
                            <p>圣斗士3D绑定微信限量礼包</p>
                            <p>剩余数量: <span>0000</span></p>
                        </div>
                        <a href="../../page/gift/gifts.html"><img class=""
                                                                  src="static/easyui/img/content/content_lqnormal_btn.png"/></a>
                    </div>
                    <div class="row news_nu">
                        <a href="../../page/gift/gifts.html"><img class="col-xs-3 col-md-3"
                                                                  src="static/easyui/img/content/images/zuixinlibao_icon.jpg"/></a>
                        <div class="col-xs-9 col-md-9">
                            <p>圣斗士3D</p>
                            <p>圣斗士3D绑定微信限量礼包</p>
                            <p>剩余数量: <span>0000</span></p>
                        </div>
                        <a href="../../page/gift/gifts.html"><img class=""
                                                                  src="static/easyui/img/content/content_lqnormal_btn.png"/></a>
                    </div>
                    <div class="row news_nu">
                        <a href="../../page/gift/gifts.html"><img class="col-xs-3 col-md-3"
                                                                  src="static/easyui/img/content/images/zuixinlibao_icon.jpg"/></a>
                        <div class="col-xs-9 col-md-9">
                            <p>圣斗士3D</p>
                            <p>圣斗士3D绑定微信限量礼包</p>
                            <p>剩余数量: <span>0000</span></p>
                        </div>
                        <a href="../../page/gift/gifts.html"><img class=""
                                                                  src="static/easyui/img/content/content_lqnormal_btn.png"/></a>
                    </div>
                    <div class="row news_nu">
                        <a href="../../page/gift/gifts.html"><img class="col-xs-3 col-md-3"
                                                                  src="static/easyui/img/content/images/zuixinlibao_icon.jpg"/></a>
                        <div class="col-xs-9 col-md-9">
                            <p>圣斗士3D</p>
                            <p>圣斗士3D绑定微信限量礼包</p>
                            <p>剩余数量: <span>0000</span></p>
                        </div>
                        <a href="../../page/gift/gifts.html"><img class=""
                                                                  src="static/easyui/img/content/content_lqnormal_btn.png"/></a>
                    </div>
                </div>
            </div>
        </div>
        <div class="col-md-8 col-xs-8 user_security">
            <div id="_user_security" class="_user_up">
                <h3>账户安全</h3>
                <div id="str_input">
                    <input style="margin-right: 30px;" onclick="ipt_str1(this)" type="button" name="" id=""
                           value="绑定手机"/>
                    <input onclick="ipt_str2(this)" type="button" name="" id="" value="绑定邮箱"/>
                    <input style="display: none;" type="button" name="" id="" value="防沉迷认证"/>
                </div>
                <div id="div_phones">
                    <h3>{{if .IsBindMobile}}手机换绑{{else}}绑定手机{{end}}</h3>
                    <p>手机号码&ensp;&ensp;<input type="text" name="" id="_phones" value=""/><img class="_inputi"
                                                                                              style="visibility:hidden"
                                                                                              src="static/easyui/img/user/cuowu.png"/>
                    </p>
                    <p>验证码&ensp;&ensp;<input style="margin-right: 5px;" type="text" name="" id="_smscode" value=""/>
                        <input type="button" onclick="settime(this)" name="mobile" id="sendvcode" class="sendvcode"
                               value="发送验证码"/></p>
                    <button><img name="mobile" id="submit_bind" class="submit_bind"
                                 src="static/easyui/img/user/account_right_queren_btn.png"/></button>
                </div>
                <div id="div_emil" style="display: none;">
                    <h3>{{if .IsBindEmail}}邮箱换绑{{else}}绑定邮箱{{end}}</h3>
                    <p>邮&ensp;&ensp;&ensp; 箱&ensp;&ensp;<input type="text" name="" id="_email" value=""/><img
                                style="visibility:hidden" src="static/easyui/img/user/cuowu.png"/></p>
                    <p>验证码&ensp;&ensp;<input style="margin-right: 5px;" type="text" name="" id="_emailcode" value=""/>
                        <input type="button" name="email" id="sendvcode" class="sendvcode" value="发送验证码"/></p>
                    <button><img name="email" id="submit_bind" class="submit_bind"
                                 src="static/easyui/img/user/account_right_queren_btn.png"/></button>
                </div>
            </div>
            <div id="_user_information" class="_user_up" style="display: none;">
                <form id="update_detail_form"><h3 style="margin-bottom: 15px;">完善个人信息</h3>
                    <p class="user_inf">温馨提示:建议你完善你的账户资料</p>
                    <div id="user_information">
                        <p>用户名称:&emsp;<input type="text" name="nickname" id="" placeholder="{{.memberinfo.NickName}}"
                                             value=""/></p>
                        <p>真实姓名:&emsp;<input type="text" name="truename" id="" placeholder="{{.memberinfo.RealName}}"
                                             value=""/></p>
                        <p>方便客服联络、实物奖品寄送</p>
                        <p>性别:&emsp;<input style="border:1px solid rgb(232,131,6);" type="radio" name="sex" id=""
                                           {{if eq .memberinfo.Sex "1" }}checked="checked"{{end}} value="1"/>男
                            <input type="radio" name="sex" id="" {{if eq .memberinfo.Sex "2" }}checked="checked"{{end}}
                                   value="2"/>女</p>
                        <div id="main">
                            <div class="demo">
                                <p>
                                    <label>出生日期：</label>
                                    <select id="sel_year" name="year"></select>年
                                    <select id="sel_month" name="month"></select>月
                                    <select id="sel_day" name="day"></select>日
                                <p>&ensp;&ensp;当前：{{dateformat .memberinfo.Birth "2006-01-02"}}</p>
                                </p>


                            </div>
                        </div>
                        <p>QQ号码:&emsp;<input type="text" name="qqnum" id="" placeholder="00000000" value=""/></p>
                        <p>手机号码:&emsp;<input type="text" name="" id=""
                                             placeholder="{{if .IsBindMobile}}{{.member.Mobile}}{{else}}请到账户安全绑定{{end}}"
                                             value="" disabled="true"/>{{if .IsBindMobile}}<span>已绑定</span>{{end}}</p>
                        <p>微信号码:&emsp;<input type="text" name="wechat" id="" placeholder="00000000" value=""/></p>
                        <p>通讯地址:&emsp;<input type="text" name="add" id="" value=""
                                             placeholder="{{.memberinfo.UserAddress}}"/></p>
                        <p>邮政号码:&emsp;<input type="text" name="zipcode" id="" placeholder="{{.memberinfo.UserPostcode}}"
                                             value=""/></p>
                        <input type="button" name="" id="submit_detail" value="提交"/>
                </form>
            </div>
        </div>
        <div id="_user_password" class="_user_up" style="display: none;">
            <form id="changepwd_form"><h3 style="margin-bottom: 147px;">修改密码</h3>
                <div id="user_password">
                    <p>旧密码:&emsp;<input type="password" name="password" id="" value=""/></p>
                    <p>新密码:&emsp;<input type="password" name="newpassword" id="paword" value=""/><img class="_word"
                                                                                                      style="visibility:hidden"
                                                                                                      src="static/easyui/img/login/1.png"/>
                    </p>
                    <p>确认新密码:&emsp;<input style="width: 52%;" type="password" name="repassword" id="pasword"
                                          value=""/><img class="_wordi" style="visibility:hidden"
                                                         src="static/easyui/img/login/2.png"/></p>
                    <input type="button" name="" id="submit_change_pwd" value="确认"/>
            </form>
        </div>
    </div>
    <div id="_user_Recharge" class="_user_up" style="display: none;">
        <h3 style="margin-bottom: 68px;">充值记录</h3>
        <div id=user_Recharges><p>查询时间:</p>
            <p class=datep><input class="datainp wicon" id=skinblue type=text placeholder="开始日期" value="" readonly><img
                        style="width: 10%" src=static/easyui/img/user/account_rili_icon.png></p>
            <p>至</p>
            <p class=datep><input class="datainp wicon" id=skinblue1 type=text placeholder="结束日期" readonly><img
                        src=static/easyui/img/user/account_rili_icon.png></p>
            <button><img id="submit_queryorder" src=static/easyui/img/user/account_chaxun_btn.png></button>
        </div>
        <div id=user_Recharge>
            <table class="table table-striped table-condensed table-bordered text-center" id="order_table">
                <tr id="head">
                    <td>充值时间</td>
                    <td>账号</td>
                    <td>充值游戏</td>
                    <td>订单状态</td>
                    <td>充值方式</td>
                    <td>充值金额</td>
                </tr>
                    {{range $key,$value := .orders}}

                    <tr>
                        <td>{{$value.OrderTime}}</td>
                        <td>{{$.member.NickName}}</td>
                        <td>{{$value.GameName}}</td>
                        <td>{{if gt $value.I.OrderStatus 6}}成功{{else if eq $value.I.OrderStatus 6}}<span
                                    style="color: red">失败</span>{{else}}未付款{{end}}</td>
                        <td>{{$value.PaymentChannel}}</td>
                        <td>¥{{$value.I.OrderRealAmount | printf "%.2f" }}</td>
                    </tr>{{end}}

            </table>
        </div>
    </div>
    <div id="_user_Package" class="_user_up" style="display: none;">
        <h3 style="margin-bottom: 68px;">我的礼包</h3>
        <div id=user_Recharges><p>查询时间:</p>
            <p class=datep><input class="datainp wicon" id=skinblue2 type=text placeholder="开始日期" value="" name="stime"
                                  readonly><img
                        src=static/easyui/img/user/account_rili_icon.png></p>
            <p>至</p>
            <p class=datep><input class="datainp wicon" id=skinblue3 type=text placeholder="结束日期" name="etime" readonly><img
                        src=static/easyui/img/user/account_rili_icon.png></p>
            <button><img src=static/easyui/img/user/account_chaxun_btn.png></button>
        </div>
        <div id=user_Recharge>
            <table class="table table-striped table-condensed table-bordered text-center">
                <tbody>
                <tr>
                    <td>礼包领取时间
                    <td>游戏名称
                    <td>礼包名称
                    <td>礼包结束时间
                    <td>激活码
                <tr>
                    <td>2016.10.26
                    <td>大话西游
                    <td>大话西游国庆礼包
                    <td>2016.11.1
                    <td>2394030
                <tr>
                    <td>2016.10.26
                    <td>大话西游
                    <td>大话西游国庆礼包
                    <td>2016.11.1
                    <td>2394030
                <tr>
                    <td>2016.10.26
                    <td>大话西游
                    <td>大话西游国庆礼包
                    <td>2016.11.1
                    <td>2394030
                <tr>
                    <td>2016.10.26
                    <td>大话西游
                    <td>大话西游国庆礼包
                    <td>2016.11.1
                    <td>2394030
                <tr>
                    <td>2016.10.26
                    <td>大话西游
                    <td>大话西游国庆礼包
                    <td>2016.11.1
                    <td>2394030
                <tr>
                    <td>2016.10.26
                    <td>大话西游
                    <td>大话西游国庆礼包
                    <td>2016.11.1
                    <td>2394030
                <tr>
                    <td>2016.10.26
                    <td>大话西游
                    <td>大话西游国庆礼包
                    <td>2016.11.1
                    <td>2394030
                <tr>
                    <td>2016.10.26
                    <td>大话西游
                    <td>大话西游国庆礼包
                    <td>2016.11.1
                    <td>2394030
                <tr>
                    <td>2016.10.26
                    <td>大话西游
                    <td>大话西游国庆礼包
                    <td>2016.11.1
                    <td>2394030
            </table>
        </div>
    </div>
    <div class="user_wd">
        <h3><b>我在玩的游戏</b></h3>
        <div class="col-md-12 col-xs-12 row">
            <div class="col-md-2  col-xs-4">
                <div><img src="static/easyui/img/gift/game_top_icon.png"/></div>
                <p>大唐游仙记</p>
                <div>
                    <ul class="all_games_ul">
                        <li>
                            <a href="">下载</a>
                        </li>
                        <li>
                            <a href="../../page/gift/gifts.html">礼包</a>
                        </li>
                    </ul>
                </div>
            </div>
            <div class="col-md-2 col-xs-4">
                <div><img src="static/easyui/img/gift/game_top_icon.png"/></div>
                <p>大唐游仙记</p>
                <div>
                    <ul class="all_games_ul">
                        <li>
                            <a href="">下载</a>
                        </li>
                        <li>
                            <a href="../../page/gift/gifts.html">礼包</a>
                        </li>
                    </ul>
                </div>
            </div>
            <div class="col-md-2 col-xs-4">
                <div><img src="static/easyui/img/gift/game_top_icon.png"/></div>
                <p>大唐游仙记</p>
                <div>
                    <ul class="all_games_ul">
                        <li>
                            <a href="">下载</a>
                        </li>
                        <li>
                            <a href="../../page/gift/gifts.html">礼包</a>
                        </li>
                    </ul>
                </div>
            </div>
            <div class="col-md-2 col-xs-4">
                <div><img src="static/easyui/img/gift/game_top_icon.png"/></div>
                <p>大唐游仙记</p>
                <div>
                    <ul class="all_games_ul">
                        <li>
                            <a href="">下载</a>
                        </li>
                        <li>
                            <a href="../../page/gift/gifts.html">礼包</a>
                        </li>
                    </ul>
                </div>
            </div>
            <div class="col-md-2 col-xs-4">
                <div><img src="static/easyui/img/gift/game_top_icon.png"/></div>
                <p>大唐游仙记</p>
                <div>
                    <ul class="all_games_ul">
                        <li>
                            <a href="">下载</a>
                        </li>
                        <li>
                            <a href="../../page/gift/gifts.html">礼包</a>
                        </li>
                    </ul>
                </div>
            </div>
            <div class="col-md-2 col-xs-4">
                <div><img src="static/easyui/img/gift/game_top_icon.png"/></div>
                <p>大唐游仙记</p>
                <div>
                    <ul class="all_games_ul">
                        <li>
                            <a href="">下载</a>
                        </li>
                        <li>
                            <a href="../../page/gift/gifts.html">礼包</a>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
    </div>
    </div>
    <div id="oDiv" class="block"></div>
</article>
<div id="footer"></div>
<script type="text/javascript" src="static/easyui/js/jquery-2.1.0.js"></script>
<script type="text/javascript" src="static/easyui/js/bootstrap.min.js"></script>

<link rel="stylesheet" href="static/easyui/css/header_or_footer/footer.css"/>
<link rel="stylesheet" href="static/easyui/css/header_or_footer/header.css"/>
<link rel="stylesheet" href="static/easyui/css/user/user.css"/>
<script type="text/javascript" src="static/easyui/js/birthday.js"></script>
<script type="text/javascript" src="static/easyui/js/jedate/jquery.jedate.js"></script>
<link type="text/css" rel="stylesheet" href="static/easyui/js/jedate/skin/jedate.css">
<script type="text/javascript" src="static/easyui/js/user.js"></script>

<script type="text/javascript">
    var start = {
        format: 'YYYY-MM-DD hh:mm:ss',
        minDate: '2014-06-16 23:59:59', //设定最小日期为当前日期
        festival: true,
        //isinitVal:true,
        maxDate: $.nowDate(0), //最大日期
        choosefun: function (elem, datas) {
            end.minDate = datas; //开始日选好后，重置结束日的最小日期
        }
    };
    var end = {
        format: 'YYYY年MM月DD日 hh:mm:ss',
        minDate: $.nowDate(0), //设定最小日期为当前日期
        festival: true,
        //isinitVal:true,
        maxDate: '2099-06-16 23:59:59', //最大日期
        choosefun: function (elem, datas) {
            start.maxDate = datas; //将结束日的初始值设定为开始日的最大日期
        }
    };
    $(function () {
        window.onload = function () {
            $('.navbar-nav>li').removeClass("active");
            $('.navbar-nav>li:nth-child(4)').addClass("active");
        }
    });
</script>
</body>

</html>

{{template "easyui/website/header_or_footer/footer.tpl" .}}