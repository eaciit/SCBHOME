var dashboard = {}
viewModel.dashboard = dashboard

dashboard.newPageData = function () {
    var page = {
        Id: "",
        ProjectName: "",
        Description: "",
        Cover: "",
        Platforms: [],
    }

    dashboard.dataMasterPlatform().forEach(function (d) {
        var platform = dashboard.newPlatform()
        platform.PlatformId = d.Id
        platform.PlatformName = d.Name
        page.Platforms.push(platform)
    })

    return page
}
dashboard.newPlatform = function () {
    return {
        PlatformId :"", 
        PlatformName : "", 
        URL : "", 
        Username : "", 
        Password : "",
    }
}
dashboard.dataMasterPlatform = ko.observableArray([])
dashboard.modalTitle = ko.observable('Insert new data')
dashboard.page = ko.mapping.fromJS(dashboard.newPageData())

dashboard.getMasterPlatformData = function (callback) {
    viewModel.isLoading(true)

    ajaxPost('/web/dashboard/getmasterplatform', {}, function (res) {
        if (res.Status !== 'OK') {
            swal("Error!", res.Message, "error")
            return
        }

        dashboard.dataMasterPlatform(res.Data)

        callback()
    }, function () {
        viewModel.isLoading(false)
        swal("Error!", "Unknown error, please try again", "error")
    })
}

dashboard.dataPage = ko.observableArray([])
dashboard.getPageData = function (callback) {
    viewModel.isLoading(true)

    ajaxPost('/web/dashboard/getpage', {}, function (res) {
        if (res.Status !== 'OK') {
            swal("Error!", res.Message, "error")
            return
        }

        res.Data.forEach(function (d) {
            dashboard.dataMasterPlatform().forEach(function (e) {
                var row = d.Platforms.find(function (k) {
                    return k.PlatformId == e.Id
                })
                if (row === undefined) {
                    d.Platforms.push(dashboard.newPlatform())
                }
            })

            d.Platforms = _.sortBy(d.Platforms, 'PlatformName')
        })

        dashboard.dataPage(_.sortBy(res.Data, 'ProjectName'))

        callback()
    }, function () {
        viewModel.isLoading(false)
        swal("Error!", "Unknown error, please try again", "error")
    })
}

dashboard.renderGrid = function () {
    var data = dashboard.dataPage()

    if (typeof $('.grid').data('kendoGrid') !== 'undefined') {
        $('.grid').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,
            pageSize: 5
        }))
        return
    }

    var columns = [
        { field: 'ProjectName', title: 'Project Name', width: 200, attributes: { class: 'bold' } },
        { field: 'Description' },
    ]
    var columnPlatforms = dashboard.dataMasterPlatform().map(function (d) {
        return {
            width: 96,
            attributes: { class: 'align-center' },
            headerTemplate: function () {
                return '<div class="align-center">' + d.Name + '</div>'
            },
            template: function (e) {
                var row = e.Platforms.find(function (k) {
                    return k.PlatformId === d.Id
                })
                if (row !== undefined) {
                    if ($.trim(row.URL) !== '') {
                        return '<i class="fa fa-check"></i>'
                    }
                }

                return ''
            }
        }
    })
    columns = columns.concat([{ 
        headerTemplate: function () {
            return '<div class="align-center">Platform Availability</div>'
        },
        columns: columnPlatforms
    }, { 
        title: '&nbsp;', 
        width: 40, 
        template: function (d) {
            return '<input class="checkboxgrid" type="checkbox" data-id="'+d.Id+'"/>'
        }, 
        attributes: { class: 'align-center' }
    }])

    $('.grid').kendoGrid({
        dataSource: {
            data: data,
            pageSize: 5
        },
        pageable: true,
        sortable: true,
        filterable: true,
        columns: columns
    })
}

dashboard.checkedData = function () {
    return $('.checkboxgrid').get().filter(function (d) {
        return d.checked
    }).map(function (d) {
        return $(d).attr('data-id')
    })
}

dashboard.addMaster = function () {
    dashboard.modalTitle("Insert new data")
    ko.mapping.fromJS(dashboard.newPageData(), dashboard.page)
    $('#modal-page').modal('show')

    $('.uploader img').attr('src')
    $('#photo').replaceWith('<input type="file" name="ProjectPhoto" id="photo" />')
}

dashboard.editMaster = function () {
    if (dashboard.checkedData().length == 0 || dashboard.checkedData().length > 1) {
        swal("Error!", "Please choose one row", "error");
        return
    }

    dashboard.modalTitle('Edit data')

    var row = _.find(dashboard.dataPage(), { Id: dashboard.checkedData()[0] })
    ko.mapping.fromJS(row, dashboard.page)
    $('#modal-page').modal('show')

    $('.uploader img').attr('src', '/web/static/main/images/upload/' + row.Cover)
    $('#photo').replaceWith('<input type="file" name="ProjectPhoto" id="photo" />')
}

dashboard.deleteMaster = function(){
    if (dashboard.checkedData().length == 0) {
        swal("Error!", "Please choose at least one row", "error");
        return
    }

    swal({
        title: "Caution",
        text: "Are you sure want to delete this Data!",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#DD6B55",
        confirmButtonText: "Yes",
        closeOnConfirm: false
    }, function () {
        setTimeout(function () {
            var url = "/web/dashboard/deletepage";
            var param = { Ids: dashboard.checkedData() }
            ajaxPost(url, param, function(data) {
                if (data.Status == "OK") {
                    dashboard.refresh();
                    swal("Deleted!", "Your Data has been deleted.", "success");
                } else {
                    swal("Error!", data.Message, "error");
                }
            });
        }, 400)
    });
}

dashboard.getUploadedFiles = function () {
    var obj = document.getElementById('photo')
    if (obj === undefined) {
        return null
    }

    var file = obj.files[0]
    if (file === undefined) {
        return null
    }

    return file
}

dashboard.saveMaster = function () {
    var validator = $("#myForm").data("kendoValidator");
    if (validator === undefined) {
       validator = $("#myForm").kendoValidator().data("kendoValidator");
    }
    if (!validator.validate()) {
        return
    }

    var formData = new FormData()

    var payload = JSON.stringify(ko.mapping.toJS(dashboard.page))
    formData.append('data', payload)

    var file = dashboard.getUploadedFiles()
    if (file !== null) {
        formData.append('file', file)
    }

    viewModel.isLoading(true)
    $.ajax({
        url: "/web/dashboard/savepage",
        data: formData,
        contentType: false,
        dataType: "json",
        mimeType: 'multipart/form-data',
        processData: false,
        type: 'POST',
        success: function (data) {       
            viewModel.isLoading(false)

            if (data.Status == "OK") {
                swal("Saved!", "Your file has been successfully Update.", "success");
                dashboard.refresh();
                $('#modal-page').modal('hide');
            } else {
                swal("Error!", data.Message, "error");
            }
        },
        error: function () {
            viewModel.isLoading(true)
        }
    });
}

dashboard.refresh = function () {
    dashboard.getPageData(function () {
        setTimeout(function () {
            viewModel.isLoading(false)
            dashboard.renderGrid()
        }, 1000)
    })
}

dashboard.registerDragImagePreview = function () {
    var imageLoader = document.getElementById('photo')
    imageLoader.addEventListener('change', function (e) {
        var reader = new FileReader();
        reader.onload = function (event) {        
            $('.uploader img').attr('src', event.target.result);
        }
        reader.readAsDataURL(e.target.files[0]);
    }, false)

    $('.uploader img, .uploader span').on('click', function () {
        $('.uploader input').trigger('click')
    })
}

dashboard.init = function () {
    dashboard.registerDragImagePreview()
    dashboard.getMasterPlatformData(function () {
        dashboard.refresh()
    })
}

$(function () {
    dashboard.init()
})