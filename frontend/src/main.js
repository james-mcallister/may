import Big from 'big.js';
import 'bulma/css/bulma.css';
import '@fortawesome/fontawesome-free/js/solid.js';
import '@fortawesome/fontawesome-free/js/fontawesome.js';
import $ from 'jquery';

// make jQuery global (for esbuild)
window.$ = $;
window.jQuery = $;

// Big.js settings
Big.DP = 2; // max decimal precision
Big.RM = Big.roundHalfUp; // rounding mode
// ex: x = Big("10.555555555");
// x.round(2) // 10.56

$(function() {
    MainModule.init();
    NavModule.init();
});

function initNext(handler) {
    const evts = {
        "entity": EntityModule,
        "form": FormModule,
        "home": HomeModule,
        "plan": PlanFormModule
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
            teardown();
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

const PlanFormModule = (function($) {
    let ele = {};

    function handleCheck() {
        if (ele.newCheck.prop("checked")) {
          ele.typeSelect.prop("disabled", true);
        } else {
          ele.typeSelect.prop("disabled", false);
        }
    }

    function handlePlanType() {
        let selected = ele.selectPlanner.val();
        if (selected == "project") {
            ele.empFields.addClass("is-hidden");
            ele.planFields.removeClass("is-hidden");
        } else {
            ele.planFields.addClass("is-hidden");
            ele.empFields.removeClass("is-hidden");
        }
    }

    function checkValidField(e) {
        e.stopPropagation();
        let ele = $(this);
        let valid = ele.get(0).checkValidity();
        if (!valid) {
            MainModule.notify("danger", "invalid form fields");
            addErrorStyle(ele);
        } else {
            removeErrorStyle(ele);
        }
    }

    function checkValidForm() {
        let formValid = true;
        ele.formInput.filter("[required]:visible").each(function() {
            let e = $(this);
            let valid = e.get(0).checkValidity();
            if (!valid) {
                MainModule.notify("danger", "invalid form fields");
                addErrorStyle(e);
                formValid = false;
            } else {
                removeErrorStyle(e);
            }
        });
        return formValid;
    }

    function addErrorStyle(ele) {
        // add is-danger class for form errors
        let p = ele.parent();
        if (p.hasClass("select")) {
            p.addClass("is-danger");
        } else {
            ele.addClass("is-danger");
        }
    }

    function removeErrorStyle(ele) {
        let p = ele.parent();
        if (p.hasClass("select")) {
            p.removeClass("is-danger");
        } else {
            ele.removeClass("is-danger");
        }
    }

    function handleSelect(e) {
        e.stopPropagation();
        e.preventDefault();
        if (!checkValidForm()) {
            MainModule.notify("danger", `Invalid Form Data: check form fields`);
            return
        }
        if (ele.newCheck.prop("checked")) {
            reqNewForm();
        } else {
            reqPlanPage();
        }
    }

    function handleSubmit(e) {
        e.stopPropagation();
        e.preventDefault();
        if (!checkValidForm()) {
            MainModule.notify("danger", `Invalid Form Data: check form fields`);
            return
        }
        reqPlanPage();
    }

    function handleCancel(e) {
        e.stopPropagation();
        e.preventDefault();
        teardown();
        HomeModule.navHome();
    }

    function reqPlanPage() {
        let formData = ele.form.serialize();
        let url = "/plan/";
        $.ajax({
            url: url,
            method: "POST",
            data: formData,
            dataType: "html",
            beforeSend: function() {
                // show loading indication
                MainModule.showProgress();
                ele.fieldset.prop("disabled", true);
            },
        }).done(function(res) {
            MainModule.endProgress();
            teardown();
            MainModule.setContent(res);
            PlanPage.init();
        }).fail(function(xhr, status, err) {
            ele.fieldset.prop("disabled", false);
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    function reqNewForm() {
        let url = "/plan/";
        $.ajax({
            url: url,
            method: "PUT",
            dataType: "html",
            beforeSend: function() {
                // show loading indication
                MainModule.showProgress();
                ele.fieldset.prop("disabled", true);
            },
        }).done(function(res) {
            MainModule.endProgress();
            teardown();
            MainModule.setContent(res);
            initForm();
        }).fail(function(xhr, status, err) {
            ele.fieldset.prop("disabled", false);
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    function initForm() {
        ele.form = $("form");
        ele.fieldset = $("fieldset");
        ele.btnSubmit = $("#btn-submit");
        ele.btnCancel = $("#btn-cancel");
        ele.empFields = $("#emp-fields");
        ele.planFields = $("#plan-fields");
        ele.typeSelect = $("#select-plan");
        ele.selectPlanner = $("#select-planner");
        ele.formInput = $("input,select");
        ele.empSelect = $("#select-employee");

        ele.form.on("blur", "input,select", checkValidField);
        ele.btnSubmit.on("click", handleSubmit);
        ele.btnCancel.on("click", handleCancel);
        ele.selectPlanner.on("change", handlePlanType);
    }

    function init() {
        ele.typeSelect = $("#select-plan");
        ele.btnSelect = $("#btn-select");
        ele.btnCancel = $("#btn-cancel");
        ele.newCheck = $("#plan-new-check");
        ele.form = $("form");
        ele.fieldset = $("fieldset");
        ele.formInput = $("input,select");

        ele.newCheck.on("click", handleCheck);
        ele.btnSelect.on("click", handleSelect);
        ele.btnCancel.on("click", handleCancel);
        ele.form.on("blur", "input,select", checkValidField);
    }

    function teardown() {
        ele.form.off();
    }

    return {
        init,
        teardown
    }

})(jQuery);

const PlanPage = (function($) {
    let ele = {};

    let currentTab;

    function checkValidField(e) {
        e.stopPropagation();
        let c = $(this);
        let valid = c.get(0).checkValidity();
        if (!valid) {
            MainModule.notify("danger", "invalid form fields");
            addErrorStyle(c);
        } else {
            removeErrorStyle(c);
        }
    }

    function checkValidForm() {
        let formValid = true;
        ele.formInput.filter("[required]:visible").each(function() {
            let e = $(this);
            let valid = e.get(0).checkValidity();
            if (!valid) {
                MainModule.notify("danger", "invalid form fields");
                addErrorStyle(e);
                formValid = false;
            } else {
                removeErrorStyle(e);
            }
        });
        return formValid;
    }

    function addErrorStyle(ele) {
        let p = ele.parent();
        if (p.hasClass("select")) {
            p.addClass("is-danger");
        } else {
            ele.addClass("is-danger");
        }
    }

    function removeErrorStyle(ele) {
        let p = ele.parent();
        if (p.hasClass("select")) {
            p.removeClass("is-danger");
        } else {
            ele.removeClass("is-danger");
        }
    }

    function handleUpdateHours() {
        removeErrorStyle(ele.targetHoursInput);
        try {
            let val = new Big(ele.targetHoursInput.val());
            let total = new Big(ele.targetHours.text());
            let delta = val.minus(total);
            ele.targetHoursDelta.text(delta);
        } catch (e) {
            addErrorStyle(ele.targetHoursInput);
        }
    }

    function handleUpdateCost() {
        removeErrorStyle(ele.targetCostInput);
        try {
            let val = new Big(ele.targetCostInput.val());
            let total = new Big(ele.targetCost.text());
            let delta = val.minus(total);
            ele.targetCostDelta.text(delta);
        } catch (e) {
            addErrorStyle(ele.targetCostInput);
        }
    }

    function handleTabClick(e) {
        // TODO: maybe use detach instead of is-hidden for performance?
        e.preventDefault();
        e.stopPropagation();
        if (currentTab) {
            currentTab.removeClass("is-active");
            currentTab.data("content").addClass("is-hidden");
        }

        let clicked = $(this);
        clicked.addClass("is-active");
        clicked.data("content").removeClass("is-hidden");
        currentTab = clicked;
    }

    function handleAddTab(tabname, content) {
        let newTab = `<li><a>${tabname}</a></li>`;
        ele.plannerTabs.append(newTab);
        // TODO: maybe add is-active class or similar for easier selection later
        // this would remove the need for a global variable
        ele.tabs.after(content);

        // might need to await these two methods calls
        let c = ele.tabs.next();
        // attach array of prodHours to table ele
        getProdHours(c);
        // attach lookup of cal_date: index to table ele
        getProdHoursLookup(c);
        // setup event listeners on the table element
        setupTableEvents(c);

        let li = $("ul li").last();
        li.data("content", c);
        li.trigger('click');
    }

    function setupTableEvents(tableEle) {
        let startDate = tableEle.data("pop-start");
        let endDate = tableEle.data("pop-end");

        tableEle.on("click", "button", function(e) {
            e.preventDefault();
            e.stopPropagation();
            let selectedEle = $(this);
            let evt = selectedEle.data("evt");

            if (!evt) {
                MainModule.notify("danger", `internal server error: invalid event`);
                return
            }

            const evts = {
                "add-row": handleAddRow,
                "adjust-col": handleAdjustCol,
                "delete-row": handleDelRow,
                "adjust-row": handleAdjustRow,
                "show-cal": handleShowCal
            };

            evts[evt](selectedEle, startDate, endDate);
        });
    }

    function handleAdjustRow(selectedEle, startDate, endDate) {
        let url = "/api/newrow";
        let planId = selectedEle.closest("tr").data("scope-id");
        let empId = selectedEle.closest("tr").data("emp-id");
        $.ajax({
            url: url,
            method: "GET",
            dataType: "html",
            beforeSend: function() {
                MainModule.showProgress();
            },
        }).done(function(res) {
            MainModule.endProgress();
            if (currentTab) {
                currentTab.data("content").addClass("is-hidden");
            }
            ele.tabs.addClass("is-hidden");
            ele.tabs.before(res);
            setupRowForm(planId, startDate, endDate);
        }).fail(function(xhr, status, err) {
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    // TODO: save logic should save the plan hours for each row, and add the
    // plan page id for load later. Since the plans and plan pages will be
    // // many to many, need an association table
    // function handleAdjustCol(selectedEle, startDate, endDate) {
    //     let fiscalPeriod = selectedEle.data("fiscal-period");
    //     if (!fiscalPeriod) {
    //         MainModule.notify("danger", "invalid element selected");
    //         return
    //     }
    //     let planId = selectedEle.closest("tr").data("scope-id");
    //     let empId = selectedEle.closest("tr").data("emp-id");
    //     let url = `/api/planrow?emp_id=${empId}&plan_id=${planId}`;
    //     let hours = serializePlanHours(selectedEle);
    //     $.ajax({
    //         url: url,
    //         method: "PUT",
    //         contentType: "application/json",
    //         data: JSON.stringify(hours),
    //         beforeSend: function() {
    //             MainModule.showProgress();
    //         },
    //     }).done(function(res) {
    //         MainModule.endProgress();
    //         ele.planPage.trigger("may:update-hours", [fiscalPeriod]);
    //     }).fail(function(xhr, status, err) {
    //         MainModule.endProgress();
    //         MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
    //     });
    // }

    function modalFormMarkup() {
        let markup = `
        <div class="modal" id="hours-modal">
            <div class="modal-background"></div>
            <div class="modal-card">
            <header class="modal-card-head">
                <p class="modal-card-title">Adjust Hours</p>
                <button class="delete" aria-label="close" id="close-modal"></button>
            </header>
            <section class="modal-card-body">
                <div class="field">
                    <label class="label">Enter a Value Between -1.0 and 1.0</label>
                    <div class="control">
                        <input class="input" type="text" placeholder="Hours" id="hours-input">
                    </div>
                </div>
            </section>
            <footer class="modal-card-foot">
                <button class="button is-success" id="btn-submit-modal">Submit</button>
                <button class="button" id="btn-cancel-modal">Cancel</button>
            </footer>
            </div>
        </div>
        `
        return markup
    }

    function setupHoursModal(fiscalPeriod, monthStart, monthEnd) {
        $("#btn-submit-modal").on("click", function(e) {
            e.preventDefault();
            e.stopPropagation();
            let m = $("#hours-input").val();
            console.log(m);
            teardownHoursModal();
            ele.planPage.trigger("may:update-hours", [m, fiscalPeriod, monthStart, monthEnd]);
        });
        $("#close-modal, #btn-cancel-modal, .modal-background").on("click", function(e) {
            e.preventDefault();
            e.stopPropagation();
            $("#hours-modal").removeClass("is-active");
            teardownHoursModal();
        });
    }

    function teardownHoursModal() {
        let modal = $("#hours-modal");
        modal.off();
        modal.remove();
    }

    function handleAdjustCol(selectedEle, startDate, endDate) {
        let markup = modalFormMarkup();
        $("#app").before(markup);
        let fiscalPeriod = selectedEle.data("fiscal-period");
        let monthStart = selectedEle.data("start-date");
        let monthEnd = selectedEle.data("end-date");
        setupHoursModal(fiscalPeriod, monthStart, monthEnd);
        $("#hours-modal").addClass("is-active");
    }

    function handleAddRow(selectedEle, startDate, endDate) {
        let url = "/api/newrow";
        let planId = selectedEle.closest("div.table-container").data("plan-id");
        $.ajax({
            url: url,
            method: "GET",
            dataType: "html",
            beforeSend: function() {
                MainModule.showProgress();
            },
        }).done(function(res) {
            MainModule.endProgress();
            if (currentTab) {
                currentTab.data("content").addClass("is-hidden");
            }
            ele.tabs.addClass("is-hidden");
            ele.tabs.before(res);
            setupRowForm(planId, startDate, endDate);
        }).fail(function(xhr, status, err) {
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    function handleDelRow(selectedEle, startDate, endDate) {
        let planId = selectedEle.closest("tr").data("scope-id");
        let empId = selectedEle.closest("tr").data("emp-id");
        let url = `/api/planrow?emp_id=${empId}&plan_id=${planId}`;
        $.ajax({
            url: url,
            method: "DELETE",
            beforeSend: function() {
                MainModule.showProgress();
            },
        }).done(function(res) {
            MainModule.endProgress();
            selectedEle.closest("tr").remove();
        }).fail(function(xhr, status, err) {
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    function handleShowCal(selectedEle, startDate, endDate) {
        let url = "";
        let planId = selectedEle.closest("tr").data("scope-id");
        let empId = selectedEle.closest("tr").data("emp-id");
        $.ajax({
            url: url,
            method: "DELETE",
            data: {
                "start_date": startDate,
                "end_date": endDate,
                "emp_id": empId,
                "plan_id": planId
            },
            dataType: "html",
            beforeSend: function() {
                MainModule.showProgress();
            },
        }).done(function(res) {
            MainModule.endProgress();
            selectedEle.closest("tr").remove();
        }).fail(function(xhr, status, err) {
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    function getTableRow(empId, planId, startDate, endDate) {
        let url = "/api/planrow";
        $.ajax({
            url: url,
            method: "GET",
            data: {
                "start_date": startDate,
                "end_date": endDate,
                "emp_ids": empId,
                "plan_ids": planId
            },
            dataType: "html",
            beforeSend: function() {
                // show loading indication
                MainModule.showProgress();
            },
        }).done(function(res) {
            MainModule.endProgress();
            teardownRowForm();
            $("#emp-multi-select").remove();
            ele.tabs.removeClass("is-hidden");
            if (currentTab) {
                currentTab.data("content").removeClass("is-hidden");
                currentTab.data("content").find("tbody").prepend(res);
                let newRow = $(`tr[data-emp-id='${empId}'][data-scope-id='${planId}']`)
                getPlanHours(newRow, startDate, endDate);
            }
        }).fail(function(xhr, status, err) {
            MainModule.endProgress();
            teardownRowForm();
            $("#emp-multi-select").remove();
            ele.tabs.removeClass("is-hidden");
            if (currentTab) {
                currentTab.data("content").removeClass("is-hidden");
                currentTab.data("content").find("tbody").prepend(res);
            }
            MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    function setupRowForm(planId, startDate, endDate) {
        ele.entityList = $("#entity-list");
        ele.entityList.data("selected", []);
        ele.entityList.on("click", "a", function(e) {
            e.preventDefault();
            e.stopPropagation();
            let d = ele.entityList.data("selected");
            let selectedEle = $(this);
            let name = selectedEle.data("id");
            if (selectedEle.hasClass("has-background-primary")) {
                selectedEle.removeClass("has-background-primary");
                let i = d.indexOf(name);
                d.splice(i, 1);
            } else {
                selectedEle.addClass("has-background-primary");
                d.push(name);
            }
            ele.entityList.data("selected", d);
        });

        let s = $("#emp-row-search-input");
        let eList = $("#entity-list a");
        s.on("keyup", function() {
            let val = s.val().toLowerCase();
            eList.filter(function() {
                $(this).toggle($(this).text().toLowerCase().indexOf(val) > -1);
            });
        });

        ele.btnPlanFormSubmit = $("#btn-add-employees");
        ele.btnPlanFormSubmit.on("click", function(e) {
            handleRowFormSubmit(e, planId, startDate, endDate);
        });
        ele.btnPlanFormCancel = $("#btn-cancel-employees");
        ele.btnPlanFormCancel.on("click", function(e) {
            e.preventDefault();
            e.stopPropagation();
            teardownRowForm();
            $("#emp-multi-select").remove();
            ele.tabs.removeClass("is-hidden");
            if (currentTab) {
                currentTab.data("content").removeClass("is-hidden");
            }
        });
    }

    async function handleRowFormSubmit(e, planId, startDate, endDate) {
        e.preventDefault();
        e.stopPropagation();

        // d is the list of empIds
        let d = ele.entityList.data("selected");

        // rows is the list of ajax calls
        let rows = [];
        for (const eId of d) {
            let call = initRow(eId, planId, startDate, endDate);
            rows.push(call);
        }
        const res = await Promise.allSettled(rows);
        for (const r of res) {
            if (r.status === "rejected") {
                MainModule.notify("danger", `error: ${r.reason.responseText}`);
                return
            }
            getTableRow(r.value["emp_id"], r.value["plan_id"], startDate, endDate);
        }
    }

    function initRow(empId, planId, startDate, endDate) {
        return $.ajax({
            url: `/api/planrow?emp_id=${empId}&plan_id=${planId}&start_date=${startDate}&end_date=${endDate}`,
            method: "POST",
            dataType: "json"
        });
    }

    function setupForm() {
        ele.form = $("form");
        ele.fieldset = $("fieldset");
        ele.btnPlanFormSubmit = $("#btn-plan-form-submit");
        ele.btnPlanFormCancel = $("#btn-plan-form-cancel");
        ele.btnPlanFormCancel.on("click", handlePlanFormCancel);
        ele.btnPlanFormSubmit.on("click", handlePlanFormSubmit);
    }

    function teardownRowForm() {
        $("#emp-row-search-input").off();
        ele.entityList.off();
    }

    function teardownForm() {
        ele.btnPlanFormCancel.off("click", handlePlanFormCancel);
        ele.btnPlanFormSubmit.off("click", handlePlanFormSubmit);
    }

    function getPlanHours(rowEle, startDate, endDate) {
        let url = "/api/planhours";
        let empId = rowEle.data("emp-id");
        let scopeId = rowEle.data("scope-id");
        $.ajax({
            url: url,
            method: "GET",
            data: {
                "start_date": startDate,
                "end_date": endDate,
                "emp_id": empId,
                "plan_id": scopeId
            },
            dataType: "json",
            beforeSend: function() {
                // show loading indication
                MainModule.showProgress();
            },
        }).done(function(res) {
            MainModule.endProgress();
            //let planHours = new Map(Object.entries(res));
            rowEle.data("planHours", res);
        }).fail(function(xhr, status, err) {
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    function getProdHours(tableEle) {
        let url = "/api/prodhours";
        let startDate = tableEle.data("pop-start");
        let endDate = tableEle.data("pop-end");
        $.ajax({
            url: url,
            method: "GET",
            data: {
                "start_date": startDate,
                "end_date": endDate
            },
            dataType: "json",
            beforeSend: function() {
                // show loading indication
                MainModule.showProgress();
            },
        }).done(function(res) {
            MainModule.endProgress();
            let numDays = res.length;
            tableEle.data("prodHours", res);
            tableEle.data("numDays", numDays);
        }).fail(function(xhr, status, err) {
            //ele.fieldset.prop("disabled", false);
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    function getProdHoursLookup(tableEle) {
        let url = "/api/prodhoursidx";
        let startDate = tableEle.data("pop-start");
        let endDate = tableEle.data("pop-end");
        $.ajax({
            url: url,
            method: "GET",
            data: {
                "start_date": startDate,
                "end_date": endDate
            },
            dataType: "json",
            beforeSend: function() {
                // show loading indication
                MainModule.showProgress();
            },
        }).done(function(res) {
            MainModule.endProgress();
            tableEle.data("lookup", res);
        }).fail(function(xhr, status, err) {
            //ele.fieldset.prop("disabled", false);
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    function handleGetNewPlanForm() {
        let url = "/evms/plan";
        $.ajax({
            url: url,
            method: "GET",
            dataType: "html",
            beforeSend: function() {
                // show loading indication
                MainModule.showProgress();
            },
        }).done(function(res) {
            MainModule.endProgress();
            if (currentTab) {
                currentTab.data("content").addClass("is-hidden");
            }
            ele.tabs.addClass("is-hidden");
            ele.tabs.before(res);
            setupForm();
        }).fail(function(xhr, status, err) {
            //ele.fieldset.prop("disabled", false);
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    function handlePlanFormSubmit(e) {
        e.preventDefault();
        e.stopPropagation();
        let tabname = $("#plan-name").val();
        let form = ele.tabs.prev("div.container");
        form.remove();
        ele.tabs.removeClass("is-hidden");
        getNewPlan(tabname);
    }

    function handlePlanFormCancel(e) {
        e.preventDefault();
        e.stopPropagation();
        teardownForm();
        let form = ele.tabs.prev("div.container");
        form.remove();
        ele.tabs.removeClass("is-hidden");
    }

    function getNewPlan(tabname) {
        let url = "/evms/plan";
        let formData = ele.form.serialize();
        $.ajax({
            url: url,
            method: "POST",
            data: formData,
            dataType: "html",
            beforeSend: function() {
                // show loading indication
                MainModule.showProgress();
                ele.fieldset.prop("disabled", true);
            },
        }).done(function(res) {
            MainModule.endProgress();
            handleAddTab(tabname, res);
        }).fail(function(xhr, status, err) {
            ele.fieldset.prop("disabled", false);
            MainModule.endProgress();
            MainModule.notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    function serializePlanHours(selectedEle) {
        let lookup = currentTab.data("lookup");
        let planHours = selectedEle.data("planHours");
        let data = {};
        let idx = 0;
        for (const k in lookup) {
            idx = lookup[k];
            data[k] = planHours[idx];
        }
        console.log(data);
    }

    function updateHours(e, multiplier, fiscalPeriod, monthStart, monthEnd) {
        e.preventDefault();
        e.stopPropagation();
        console.log("update hours called");
        console.log(multiplier);
        console.log(fiscalPeriod);
        if (!currentTab) {
            console.log("no current tab");
            return
        }

        try {
            let m = new Big(multiplier);
            let selector = `button[data-fiscal-period='${fiscalPeriod}']`
            let tBody = currentTab.data("content").find("tbody");
            let tHead = currentTab.data("content").find("thead");
            let hoursOp = tHead.find("select").val();
            console.log(hoursOp);
            console.log(selector);
            console.log(tBody);
            if (hoursOp === "Adjust") {
                tBody.find(selector).each(function(i, e) {
                    let mul = m.plus("1.0");
                    let val = $(this).text();
                    let newVal = mul.times(val);
                    $(this).text(newVal.toString());
                });
            } else if (hoursOp === "Reset") {
                let hours = resetMonthHours(m, monthStart, monthEnd);
                tBody.find(selector).each(function(i, e) {
                    $(this).text(hours);
                });
            }
        } catch (e) {
            console.log(e);
            return
        }
    }

    function resetMonthHours(multiplier, monthStart, monthEnd) {
        let prodHours = currentTab.data("content").data("prodHours");
        let lookup = currentTab.data("content").data("lookup");
        let startIdx = lookup[monthStart];
        let endIdx = lookup[monthEnd];
        let newHours = new Big("0");

        for (let i=0; i <= endIdx; i++) {
            let idx = startIdx + i;
            let h = new Big(prodHours[idx]);
            let newVal = h.times(multiplier);
            newHours = newHours.plus(newVal);
        }
        // this is coming out wrong
        console.log(newHours.toString());
        return newHours.toString()
    }

    function init() {
        // TODO: on table add/load need to make three requests:
        // table markup
        // prod hours data (json object)
        // plan Hours data (json object)
        // the math can be done on the frontend. Maybe do a custom
        // event to calculate all the totals once the requests are
        // complete.
        ele.btnAdd = $("#btn-add-plan");
        ele.btnDel = $("#btn-delete-plan");
        ele.btnLoad = $("#btn-load-plan");
        ele.targetHours = $("#target-hours");
        ele.targetHoursInput = $("#target-hours-input");
        ele.targetHoursDelta = $("#target-hours-delta");
        ele.targetCost = $("#target-cost");
        ele.targetCostInput = $("#target-cost-input");
        ele.targetCostDelta = $("#target-cost-delta");
        ele.tabs = $("div.tabs");
        ele.plannerTabs = $("#planner-tabs");
        ele.planPage = $("#plan-page");

        ele.targetHoursInput.on("blur", handleUpdateHours);
        ele.targetCostInput.on("blur", handleUpdateCost);
        ele.plannerTabs.on("click", "li", handleTabClick);
        ele.btnAdd.on("click", handleGetNewPlanForm);
        ele.planPage.on("may:update-hours", updateHours);
        currentTab = null;
    }

    function teardown() {
        return
    }

    function reqMarkup() {
        // request for calendar page and main plan page
        return
    }

    function reqHours() {
        // request json object with prod/plan hours
        return
    }

    function calcTotals() {
        // custom event to 
        return
    }

    // const served = `{"a":"1", "b":"2", "c":"3", "d":"4", "e":"5", "f":"6", "g":"7", "h":"8"}`;
    // const data = JSON.parse(served);
    // const keys = Object.keys(data).sort();

    function* dateIt(startDate, endDate, data, keys) {
        // data should be a map of "YYYY-MM-DD": <float>
        // this is for iterating over a date range
        // const keys = Object.keys(data).sort(); // need to move this out of the function
        let max = keys.length;
        let i = 0;
        while (i < max) {
            if (keys[i] === startDate) {
                let j = i;
                let v;
                do {
                    v = keys[j];
                    yield data[v];
                    j++;
                } while(v !== endDate && j < max);
                i = max;
            }
            i++;
        }
    }

// const d = {"a": 1, "c": 2, "b": 3, "e": 4, "d": 5};
// const it = dateIt("c", "e", d);
// let result = it.next();
// while(!result.done) {
//   let v = result.value;
//   console.log(v);
//   result = it.next();
// }

    return {
        init,
        teardown
    }
})(jQuery);
