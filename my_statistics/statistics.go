// Copyright © 2010-12 Qtrac Ltd.
// 
// This program or package and any associated files are licensed under the
// Apache License, Version 2.0 (the "License"); you may not use these files
// except in compliance with the License. You can get a copy of the License
// at: http://www.apache.org/licenses/LICENSE-2.0.
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    "os"
    "fmt"
    //"log"
    "net/http"
    "net"
    "sort"
    "strconv"
    "strings"
    "math"
    "github.com/op/go-logging"
)

const (
    pageTop    = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head><title>Statistics</title>
<body><h3>Statistics</h3>
<p>Computes basic statistics for a given list of numbers</p>`
    form       = `<form action="/" method="POST">
<label for="numbers">Numbers (comma or space-separated):</label><br />
<input type="text" name="numbers" size="30"><br />
<input type="submit" value="Calculate">
</form>`
    pageBottom = `</body></html>`
    anError    = `<p class="error">%s</p>`
)

type statistics struct {
    numbers []float64
    mean    float64
    median  float64
    mode []float64
    stddev float64
}

var log = logging.MustGetLogger("statistics")
var format = logging.MustStringFormatter(
    `%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {

    // For demo purposes, create two backend for os.Stderr.
    backend1 := logging.NewLogBackend(os.Stderr, "", 0)

    // For messages written to backend2 we want to add some additional
    // information to the output, including the used log level and the name of
    // the function.
    backend1Formatter := logging.NewBackendFormatter(backend1, format)

    backend1Leveled := logging.AddModuleLevel(backend1)
    backend1Leveled.SetLevel(logging.INFO, "")

    // Set the backends to be used.
    logging.SetBackend(backend1Leveled, backend1Formatter)


    http.HandleFunc("/", homePage)
    if err := http.ListenAndServe(":9001", nil); err != nil {
        log.Error("failed to start server", err)
    }
}

func homePage(writer http.ResponseWriter, request *http.Request) {

    ip, port, err := net.SplitHostPort(request.RemoteAddr)
    log.Info("Request incoming from %s : %s", ip, port)

    err = request.ParseForm() // Must be called before writing response
    fmt.Fprint(writer, pageTop, form)
    if err != nil {
        fmt.Fprintf(writer, anError, err)
    } else {
        if numbers, message, ok := processRequest(request); ok {
            stats := getStats(numbers)
            fmt.Fprint(writer, formatStats(stats))
        } else if message != "" {
            fmt.Fprintf(writer, anError, message)
        }
    }
    fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) ([]float64, string, bool) {
    var numbers []float64
    if slice, found := request.Form["numbers"]; found && len(slice) > 0 {
        text := strings.Replace(slice[0], ",", " ", -1)
        for _, field := range strings.Fields(text) {
            if x, err := strconv.ParseFloat(field, 64); err != nil {
                return numbers, "'" + field + "' is invalid", false
            } else {
                numbers = append(numbers, x)
            }
        }
    }
    if len(numbers) == 0 {
        return numbers, "", false // no data first time form is shown
    }
    return numbers, "", true
}

func formatStats(stats statistics) string {
    return fmt.Sprintf(`<table border="1">
<tr><th colspan="2">Results</th></tr>
<tr><td>Numbers</td><td>%v</td></tr>
<tr><td>Count</td><td>%d</td></tr>
<tr><td>Mean</td><td>%f</td></tr>
<tr><td>Median</td><td>%f</td></tr>
<tr><td>Mode</td><td>%f</td></tr>
<tr><td>Std deviation</td><td>%f</td></tr>
</table>`, stats.numbers, len(stats.numbers), stats.mean, stats.median, stats.mode, stats.stddev)
}

func getStats(numbers []float64) (stats statistics) {
    stats.numbers = numbers
    sort.Float64s(stats.numbers)
    stats.mean = sum(numbers) / float64(len(numbers))
    stats.median = median(numbers)
    stats.mode = mode(numbers)
    stats.stddev = stddev(numbers)
    return stats
}

func sum(numbers []float64) (total float64) {
    for _, x := range numbers {
        total += x
    }
    return total
}

func median(numbers []float64) float64 {
    middle := len(numbers) / 2
    result := numbers[middle]
    if len(numbers)%2 == 0 {
        result = (result + numbers[middle-1]) / 2
    }
    return result
}

func mode(numbers []float64) [] float64 {
    var fq map[float64]int = make(map[float64]int)
    for _, x := range numbers {
        fq[x] = fq[x] + 1
    }
    min := -1
    max := 0
    for _,v := range fq {
        if (min == -1) {
            min = v
            max = v
        } else {
            if min > v {
                min = v
            }
            if max < v {
                max = v
            }
        }
    }
    var result []float64
    log.Info("Min max %d %d", min, max)
    if min != max {
        for k,v := range fq {
           if v == max {
               result = append(result, k)
           }
        }
    }
    return result;
}

func stddev(numbers []float64) float64 {
    mean := median(numbers)
    var diff float64 = 0.0
    for _, x := range numbers {
        diff += math.Pow(x - mean, 2)
    }
    diff /= float64(1.0 * len(numbers) - 1.0)
    return math.Sqrt(diff)
}

