import Big from 'big.js';
import { PlanRow, PlanTable } from './plan';
import { init as msgInit, notify, showProgress, endProgress } from './messages';
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

const formatterUSD = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
});

const formatNumber = new Intl.NumberFormat("en-US");

$(function() {
    MainModule.init();
    NavModule.init();
    msgInit();
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
                showProgress();
            }
        }).done(function(markup) {
            endProgress();
            MainModule.setContent(markup);
            if (showNotify) {
                notify(nColor, msg);
            }
        }).fail(function(xhr, status, err) {
            endProgress();
            notify("danger", `request failure: /home/ ${xhr.responseText}`);
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

    return {
        init() {
            ele.app = $("#app");
        },
        setContent(content) {
            ele.app.empty();
            ele.app.append(content);
        }
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
            notify("danger", `request failure: ${route} ${xhr.responseText}`);
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
                showProgress();
            }
        }).done(function(markup) {
            endProgress();
            teardown();
            MainModule.setContent(markup);
            initNext(handler);
        }).fail(function(xhr, status, err) {
            endProgress();
            notify("danger", `request failure: /${ent}/${id}/ ${xhr.responseText}`);
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
            notify("danger", "invalid form fields");
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
                notify("danger", "invalid form fields");
                ele.addClass("is-danger");
            } else {
                ele.removeClass("is-danger");
            }
        });

        // delete button should be only element with is-danger class
        return $(".is-danger").length === 1;
    }

    function teardown() {
        ele.btnClear.off("click", clearDropdown);
        ele.btnDelete.off("click", handleDelete);
        ele.btnCancel.off("click", handleCancel);
        ele.btnSubmit.off("click", handleUpdate);
        ele.form.off("blur", "input", checkValidField);
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
                showProgress();
                ele.fieldset.prop("disabled", true);
            },
        }).done(function(res) {
            endProgress();
            teardown();
            HomeModule.navHome("primary", res);
        }).fail(function(xhr, status, err) {
            ele.fieldset.prop("disabled", false);
            endProgress();
            notify("danger", `request failure: ${url} ${xhr.responseText}`);
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
        teardown
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
            notify("danger", "invalid form fields");
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
                notify("danger", "invalid form fields");
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
            notify("danger", `Invalid Form Data: check form fields`);
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
            notify("danger", `Invalid Form Data: check form fields`);
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
                showProgress();
                ele.fieldset.prop("disabled", true);
            },
        }).done(function(res) {
            endProgress();
            teardown();
            MainModule.setContent(res);
            PlanPage.init();
        }).fail(function(xhr, status, err) {
            ele.fieldset.prop("disabled", false);
            endProgress();
            notify("danger", `request failure: ${url} ${xhr.responseText}`);
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
                showProgress();
                ele.fieldset.prop("disabled", true);
            },
        }).done(function(res) {
            endProgress();
            teardown();
            MainModule.setContent(res);
            initForm();
        }).fail(function(xhr, status, err) {
            ele.fieldset.prop("disabled", false);
            endProgress();
            notify("danger", `request failure: ${url} ${xhr.responseText}`);
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
            notify("danger", "invalid form fields");
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
                notify("danger", "invalid form fields");
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
            let total = new Big(ele.targetHours.data("val"));
            let delta = val.minus(total);
            ele.targetHoursDelta.text(formatNumber.format(delta.toString()));
        } catch (e) {
            addErrorStyle(ele.targetHoursInput);
        }
    }

    function handleUpdateCost() {
        removeErrorStyle(ele.targetCostInput);
        try {
            let val = new Big(ele.targetCostInput.val());
            let total = new Big(ele.targetCost.data("val"));
            let delta = val.minus(total);
            ele.targetCostDelta.text(formatterUSD.format(delta));
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


        let c = ele.tabs.next();
        let li = $("ul li").last();
        initTableData(c, li);

        // setup event listeners on the table element
        setupTableEvents(c);

        li.data("content", c);
        li.trigger('click');
    }

    function initTableData(tableEle, listEle) {
        let oldTab = currentTab;
        let startDate = tableEle.data("pop-start");
        let endDate = tableEle.data("pop-end");
        let tableId = tableEle.data("plan-id");
        let t = new PlanTable(tableId, startDate, endDate);
        t.init().then(() => {
            tableEle.data("tableData", t);
        }).catch((err) => {
            notify("danger", `initTableData error: ${err.message}`);
            tableEle.remove();
            listEle.remove();
            if (oldTab) {
                currentTab = oldTab;
            } else {
                currentTab = null;
            }
        });
    }

    function setupTableEvents(tableEle) {
        let startDate = tableEle.data("pop-start");
        let endDate = tableEle.data("pop-end");

        tableEle.on("may:calc-totals", function(e) {
            e.preventDefault();
            e.stopPropagation();
            if (currentTab) {
                let t = currentTab.data("content");
                let tHead = t.find("thead");
                let tBody = t.find("tbody");
                let headBtns = tHead.find("button.col");
                headBtns.each(function() {
                    let cost = new Big(0);
                    let hours = new Big(0);
                    let fp = $(this).data("fiscal-period");
                    let mHours = $(this).data("month-hours");
                    let hBtns = tBody.find(`button.hours[data-fiscal-period="${fp}"]`);
                    hBtns.each(function() {
                        let btnHours = new Big($(this).text());
                        let rate = $(this).parent().siblings("td.rate").text();
                        let btnCost = btnHours.times(rate);
                        hours = hours.plus(btnHours);
                        cost = cost.plus(btnCost);
                    });
                    let fte = hours.div(mHours);
                    tBody.find(`td.fte[data-fiscal-period="${fp}"]`).text(fte);
                    tBody.find(`td.cost[data-fiscal-period="${fp}"]`).text(formatterUSD.format(cost));
                    tBody.find(`td.cost[data-fiscal-period="${fp}"]`).data("val", cost);
                    tBody.find(`td.hours[data-fiscal-period="${fp}"]`).text(formatNumber.format(hours.toString()));
                    tBody.find(`td.hours[data-fiscal-period="${fp}"]`).data("val", hours);
                });
            }
        });

        tableEle.on("click", "button", function(e) {
            e.preventDefault();
            e.stopPropagation();
            let selectedEle = $(this);
            let evt = selectedEle.data("evt");

            if (!evt) {
                notify("danger", `internal server error: invalid event`);
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

    function calcTotals() {
        if (currentTab) {
            let t = currentTab.data("content");
            t.trigger("may:calc-totals", []);

            let costTotal = new Big(0);
            t.find("td.cost").each(function() {
                let colCost = $(this).data("val");
                costTotal = costTotal.plus(colCost);
            });
            ele.targetCost.text(formatterUSD.format(costTotal));
            ele.targetCost.data("val", costTotal);

            let hoursTotal = new Big(0);
            t.find("td.hours").each(function() {
                let colHours = $(this).data("val");
                hoursTotal = hoursTotal.plus(colHours);
            });
            ele.targetHours.text(formatNumber.format(hoursTotal.toString()));
            ele.targetHours.data("val", hoursTotal);
        }
    }

    function handleAdjustRow(selectedEle) {
        let row = selectedEle.closest("tr");
        let planId = row.data("scope-id");
        let empId = row.data("emp-id");
        let multiplier = row.find("input").val();

        // loop over all buttons in the row and trigger update
        let btns = $(`tr[data-emp-id='${empId}'][data-scope-id='${planId}'] button.hours`);
        btns.each(function() {
            let btn = $(this);
            let fiscalPeriod = btn.data("fiscal-period");
            let monthStart = btn.data("start-date");
            let monthEnd = btn.data("end-date");
            let selector = `tr[data-emp-id='${empId}'][data-scope-id='${planId}'] button[data-fiscal-period='${fiscalPeriod}']`;
            ele.planPage.trigger("may:update-hours", [multiplier, selector, monthStart, monthEnd]);
        });
        calcTotals();
    }

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
            teardownHoursModal();
            let selector = `button[data-fiscal-period='${fiscalPeriod}']`;
            ele.planPage.trigger("may:update-hours", [m, selector, monthStart, monthEnd]);
            calcTotals();
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
                showProgress();
            },
        }).done(function(res) {
            endProgress();
            if (currentTab) {
                currentTab.data("content").addClass("is-hidden");
            }
            ele.tabs.addClass("is-hidden");
            ele.tabs.before(res);
            setupRowForm(planId, startDate, endDate);
        }).fail(function(xhr, status, err) {
            endProgress();
            notify("danger", `request failure: ${url} ${xhr.responseText}`);
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
                showProgress();
            },
        }).done(function(res) {
            endProgress();
            selectedEle.closest("tr").remove();
            calcTotals();
        }).fail(function(xhr, status, err) {
            endProgress();
            notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    function handleShowCal(selectedEle, startDate, endDate) {
        let url = "/evms/cal";
        let fiscalPeriod = selectedEle.data("fiscal-period");
        $.ajax({
            url: url,
            method: "GET",
            data: {
                "start_date": startDate,
                "end_date": endDate,
                "fiscal_period": fiscalPeriod
            },
            dataType: "html",
            beforeSend: function() {
                showProgress();
            },
        }).done(function(res) {
            endProgress();
            if (currentTab) {
                currentTab.data("content").addClass("is-hidden");
            }
            ele.tabs.addClass("is-hidden");
            ele.tabs.before(res);
            setupCal(selectedEle);
        }).fail(function(xhr, status, err) {
            endProgress();
            notify("danger", `request failure: ${url} ${xhr.responseText}`);
        });
    }

    function setupCal(selectedEle) {
        let row = selectedEle.closest("tr").data("rowData");
        let startDate = selectedEle.data("start-date");
        let endDate = selectedEle.data("end-date");
        let hours = currentTab.data("content").data("tableData").getPlanHours(startDate, endDate, row);
        // skip disabled calendar tiles (needed to ensure correct planHours index)
        let offset = 0;
        $(".calendar-tile").find("input").each(function(i) {
            let isDisabled = $(this).prop("disabled");
            if (!isDisabled) {
                let idx = i - offset;
                $(this).val(hours[idx]);
            } else {
                offset++;
            }
        });

        let slider = $("#capacity-slider");
        let cap = $("#capacity-value");
        slider.on("input", function(e) {
            e.preventDefault();
            e.stopPropagation();
            let v = slider.val();
            slider.data("val", v);
            cap.text(formatAsPercentage(v));
        });

        let btnCalc = $("#btn-cal-calculate");
        let calcControl = $("#calc-control");
        btnCalc.on("click", function(e) {
            e.preventDefault();
            e.stopPropagation();
            let cap = slider.data("val") || "1";
            let control = calcControl.val();
            let tiles = $(".calendar-tile > p");
            tiles.each(function() {
                let prodHours = $(this).text();
                let planHours = $(this).prev().val() || "0";
                let newPlanHours = "error";
                if (control === "adjust") {
                    let adjustCap = addNumber(cap, "1.00");
                    newPlanHours = multiplyNumber(planHours, adjustCap);
                } else {
                    newPlanHours = multiplyNumber(prodHours, cap);
                }
                $(this).prev().val(newPlanHours);
            });
        });
  
        calcControl.on("change", function(e) {
            e.preventDefault();
            e.stopPropagation();
            let control = calcControl.val();
            if (control === "adjust") {
                slider.attr("min", "-1");
            } else {
                slider.attr("min", "0");
                let v = parseFloat(slider.data("val"));
                if (v && v < 0) {
                    slider.data("val", "0");
                    cap.text(formatAsPercentage("0"));
                }
            }
        });

        // TODO: save event should update planHours array
        let btnSave = $("#btn-cal-save");
        btnSave.on("click", function(e) {
            e.preventDefault();
            e.stopPropagation();
            let saved = saveHours(selectedEle);
            if (saved) {
                calcTotals();
            }
        });
        // TODO: cancel should remove cal with no changes (teardown)
        let btnCancel = $("#btn-cal-cancel");
        btnCancel.on("click", function(e) {
            e.preventDefault();
            e.stopPropagation();
            teardownCal();
            calcTotals();
        });
        // TODO: prev should trigger the show cal event with previous fiscal period
        let btnPrev = $("#btn-cal-prev");
        btnPrev.on("click", function(e) {
            e.preventDefault();
            e.stopPropagation();
            let prevBtn = selectedEle.closest("td").prev().find("button.hours");
            if (prevBtn.length !== 0) {
                let t = currentTab.data("content");
                let startDate = t.data("pop-start");
                let endDate = t.data("pop-end");
                let saved = saveHours(selectedEle);
                if (saved) {
                    handleShowCal(prevBtn, startDate, endDate);
                }
            }
        });
        // TODO: next should trigger the show cal event with next fiscal period
        let btnNext = $("#btn-cal-next");
        btnNext.on("click", function(e) {
            e.preventDefault();
            e.stopPropagation();
            let nextBtn = selectedEle.closest("td").next().find("button.hours");
            if (nextBtn.length !== 0) {
                let t = currentTab.data("content");
                let startDate = t.data("pop-start");
                let endDate = t.data("pop-end");
                let saved = saveHours(selectedEle);
                if (saved) {
                    handleShowCal(nextBtn, startDate, endDate);
                }
            }
        });
    }

    function saveHours(selectedEle) {
        let startDate = selectedEle.data("start-date");
        let endDate = selectedEle.data("end-date");
        let row = selectedEle.closest("tr").data("rowData");
        let t = currentTab.data("content").data("tableData");
        let vals = [];

        $(".calendar-tile").find("input").each(function() {
            let isDisabled = $(this).prop("disabled");
            if (!isDisabled) {
                let val = $(this).val();
                vals.push(val);
            }
        });

        try {
            t.updatePlanHours(startDate, endDate, row, vals);
            let s = t.sum(startDate, endDate, row);
            selectedEle.text(s);
            teardownCal();
        } catch (e) {
            notify("danger", `${e.message}`);
            return false;
        }
        return true;
    }

    function addNumber(val1, val2) {
        let v1 = new Big(val1);
        let v2 = new Big(val2);
        return v1.add(v2)
    }

    function multiplyNumber(val1, val2) {
        let v1 = new Big(val1);
        let v2 = new Big(val2);
        return v1.times(v2)
    }

    function formatAsPercentage(number, minimumFractionDigits = 0, maximumFractionDigits = 2) {
        const formatter = new Intl.NumberFormat('default', {
            style: 'percent',
            minimumFractionDigits,
            maximumFractionDigits,
        });
        return formatter.format(number);
    }

    function teardownCal() {
        $("#plan-cal").off();
        $("#plan-cal").remove();
        ele.tabs.removeClass("is-hidden");
        if (currentTab) {
            currentTab.data("content").removeClass("is-hidden");
        }
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
                showProgress();
            },
        }).done(function(res) {
            endProgress();
            teardownRowForm();
            $("#emp-multi-select").remove();
            ele.tabs.removeClass("is-hidden");
            if (currentTab) {
                currentTab.data("content").removeClass("is-hidden");
                currentTab.data("content").find("tbody").prepend(res);
                let newRow = $(`tr[data-emp-id='${empId}'][data-scope-id='${planId}']`)
                initRowData(newRow, startDate, endDate);
            }
        }).fail(function(xhr, status, err) {
            endProgress();
            teardownRowForm();
            $("#emp-multi-select").remove();
            ele.tabs.removeClass("is-hidden");
            if (currentTab) {
                currentTab.data("content").removeClass("is-hidden");
            }
            notify("danger", `request failure: ${url} ${xhr.responseText}`);
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
            teardownRowForm();
            $("#emp-multi-select").remove();
            ele.tabs.removeClass("is-hidden");
            if (currentTab) {
                currentTab.data("content").removeClass("is-hidden");
            }
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
                notify("danger", `error: ${r.reason.responseText}`);
                continue
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

    function initRowData(rowEle, startDate, endDate) {
        let empId = rowEle.data("emp-id");
        let scopeId = rowEle.data("scope-id");
        let r = new PlanRow(empId, scopeId);
        r.init(startDate, endDate).then(() => {
            rowEle.data("rowData", r);
        }).catch((err) => {
            notify("danger", `row init error: ${err.message}`);
            rowEle.remove();
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
                showProgress();
            },
        }).done(function(res) {
            endProgress();
            if (currentTab) {
                currentTab.data("content").addClass("is-hidden");
            }
            ele.tabs.addClass("is-hidden");
            ele.tabs.before(res);
            setupForm();
        }).fail(function(xhr, status, err) {
            //ele.fieldset.prop("disabled", false);
            endProgress();
            notify("danger", `request failure: ${url} ${xhr.responseText}`);
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
        if (currentTab) {
            currentTab.data("content").removeClass("is-hidden");
        }
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
                showProgress();
                ele.fieldset.prop("disabled", true);
            },
        }).done(function(res) {
            endProgress();
            handleAddTab(tabname, res);
        }).fail(function(xhr, status, err) {
            ele.fieldset.prop("disabled", false);
            endProgress();
            notify("danger", `request failure: ${url} ${xhr.responseText}`);
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

    function updateHours(e, multiplier, selector, monthStart, monthEnd) {
        e.preventDefault();
        e.stopPropagation();
        if (!currentTab) {
            return
        }

        try {
            let m = new Big(multiplier);
            let t = currentTab.data("content");
            let tBody = t.find("tbody");
            let tHead = t.find("thead");
            let tData = t.data("tableData");
            let hoursOp = tHead.find("select").val();
            if (hoursOp === "Adjust") {
                tBody.find(selector).each(function(i, e) {
                    let mul = m.plus("1.0");
                    let row = $(this).closest("tr").data("rowData");
                    tData.adjustHours(mul, monthStart, monthEnd, row);
                    let newVal = tData.sum(monthStart, monthEnd, row);
                    $(this).text(newVal.toString());
                });
            } else if (hoursOp === "Reset") {
                tBody.find(selector).each(function(i, e) {
                    let row = $(this).closest("tr").data("rowData");
                    tData.resetHours(m, monthStart, monthEnd, row);
                    let newVal = tData.sum(monthStart, monthEnd, row);
                    $(this).text(newVal.toString());
                });
            }
        } catch (e) {
            // TODO: undecided if error should be ignored or displayed
            console.log(e);
            return
        }
    }

    function loadPlanTable() {
        return
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
        ele.btnLoad.on("click", loadPlanTable);
        currentTab = null;
    }

    function teardown() {
        ele.planPage.off();
    }

    return {
        init,
        teardown
    }
})(jQuery);
