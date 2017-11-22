{{template "../manage/header.tpl"}}

<script type="text/javascript">
var gametypes=$.parseJSON({{.gametypes | stringsToJson}})
var systemlist = [
    {id:'1',text:'android'},
    {id:'2',text:'ios'}
];

var URL="/rbac/game";
$(function(){
    //用户列表
    $("#datagrid").datagrid({
        title:'游戏列表',
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
            {field:"Gamename",title:"游戏名",align:"center",editor:"text"},
            {field:"GameType",title:"游戏类型",formatter:function(value){
                    for(var i=0; i<gametypes.length; i++){
                        if (gametypes[i].Id == value){
                            return gametypes[i].Typename;
                        } 
                    }
                    return value;
                },
                editor:{
                    type:'combobox',
                    options:{
                        valueField:'Id',
                        textField:'Typename',
                        data:gametypes,
                        required:true
                    }
                }
            },
            {field:"Size",title:"游戏大小",align:"center",editor:"numberbox"},
            {field:"System",title:"操作系统",align:"center",
                formatter:function(value){
                    for(var i=0; i<systemlist.length; i++){
                        if (systemlist[i].id == value) return systemlist[i].text;
                    }
                    return value;
                },
                editor:{
                    type:'combobox',
                    options:{
                        valueField:'id',
                        textField:'text',
                        data:systemlist,
                        required:true
                    }
                }
            },
            {field:"Version",title:"版本",align:"center",editor:"text"},
            {field:"Icon",title:"Icon",align:"center",
                formatter:function (value) {
                          return " <a href='" + value + "' target='_blank'><img src='" + value + "' style='width: 60px; height:60px;'/> </a>";
                },
                editor: "text"
            },
            {field:"Site",title:"官网地址",align:"center",editor:"text"},
            {field:"Download",title:"下载链接",align:"center",editor:"text"},
            {field:"Desc",title:"描述信息",align:"center",editor:"text"}

        ]],
        // onAfterEdit:function(index, data, changes){
        //     new_icon = $("#gameIcon-src").val();
        //     alert("new_icon:", new_icon)
        //     alert("old_icon:", data.Icon)
        //     if(vac.isEmpty(changes)){
        //         return;
        //     }
        //     changes.Id = data.Id;
        //     alert(changes.Icon);
        //     vac.ajax(URL+'/UpdateGame', changes, 'POST', function(r){
        //         if(!r.status){
        //             vac.alert(r.info);
        //         }else{
        //             $("#datagrid").datagrid("reload");
        //         }
        //     })
        // },
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

function editrow(S){
    if(!$("#datagrid").datagrid("getSelected")){
        vac.alert("请选择要编辑的行");
        return;
    }
    $('#datagrid').datagrid('beginEdit', vac.getindex("datagrid"));
    var str_input =  $(".datagrid-view .datagrid-body>table>tbody>tr>td:nth-child(7)>div>table>tbody>tr>td");
    // console.log(str_input.html())
    str_input.html('<div id="gameIcon-upload" class="upload-box" data-api="/rbac/ueditor?action=uploadimage"><div id="gameIcon-fileList" class="uploader-list"></div><div id="gameIcon-filePicker" class="upload-button">选择图片</div><input id="gameIcon-src" name="Icon" type="hidden" value=""><span class="hint"></span></div>')
    HJR.Uploader.init();
    // str_input.prop("type","button");
    // str_input.removeClass("datagrid-editable-input")
    // str_input.val("更换Icon")
    // str_input.css("height", "20px")
    // console.log(str_input)
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
    var str_input =  $(".datagrid-view .datagrid-body>table>tbody>tr>td:nth-child(7)>div>table>tbody>tr>td");
    // console.log(str_input.html())
    str_input.html('<div id="gameIcon-upload" class="upload-box" data-api="/rbac/ueditor?action=uploadimage"><div id="gameIcon-fileList" class="uploader-list"></div><div id="gameIcon-filePicker" class="upload-button">选择图片</div><input id="gameIcon-src" name="Icon" type="hidden" value=""><span class="hint"></span></div>')
    HJR.Uploader.init();
}

function saverow(index){
    if(!$("#datagrid").datagrid("getSelected")){
        vac.alert("请选择要保存的行");
        return;
    }
    wins = $("#gameIcon-src").val();
    console.log(wins);
    // var testData = $("#datagrid").datagrid('getChanges');
    // console.log("Data:", testData)
    $('#datagrid').datagrid('endEdit', vac.getindex("datagrid"));
    var changes = $("#datagrid").datagrid("getSelected");
    console.log("selected: ",changes);
    changes.Icon = wins
    if(changes.Id == undefined){
        changes.Id = 0;
    }
    vac.ajax(URL+'/AddAndEdit', changes, 'POST', function(r){
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

function loadData() {
    var datas = new Array();
    for(var i=0; i<gametypes.length; i++){
        console.log(gametypes[i].Id, gametypes[i].Typename)
        datas.push({
            value: gametypes[i].Id,
            text: gametypes[i].Typename
        });
    }
    console.log(datas);
    //Reload the data
    $("#gameTypeSelect").combobox("loadData", datas);
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
            vac.ajax(URL+'/DelGame', {Id:row.Id}, 'POST', function(r){
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
</body>
<script type="text/javascript" src="/static/easyui/js/webuploader.min.js"></script>
<script type="text/javascript" src="/static/easyui/js/webuploader_config.js"></script>
</html>