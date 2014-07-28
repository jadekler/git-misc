package home

import (
    "fmt"
    "strconv"
    "net/http"
    "html/template"

    "git-misc/logic-tree/app/common"
)

type Condition struct {
    Text string
    Other string
}

func GetHomePage(rw http.ResponseWriter, req *http.Request) {
    type Page struct {
        Title string
        Conditions []Condition
    }
    
    p := Page{
        Title: "home",
        Conditions: getConditions(),
    }

    common.Templates = template.Must(template.ParseFiles("templates/home/home.html", common.LayoutPath))
    err := common.Templates.ExecuteTemplate(rw, "base", p)
    common.CheckError(err, 2)
}

func SaveState(rw http.ResponseWriter, req *http.Request) {
    field := req.FormValue("field")
    operator := req.FormValue("operator")
    value, err := strconv.Atoi(req.FormValue("value"))
    common.CheckError(err, 2)

    _, err = common.DB.Query(fmt.Sprintf("INSERT INTO logictree.equality(field, operator, value) VALUES ('%s', '%s', %d)", field, operator, value))
    common.CheckError(err, 2)

    _, err = common.DB.Query("INSERT INTO logictree.logic(operator) VALUES ('AND')")
    common.CheckError(err, 2)

    GetHomePage(rw, req)
}

func Truncate(rw http.ResponseWriter, req *http.Request) {
    _, err := common.DB.Query("TRUNCATE TABLE logictree.equality")
    common.CheckError(err, 2)

    _, err = common.DB.Query("TRUNCATE TABLE logictree.logic")
    common.CheckError(err, 2)

    GetHomePage(rw, req)
}

func getConditions() []Condition {
    conditions := make([]Condition, 0)

    rows, err := common.DB.Query("SELECT field, operator, value FROM logictree.equality")
    common.CheckError(err, 2)

    var field, operator, value string

    i := 0

    for rows.Next() {
        rows.Scan(&field, &operator, &value)
        common.CheckError(err, 2)

        if i != 0 {
            conditions = append(conditions, Condition{
                Text: "AND",
                Other: "data-type='logic'",
            })
        }

        conditions = append(conditions, Condition{
            Text: fmt.Sprintf("%s %s %s", field, operator, value),
            Other: fmt.Sprintf("data-type='equality' data-field='%s' data-operator='%s' data-value='%s'", field, operator, value)},
        )

        i++
    }

    return conditions
}







