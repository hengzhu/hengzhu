{{template "../manage/header.tpl"}}

<script type="text/javascript" charset="utf-8" src="/static/ueditor/ueditor.config.js"></script>
<script type="text/javascript" charset="utf-8" src="/static/ueditor/ueditor.all.js"> </script>
<script type="text/javascript" charset="utf-8" src="/static/ueditor/lang/zh-cn/zh-cn.js"></script>
<script type="text/javascript">
    var subtypelist=$.parseJSON({{.newsTypes | stringsToJson}})
    var maintypelist = [
        {id:'1',text:'1'},
        {id:'2',text:'2'},
        {id:'3',text:'3'}
    ];
    var URL="/rbac/news";
    var newsId = 0;
    $(".datagrid-btable tbody tr").each(function () {
        $(this).children("td:first").attr("style","height:70px;")
    })
    $(function(){
        //新闻列表
        $("#datagrid").datagrid({
            title:'新闻管理',
            url:URL+"/index",
            method:'POST',
            pagination:true,
            fitColumns:true,
            striped:true,
            rownumbers:true,
            singleSelect:true,
            idField:'Id',
            pagination:true,
            pageSize:20,
            pageList:[10,20,30,50,100],
            columns:[[
                {field:'Id',title:'ID',width:35,align:'center'},
                {field:'Title',title:'标题',width:105,align:'center',editor:'text'},
                {field:'Maintype',title:'等级',width:35,align:'center',
                    formatter:function(value){
                        for(var i=0; i<maintypelist.length; i++){
                            if (maintypelist[i].id == value) return maintypelist[i].text;
                        }
                        return value;
                    },
                    editor:{
                        type:'combobox',
                        options:{
                            valueField:'id',
                            textField:'text',
                            data:maintypelist,
                            required:true
                        }
                    }
                },
                {field:'Subtype',title:'类型',width:35,align:'center',
                    formatter:function(value){
                        for(var i=0; i<subtypelist.length; i++){
                            if (subtypelist[i].Id == value) return subtypelist[i].Typename;
                        }
                        return value;
                    },
                    editor:{
                        type:'combobox',
                        options:{
                            valueField:'Id',
                            textField:'Typename',
                            data:subtypelist,
                            required:true
                        }
                    }
                },
                {field:'Sourcefrom',title:'文章来源',width:55,align:'center',editor:'text'},
                {field:'Author',title:'作者',width:55,align:'center',editor:'text'},
                {field:'Banner',title:'图片链接',width:70,align:'center',
                    formatter:function (value) {
                      return " <a href='" + value + "' target='_blank'><img src='" + value + "' style='width: 60px; height:60px;'/> </a>";
                    }/*,
                    editor:'text'*/},
//                {field:'Content',title:'内容',width:210,height:60,align:'center', editor:'text'},
                {field:'Updatetime',title:'更新时间',width:75,align:'center',
                    formatter:function(value,row,index){
                        if(value) return phpjs.date("Y-m-d H:i:s",phpjs.strtotime(value));
                        return value;
                    }
                },
                {field:'Createtime',title:'创建时间',width:75,align:'center',
                    formatter:function(value,row,index){
                        if(value) return phpjs.date("Y-m-d H:i:s",phpjs.strtotime(value));
                        return value;
                    }
                },
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
                vac.ajax(URL+'/AddAndEditNews', changes, 'POST', function(r){
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
        //创建添加新闻窗口
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
//                        if()
                        url:URL+'/AddAndEditNews?Id=' + newsId,

                        onSubmit:function(){
//                            vac.alert(flag);
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

    function loadData() {
        var datas = new Array();
        for(var i=0; i<subtypelist.length; i++){
            console.log(subtypelist[i].Id, subtypelist[i].Typename)
            datas.push({
                value: subtypelist[i].Id,
                text: subtypelist[i].Typename
            });
        }
        console.log(datas);
        //Reload the data
        $("#newsTypeSelect").combobox("loadData", datas);
    }

    //新增新闻
    function addrow(){
        newsId=0;
        loadData();
        HJR.Uploader.init();
        $("#dialog").dialog('open');
        $("#form1").form('clear');
    }

    function editrow(){
        var row = $("#datagrid").datagrid("getSelected");
        newsId=row.Id;
        if(!row){
            vac.alert("请选择要编辑的行");
            return;
        }
//        if(!$("#datagrid").datagrid("getSelected")){
//            vac.alert("请选择要编辑的行");
//            return;
//        }
//        $('#datagrid').datagrid('beginEdit', vac.getindex("datagrid"));

        loadData(row)
        $("#dialog").dialog('open');
//        $("#form1").form('clear');
        $("#form1").form('load',{
            Title:row.Title,
            Maintype:row.Maintype,
            Subtype:row.Subtype,
            Sourcefrom:row.Sourcefrom,
            Author:row.Author,
            Banner:row.Banner
        });
        ue.setContent(row.Content)
        HJR.Uploader.init();
    }

    function saverow(index){
        if(!$("#datagrid").datagrid("getSelected")){
            vac.alert("请选择要保存的行");
            return;
        }
        $('#datagrid').datagrid('endEdit', vac.getindex("datagrid"));
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
                vac.ajax(URL+'/DelNews', {Id:row.Id}, 'POST', function(r){
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
        <div iconCls='icon-save' onclick="saverow()">保存</div>.
        <!--<div iconCls='icon-cancel' onclick="cancelrow()">取消</div>-->
        <div class="menu-sep"></div>
        <div iconCls='icon-cancel' onclick="delrow()">删除</div>
        <div iconCls='icon-reload' onclick="reloadrow()">刷新</div>
        <div class="menu-sep"></div>
        <div>Exit</div>
    </div>
    <!--表头的右键菜单-->
    <!--<div id="mm1" class="easyui-menu" style="width:120px;display: none"  >
        <div icon='icon-add' onclick="addrow()">新增</div>
    </div>-->
    <div id="dialog" title="编辑新闻" style="width:900px;height:600px;">
        <div style="padding:20px 20px 40px 80px;" >
            <form id="form1" method="post">
                <table>
                    <tr>
                        <td style="width: 65px;">标&nbsp;&nbsp;题：</td>
                        <td colspan="3"><input id="_Title" value="" name="Title" type="text" class="easyui-validatebox" style="width: 534px;" required="true"/></td>
                    </tr>
                    <tr>
                        <td>等&nbsp;&nbsp;级：</td>
                        <td>
                            <select id="Maintype" name="Maintype"  style="width:158px;" class="easyui-combobox " data-options="value:3" editable="false" >
                                <option value="1" selected="selected">1</option>
                                <option value="2">2</option>
                                <option value="3">3</option>
                            </select>
                        </td>
                    <!--</tr>-->
                    <!--<tr>-->
                        <td style="width: 61px;">类&nbsp;&nbsp;型：</td>
                        <td>
                            <select id="newsTypeSelect" name="Subtype"  style="width:158px;" class="easyui-combobox " data-options="value:3" editable="false" >
                                <!--<option value="1" selected="selected">新闻</option>-->
                                <!--<option value="2">公告</option>-->
                                <!--<option value="3">活动</option>-->
                            </select>
                        </td>
                    </tr>
                    <tr>
                        <td>文章来源：</td>
                        <td><input id="Sourcefrom" name="Sourcefrom" class="easyui-validatebox"/>
                        </td>
                    <!--</tr>-->
                    <!--<tr>-->
                        <td>作&nbsp;&nbsp;者：</td>
                        <td><input id="Author" name="Author" class="easyui-validatebox"/>
                        </td>
                    </tr>
                    <tr>
                        <td>内&nbsp;&nbsp;容：</td>
                        <td colspan="3">
                            <div>
                                <script id="Content" type="text/plain" style="height:500px;" name="Content"></script>
                            </div>
                        </td>
                    </tr>

                    <tr>
                        <td>图片链接：</td>
                        <td colspan="3">
                            <div id="gameIcon-upload" class="upload-box" data-api="/rbac/ueditor?action=uploadimage">
                                <div id="gameIcon-fileList" class="uploader-list"></div>
                                <div id="gameIcon-filePicker" class="upload-button">选择图片</div>
                                <input id="gameIcon-src" name="Banner" type="hidden" value=""><span class="hint"></span>
                            </div>
                        </td>
                    </tr>
                </table>
            </form>
        </div>
    </div>
</body>


<script type="text/javascript" src="/static/easyui/js/webuploader.min.js"></script>
<script type="text/javascript" src="/static/easyui/js/webuploader_config.js"></script>
<script type="text/javascript">
    var ue = UE.getEditor('Content');

//    $(function () {
//        $('#edui149').css("z-index","9199");
//    });
</script>
</html>