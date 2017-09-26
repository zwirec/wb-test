package counter

import (
	"io"
	"sync"
	"bufio"
	"strconv"
	"net/http"
	"io/ioutil"
	"sync/atomic"
	"strings"
	"os"
	"bytes"
	"log"
)

var MaxNumWorkers int32 = 5

/**
A Counter defines parameters for running count
Substring in urls that received from In and send to Out
 */
type Counter struct {
	In         io.Reader
	Out        io.Writer
	buffer     *bytes.Buffer
	wg         sync.WaitGroup
	result     chan string
	doneWrite  chan struct{}
	doneWork   chan struct{}
	Substring  string
	totalCount uint64
	WorkersNum int32
	wCurr      int32
}

func NewCounter(in io.Reader, out io.Writer, substr string) *Counter {
	return &Counter{In: in,
		Out: out,
		Substring: substr,
		WorkersNum: MaxNumWorkers}
}

func (c *Counter) SetMaxNumWorkers(num int32) {
	c.WorkersNum = num
}

func (c *Counter) init() {
	if c.WorkersNum == 0 {
		c.WorkersNum = MaxNumWorkers
	}
	c.result = make(chan string, c.WorkersNum)
	c.buffer = new(bytes.Buffer)
	c.doneWork = make(chan struct{})
	c.doneWrite = make(chan struct{})
	return
}

func (c *Counter) Count() error {
	c.init()

	go c.writeWorkerRun()

	scanner := bufio.NewScanner(c.In)
	writer := bufio.NewWriter(c.Out)

	for scanner.Scan() {
		if c.WorkersNum == c.wCurr {
			<-c.doneWork
		}
		c.wg.Add(1)
		go c.taskWorkerRun(scanner.Text())
	}
	c.wg.Wait()

	close(c.result)

	<-c.doneWrite

	if _, err := writer.Write(c.buffer.Bytes()); err != nil {
		return err
	}

	writer.Flush()
	if _, err := writer.WriteString("Total: " + strconv.FormatUint(c.totalCount, 10) + "\n"); err != nil {
		log.Fatal(err)
	}

	writer.Flush()
	return nil
}

func (c *Counter) writeWorkerRun() {
	for {
		res, ok := <-c.result
		if ok {
			if _, err := c.buffer.WriteString(res); err != nil {
				log.Fatal(err)
			}
		} else {
			c.doneWrite <- struct{}{}
			return
		}
	}
}

func (c *Counter) taskWorkerRun(url string) {
	atomic.AddInt32(&c.wCurr, 1)

	defer func() {
		c.wg.Done()
		atomic.AddInt32(&c.wCurr, -1)
		select {
		case c.doneWork <- struct{}{}:
			{
			}
		default:

		}
	}()
	if response, err := http.Get(url); err == nil {

		defer response.Body.Close()

		if body_bytes, err := ioutil.ReadAll(response.Body); err == nil {
			cnt := strings.Count(string(body_bytes), c.Substring)
			atomic.AddUint64(&c.totalCount, uint64(cnt))
			c.result <- "Count for " + url + ": " + strconv.Itoa(cnt) + "\n"
			return
		} else {
			c.result <- `Error in url '` + url + `'` + "error: " + err.Error() + "\n"
			return
		}
	} else {
		c.result <- `Error in url '` + url + `'` + " error: " + err.Error() + "\n"
		return
	}

}

func Count() {
	c := Counter{In: os.Stdin,
		Out: os.Stdout,
		Substring: "Go",
	}
	c.Count()
	return
}

func CountFromTo(input io.Reader, output io.Writer, substring string) {
	c := Counter{In: input,
		Out: output,
		Substring: substring,
		WorkersNum: MaxNumWorkers,
	}
	c.Count()
	return
}
