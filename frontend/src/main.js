import Big from 'big.js';
import 'bulma/css/bulma.css';
import '@fortawesome/fontawesome-free/js/solid.js';
import '@fortawesome/fontawesome-free/js/fontawesome.js';
import $ from 'jquery';

// make jQuery global (for esbuild)
window.$ = $;
window.jQuery = $;

$(function() {
    MainModule.init();
    NavModule.init();
});

function initNext(handler) {
    const evts = {
        "entity": EntityModule,
        "form": FormModule,
        "home": HomeModule
    };
    evts[handler].init();
}

function showModalConfirm() {
    let modal = $("#modal-confirm");
    let modalBtn = $("#btn-modal-confirm");
    let modalCancel = $("#btn-modal-cancel");

    return new Promise((resolve, reject) => {
        modal.addClass("is-active");
        modalBtn.on("click", function() {
            modal.removeClass("is-active");
            resolve();
        });
        modalCancel.on("click", function() {
            modal.removeClass("is-active");
            reject();
        });
    });
}

const HomeModule = (function($) {
    function navHome(nColor, msg) {
        let showNotify = arguments.length === 2
        $.ajax({
            url: "/home/",
            method: "GET",
            dataType: "html",
            beforeSend: function(xhr) {
                // show loading indication
                MainModule.showProgress();
            }
        }).done(function(markup) {
            MainModule.endProgress();
            MainModule.setContent(markup);
            if (showNotify) {
                MainModule.notify(nColor, msg);
            }
        }).fail(function(xhr, status, err) {
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: /home/ ${xhr.responseText}`);
        });
    }

    function init() {
        return
    }

    function teardown() {
        return
    }

    return {
        navHome,
        init,
        teardown
    }
})(jQuery);

const MainModule = (function($) {
    let ele = {};

    function displayMessage(e, nColor, msg) {
        e.stopPropagation();
        e.preventDefault();
        ele.notify.removeClass("is-hidden");
        let markup = `
        <div class="notification is-${nColor} is-light">${msg}</div>
        `
        
        // 500ms animations with a 5 second delay between reveal and hide.
        ele.notify.html(markup).slideDown(500, function() {
          setTimeout(function() {
            ele.notify.slideUp(500, function() {
              ele.notify.empty();
              ele.notify.addClass("is-hidden");
            });
          }, 5000);
        });
    }

    function showProgress() {
        let markup = `<progress class="progress is-small" max="100"></progress>`
        ele.notify.html(markup);
    }

    function endProgress() {
        ele.notify.empty();
    }

    return {
        init() {
            ele.app = $("#app");
            ele.notify = $("#notification");

            ele.notify.on("may:notify", displayMessage);
        },
        setContent(content) {
            ele.app.empty();
            ele.app.append(content);
        },
        notify(nColor, msg) {
            ele.notify.trigger("may:notify", [nColor, msg]);
        },
        showProgress,
        endProgress
    }
})(jQuery);

const NavModule = (function($) {
    let ele = {
        navbar: $("nav")
    };

    function onClick(e) {
        e.stopPropagation();
        e.preventDefault();
        e.target.blur();
        let route = $(this).attr("href");
        let handler = $(this).data("handler");
        $.ajax({
            url: `/${route}/`,
            method: "GET",
            dataType: "html",
            beforeSend: function(xhr) {
                // no loader for nav
                return
            },
        }).done(function(markup) {
            MainModule.setContent(markup);
            initNext(handler);
        }).fail(function(xhr, status, err) {
            MainModule.notify("danger", `request failure: ${route} ${xhr.responseText}`);
        });
    }

    return {
        init() {
            ele.navbar.on("click", "a", onClick);
        },
        teardown() {
            ele.navbar.off("click", "a", onClick);
        }
    }
})(jQuery);

const EntityModule = (function($) {
    let ele = {};

    function filterTable() {
        let val = ele.searchInput.val().toLowerCase();
        ele.eList.filter(function() {
            $(this).toggle($(this).text().toLowerCase().indexOf(val) > -1);
        });
    }

    function handleUpdate(e) {
        e.stopPropagation();
        e.preventDefault();
        let ent = ele.btnNew.attr("href");
        let entId = $(this).find("td:eq(0)").text();
        makeRequest(ent, entId, "GET");
    }

    function handleNew(e) {
        e.stopPropagation();
        e.preventDefault();
        let ent = ele.btnNew.attr("href");
        // id=0 will request a blank form with defaults filled in
        makeRequest(ent, "0", "GET");
    }

    function makeRequest(ent, id, reqMethod) {
        let handler = ele.tBody.data("handler");
        $.ajax({
            url: `/${ent}/${id}/`,
            method: reqMethod,
            dataType: "html",
            beforeSend: function(xhr) {
                // show loading indication
                MainModule.showProgress();
            }
        }).done(function(markup) {
            MainModule.endProgress();
            teardown();
            MainModule.setContent(markup);
            initNext(handler);
        }).fail(function(xhr, status, err) {
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: /${ent}/${id}/ ${xhr.responseText}`);
        });
    }

    function init() {
        // setup events for click on table row, add new button, and cancel button
        ele.searchInput = $("#search-input");
        ele.btnNew = $("#btn-new");
        ele.tBody = $("tbody");
        ele.eList = $("tbody tr");

        ele.searchInput.on("keyup", filterTable);
        ele.tBody.on("click", "tr", handleUpdate);
        ele.btnNew.on("click", handleNew);
    }

    function teardown() {
        ele.searchInput.off("keyup", filterTable);
        ele.tBody.off("click", "tr", handleUpdate);
        ele.btnNew.off("click", handleNew);
    }

    return {
        init,
        teardown
    }
})(jQuery);

const FormModule = (function($) {
    let ele = {};

    function clearDropdown(e) {
        e.stopPropagation();
        e.preventDefault();
        let dropdown = $(this).data("select-id");
        $(`#${dropdown}`).prop("selectedIndex", 0);
    }

    function handleUpdate(e) {
        e.stopPropagation();
        e.preventDefault();
        if (!checkValidForm()) {
            return
        }
        let ent = ele.btnSubmit.attr("href");
        let id = parseInt(ele.entId.val(), 10);
        if (id === 0) {
            makeRequest(`/${ent}/`, "POST")
        } else {
            makeRequest(`/${ent}/${id}/`, "PUT");
        }
    }

    function handleDelete(e) {
        e.stopPropagation();
        e.preventDefault();
        if (!checkValidForm()) {
            return
        }
        ele.fieldset.prop("disabled", true);
        showModalConfirm().then(() => {
            let ent = ele.btnSubmit.attr("href");
            let id = parseInt(ele.entId.val(), 10);
            makeRequest(`/${ent}/${id}/`, "DELETE");
        }).catch(() => {
            ele.fieldset.prop("disabled", false);
        });
    }

    function handleCancel(e) {
        e.stopPropagation();
        e.preventDefault();
        HomeModule.navHome();
    }

    function checkValidField(e) {
        e.stopPropagation();
        let ele = $(this);
        let valid = ele.get(0).checkValidity();
        if (!valid) {
            MainModule.notify("danger", "invalid form fields");
            ele.addClass("is-danger");
        } else {
            ele.removeClass("is-danger");
        }
    }

    function checkValidForm() {
        ele.formInput.each(function() {
            let ele = $(this);
            let valid = ele.get(0).checkValidity();
            if (!valid) {
                MainModule.notify("danger", "invalid form fields");
                ele.addClass("is-danger");
            } else {
                ele.removeClass("is-danger");
            }
        });

        // delete button should be only element with is-danger class
        return $(".is-danger").length === 1;
    }

    function makeRequest(url, reqMethod) {
        let formData = ele.form.serialize();
        $.ajax({
            url: url,
            method: reqMethod,
            data: formData,
            dataType: "text",
            beforeSend: function() {
                // show loading indication
                MainModule.showProgress();
                ele.fieldset.prop("disabled", true);
            },
        }).done(function(res) {
            MainModule.endProgress();
            HomeModule.navHome("primary", res);
        }).fail(function(xhr, status, err) {
            ele.fieldset.prop("disabled", false);
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    return {
        init() {
            ele.btnClear = $("button.clear");
            ele.btnSubmit = $("#btn-submit");
            ele.btnDelete = $("#btn-delete");
            ele.btnCancel = $("#btn-cancel");
            ele.fieldset = $("fieldset");
            ele.form = $("form");
            ele.formInput = $("div.control input");
            ele.entId = $("#entity-id");

            ele.btnClear.on("click", clearDropdown);
            ele.btnSubmit.on("click", handleUpdate);
            ele.btnCancel.on("click", handleCancel);
            ele.form.on("blur", "input", checkValidField);

            let id = parseInt(ele.entId.val(), 10);
            if (id === 0) {
                ele.btnDelete.addClass("is-hidden");
            } else {
                ele.btnDelete.on("click", handleDelete);
            }
        },
        teardown() {
            ele.btnClear.off("click", clearDropdown);
            ele.btnDelete.off("click", handleDelete);
            ele.btnCancel.off("click", handleCancel);
            ele.btnSubmit.off("click", handleUpdate);
            ele.form.off("blur", "input", checkValidField);
        }
    }
})(jQuery);
