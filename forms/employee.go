package forms

import (
	"fmt"
	"html/template"

	"github.com/james-mcallister/may/database"
)

// init template and any persistent data
type EmployeeForm struct {
	template *template.Template
}

// data to use in the template.Execute method. Called every request
type EmployeeFormData struct {
	Emp          database.Employee
	EmpDropdown  database.Dropdown
	IptDropdown  database.Dropdown
	CompDropdown database.Dropdown
}

func (f *EmployeeForm) InitTemplate() error {
	markup := `
<div class="block mt-1">
    <form>
        <fieldset>

            <div class="field">
                <label class="label">ID</label>
                <div class="control">
                    <input class="input" type="text" value="{{ .Emp.Id }}" disabled />
                </div>
                <p class="help">ID is managed by the database</p>
            </div>

            <div class="field">
                <label class="label">My ID</label>
                <div class="control">
                    <input class="input" type="text" value="{{ .Emp.Myid }}" />
                </div>
                <p class="help">Required Field (ex: m33445)</p>
            </div>

            <div class="field">
                <label class="label">First Name</label>
                <div class="control">
                    <input class="input" type="text" value="{{ .Emp.FirstName }}" />
                </div>
                <p class="help"></p>
            </div>

            <div class="field">
                <label class="label">Last Name</label>
                <div class="control">
                    <input class="input" type="text" value="{{ .Emp.LastName }}" />
                </div>
                <p class="help"></p>
            </div>

            <div class="field">
                <label class="label">Employee ID</label>
                <div class="control">
                    <input class="input" type="text" value="{{ .Emp.Empid }}" />
                </div>
                <p class="help">From workday or COBRA (ex: 1002456 or mcallja)</p>
            </div>

            <div class="field">
                <label class="label">Labor Capacity</label>
                <div class="control">
                    <input class="input" type="number" value="{{ .Emp.LaborCapacity }}" />
                </div>
                <p class="help">1.0 for full-time workers</p>
            </div>

            <div class="field">
                <label class="label">Desk</label>
                <div class="control">
                    <input class="input" type="text" value="{{ .Emp.Desk }}" />
                </div>
                <p class="help">Where they sit on campus (cube, office or lab room number)</p>
            </div>

            <div class="field">
                <div class="control">
                  <label class="radio">
                    <input type="radio" name="active" {{ if .Emp.Active }}checked{{ end }}>
                    Active
                  </label>
                  <label class="radio">
                    <input type="radio" name="active"{{ if .Emp.Active }}{{ else }}checked{{ end }}>
                    Inactive
                  </label>
                </div>
              </div>

              <div class="field">
                  <label class="label">Coverage Start Date</label>
                  <div class="control">
                      <input class="input" type="date" value="{{ .Emp.CoverageStart }}" />
                  </div>
                  <p class="help"></p>
              </div>

              <div class="field">
                  <label class="label">Coverage End Date</label>
                  <div class="control">
                      <input class="input" type="date" value="{{ .Emp.CoverageEnd }}" />
                  </div>
                  <p class="help"></p>
              </div>

            <div class="field has-addons">
                <div class="control is-expanded">
                    <div class="select is-fullwidth">
                        <select id="select-grade">
                            <option value="0" disabled selected>Select Grade...</option>
                            {{ range .CompDropdown }}
                            <option value="{{.Id}}">{{ .Name }}</option>
                            {{ end }}
                        </select>
                    </div>
                </div>
                <div class="control">
                    <button class="button clear" data-select-id="select-grade">Clear</button>
                </div>
            </div>

            <div class="field has-addons">
                <div class="control is-expanded">
                    <div class="select is-fullwidth">
                        <select id="select-manager">
                            <option value="0" disabled selected>Select Manager...</option>
                            {{ range .EmpDropdown }}
                            <option value="{{.Id}}">{{ .Name }}</option>
                            {{ end }}
                        </select>
                    </div>
                </div>
                <div class="control">
                    <button class="button clear" data-select-id="select-manager">Clear</button>
                </div>
            </div>

            <div class="field has-addons">
                <div class="control is-expanded">
                    <div class="select is-fullwidth">
                        <select id="select-ipt">
                            <option value="0" disabled selected>Select IPT...</option>
                            {{ range .IptDropdown }}
                            <option value="{{.Id}}">{{ .Name }}</option>
                            {{ end }}
                        </select>
                    </div>
                </div>
                <div class="control">
                    <button class="button clear" data-select-id="select-ipt">Clear</button>
                </div>
            </div>

            <div class="field is-grouped">
                <div class="control">
                    <button class="button is-link">Submit</button>
                </div>
                <div class="control">
                    <button class="button is-danger">Delete</button>
                </div>
                <div class="control">
                    <button class="button is-link is-light">Cancel</button>
                </div>
            </div>
        </fieldset>
    </form>
</div>
	`
	var err error
	f.template, err = template.New("employee-form").Parse(markup)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}
	return nil
}

// handler functions here for the 5 requests against the form
