{{template "../manage/header.tpl"}}
<script type="text/javascript">
    var URL = "/state";
    $(function () {
        $("#datagrid").datagrid({
            title: '状态',
            url: URL + "/index",
            method: 'POST',
            pagination: true,
            fitColumns: true,
            striped: true,
            rownumbers: true,
            singleSelect: true,
            idField: 'Id',
            columns: [[
                {
                    field: 'CabinetID', title: 'ID', width: 50, align: 'center',
                    formatter: function (value, row) {
                        if (row.IsOnline == "否") {
                            return "<span style='color: gray'>" + value + "</span>"
                        }
                        return value
                    }
                },
                {
                    field: 'IsOnline', title: '在线', width: 150, align: 'center',
                    formatter: function (value, row) {
                        if (row.IsOnline == "否") {
                            return "<span style='color: gray'>" + value + "</span>"
                        }
                        return value
                    }
                },
                {
                    field: 'Doors', title: '门数', width: 250, align: 'center',
                    formatter: function (value, row) {
                        if (row.IsOnline == "否") {
                            return "<span style='color: gray'>" + value + "</span>"
                        }
                        return value
                    }
                },
                {
                    field: 'OnUse', title: '使用情况', width: 100, align: 'center',
                    formatter: function (value, row) {
                        if (row.IsOnline == "否") {
                            return "<span style='color: gray'>" + value + "/" + row.Doors + "</span>"
                        }
                        return value + "/" + row.Doors
                    }

                },
                {
                    field: 'Close', title: '开关状态', width: 200, align: 'center',
                    formatter: function (value, row) {
                        if (row.IsOnline == "否") {
                            return "<span style='color: gray'>" + value + "/" + row.Doors + "</span>"
                        }
                        return value + "/" + row.Doors
                    }
                },
                {
                    field: 'TypeName', title: '类型', width: 200, align: 'center',
                    formatter: function (value, row) {
                        if (row.IsOnline == "否") {
                            return "<span style='color: gray'>" + value + "</span>"
                        }
                        return value
                    }
                },
                {
                    field: 'action', title: '查看', width: 200, align: 'center',
                    formatter:function(value,row){
                        var a = '<a href="/state/detail/'+row.Id+'" icon="icon-add" plain="true" class="easyui-linkbutton" target="_self">详情</a>'
                        return a;
                    }
                }
            ]],
        });

    })

    //创建状态详情窗口
    /*$("#dialog").dialog({
        modal:true,
        resizable:true,
        top:50,
        closed:true,
        buttons:[{
            text:'保存',
            iconCls:'icon-save',
            handler:function(){
                $("#form1").form('submit',{
                    url:URL+'/AddUser',
                    onSubmit:function(){
                        return $("#form1").form('validate');
                    },
                    success:function(r){
                        var r = $.parseJSON( r );
                        if(r.status){
                            $("#dialog").dialog("close");
                            $("#datagrid").datagrid('reload');
                        }else{
                            vac.alert(r.info);
                        }
                    }
                });
            }
        },{
            text:'取消',
            iconCls:'icon-cancel',
            handler:function(){
                $("#dialog").dialog("close");
            }
        }]
    });*/

    //状态详情弹窗
    // function detail(id){
    //     $("#dialog").dialog('open');
    //     $("#form1").form('clear');
    // }

    //刷新
    function reloadrow() {
        $("#datagrid").datagrid("reload");
    }
    
</script>
<body>
<table id="datagrid" toolbar="#tb"></table>
<div id="tb" style="padding:5px;height:auto">
    <a href="#" icon='icon-reload' plain="true" onclick="reloadrow()" class="easyui-linkbutton">刷新</a>
</div>
{{/*<div id="dialog" title="状态详情" style="width:400px;height:400px;">
    <div style="padding:20px 20px 40px 80px;" >
    <form id="form1" method="post">
        <table>
            <tr>
                <td>ID：</td>
                <td>{{.cabinet.CabinetID}}</td>
            </tr>
            <tr>
                <td>门数：</td>
                <td>{{.cabinet.Doors}}</td>
            </tr>
            <tr>
                <td>使用：</td>
                <td>{{.cabinet.OnUse}}/{{.cabinet.Doors}}</td>
            </tr>
            <tr
                <td>类型：</td>
                <td><input id="type" type="text" readonly value="{{.cabinet.TypeName}}"/>{{.cabinet.TypeName}}</td>
                <td><input type="button" onclick="$('#type').removeAttr('readonly')" value="修改"></td>
            </tr>
            <tr>
                <td>位置：</td>
                <td><input id="address" type="text" readonly value="{{.cabinet.TypeName}}"/>{{.cabinet.TypeName}}</td>
                <td><input type="button" onclick="$('#address').removeAttr('readonly')" value="修改"></td>
            </tr>
            <tr>
                <td>状态：</td>
                <td>
                    <select name="Status"  style="width:153px;" class="easyui-combobox " editable="false" required="true"  >
                        <option value="2" >启用</option>
                        <option value="1">禁用</option>
                    </select>
                </td>
            </tr>
            <tr>
                <td>备注：</td>
                <td><textarea name="Remark" class="easyui-validatebox" validType="length[0,200]"></textarea></td>
            </tr>
        </table>
    </form>
    </div>
</div>*/}}
</body>
</html>