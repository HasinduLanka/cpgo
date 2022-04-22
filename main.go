package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"sync"
)

// --------------------- ReaderLogic.go --------------------- //

// Read a list of strings separated by sep
func (rd *Reader) ReadArray(sep string) []string {
	S := rd.ReadLine()
	A := strings.Split(S, sep)
	return A
}

// Read a list of strings seperated by sep and clean the strings by removing empty/sep values
func (rd *Reader) ReadArrayClean(sep string) []string {
	A := rd.ReadArray(sep)
	Out := make([]string, 0, len(A))

	for _, S := range A {
		if (len(S) > 0) && (S != sep) {
			Out = append(Out, S)
		}
	}

	return Out
}

// Read a list of strings seperated by sep, and return a HashSet of the strings
func (rd *Reader) ReadHashset(sep string) map[string]struct{} {
	A := rd.ReadArray(sep)
	Out := make(map[string]struct{}, len(A))

	for _, S := range A {
		Out[S] = struct{}{}
	}

	delete(Out, sep)
	delete(Out, "")

	return Out
}

// Read pairs of key/value in format: key1 value1 key2 value2
//
// sep : separator
// skip : skip this many entries before reading key/value pairs
func (rd *Reader) ReadPairs(sep string, skip int) map[string]string {
	A := rd.ReadArrayClean(sep)

	Out := make(map[string]string, (len(A)-skip)/2)

	lenA := len(A)

	if (len(A)-skip)%2 == 1 {
		// odd number of entries
		lenA--

		Out[A[lenA]] = ""
	}

	for i := skip; i < lenA; i += 2 {
		Out[A[i]] = A[i+1]
	}

	return Out
}

// Read pairs of key/value in format: key1 111 key2 222
//
// sep : separator
// skip : skip this many entries before reading key/value pairs
func (rd *Reader) ReadStringIntPairs(sep string, skip int) (map[string]int, error) {
	A := rd.ReadArrayClean(sep)

	Out := make(map[string]int, (len(A)-skip)/2)

	lenA := len(A)

	if (len(A)-skip)%2 == 1 {
		// odd number of entries
		lenA--

		Out[A[lenA]] = 0
	}

	for i := skip; i < lenA; i += 2 {
		n, err := strconv.Atoi(A[i+1])
		if err != nil {
			return nil, err
		}
		Out[A[i]] = n
	}

	return Out, nil
}

func (rd *Reader) ReadBoolean(possitiveResponse string) bool {
	S := strings.ToLower(rd.ReadLine())
	possitiveResponse = strings.ToLower(possitiveResponse)

	return S == possitiveResponse
}

// Read one integer in one line
func (rd *Reader) ReadInt() (int, error) {
	S := rd.ReadLine()
	I, err := strconv.Atoi(S)
	return I, err
}

// Read a list of integers separated by sep
func (rd *Reader) ReadIntArray(sep string) ([]int, error) {
	A := rd.ReadArrayClean(sep)
	Out := make([]int, len(A))

	for i, S := range A {
		n, err := strconv.Atoi(S)
		if err != nil {
			return nil, err
		}
		Out[i] = n
	}

	return Out, nil
}

// --------------------- WriterLogic.go --------------------- //

// Print a string and a new line
func (wr *Writer) Print(S string) {
	wr.PrintInline(S + "\n")
}

// Print an integer and a new line
func (wr *Writer) PrintInt(I int) {
	wr.Print(strconv.Itoa(I))
}

// Print an integer without a new line
func (wr *Writer) PrintIntInline(I int) {
	wr.PrintInline(strconv.Itoa(I))
}

// Print an array of strings, separated by sep
func (wr *Writer) PrintArray(A []string, sep string) {
	wr.Print(strings.Join(A, sep))
}

// Print an array of integers, separated by sep
func (wr *Writer) PrintIntArray(A []int, sep string) {
	SA := make([]string, len(A))

	for i, v := range A {
		SA[i] = strconv.Itoa(v)
	}

	wr.PrintArray(SA, sep)
}

// Print any object as a json, spaced and indented
func (wr *Writer) Log(Obj interface{}) {
	B, JErr := json.MarshalIndent(Obj, "", "\t")

	if JErr != nil {
		return
	}

	S := string(B)

	wr.Print(S)
}

// Print any object as a single line json
func (wr *Writer) LogLine(Obj interface{}) {
	B, JErr := json.Marshal(Obj)

	if JErr != nil {
		return
	}

	S := string(B)

	wr.Print(S)
}

// --------------------- CP Adapter --------------------- //

type Reader struct {
	// Scanner   *bufio.Scanner
	BufReader *bufio.Reader
}

func NewReaderFromSTDIn() *Reader {
	rd := new(Reader)
	rd.ScanFromStandardInput()
	return rd
}

func (rd *Reader) ScanFromStandardInput() {
	// rd.Scanner = bufio.NewScanner(os.Stdin)
	rd.BufReader = bufio.NewReader(os.Stdin)
}

// Read a line from STDIN
func (rd *Reader) ReadLine() string {

	line, _ := rd.BufReader.ReadString('\n')

	line = strings.TrimSpace(line)
	return line
}

type Writer struct {
	channel   chan string
	Buff      *bufio.Writer
	waitGroup *sync.WaitGroup
}

func NewWriterToStandardOutput() *Writer {
	wr := new(Writer)
	wr.WriteToStandardOutput()
	return wr
}

func (wr *Writer) WriteToStandardOutput() {
	wr.Buff = bufio.NewWriter(os.Stdout)
	wr.channel = make(chan string, 1024)
	wr.waitGroup = new(sync.WaitGroup)

	go func() {
		for output := range wr.channel {
			wr.Buff.WriteString(output)
			wr.Buff.Flush()
			wr.waitGroup.Done()
		}
	}()
}

// Print a string without a new line
func (wr *Writer) PrintInline(S string) {
	wr.waitGroup.Add(1)
	wr.channel <- S
}

// Wait until all writes are done
func (wr *Writer) Flush() {
	wr.waitGroup.Wait()
}

type RWConsole struct {
	Reader
	Writer
}

var console = RWConsole{
	Reader: *NewReaderFromSTDIn(),
	Writer: *NewWriterToStandardOutput(),
}

// --------------------- ------- --------------------- //
// --------------------- Program --------------------- //
// --------------------- ------- --------------------- //

func main() {

	console.Print("Hello World")

	// Remeber to flush the buffer
	console.Flush()
}
