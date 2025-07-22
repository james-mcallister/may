// the lookup and prod hours are associated with the table, and then
// each row should have a Map() of its plan hours. Instead of associating
// the planHours array in each row, I think we should just associate
// an object of this type to simiplify the calculate operations.

import { notify } from './messages';
import Big from 'big.js';
Big.DP = 2; // max decimal precision
Big.RM = Big.roundHalfUp; // rounding mode

class PlanTable {
    constructor(planId, popStart, popEnd) {
        this.id = planId;
        this.popStart = popStart;
        this.popEnd = popEnd;
        this.prodHours = null;
        this.lookup = null;
        this.dates = null;
        this.lookupLength = 0;
    }

    async init() {
        return fetchProdHours(this.popStart, this.popEnd)
          .then((res) => {
              this.setProdHours(res);
              return fetchProdHoursLookup(this.popStart, this.popEnd);
          })
          .then((res) => {
              this.setProdHoursLookup(res);
          });
    }

    setProdHours(prodHours) {
        this.prodHours = prodHours;
    }

    setProdHoursLookup(prodHoursLookup) {
        this.lookup = prodHoursLookup;
        this.dates = Object.keys(prodHoursLookup);
        this.lookupLength = this.dates.length;
    }

    // multiply the prodHours by the multiplier to calculate the planHours
    // pass the planHours array here or the PlanRow object?
    // m: multiplier of type Big
    // startDate: string
    // endDate: string
    // row: PlanRow
    resetHours(m, startDate, endDate, row) {
        let startIdx = this.getStartIndex(startDate);
        let endIdx = this.getEndIndex(endDate);

        let i = startIdx;
        while (i <= endIdx) {
            let newVal = m.times(this.prodHours[i]);
            row.updateHours(newVal, i);
            i++;
        }
    }

    // same as reset except multiply by existing plan val
    // m should be m + 1.0 and of type Big before it is passed here
    adjustHours(m, startDate, endDate, row) {
        let startIdx = this.getStartIndex(startDate);
        let endIdx = this.getEndIndex(endDate);

        let i = startIdx;
        while (i <= endIdx) {
            let val = row.getHours(i);
            let newVal = m.times(val);
            row.updateHours(newVal, i);
            i++;
        }
    }

    // return a slice of the plan hours (28 or 35 days for the calendar)
    getPlanHours(startDate, endDate, row) {
        let startIdx = this.getStartIndex(startDate);
        let endIdx = this.getEndIndex(endDate);
        let hours = [];

        let i = startIdx;
        while (i <= endIdx) {
            let val = row.getHours(i);
            hours.push(val);
            i++;
        }
        return hours;
    }

    // update the plan hours for a date range (calendar update)
    updatePlanHours(startDate, endDate, row, vals) {
        let startIdx = this.getStartIndex(startDate);
        let endIdx = this.getEndIndex(endDate);

        let updateCount = endIdx - startIdx + 1;
        if (updateCount != vals.length) {
            throw new Error("server error: incorrect update logic");
        }

        let i = startIdx;
        let j = 0;
        while (i <= endIdx) {
            let v = vals[j];
            row.updateHours(v, i);
            i++;
            j++;
        }
    }

    sum(startDate, endDate, row) {
        let startIdx = this.getStartIndex(startDate);
        let endIdx = this.getEndIndex(endDate);
        return row.sumHours(startIdx, endIdx);
    }

    saveHours(row) {
        let startIdx = this.getStartIndex(this.popStart);
        let endIdx = this.getEndIndex(this.popEnd);
        let hours = {};
        let i = startIdx;
        while (i <= endIdx) {
            let val = row.getHours(i);
            let d = this.dates[i];
            hours[d] = val;
            i++;
        }
        savePlanHours(JSON.stringify(hours));
    }

    getStartIndex(monthStart) {
        // edge case where the beginning and end of the PoP might be after
        // the month start date
        let c = monthStart in this.lookup
        if (!c) {
            return 0; // first index
        }
        return this.lookup[monthStart];
    }

    getEndIndex(monthEnd) {
        // edge case where the beginning and end of the PoP might be after
        // the month end date
        let c = monthEnd in this.lookup
        if (!c) {
            return this.lookupLength - 1; // last index
        }
        return this.lookup[monthEnd];
    }
}

class PlanRow {
    constructor(empId, planId) {
        this.empId = empId;
        this.planId = planId;
        this.planHours = null;
    }

    async init(startDate, endDate) {
        return fetchPlanHours(this.empId, this.planId, startDate, endDate)
            .then((res) => {
                this.setPlanHours(res);
            });
    }

    setPlanHours(planHours) {
        this.planHours = planHours;
    }

    sumHours(startIdx, endIdx) {
        let sum = new Big(0);
        let i = startIdx;
        while (i <= endIdx) {
            sum = sum.plus(this.planHours[i]);
            i++;
        }
        return sum
    }

    updateHours(v, idx) {
        this.planHours[idx] = v;
    }

    getHours(idx) {
        return this.planHours[idx];
    }
}

function savePlanHours(hours) {
    let url = "";
    $.ajax({
        url: url,
        method: "PUT",
        contentType: "application/json",
        data: hours,
    }).done(function(res) {
        notify("success", `Hours Updated`);
    }).fail(function(xhr, status, err) {
        // need to import notify somehow here or figure out error handling
        notify("danger", `request failure: ${url} ${xhr.responseText}`);
    });
}

function fetchProdHours(popStart, popEnd) {
    let url = "/api/prodhours";
    return $.ajax({
        url: url,
        method: "GET",
        data: {
            "start_date": popStart,
            "end_date": popEnd
        },
        dataType: "json",
    });
}

function fetchProdHoursLookup(popStart, popEnd) {
    let url = "/api/prodhoursidx";
    return $.ajax({
        url: url,
        method: "GET",
        data: {
            "start_date": popStart,
            "end_date": popEnd
        },
        dataType: "json",
    });
}

function fetchPlanHours(empId, planId, popStart, popEnd) {
    let url = "/api/planhours";
    return $.ajax({
        url: url,
        method: "GET",
        data: {
            "start_date": popStart,
            "end_date": popEnd,
            "emp_id": empId,
            "plan_id": planId
        },
        dataType: "json",
    });
}

export { PlanRow, PlanTable };
