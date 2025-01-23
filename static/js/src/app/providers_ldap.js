var ldap_servers = []

// Save attempts to POST to /providers/ldap/
function saveLdap(idx) {
    let server = {}
    server.name = $("#ldap_name").val()
    server.protocol = $("#ldap_protocol").val()
    server.host = $("#ldap_host").val()
    server.username = $("#ldap_username").val()
    server.password = $("#ldap_password").val()
    server.base_dn = $("#ldap_base_dn").val()
    server.query = $("#ldap_query").val()
    server.attributes = $("#ldap_attributes").val()
    server.ignore_cert_errors = $("#ldap_ignore_cert_errors").prop("checked")
    if (idx != -1) {
        server.id = ldap_servers[idx].id
        api.LDAPId.put(server)
            .success(function (data) {
                successFlash("Server edited successfully!")
                loadLdaps()
                dismissLdap()
            })
            .error(function (data) {
                modalError(data.responseJSON.message)
            })
    } else {
        // Submit the server
        api.LDAP.post(server)
            .success(function (data) {
                successFlash("Server added successfully!")
                loadLdaps()
                dismissLdap()
            })
            .error(function (data) {
                modalError(data.responseJSON.message)
            })
    }
}

function dismissLdap() {
    $("#ldap_modal\\.flashes").empty()
    $("#ldap_name").val("")
    $("#ldap_protocol").val("ldap")
    $("#ldap_host").val("")
    $("#ldap_username").val("")
    $("#ldap_password").val("")
    $("#ldap_base_dn").val("")
    $("#ldap_query").val("")
    $("#ldap_attributes").val("")
    $("#ldap_ignore_cert_errors").prop("checked", true)
    $("#ldap_modal").modal('hide')
}

var deleteLdap = function (idx) {
    Swal.fire({
        title: "Are you sure?",
        text: "This will delete the LDAP config. This can't be undone!",
        type: "warning",
        animation: false,
        showCancelButton: true,
        confirmButtonText: "Delete " + escapeHtml(ldap_servers[idx].name),
        confirmButtonColor: "#428bca",
        reverseButtons: true,
        allowOutsideClick: false,
        preConfirm: function () {
            return new Promise(function (resolve, reject) {
                api.LDAPId.delete(ldap_servers[idx].id)
                    .success(function (msg) {
                        resolve()
                    })
                    .error(function (data) {
                        reject(data.responseJSON.message)
                    })
            })
        }
    }).then(function (result) {
        if (result.value){
            Swal.fire(
                'LDAP config Deleted!',
                'This LDAP config has been deleted!',
                'success'
            );
        }
        $('button:contains("OK")').on('click', function () {
            location.reload()
        })
    })
}

function editLdap(idx) {
    $("#ldap_modalSubmit").unbind('click').click(function () {
        saveLdap(idx)
    })
    var ldap = {}
    if (idx != -1) {
        $("#ldapModalLabel").text("Edit LDAP config")
        let server = ldap_servers[idx]
        $("#ldap_name").val(server.name)
        $("#ldap_protocol").val(server.protocol)
        $("#ldap_host").val(server.host)
        $("#ldap_username").val(server.username)
        $("#ldap_password").val(server.password)
        $("#ldap_base_dn").val(server.base_dn)
        $("#ldap_query").val(server.query)
        $("#ldap_attributes").val(server.attributes)
        $("#ldap_ignore_cert_errors").prop("checked", server.ignore_cert_errors)
    } else {
        $("#ldapModalLabel").text("New LDAP config")
    }
}

function loadLdaps() {
    $("#ldapTable").hide()
    $("#ldap_emptyMessage").hide()
    $("#ldap_loading").show()
    api.LDAP.get()
        .success(function (ls) {
            ldap_servers = ls
            $("#ldap_loading").hide()
            if (ldap_servers.length > 0) {
                $("#ldapTable").show()
                ldapTable = $("#ldapTable").DataTable({
                    destroy: true,
                    columnDefs: [{
                        orderable: false,
                        targets: "no-sort"
                    }]
                });
                ldapTable.clear()
                ldapRows = []
                $.each(ldap_servers, function (i, server) {
                    ldapRows.push([
                        escapeHtml(server.name),
                        escapeHtml(server.host),
                        moment(server.modified_date).format('MMMM Do YYYY, h:mm:ss a'),
                        "<div class='pull-right'><span data-toggle='modal' data-backdrop='static' data-target='#ldap_modal'><button class='btn btn-primary' data-toggle='tooltip' data-placement='left' title='Edit LDAP config' onclick='editLdap(" + i + ")'>\
                    <i class='fa fa-pencil'></i>\
                    </button></span>\
                    <button class='btn btn-danger' data-toggle='tooltip' data-placement='left' title='Delete LDAP config' onclick='deleteLdap(" + i + ")'>\
                    <i class='fa fa-trash-o'></i>\
                    </button></div>"
                    ])
                })
                ldapTable.rows.add(ldapRows).draw()
                $('[data-toggle="tooltip"]').tooltip()
            } else {
                $("#ldap_emptyMessage").show()
            }
        })
        .error(function () {
            $("#ldap_loading").hide()
            errorFlash("Error fetching LDAP configs")
        })
}

$(document).ready(function () {
    $('#ldap_modal').on('hidden.bs.modal', function (event) {
        dismissLdap()
    });
    loadLdaps()
})