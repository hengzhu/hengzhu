{{template "../manage/header.tpl"}}

<script type="text/javascript">
var URL="/rbac/gametype";
$(function(){
    //用户列表
    $("#datagrid").datagrid({
        title:'分组列表',
        url:URL+'/index',
        method:'POST',
        pagination:true,
        fitColumns:true,
        striped:true,
        rownumbers:true,
        singleSelect:true,
        idField:'id',
        pageSize:20,
        pageList:[10,20,30,50,100],
        columns:[[
            {field:'Id',title:'ID'},
            {field:"Typename",title:"类型名称",align:"center",editor:"text"},
            {field:"Desc",title:"描述信息",align:"center",editor:"text"}

        ]],
        onAfterEdit:function(index, data, changes){
             if(vac.isEmpty(changes)){
                return;
            }
            if(data.Id == undefined){
                changes.Id = 0;
            }else{
                changes.Id = data.Id;
            }
            vac.ajax(URL+'/AddAndEdit', changes, 'POST', function(r){
                if(!r.status){
                    vac.alert(r.info);
                }else{
                    $("#datagrid").datagrid("reload");
                }
            })
        },
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

    $("#dialog").dialog({
        modal:true,
        resizable:true,
        top:150,
        closed:true,
        buttons:[{
            text:'保存',
            iconCls:'icon-save',
            handler:function(){
                $("#form1").form('submit',{
                    url:URL+'/AddGameType',
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
    });

})
function editrow(){
    if(!$("#datagrid").datagrid("getSelected")){
        vac.alert("请选择要编辑的行");
        return;
    }
    $('#datagrid').datagrid('beginEdit', vac.getindex("datagrid"));
}
function saverow(index){
    if(!$("#datagrid").datagrid("getSelected")){
        vac.alert("请选择要保存的行");
        return;
    }
    $('#datagrid').datagrid('endEdit', vac.getindex("datagrid"));
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

function addrow(){
    var getRows = $("#datagrid").datagrid("getRows");

    //如果没有数据，则从0行开始新增
    if(vac.isEmpty(getRows)){
        var length = 0;
    }else{
        var length = getRows.length;
    }
    $("#datagrid").datagrid("appendRow",{Status:2});//插入
    $("#datagrid").datagrid("selectRow",length);//选中
    $("#datagrid").datagrid("beginEdit",length);//编辑输入
}

//删除
function delrow(){
    $.messager.confirm('Confirm','你确定要删除?',function(r){
        if (r){
            var row = $("#datagrid").datagrid("getSelected");
            if(!row){
                vac.alert("请选择要删除的行");
                return;
            }
            vac.ajax(URL+'/DelGameType', {Id:row.Id}, 'POST', function(r){
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
<!-- <div id="dialog" title="添加游戏类型" style="width:400px;height:400px;">
    <div style="padding:20px 20px 40px 80px;" >
        <form id="form1" method="post">
            <table>
                <tr>
                    <td>类型名称：</td>
                    <td><input name="Typename" class="easyui-validatebox" required="true"/></td>
                </tr>
                <tr>
                	<td>描述：</td>
                	<td><input name="Desc" class="easyui-validatebox"/></td>
                </tr>
            </table>
        </form>
    </div>
</div> -->
</body>
</html>