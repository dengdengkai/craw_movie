
Golang学习 - regexp 包

------------------------------------------------------------

// 函数

// 判断在 b（s、r）中能否找到 pattern 所匹配的字符串
func Match(pattern string, b []byte) (matched bool, err error)
func MatchString(pattern string, s string) (matched bool, err error)
func MatchReader(pattern string, r io.RuneReader) (matched bool, err error)

// 将 s 中的正则表达式元字符转义成普通字符。
func QuoteMeta(s string) string

------------------------------

// 示例：MatchString、QuoteMeta
func main() {
	pat := `(((abc.)def.)ghi)`
	src := `abc-def-ghi abc+def+ghi`

	fmt.Println(regexp.MatchString(pat, src))
	// true <nil>

	fmt.Println(regexp.QuoteMeta(pat))
	// \(\(\(abc\.\)def\.\)ghi\)
}

------------------------------------------------------------

// Regexp 代表一个编译好的正则表达式，我们这里称之为正则对象。正则对象可以
// 在文本中查找匹配的内容。
//
// Regexp 可以安全的在多个例程中并行使用。
type Regexp struct { ... }

------------------------------

// 编译

// 将正则表达式编译成一个正则对象（使用 PERL 语法）。
// 该正则对象会采用“leftmost-first”模式。选择第一个匹配结果。
// 如果正则表达式语法错误，则返回错误信息。
func Compile(expr string) (*Regexp, error)

// 将正则表达式编译成一个正则对象（正则语法限制在 POSIX ERE 范围内）。
// 该正则对象会采用“leftmost-longest”模式。选择最长的匹配结果。
// POSIX 语法不支持 Perl 的语法格式：\d、\D、\s、\S、\w、\W
// 如果正则表达式语法错误，则返回错误信息。
func CompilePOSIX(expr string) (*Regexp, error)

// 功能同上，但会在解析失败时 panic
func MustCompile(str string) *Regexp
func MustCompilePOSIX(str string) *Regexp

// 让正则表达式在之后的搜索中都采用“leftmost-longest”模式。
func (re *Regexp) Longest()

// 返回编译时使用的正则表达式字符串
func (re *Regexp) String() string

// 返回正则表达式中分组的数量
func (re *Regexp) NumSubexp() int

// 返回正则表达式中分组的名字
// 第 0 个元素表示整个正则表达式的名字，永远是空字符串。
func (re *Regexp) SubexpNames() []string

// 返回正则表达式必须匹配到的字面前缀（不包含可变部分）。
// 如果整个正则表达式都是字面值，则 complete 返回 true。
func (re *Regexp) LiteralPrefix() (prefix string, complete bool)

------------------------------

// 示例：第一匹配和最长匹配
func main() {
	b := []byte("abc1def1")
	pat := `abc1|abc1def1`
	reg1 := regexp.MustCompile(pat)      // 第一匹配
	reg2 := regexp.MustCompilePOSIX(pat) // 最长匹配
	fmt.Printf("%s\n", reg1.Find(b))     // abc1
	fmt.Printf("%s\n", reg2.Find(b))     // abc1def1

	b = []byte("abc1def1")
	pat = `(abc|abc1def)*1`
	reg1 = regexp.MustCompile(pat)      // 第一匹配
	reg2 = regexp.MustCompilePOSIX(pat) // 最长匹配
	fmt.Printf("%s\n", reg1.Find(b))    // abc1
	fmt.Printf("%s\n", reg2.Find(b))    // abc1def1
}

------------------------------

// 示例：正则信息
func main() {
	pat := `(abc)(def)(ghi)`
	reg := regexp.MustCompile(pat)

	// 获取正则表达式字符串
	fmt.Println(reg.String())    // (abc)(def)(ghi)

	// 获取分组数量
	fmt.Println(reg.NumSubexp()) // 3

	fmt.Println()

	// 获取分组名称
	pat = `(?P<Name1>abc)(def)(?P<Name3>ghi)`
	reg = regexp.MustCompile(pat)

	for i := 0; i <= reg.NumSubexp(); i++ {
		fmt.Printf("%d: %q\n", i, reg.SubexpNames()[i])
	}
	// 0: ""
	// 1: "Name1"
	// 2: ""
	// 3: "Name3"

	fmt.Println()

	// 获取字面前缀
	pat = `(abc1)(abc2)(abc3)`
	reg = regexp.MustCompile(pat)
	fmt.Println(reg.LiteralPrefix()) // abc1abc2abc3 true

	pat = `(abc1)|(abc2)|(abc3)`
	reg = regexp.MustCompile(pat)
	fmt.Println(reg.LiteralPrefix()) //  false

	pat = `abc1|abc2|abc3`
	reg = regexp.MustCompile(pat)
	fmt.Println(reg.LiteralPrefix()) // abc false
}

------------------------------

// 判断

// 判断在 b（s、r）中能否找到匹配的字符串
func (re *Regexp) Match(b []byte) bool
func (re *Regexp) MatchString(s string) bool
func (re *Regexp) MatchReader(r io.RuneReader) bool

------------------------------

// 查找

// 返回第一个匹配到的结果（结果以 b 的切片形式返回）。
func (re *Regexp) Find(b []byte) []byte

// 返回第一个匹配到的结果及其分组内容（结果以 b 的切片形式返回）。
// 返回值中的第 0 个元素是整个正则表达式的匹配结果，后续元素是各个分组的
// 匹配内容，分组顺序按照“(”的出现次序而定。
func (re *Regexp) FindSubmatch(b []byte) [][]byte

// 功能同 Find，只不过返回的是匹配结果的首尾下标，通过这些下标可以生成切片。
// loc[0] 是结果切片的起始下标，loc[1] 是结果切片的结束下标。
func (re *Regexp) FindIndex(b []byte) (loc []int)

// 功能同 FindSubmatch，只不过返回的是匹配结果的首尾下标，通过这些下标可以生成切片。
// loc[0] 是结果切片的起始下标，loc[1] 是结果切片的结束下标。
// loc[2] 是分组1切片的起始下标，loc[3] 是分组1切片的结束下标。
// loc[4] 是分组2切片的起始下标，loc[5] 是分组2切片的结束下标。
// 以此类推
func (re *Regexp) FindSubmatchIndex(b []byte) (loc []int)

------------------------------

// 示例：Find、FindSubmatch
func main() {
	pat := `(((abc.)def.)ghi)`
	reg := regexp.MustCompile(pat)

	src := []byte(`abc-def-ghi abc+def+ghi`)

	// 查找第一个匹配结果
	fmt.Printf("%s\n", reg.Find(src)) // abc-def-ghi

	fmt.Println()

	// 查找第一个匹配结果及其分组字符串
	first := reg.FindSubmatch(src)
	for i := 0; i < len(first); i++ {
		fmt.Printf("%d: %s\n", i, first[i])
	}
	// 0: abc-def-ghi
	// 1: abc-def-ghi
	// 2: abc-def-
	// 3: abc-
}

------------------------------

// 示例：FindIndex、FindSubmatchIndex
func main() {
	pat := `(((abc.)def.)ghi)`
	reg := regexp.MustCompile(pat)

	src := []byte(`abc-def-ghi abc+def+ghi`)

	// 查找第一个匹配结果
	matched := reg.FindIndex(src)
	fmt.Printf("%v\n", matched) // [0 11]
	m := matched[0]
	n := matched[1]
	fmt.Printf("%s\n\n", src[m:n]) // abc-def-ghi

	// 查找第一个匹配结果及其分组字符串
	matched = reg.FindSubmatchIndex(src)
	fmt.Printf("%v\n", matched) // [0 11 0 11 0 8 0 4]
	for i := 0; i < len(matched)/2; i++ {
		m := matched[i*2]
		n := matched[i*2+1]
		fmt.Printf("%s\n", src[m:n])
	}
	// abc-def-ghi
	// abc-def-ghi
	// abc-def-
	// abc-
}

------------------------------

// 功能同上，只不过返回多个匹配的结果，而不只是第一个。
// n 是查找次数，负数表示不限次数。
func (re *Regexp) FindAll(b []byte, n int) [][]byte
func (re *Regexp) FindAllSubmatch(b []byte, n int) [][][]byte

func (re *Regexp) FindAllIndex(b []byte, n int) [][]int
func (re *Regexp) FindAllSubmatchIndex(b []byte, n int) [][]int

------------------------------

// 示例：FindAll、FindAllSubmatch
func main() {
	pat := `(((abc.)def.)ghi)`
	reg := regexp.MustCompile(pat)

	s := []byte(`abc-def-ghi abc+def+ghi`)

	// 查找所有匹配结果
	for _, one := range reg.FindAll(s, -1) {
		fmt.Printf("%s\n", one)
	}
	// abc-def-ghi
	// abc+def+ghi

	// 查找所有匹配结果及其分组字符串
	all := reg.FindAllSubmatch(s, -1)
	for i := 0; i < len(all); i++ {
		fmt.Println()
		one := all[i]
		for i := 0; i < len(one); i++ {
			fmt.Printf("%d: %s\n", i, one[i])
		}
	}
	// 0: abc-def-ghi
	// 1: abc-def-ghi
	// 2: abc-def-
	// 3: abc-

	// 0: abc+def+ghi
	// 1: abc+def+ghi
	// 2: abc+def+
	// 3: abc+
}

------------------------------

// 功能同上，只不过在字符串中查找
func (re *Regexp) FindString(s string) string
func (re *Regexp) FindStringSubmatch(s string) []string

func (re *Regexp) FindStringIndex(s string) (loc []int)
func (re *Regexp) FindStringSubmatchIndex(s string) []int

func (re *Regexp) FindAllString(s string, n int) []string
func (re *Regexp) FindAllStringSubmatch(s string, n int) [][]string

func (re *Regexp) FindAllStringIndex(s string, n int) [][]int
func (re *Regexp) FindAllStringSubmatchIndex(s string, n int) [][]int

// 功能同上，只不过在 io.RuneReader 中查找。
func (re *Regexp) FindReaderIndex(r io.RuneReader) (loc []int)
func (re *Regexp) FindReaderSubmatchIndex(r io.RuneReader) []int

------------------------------

// 替换（不会修改参数，结果是参数的副本）

// 将 src 中匹配的内容替换为 repl（repl 中可以使用 $1 $name 等分组引用符）。
func (re *Regexp) ReplaceAll(src, repl []byte) []byte

// 将 src 中匹配的内容经过 repl 函数处理后替换回去。
func (re *Regexp) ReplaceAllFunc(src []byte, repl func([]byte) []byte) []byte

// 将 src 中匹配的内容替换为 repl（repl 为字面值，不解析其中的 $1 $name 等）。
func (re *Regexp) ReplaceAllLiteral(src, repl []byte) []byte

// 功能同上，只不过在字符串中查找。
func (re *Regexp) ReplaceAllString(src, repl string) string
func (re *Regexp) ReplaceAllStringFunc(src string, repl func(string) string) string
func (re *Regexp) ReplaceAllLiteralString(src, repl string) string

// Expand 要配合 FindSubmatchIndex 一起使用。FindSubmatchIndex 在 src 中进行
// 查找，将结果存入 match 中。这样就可以通过 src 和 match 得到匹配的字符串。
// template 是替换内容，可以使用分组引用符 $1、$2、$name 等。Expane 将其中的分
// 组引用符替换为前面匹配到的字符串。然后追加到 dst 的尾部（dst 可以为空）。
// 说白了 Expand 就是一次替换过程，只不过需要 FindSubmatchIndex 的配合。
func (re *Regexp) Expand(dst []byte, template []byte, src []byte, match []int) []byte

// 功能同上，参数为字符串。
func (re *Regexp) ExpandString(dst []byte, template string, src string, match []int) []byte

------------------------------

// 示例：Expand
func main() {
	pat := `(((abc.)def.)ghi)`
	reg := regexp.MustCompile(pat)

	src := []byte(`abc-def-ghi abc+def+ghi`)
	template := []byte(`$0   $1   $2   $3`)

	// 替换第一次匹配结果
	match := reg.FindSubmatchIndex(src)
	fmt.Printf("%v\n", match) // [0 11 0 11 0 8 0 4]
	dst := reg.Expand(nil, template, src, match)
	fmt.Printf("%s\n\n", dst)
	// abc-def-ghi   abc-def-ghi   abc-def-   abc-

	// 替换所有匹配结果
	for _, match := range reg.FindAllSubmatchIndex(src, -1) {
		fmt.Printf("%v\n", match)
		dst := reg.Expand(nil, template, src, match)
		fmt.Printf("%s\n", dst)
	}
	// [0 11 0 11 0 8 0 4]
	// abc-def-ghi   abc-def-ghi   abc-def-   abc-
	// [12 23 12 23 12 20 12 16]
	// abc+def+ghi   abc+def+ghi   abc+def+   abc+
}

------------------------------

// 其它

// 以 s 中的匹配结果作为分割符将 s 分割成字符串列表。
// n 是分割次数，负数表示不限次数。
func (re *Regexp) Split(s string, n int) []string

// 将当前正则对象复制一份。在多例程中使用同一正则对象时，给每个例程分配一个
// 正则对象的副本，可以避免多例程对单个正则对象的争夺锁定。
func (re *Regexp) Copy() *Regexp

