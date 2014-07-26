- **add** *func(args ...interface{}) (interface{}, error)*

	Add all the arguments, returning either an int (if the result does not have a decimal part) or a float64.

- **addf** *func(args ...interface{}) (float64, error)*

	Add all the arguments, returning a float64.

- **addi** *func(args ...interface{}) (int, error)*

	Add all the arguments, returning a int. If the result has a decimal part, it's truncated.

- **and** *func(args ...interface{}) interface{}*

	Return the last argument of the given ones if all of them are true. Otherwise, return the first non-true one.

- **concat** *func(args ...interface{}) string*

	Return the result of concatenating all the arguments.

- **divisible** *func(n interface{}, d interface{}) (bool, error)*

	Returns true if the first argument is divisible by the second one.

- **eq** *func(args ...interface{}) bool*

	Returns true iff the first argument is equal to any of the following ones.

- **even** *func(arg interface{}) (bool, error)*

	An alias for divisible(arg, 2)

- **float** *func(val interface{}) (float64, error)*

	float is an alias for [gnd.la/util/types.ToFloat](/doc/pkg/gnd.la/util/types#func-ToFloat)

- **gt** *func(arg1 interface{}, arg2 interface{}) (bool, error)*

	Returns true iff arg1 > arg2. Produces an error if arguments are of different types of if its type is not comparable.

- **gte** *func(arg1 interface{}, arg2 interface{}) (bool, error)*

	Returns true iff arg1 >= arg2. Produces an error if arguments are of different types of if its type is not comparable.

- **html** *func(args ...interface{}) string*

	html is an alias for [html/template.HTMLEscaper](/doc/pkg/html/template#func-HTMLEscaper)

- **index** *func(item interface{}, indices ...interface{}) (interface{}, error)*

	Return the result of indexing into the first argument, which must be a map or slice, using the second one (i.e. item[idx]).

- **int** *func(val interface{}) (int, error)*

	int is an alias for [gnd.la/util/types.ToInt](/doc/pkg/gnd.la/util/types#func-ToInt)

- **join** *func(a []string, sep string) string*

	join is an alias for [strings.Join](/doc/pkg/strings#func-Join)

- **js** *func(args ...interface{}) string*

	js is an alias for [html/template.JSEscaper](/doc/pkg/html/template#func-JSEscaper)

- **json** *func(arg interface{}) (html/template.JS, error)*

	Returns the JSON representation of the given argument as a string. Produces an error in the argument can't be converted to JSON.

- **jsons** *func(arg interface{}) (string, error)*

	Same as jsons, but returns a template.JS, which can be embedded in script sections of an HTML template without further escaping.

- **len** *func(item interface{}) (int, error)*

	Return the length of the argument, which must be map, slice or array.

- **lt** *func(arg1 interface{}, arg2 interface{}) (bool, error)*

	Returns true iff arg1 < arg2. Produces an error if arguments are of different types of if its type is not comparable.

- **lte** *func(arg1 interface{}, arg2 interface{}) (bool, error)*

	Returns true iff arg1 <= arg2. Produces an error if arguments are of different types of if its type is not comparable.

- **mul** *func(args ...interface{}) (interface{}, error)*

	Multiply all the arguments, returning either an int (if the result does not have a decimal part) or a float64.

- **mulf** *func(args ...interface{}) (float64, error)*

	Multiply all the arguments, returning a float64.

- **muli** *func(args ...interface{}) (int, error)*

	Multiply all the arguments, returning a int. If the result has a decimal part, it's truncated.

- **neq** *func(args ...interface{}) bool*

	Returns true iff the first argument is different to all the following ones.

- **not** *func(arg interface{}) bool*

	Return the negation of the truth value of the given argument.

- **now** *func() time.Time*

	Return the current time.Time in the local timezone.

- **odd** *func(arg interface{}) (bool, error)*

	An alias for not divisible(arg, 2)

- **or** *func(args ...interface{}) interface{}*

	Return the first true argument of the given ones. If none of them is true, return false.

- **print** *func(a ...interface{}) string*

	print is an alias for [fmt.Sprint](/doc/pkg/fmt#func-Sprint)

- **printf** *func(format string, a ...interface{}) string*

	printf is an alias for [fmt.Sprintf](/doc/pkg/fmt#func-Sprintf)

- **println** *func(a ...interface{}) string*

	println is an alias for [fmt.Sprintln](/doc/pkg/fmt#func-Sprintln)

- **slice** *func(args ...interface{}) *[]interface{}*

	Returns a slice with the given arguments.

- **split** *func(s string, sep string) []string*

	split is an alias for [strings.Split](/doc/pkg/strings#func-Split)

- **split_n** *func(s string, sep string, n int) []string*

	split_n is an alias for [strings.SplitN](/doc/pkg/strings#func-SplitN)

- **sub** *func(args ...interface{}) (interface{}, error)*

	Substract all the arguments in the given order, from left to right, returning either an int (if the result does not have a decimal part) or a float64.

- **subf** *func(args ...interface{}) (float64, error)*

	Substract all the arguments in the given order, from left to right, returning a float64.

- **subi** *func(args ...interface{}) (int, error)*

	Substract all the arguments in the given order, from left to right, returning a int. If the result has a decimal part, it's truncated.

- **to_html** *func(s string) html/template.HTML*

	Converts plain text to HTML by escaping it and replacing newlines with <br> tags.

- **to_lower** *func(s string) string*

	to_lower is an alias for [strings.ToLower](/doc/pkg/strings#func-ToLower)

- **to_title** *func(s string) string*

	to_title is an alias for [strings.ToTitle](/doc/pkg/strings#func-ToTitle)

- **to_upper** *func(s string) string*

	to_upper is an alias for [strings.ToUpper](/doc/pkg/strings#func-ToUpper)

- **urlquery** *func(args ...interface{}) string*

	urlquery is an alias for [html/template.URLQueryEscaper](/doc/pkg/html/template#func-URLQueryEscaper)

- **var** *func(s *gnd.la/template.State, name string) interface{}*

	Return the value of the given variable, or an empty value if no such variable exists.




[title] = Template Function Reference
[synopsis] = All the standard functions available in Gondola templates.
[updated] = 2014-07-26 17:23:08
[updated] = 2014-07-26 17:26:36
[priority] = 2
