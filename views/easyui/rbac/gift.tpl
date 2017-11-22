{{template "../manage/header.tpl"}}

<script type="text/javascript">
    var gamelist=$.parseJSON({{.game | stringsToJson}})
    var URL="/rbac/gift";
    $(function(){
        //角色列表
        $("#datagrid").datagrid({
            title:'角色管理',
            url:URL+"/index",
            method:'POST',
            pagination:true,
            fitColumns:true,
            striped:true,
            rownumbers:true,
            singleSelect:true,
            idField:'Id',
            columns:[[
                {field:'Id',title:'ID',width:30,align:'center'},
                {field:'Packagename',title:'礼包名',width:130,align:'center',editor:'text'},
                {field:'Game',title:'游戏',width:130,align:'center',
                    formatter:function(value){
                        for(var i=0; i<gamelist.length; i++){
                            if (gamelist[i].Id == value) return gamelist[i].Gamename;
                        }
                        return value;
                    },
                    editor:{
                        type:'combobox',
                        options:{
                            valueField:'Id',
                            textField:'Gamename',
                            data:gamelist,
                            required:true
                        }
                    }},
                {field:'Quantity',title:'礼包数量',width:70,align:'center',editor:'text'},
                {field:'Code',title:'礼包码',width:70,align:'center',
                    formatter:function (value) {
                        var row = $("#datagrid").datagrid("getSelected")
                        return "<input type='submit' value='下载礼包码' onclick='downLoadGiftCode()'/>"
                        return " <a href='" + value + "' target='_blank'><img src='" + value + "' style='width: 60px; height:60px;'/> </a>";
                    },editor:'text'},
                {field:'Starttime',title:'开始时间',width:100,align:'center',
                    formatter:function(value){
                        if(value) return phpjs.date("Y-m-d H:i:s",phpjs.strtotime(value));
                        return value;
                    },
                    editor:'text'},
                {field:'Endtime',title:'结束时间',width:100,align:'center',
                    formatter:function(value){
                        if(value) return phpjs.date("Y-m-d H:i:s",phpjs.strtotime(value));
                        return value;
                    },
                    editor:'text'},
                {field:'Createtime',title:'发放时间',width:100,align:'center',
                    formatter:function(value){
                        if(value) return phpjs.date("Y-m-d H:i:s",phpjs.strtotime(value));
                        return value;
                    }},
                {field:'Desc',title:'描述',width:190,align:'center',editor:'text'}
            ]],
//            onAfterEdit:function(index, data, changes){
//                if(vac.isEmpty(changes)){
//                    return;
//                }
//                if(data.Id == undefined){
//                    changes.Id = 0;
//                }else{
//                    changes.Id = data.Id;
//                }
//                wins = $("#datetimepicker").val();
//                vac.alert(wins);
//                vac.ajax(URL+'/UpdateGift', changes, 'POST', function(r){
//                    if(!r.status){
//                        vac.alert(r.info);
//                    }else{
//                        $("#datagrid").datagrid("reload");
//                    }
//                })
//            },
            onDblClickRow:function(index,row){
                editrow();

            },
            onRowContextMenu:function(e, index, row){
                e.preventDefault();
                $(this).datagrid("selectRow",index);
                $('#mm').menu('show',{
                    left: e.clientX,
                    top: e.clientY
                });
            },
            onHeaderContextMenu:function(e, field){
                e.preventDefault();
                $('#mm1').menu('show',{
                    left: e.clientX,
                    top: e.clientY
                });
            }
        });

    })

    //下载礼包码
    function downLoadGiftCode() {
//        $.messager.confirm('Confirm','你确定要删除?',function(r){
//            if (r){
                var row = $("#datagrid").datagrid("getSelected");
                if(! row){
                    vac.alert("请选择要下载的礼包码");
                    return;
                }
                vac.alert(row.Id)
                vac.ajax(URL+'/downLoadGiftCode', {Id:row.Id}, 'POST', function(r){
                    if(r.status){
                        $("#datagrid").datagrid('reload');
                    }else{
                        vac.alert(r.info);
                    }
                })
//            }
//        });
    }

    //新增行
    function addrow(){
        var getRows = $("#datagrid").datagrid("getRows");

        //如果没有数据，则从0行开始新增
        if(vac.isEmpty(getRows)){
            var lenght = 0;
        }else{
            var lenght = getRows.length;
        }
        $("#datagrid").datagrid("appendRow",{Status:2});//插入
        $("#datagrid").datagrid("selectRow",lenght);//选中
        $("#datagrid").datagrid("beginEdit",lenght);//编辑输入
        var str_input =  $(".datagrid-view .datagrid-body>table>tbody>tr>td:nth-child(5)>div>table>tbody>tr>td");
        var str_input_start =  $(".datagrid-view .datagrid-body>table>tbody>tr>td:nth-child(6)>div>table>tbody>tr>td");
        var str_input_end =  $(".datagrid-view .datagrid-body>table>tbody>tr>td:nth-child(7)>div>table>tbody>tr>td");

        str_input.html('<div id="giftCode-upload" class="upload-box" data-api="/rbac/ueditor?action=uploadfile">' +
                '<div id="giftCode-fileList" class="uploader-list"></div>' +
                '<div id="giftCode-filePicker" class="upload-button">上传礼包码</div>' +
                '<input id="giftCode-src" name="Code" type="hidden" value=""><span class="hint"></span></div>')
        HJR2.Uploader.init();

        str_input_start.html('<input name="Starttime" type="text" value="" id="datetimepicker"/>');
        str_input_end.html('<input name="Endtime" type="text" value="" id="datetimepicker1"/>');
        $.datetimepicker.setLocale('ch');
        $('#datetimepicker').datetimepicker({
            format:"Y-m-d H:i:00",      //格式化日期
            step:10,
            yearStart:2010,     //设置最小年份
            yearEnd:2050,        //设置最大年份
        });

        $('#datetimepicker1').datetimepicker({
            format:"Y-m-d H:i:00",      //格式化日期
            step:10,
            yearStart:2010,     //设置最小年份
            yearEnd:2050,        //设置最大年份
        });

        jQuery(function() {
            jQuery('#datetimepicker').datetimepicker({
                onShow: function(ct) {
                    this.setOptions({
                        maxDate: jQuery('#datetimepicker1').val() ? jQuery('#datetimepicker1').val() : false
                    })

                }
            });

            jQuery('#datetimepicker1').datetimepicker({
                onShow: function(ct) {
                    this.setOptions({
                        minDate: jQuery('#datetimepicker').val() ? jQuery('#datetimepicker').val() : false
                    })
                }
            });
        });
    }

    function editrow(){
        if(!$("#datagrid").datagrid("getSelected")){
            vac.alert("请选择要编辑的行");
            return;
        }
        $('#datagrid').datagrid('beginEdit', vac.getindex("datagrid"));
        var str_input =  $(".datagrid-view .datagrid-body>table>tbody>tr>td:nth-child(5)>div>table>tbody>tr>td");
        var str_input_start =  $(".datagrid-view .datagrid-body>table>tbody>tr>td:nth-child(6)>div>table>tbody>tr>td");
        var str_input_end =  $(".datagrid-view .datagrid-body>table>tbody>tr>td:nth-child(7)>div>table>tbody>tr>td");

        str_input.html('<div id="giftCode-upload" class="upload-box" data-api="/rbac/ueditor?action=uploadfile">' +
                '<div id="giftCode-fileList" class="uploader-list"></div>' +
                '<div id="giftCode-filePicker" class="upload-button">上传礼包码</div>' +
                '<input id="giftCode-src" name="Code" type="hidden" value=""><span class="hint"></span></div>')
        HJR2.Uploader.init();

        str_input_start.html('<input name="Starttime" type="text" value="" id="datetimepicker"/>');
        str_input_end.html('<input name="Endtime" type="text" value="" id="datetimepicker1"/>');
        $.datetimepicker.setLocale('ch');
        $('#datetimepicker').datetimepicker({
            format:"Y-m-d H:i:00",      //格式化日期
            step:10,
            yearStart:2010,     //设置最小年份
            yearEnd:2050,        //设置最大年份
//            todayButton:false,    //关闭选择今天按钮
//            onSelectDate: function () {//选择完时间后执行
//                console.log($("#datetimepicker").val())
////                $("#Starttime")
//            },
//            onSelectTime:function(){
//                console.log($("#datetimepicker").val()) //选择完时分秒后执行
//            }
        });

        $('#datetimepicker1').datetimepicker({
            format:"Y-m-d H:i:00",      //格式化日期
            step:10,
            yearStart:2010,     //设置最小年份
            yearEnd:2050,        //设置最大年份
//            todayButton:false,    //关闭选择今天按钮
//            onSelectDate: function () {//选择完时间后执行
//                console.log($("#datetimepicker1").val())
//            },
//            onSelectTime:function(){
//                console.log($("#datetimepicker1").val()) //选择完时分秒后执行
//            }
        });

        jQuery(function() {
            jQuery('#datetimepicker').datetimepicker({
                onShow: function(ct) {
                    this.setOptions({
                        maxDate: jQuery('#datetimepicker1').val() ? jQuery('#datetimepicker1').val() : false
                    })

                }
            });

            jQuery('#datetimepicker1').datetimepicker({
                onShow: function(ct) {
                    this.setOptions({
                        minDate: jQuery('#datetimepicker').val() ? jQuery('#datetimepicker').val() : false
                    })
                }
            });
        });
    }


    function saverow(index){
        if(!$("#datagrid").datagrid("getSelected")){
            vac.alert("请选择要保存的行");
            return;
        }
        start_time = $("#datetimepicker").val();
        end_time = $("#datetimepicker1").val();
        Code = $("#giftCode-src").val();
        $('#datagrid').datagrid('endEdit', vac.getindex("datagrid"));
        var changes = $("#datagrid").datagrid("getSelected");
//        console.log("selected: ",changes);
        changes.Starttime = start_time;
        changes.Endtime = end_time;
        changes.Code = Code;
        if(changes.Id == undefined){
            changes.Id = 0;
        }
        vac.ajax(URL+'/AddOrUpdateGift?Code=' + Code, changes, 'POST', function(r){
            if(!r.status){
                vac.alert(r.info);
            }else{
                $("#datagrid").datagrid("reload");
            }
        })

    }
    //取消
    function cancelrow(){
        if(! $("#datagrid").datagrid("getSelected")){
            vac.alert("请选择要取消的行");
            return;
        }
        $("#datagrid").datagrid("cancelEdit",vac.getindex("datagrid"));
    }
    //刷新
    function reloadrow(){
        $("#datagrid").datagrid("reload");
    }

    //删除
    function delrow(){
        $.messager.confirm('Confirm','你确定要删除?',function(r){
            if (r){
                var row = $("#datagrid").datagrid("getSelected");
                if(! row){
                    vac.alert("请选择要删除的行");
                    return;
                }
                vac.ajax(URL+'/DelGift', {Id:row.Id}, 'POST', function(r){
                    if(r.status){
                        $("#datagrid").datagrid('reload');
                    }else{
                        vac.alert(r.info);
                    }
                })
            }
        });
    }
</script>
<script type="text/javascript" src="/static/easyui/js/webuploader.min.js"></script>
<script type="text/javascript" src="/static/easyui/js/webuploader_config.js"></script>
<body>
    <table id="datagrid" toolbar="#tb"></table>
    <div id="tb" style="padding:5px;height:auto">
        <a href="#" icon='icon-add' plain="true" onclick="addrow()" class="easyui-linkbutton" >新增</a>
        <a href="#" icon='icon-edit' plain="true" onclick="editrow()" class="easyui-linkbutton" >编辑</a>
        <a href="#" icon='icon-save' plain="true" onclick="saverow()" class="easyui-linkbutton" >保存</a>
        <a href="#" icon='icon-cancel' plain="true" onclick="delrow()" class="easyui-linkbutton" >删除</a>
        <a href="#" icon='icon-reload' plain="true" onclick="reloadrow()" class="easyui-linkbutton" >刷新</a>
    </div>
    <!--表格内的右键菜单-->
    <div id="mm" class="easyui-menu" style="width:120px;display: none" >
        <div iconCls='icon-add' onclick="addrow()">新增</div>
        <div iconCls="icon-edit" onclick="editrow()">编辑</div>
        <div iconCls='icon-save' onclick="saverow()">保存</div>
        <div iconCls='icon-cancel' onclick="cancelrow()">取消</div>
        <div class="menu-sep"></div>
        <div iconCls='icon-cancel' onclick="delrow()">删除</div>
        <div iconCls='icon-reload' onclick="reloadrow()">刷新</div>
        <div class="menu-sep"></div>
        <div>Exit</div>
    </div>
    <!--表头的右键菜单-->
    <div id="mm1" class="easyui-menu" style="width:120px;display: none"  >
        <div icon='icon-add' onclick="addrow()">新增</div>
    </div>
</body>
</html>