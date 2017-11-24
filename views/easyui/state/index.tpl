{{template "../manage/header.tpl"}}
<script type="text/javascript">
    var URL = "/state";
    $(function () {
        //角色列表
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
                {field: 'CabinetID', title: 'ID', width: 50, align: 'center'},
                {field: 'isonline', title: '组名', width: 150, align: 'center'},
                {field: 'Remark', title: '描述', width: 250, align: 'center'},
                {field: 'Status', title: '状态', width: 100, align: 'center'},
                {field: 'action', title: '操作', width: 200, align: 'center'}
            ]],
        });
    })
</script>
<body>

<table id="datagrid" toolbar="#tb"></table>
{{print .data}}
{{range .data.rows}}
<div class="col-md-2">
    <p>{{.Gamename}}</p>
    <div>
        <ul class="all_games_ul">
            <li><a href="{{.Download}}">下载</a></li>
            <li><a href="static/easyui/page/gift/gift.html">礼包</a></li>
        </ul>
    </div>
</div>
{{end}}
</body>
</html>