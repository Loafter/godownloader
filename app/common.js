
        function UpdateTable() {
            $("#jqGrid")
                .jqGrid({
                    url: 'http://' + location.hostname + ':9981/progress.json',
                    mtype: "GET",
                    ajaxSubgridOptions: {
                        async: false
                    },
                    styleUI: 'Bootstrap',
                    datatype: "json",
                    colModel: [{
                        label: '#',
                        name: 'Id',
                        key: true,
                        width: 5
                    }, {
                        label: 'File Name',
                        name: 'FileName',
                        width: 15
                    }, {
                        label: 'Size',
                        name: 'Size',
                        width: 20,
                        formatter: FormatByte
                    }, {
                        label: 'Downloaded',
                        name: 'Downloaded',
                        width: 20,
                        formatter: FormatByte
                    }, {
                        label: '%',
                        name: 'Progress',
                        width: 5
                    }, {
                        label: 'Speed',
                        name: 'Speed',
                        width: 15,
                        formatter: FormatSpeedByte
                    }, {
                        label: 'Progress',
                        name: 'Progress',
                        formatter: FormatProgressBar
                    }],
                    viewrecords: true,
                    rowNum: 20,
                    pager: "#jqGridPager"
                });
        }

        function FixTable() {
            $.extend($.jgrid.ajaxOptions, {
                async: false
            })
            $("#jqGrid")
                .setGridWidth($(window)
                    .width() - 5)
            $("#jqGrid")
                .setGridHeight($(window)
                    .height())
            $(window)
                .bind('resize', function() {
                    $("#jqGrid")
                        .setGridWidth($(window)
                            .width() - 5);
                    $("#jqGrid")
                        .setGridHeight($(window)
                            .height())
                })
        }

        function UpdateData() {
            var grid = $("#jqGrid");
            var rowKey = grid.jqGrid('getGridParam', "selrow");
            $("#jqGrid").trigger("reloadGrid");
            if (rowKey) {
                $('#jqGrid').jqGrid("resetSelection")
                $('#jqGrid').jqGrid('setSelection', rowKey);
            }
        }

        function FormatProgressBar(cellValue, options, rowObject) {
            var intVal = parseInt(cellValue);

            var cellHtml = '<div class="progress"><div class="progress-bar" style="width: ' + intVal + '%;"></div></div>'

            return cellHtml;
        }

        function FormatByte(cellValue, options, rowObject) {
            var intVal = parseInt(cellValue);
            var ras = " B."
            if (intVal > 1024) {
                intVal /= 1024
                ras = " KB."
            }
            if (intVal > 1024) {
                intVal /= 1024
                ras = " MB."
            }
            if (intVal > 1024) {
                intVal /= 1024
                ras = " GB."
            }

            if (intVal > 1024) {
                intVal /= 1024
                ras = " TB."
            }
            var cellHtml = (intVal).toFixed(1) + ras;
            return cellHtml;
        }

        function FormatSpeedByte(cellValue, options, rowObject) {
            var intVal = parseInt(cellValue);
            var ras = " B/sec."
            if (intVal > 1024) {
                intVal /= 1024
                ras = " KB/sec."
            }
            if (intVal > 1024) {
                intVal /= 1024
                ras = " MB/sec."
            }
            if (intVal > 1024) {
                intVal /= 1024
                ras = " GB/sec"
            }

            if (intVal > 1024) {
                intVal /= 1024
                ras = " TB."
            }
            var cellHtml = (intVal).toFixed(1) + ras;
            return cellHtml;
        }

        function OnLoad() {

            UpdateTable()
            FixTable()
            setInterval(UpdateData, 500);
        }

        function AddDownload() {
            var req = {
                PartCount: parseInt($("#part_count_id").val()),
                FilePath: $("#save_path_id").val(),
                Url: $("#url_id").val()
            };
            $.ajax({
                    url: "/add_task",
                    type: "POST",
                    data: JSON.stringify(req),
                    dataType: "text"
                })
                .error(function(jsonData) {
                    console.log(jsonData)
                })
        }

        function RemoveDownload() {
            var grid = $("#jqGrid");
            var rowKey = parseInt(grid.jqGrid('getGridParam', "selrow"));
            var req = rowKey;
            $.ajax({
                    url: "/remove_task",
                    type: "POST",
                    data: JSON.stringify(req),
                    dataType: "text"
                })
                .error(function(jsonData) {
                    console.log(jsonData)
                })
        }

        function StartDownload() {
            var grid = $("#jqGrid");
            var rowKey = parseInt(grid.jqGrid('getGridParam', "selrow"));
            var req = rowKey;
            $.ajax({
                    url: "/start_task",
                    type: "POST",
                    data: JSON.stringify(req),
                    dataType: "text"
                })
                .error(function(jsonData) {
                    console.log(jsonData)
                })
        }

        function StopDownload() {
            var grid = $("#jqGrid");
            var rowKey = parseInt(grid.jqGrid('getGridParam', "selrow"));
            var req = rowKey;
            $.ajax({
                    url: "/stop_task",
                    type: "POST",
                    data: JSON.stringify(req),
                    dataType: "text"
                })
                .error(function(jsonData) {
                    console.log(jsonData)
                })
        }

        function StartAllDownload() {
            $.ajax({
                    url: "/start_all_task",
                    type: "POST",
                    dataType: "text"
                })
                .error(function(jsonData) {
                    console.log(jsonData)
                })
        }

        function StopAllDownload() {
            $.ajax({
                    url: "/stop_all_task",
                    type: "POST",
                    dataType: "text"
                })
                .error(function(jsonData) {
                    console.log(jsonData)
                })
        }


        function OnChangeUrl() {
            var filename = $("#url_id").val().split('/').pop()
            $("#save_path_id").val(filename)
        }
