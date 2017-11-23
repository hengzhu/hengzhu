/*** 上传相关 ***/
var HJR = {};
(function(S) {
    var Uploader = {
        init: function() {
            /**上传图片 **/
            this._upload('#gameIcon');
        },
        _upload:function(id){
            var api ,picker ,list,res;
            api = $(id+'-upload').data('api');
            picker = id+'-filePicker';
            list =  $(id+'-fileList');
            res =  $(id+'-src');
            var thumbnailWidth = 60,
                thumbnailHeight = 60;

            // Web Uploader实例
            var uploader = WebUploader.create({
                auto: true,
                swf: '/static/easyui/js/Uploader.swf',
                server: api,
                pick: picker,
                fileVal:'upfile',
                accept: {
                    title: 'image',
                    extensions: 'gif,jpg,jpeg,png',
                    mimeTypes: 'image/*'
                },
                resize:false
            });

            // 当有文件添加进来的时候
            uploader.on('fileQueued', function(file) {
                var $li = $(
                        '<div id="' + file.id + '" class="file-item thumbnail">' +
                        '<img>' +
                        '</div>'
                    ),
                    $img = $li.find('img');

                list.html($li);
                res.val('');
                // 创建缩略图
                uploader.makeThumb(file, function(error, src) {
                    if (error) {
                        $img.replaceWith('<span>不能预览</span>');
                        return;
                    }

                    $img.attr('src', src);

                }, thumbnailWidth, thumbnailHeight);
            });

            // 文件上传过程中创建进度条实时显示。
            uploader.on('uploadProgress', function(file, percentage) {
                var $li = $('#' + file.id),
                    $percent = $li.find('.progress span');

                // 避免重复创建
                if (!$percent.length) {
                    $percent = $('<p class="progress"><span></span></p>')
                        .appendTo($li)
                        .find('span');
                }

                $percent.css('width', percentage * 100 + '%');
            });

            // 文件上传成功，给item添加成功class, 用样式标记上传成功。
            uploader.on('uploadSuccess', function(file, response) {
            	if(response.state=='SUCCESS'){
                $('#' + file.id).addClass('upload-state-done');
                res.val(response.url);
                }
            });

            // 文件上传失败，现实上传出错。
            uploader.on('uploadError', function(file) {
                var $li = $('#' + file.id),
                    $error = $li.find('div.error');

                // 避免重复创建
                if (!$error.length) {
                    $error = $('<div class="error"></div>').appendTo($li);
                }

                $error.text('上传失败');
            });

            // 完成上传完了，成功或者失败，先删除进度条。
            uploader.on('uploadComplete', function(file) {
                $('#' + file.id).find('.progress').remove();

            });
        },
    }

    // $(document).ready(function() {
    //     Uploader.init();
    // });

    return S.Uploader = Uploader;

})(HJR);

/*** 上传相关 ***/
var HJR2 = {};
(function(T) {
    var Uploader = {
        init: function() {
            /**上传excel文件 **/
            this._upload('#giftCode');
        },
        _upload:function(id){
            var api ,picker ,list,res;
            api = $(id+'-upload').data('api');
            picker = id+'-filePicker';
            list =  $(id+'-fileList');
            res =  $(id+'-src');
            // var thumbnailWidth = 60,
            //     thumbnailHeight = 60;

            // Web Uploader实例
            var uploader = WebUploader.create({
                auto: true,
                swf: '/static/easyui/js/Uploader.swf',
                server: api,
                pick: picker,
                fileVal:'upfile',
                accept: {
                    title: 'file',
                    extensions: 'xlsx,xls',
                    mimeTypes: 'file/*'
                },
                resize:false
            });

            // // 当有文件添加进来的时候
            // uploader.on('fileQueued', function(file) {
            //     var $li = $(
            //             '<div id="' + file.id + '" class="file-item thumbnail">' +
            //             '<img>' +
            //             '</div>'
            //         ),
            //         $img = $li.find('img');
            //
            //     list.html($li);
            //     res.val('');
            //     创建缩略图
            //     uploader.makeThumb(file, function(error, src) {
            //         if (error) {
            //             $img.replaceWith('<span>不能预览</span>');
            //             return;
            //         }
            //
            //         $img.attr('src', src);
            //
            //     }, thumbnailWidth, thumbnailHeight);
            // });

            // 文件上传过程中创建进度条实时显示。
            uploader.on('uploadProgress', function(file, percentage) {
                var $li = $('#' + file.id),
                    $percent = $li.find('.progress span');

                // 避免重复创建
                if (!$percent.length) {
                    $percent = $('<p class="progress"><span></span></p>')
                        .appendTo($li)
                        .find('span');
                }

                $percent.css('width', percentage * 100 + '%');
            });

            // 文件上传成功，给item添加成功class, 用样式标记上传成功。
            uploader.on('uploadSuccess', function(file, response) {
            	if(response.state=='SUCCESS'){
                $('#' + file.id).addClass('upload-state-done');
                res.val(response.url);
                }
            });

            // 文件上传失败，现实上传出错。
            uploader.on('uploadError', function(file) {
                var $li = $('#' + file.id),
                    $error = $li.find('div.error');

                // 避免重复创建
                if (!$error.length) {
                    $error = $('<div class="error"></div>').appendTo($li);
                }

                $error.text('上传失败');
            });

            // 完成上传完了，成功或者失败，先删除进度条。
            uploader.on('uploadComplete', function(file) {
                $('#' + file.id).find('.progress').remove();

            });
        },
    }

    // $(document).ready(function() {
    //     Uploader.init();
    // });

    return T.Uploader = Uploader;

})(HJR2);