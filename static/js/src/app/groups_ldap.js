var ldap_servers_loaded = false;

function loadLdapServers() {
    api.LDAP.get()
        .success(function (ls) {
            if (ls.length === 0) {
                $('#ldap_select_group').hide()
                $('#no_ldap').show()
                return false
            } else {
                let ldap_s2 = $.map(ls, function (obj) {
                    obj.text = obj.name + " -- " + obj.protocol + "://" + obj.host
                    return obj
                });
                let ldap_select = $("#ldap_servers")
                ldap_select.select2({
                    placeholder: "Select a server",
                    data: ldap_s2,
                });
                if (ls.length === 1) {
                    ldap_select.val(ldap_s2[0].id)
                    ldap_select.trigger('change.select2')
                }
            }
        });
}

function importLdapUsers(ldap_id) {
    /*targets = $("#targetsTable").dataTable({
        destroy: true, // Destroy any other instantiated table - http://datatables.net/manual/tech-notes/3#destroy
        columnDefs: [{
            orderable: false,
            targets: "no-sort"
        }]
    })*/
    api.import.ldap(ldap_id)
        .success(function (ts) {
            console.log(ts)
            $.each(ts, function (i, record) {
                addTarget(
                    record.first_name,
                    record.last_name,
                    record.email,
                    record.position);
            });
            targets.DataTable().draw();
        })
        .error(function (err) {
            modalError(err.responseJSON.message)
        })
}


$(document).ready(function () {
    $('#ldap_btn').click(function(e) {
        e.preventDefault()
        let ldap_form = $('#ldap_form')
        if(ldap_form.is(":visible")) {
            ldap_form.hide()
        } else {
            ldap_form.show()
            if(!ldap_servers_loaded) {
                loadLdapServers();
                ldap_servers_loaded = true
            }
        }
    })

    $('#ldap_import_btn').click(function(e) {
        e.preventDefault()
        $("#modal\\.flashes").empty()
        let btn = $(this).button('loading')
        const ldap_id = $('#ldap_servers').val()
        if(ldap_id) {
            importLdapUsers(ldap_id)
        } else {
            modalError("Please select a server")
        }
        btn.button('reset')
    })
})
